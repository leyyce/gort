package argParser

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/scanUtils"
	"strconv"
	"strings"
)

func ParseHostArgs(hostArgs string) []*scanUtils.TargetInfo {
	var tgtHosts []*scanUtils.TargetInfo
	hosts := strings.Split(hostArgs, ",")
	for _, hostArg := range hosts {
		if helper.ValidateIPOrRange(hostArg) {
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
					tgtHosts = append(tgtHosts, scanUtils.NewHostInfo(t))
				}
			} else {
				tgtHosts = append(tgtHosts, scanUtils.NewHostInfo(hostArg))
			}
		} else {
			tgtHosts = append(tgtHosts, scanUtils.NewHostInfo(hostArg))
		}
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
