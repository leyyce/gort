package main

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/argParser"
	"github.com/ElCap1tan/gort/internal/csvParser"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	var hostArgs string
	var portArgs string

	args := os.Args

	// Try to update port-numbers.xml and port_open_freq.csv if necessary
	err := updateKnownPorts(5)
	if err != nil {
		println("[!] Error while updating list of known ports. Using old list...")
	}
	err = updatePortFreq(5)
	if err != nil {
		println("[!] Error while updating list of most common open ports. Using old list...")
	}

	if len(args) == 2 {
		fmt.Println("[!] No port arguments provided assuming 100 most common open ports...")
		hostArgs = args[1]
		mostScanned := csvParser.NewMostScannedPorts()
		var i int
		for _, mS := range *mostScanned {
			if i == 99 {
				break
			}
			if mS.Protocol == "tcp" {
				portArgs += strconv.Itoa(int(mS.Number)) + ","
				i++
			}
		}
		portArgs = portArgs[:len(portArgs)-2]
	} else if len(args) == 3 {
		hostArgs = args[1]
		portArgs = args[2]
	} else {
		fmt.Printf("[ERROR] Called with the wrong number of arguments: Provided|Needed [%d]|[1-2]", len(args)-1)
		return
	}

	targets := argParser.ParseHostArgs(hostArgs, argParser.ParsePortArgs(portArgs, "tcp"))
	fmt.Println("[!] STARTING SCAN")
	multiScanRes := targets.Scan()
	tFinished := time.Now()
	fmt.Println(multiScanRes.String())
	file, err := os.Create(fmt.Sprintf(
		"scanlog_%d-%02d-%02d_%02d-%02d-%02d.txt",
		tFinished.Year(), tFinished.Month(), tFinished.Day(),
		tFinished.Hour(), tFinished.Minute(), tFinished.Second()))
	if err == nil {
		_, err = file.WriteString(multiScanRes.String())
		if err != nil {
			println(err.Error())
		}
	}
}

func updateKnownPorts(maxAgeDays int) error {
	pnPath := "data/port-numbers.xml"
	url := "https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xml"
	pnStats, err := os.Stat(pnPath)
	if err != nil {
		fmt.Printf("[!] Known ports list under %s not found. Trying to download it from %s...\n", pnPath, url)
		err = fetchFile(url, pnPath)
		return err
	}
	lastModTime := pnStats.ModTime()
	if lastModTime.Add(time.Hour * 24 * time.Duration(maxAgeDays)).Before(time.Now()) {
		fmt.Printf("[!] List of known ports not updated since %d days. Trying to update now...\n", maxAgeDays)
		err = fetchFile(url, pnPath)
		return err
	}
	return err
}

func updatePortFreq(maxAgeDays int) error {
	pfPath := "data/port_open_freq.csv"
	url := "https://docs.google.com/spreadsheets/d/1r_IriqmkTNPSTiUwii_hQ8Gwl2tfTUz8AGIOIL-wMIE/export?format=csv"
	pnStats, err := os.Stat(pfPath)
	if err != nil {
		fmt.Printf("[!] Most common open ports list under %s not found. Trying to download it from %s...\n", pfPath, url)
		err = fetchFile(url, pfPath)
		return err
	}
	lastModTime := pnStats.ModTime()
	if lastModTime.Add(time.Hour * 24 * time.Duration(maxAgeDays)).Before(time.Now()) {
		fmt.Printf("[!] List of most common open ports not updated since %d days. Trying to update now...\n", maxAgeDays)
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
		fmt.Println("[🗸] Success!")
	}
	return err
}
