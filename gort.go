package main

import (
	"flag"
	"fmt"
	"github.com/ElCap1tan/gort/internal/colorFmt"
	"github.com/ElCap1tan/gort/internal/csvParser"
	"github.com/ElCap1tan/gort/internal/symbols"
	"github.com/ElCap1tan/gort/netUtil"
	"github.com/ElCap1tan/gort/netUtil/pScan"
	"github.com/fatih/color"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

const dataFolder, resultFolder = "data", "scans"

func main() {
	var usage = "" +
		"Usage:\n" +
		"\tgort [-p ports] [-mc count] [-closed] [-online] [-file] [-elevated] hosts\n" +
		"\tMandatory argument:\n" +
		"\thosts are comma separated values that can either be\n" +
		"\t\tA single host : 192.88.99.1 or example.com\n" +
		"\t\tA range of hosts : 192.88.99.1-50 or 192.88.99-100.1-50\n" +
		"\t\tA CIDR formatted host range : 192.88.99.1/24\n" +
		"\tOptional arguments\n" +
		"\t\t-p\n" +
		"\t\t\tports are comma separated values that either can be\n" +
		"\t\t\t\tA single port : 80\n" +
		"\t\t\t\tA range of ports : 80-100\n" +
		"\t\t-mc [int]\n" +
		"\t\t\tSets the number of most common open ports to scan. If omitted defaults to 1000.\n" +
		"\t\t-closed\n" +
		"\t\t\tIf this flag is passed ports with closed and unknown/filtered state are also shown in the console output.\n" +
		"\t\t-online\n" +
		"\t\t\tIf this flag is passed only hosts confirmed as online are shown in the console output.\n" +
		"\t\t-file\n" +
		"\t\t\tIf this flag is passed the scan result will be saved to a file.\n" +
		"\t\t-elevated\n" +
		"\t\t\tOnly important for Linux. If this flag is passed the ICMP echo requests will be send via raw sockets.\n" +
		"\t\t\tYou might want to try in unprivileged mode first.\n" +
		"\t\t\tImportant: Must be run as a super-user when this flag is used or else ping tests won't work!\n" +
		"Examples:\n" +
		"\t# scan the 1000 most common open ports of example.com\n" +
		"\t\tgort example.com\n" +
		"\t# scan the 500 most common open ports of example.com\n" +
		"\t\tgort -mc 500 example.com\n" +
		"\t# scan a custom list of ports for example.com and also show closed or unknown ports in result\n" +
		"\t\tgort -p 80,443,1000-1024 -closed example.com\n" +
		"\t# scan the subnet 192.88.99.0/24 for the 100 most common open ports and and a custom list of ports\n" +
		"\t# and only show targets confirmed as online in the scan result (Some ports could be scanned double).\n" +
		"\t\tgort -mc 100 -p 10334,12012 -online 192.88.99.0/24\n" +
		"\t\tor\n" +
		"\t\tgort -mc 100 -p 10334,12012 -online 192.88.99.0-255\n"
	flag.Usage = func() {
		fmt.Printf(usage)
	}
	var hostArgs string
	mostCommonCount := flag.Int("mc", 1000, "")
	portArgs := flag.String("p", "", "")
	onlineOnly := flag.Bool("online", false, "")
	showClosed := flag.Bool("closed", false, "")
	writeFile := flag.Bool("file", false, "")
	privileged := flag.Bool("elevated", false, "")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	hostArgs = flag.Arg(0)

	// Try to update port-numbers.xml and port_open_freq.csv if necessary
	err := ensureDir(dataFolder)
	if err != nil {
		colorFmt.Fatalf("%s Error creating data dir '%s': %s", symbols.FAILURE, dataFolder, err.Error())
		return
	}
	err = updateKnownPorts(5)
	if err != nil {
		colorFmt.Warnf("%s Error while updating list of known ports. Using old list...\n", symbols.INFO)
	}
	err = updatePortFreq(5)
	if err != nil {
		colorFmt.Warnf("%s Error while updating list of most common open ports. Using old list...\n", symbols.INFO)
	}

	if *portArgs == "" || *portArgs != "" && *mostCommonCount != 1000 {
		mostCommon := csvParser.NewMostCommonPorts()
		maxAvailable := 0
		for _, p := range *mostCommon {
			if p.Protocol == "tcp" {
				maxAvailable++
			}
		}

		if *mostCommonCount > maxAvailable {
			colorFmt.Infof("%s Can't start scan for the %d most common open ports because stats are only available for %d ports. Using that number of ports instead...\n",
				symbols.INFO, *mostCommonCount, maxAvailable)
			*mostCommonCount = maxAvailable
		}
		if *portArgs == "" {
			colorFmt.Infof("%s No port arguments provided assuming %d most common open ports...\n", symbols.INFO, *mostCommonCount)
		} else {
			colorFmt.Infof("%s Adding the %d most common open ports to the provided list of port arguments...\n", symbols.INFO, *mostCommonCount)
		}
		p := *portArgs + mostCommon.GetMostCommonString(*mostCommonCount)
		portArgs = &p
	}

	colorFmt.Infof("%s Parsing and resolving host arguments...\n", symbols.INFO)
	targets := pScan.ParseHostString(hostArgs, netUtil.ParsePortString(*portArgs, "tcp"), *privileged)
	colorFmt.Infof("%s STARTING SCAN...\n", symbols.INFO)
	multiScanRes := targets.Scan()
	tFinished := time.Now()
	if runtime.GOOS == "windows" {
		_, err = color.Output.Write([]byte(multiScanRes.CustomColorString(*onlineOnly, *showClosed)))
		if err != nil {
			colorFmt.Infof("Error writing colored scan result to the console. Trying uncolored...")
			fmt.Println(multiScanRes.String())
		}
	} else {
		fmt.Println(multiScanRes.CustomColorString(*onlineOnly, *showClosed))
	}

	if *writeFile {
		fileName := fmt.Sprintf("scanlog_%d-%02d-%02d_%02d-%02d-%02d.txt",
			tFinished.Year(), tFinished.Month(), tFinished.Day(),
			tFinished.Hour(), tFinished.Minute(), tFinished.Second())
		err = ensureDir(resultFolder)
		if err != nil {
			colorFmt.Warnf("%s Failed to create results dir under '%s'. Trying to save in the current working directory.", symbols.INFO, resultFolder)
			file, err := os.Create(fileName)
			if err == nil {
				_, err = file.WriteString(multiScanRes.String())
				if err != nil {
					println(err.Error())
				} else {
					colorFmt.Infof("%s Scan result saved as '%s'\n\n", symbols.INFO, fileName)
				}
			} else {
				println(err.Error())
			}
		} else {
			filePath := path.Join(resultFolder, fileName)
			file, err := os.Create(filePath)
			if err == nil {
				_, err = file.WriteString(multiScanRes.String())
				if err != nil {
					println(err.Error())
				} else {
					colorFmt.Infof("%s Scan result saved as '%s'\n\n", symbols.INFO, filePath)
				}
			} else {
				println(err.Error())
			}
		}
	}
}

func updateKnownPorts(maxAgeDays int) error {
	pnPath := path.Join(dataFolder, "port-numbers.xml")
	url := "https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xml"
	pnStats, err := os.Stat(pnPath)
	if err != nil {
		colorFmt.Warnf("%s Known ports list under %s not found. Trying to download it from %s...\n", symbols.INFO, pnPath, url)
		err = fetchFile(url, pnPath)
		return err
	}
	lastModTime := pnStats.ModTime()
	if lastModTime.Add(time.Hour * 24 * time.Duration(maxAgeDays)).Before(time.Now()) {
		colorFmt.Infof("%s List of known ports not updated since %d days. Trying to update now...\n", symbols.INFO, maxAgeDays)
		err = fetchFile(url, pnPath)
		return err
	}
	return err
}

func updatePortFreq(maxAgeDays int) error {
	pfPath := path.Join(dataFolder, "port_open_freq.csv")
	url := "https://docs.google.com/spreadsheets/d/1r_IriqmkTNPSTiUwii_hQ8Gwl2tfTUz8AGIOIL-wMIE/export?format=csv"
	pnStats, err := os.Stat(pfPath)
	if err != nil {
		colorFmt.Warnf("%s Most common open ports list under %s not found. Trying to download it from %s...\n", symbols.INFO, pfPath, url)
		err = fetchFile(url, pfPath)
		return err
	}
	lastModTime := pnStats.ModTime()
	if lastModTime.Add(time.Hour * 24 * time.Duration(maxAgeDays)).Before(time.Now()) {
		colorFmt.Infof("%s List of most common open ports not updated since %d days. Trying to update now...\n", symbols.INFO, maxAgeDays)
		err = fetchFile(url, pfPath)
		return err
	}
	return err
}

func fetchFile(url, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, err = io.Copy(file, resp.Body)
	if err == nil {
		colorFmt.Successf("%s Success!\n", symbols.SUCCESS)
	}
	return err
}

func ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, 0766)
	if err == nil || os.IsExist(err) {
		return nil
	}
	return err
}
