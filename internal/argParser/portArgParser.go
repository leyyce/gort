package argParser

import (
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/netUtil"
	"strconv"
	"strings"
)

func ParsePortArgs(portArgs string, proto string) netUtil.Ports {
	var tgtPorts netUtil.Ports
	ports := strings.Split(portArgs, ",")
	for _, portArg := range ports {
		if strings.Contains(portArg, "-") {
			rawPorts := helper.StrRangeToArray(portArg)
			for _, rawPort := range rawPorts {
				if helper.ValidatePort(strconv.Itoa(rawPort)) {
					tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto))
				}
			}
		} else {
			if helper.ValidatePort(portArg) {
				p, _ := strconv.ParseUint(portArg, 10, 16)
				tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(p), proto))
			}
		}
	}
	return tgtPorts
}
