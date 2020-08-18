package main

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/argParser"
	"github.com/ElCap1tan/gort/scanUtils"
	"os"
)

type multiScanResult struct {
	resolved   []*scanUtils.ScanResult
	unresolved []*scanUtils.TargetInfo
}

func printMultiScanResult(multiScanRes multiScanResult) {
	fmt.Println("--------------- SCAN RESULTS ---------------")
	for _, scanResult := range multiScanRes.resolved {
		fmt.Println(scanResult.String())
	}
}

func main() {
	args := os.Args

	hostArgs := args[1]
	portArgs := args[2]

	hosts := argParser.ParseHostArgs(hostArgs)
	multiScanRes := multiScanResult{}
	ports := argParser.ParsePortArgs(portArgs)
	out := make(chan *scanUtils.ScanResult)

	for _, h := range hosts {
		if h.IPAddr == nil {
			multiScanRes.unresolved = append(multiScanRes.unresolved, h)
		} else {
			go scanUtils.NewFCPortScanner(h, ports).Scan(out)
		}
	}

	for i := 0; i < len(hosts)-len(multiScanRes.unresolved); i++ {
		multiScanRes.resolved = append(multiScanRes.resolved, <-out)
	}
	printMultiScanResult(multiScanRes)
}
