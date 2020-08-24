package main

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/argParser"
	"github.com/ElCap1tan/gort/internal/csvParser"
	"os"
	"strconv"
	"time"
)

func main() {
	var hostArgs string
	var portArgs string

	args := os.Args

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
