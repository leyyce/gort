package argParser

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/internal/helper/ulimit"
	"github.com/ElCap1tan/gort/netUtil"
	"github.com/ElCap1tan/gort/netUtil/pScan"
	"golang.org/x/sync/semaphore"
	"net"
	"strconv"
	"strings"
)

func ParseHostArgs(hostArgs string, ports netUtil.Ports) pScan.Targets {
	var tgtHosts pScan.Targets
	hostCount := 0
	out := make(chan *pScan.Target)

	var limit int64
	l, err := ulimit.GetUlimit()
	if err != nil {
		limit = 1024
	} else {
		limit = int64(l)
	}

	lock := semaphore.NewWeighted(limit)

	hosts := strings.Split(hostArgs, ",")
	for _, hostArg := range hosts {
		if ip, ipNet, err := net.ParseCIDR(hostArg); err == nil {
			for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); helper.IncIp(ip) {
				go pScan.AsyncNewTarget(ip.String(), ports, out, lock)
				hostCount++
			}
		} else if helper.ValidateIPOrRange(hostArg) {
			if strings.Contains(hostArg, "-") {
				addrParts := strings.Split(hostArg, ".")
				var octets [4][]int
				for i, addrPart := range addrParts {
					if strings.Contains(addrPart, "-") {
						octets[i] = append(octets[i], helper.StrRangeToArray(addrPart)...)
					} else {
						p, _ := strconv.Atoi(addrPart)
						octets[i] = append(octets[i], p)
					}
				}
				for _, t := range octetsToTargets(octets) {
					go pScan.AsyncNewTarget(t, ports, out, lock)
					hostCount++
				}
			} else {
				go pScan.AsyncNewTarget(hostArg, ports, out, lock)
				hostCount++
			}
		} else {
			go pScan.AsyncNewTarget(hostArg, ports, out, lock)
			hostCount++
		}
	}
	for i := 1; i <= hostCount; i++ {
		tgtHosts = append(tgtHosts, <-out)
	}
	return tgtHosts
}

func octetsToTargets(octets [4][]int) []string {
	var targets []string
	for _, oc0 := range octets[0] {
		for _, oc1 := range octets[1] {
			for _, oc2 := range octets[2] {
				for _, oc3 := range octets[3] {
					targets = append(targets, fmt.Sprintf("%d.%d.%d.%d", oc0, oc1, oc2, oc3))
				}
			}
		}
	}
	return targets
}
