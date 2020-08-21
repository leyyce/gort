package main

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/argParser"
	"os"
	"time"
)

func main() {
	args := os.Args

	hostArgs := args[1]
	portArgs := args[2]

	targets := argParser.ParseHostArgs(hostArgs, argParser.ParsePortArgs(portArgs, "tcp"))
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
