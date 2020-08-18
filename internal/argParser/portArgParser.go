package argParser

import (
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/scanUtils"
	"strconv"
	"strings"
)

func ParsePortArgs(portArgs string) []scanUtils.Port {
	var tgtPorts []scanUtils.Port
	ports := strings.Split(portArgs, ",")
	for _, portArg := range ports {
		if strings.Contains(portArg, "-") {
			rawPorts := helper.StrRangeToArray(portArg)
			for _, rawPort := range rawPorts {
				if helper.ValidatePort(strconv.Itoa(rawPort)) {
					tgtPorts = append(tgtPorts, scanUtils.Port(rawPort))
				}
			}
		} else {
			if helper.ValidatePort(portArg) {
				p, _ := strconv.ParseUint(portArg, 10, 16)
				tgtPorts = append(tgtPorts, scanUtils.Port(p))
			}
		}
	}
	return tgtPorts
}
