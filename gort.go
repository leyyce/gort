package main

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/argParser"
	"os"
)

func main() {
	args := os.Args

	hostArgs := args[1]
	portArgs := args[2]

	targets := argParser.ParseHostArgs(hostArgs, argParser.ParsePortArgs(portArgs, "tcp"))
	multiScanRes := targets.Scan()
	fmt.Println(multiScanRes.String())
}
