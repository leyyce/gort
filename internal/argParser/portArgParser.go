package argParser

import (
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/internal/xmlParser"
	"github.com/ElCap1tan/gort/netUtil"
	"strconv"
	"strings"
)

func ParsePortArgs(portArgs string, proto string) netUtil.Ports {
	var tgtPorts netUtil.Ports
	portRegistry, err := xmlParser.NewPortRegistry()
	if err == nil {
		ports := strings.Split(portArgs, ",")
		for _, portArg := range ports {
			if strings.Contains(portArg, "-") {
				rawPorts := helper.StrRangeToArray(portArg)
				for _, rawPort := range rawPorts {
					if helper.ValidatePort(strconv.Itoa(rawPort)) {
						var found bool
						for _, p := range portRegistry.Records {
							found = false
							if p.Protocol == proto {
								pNum, err := strconv.Atoi(p.Number)
								if err == nil && pNum == rawPort {
									tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, p.Service, p.Description))
									found = true
									break
								} else if strings.Contains(p.Number, "-") {
									pRange := helper.StrRangeToArray(p.Number)
									for _, n := range pRange {
										if n == rawPort {
											tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, p.Service, p.Description))
											found = true
											break
										}
									}
								}
							}
						}
						if !found {
							tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, "N/A",
								"No description available"))
						}
					}
				}
			} else {
				if helper.ValidatePort(portArg) {
					rawPort, err := strconv.ParseUint(portArg, 10, 16)
					if err == nil {
						var found bool
						for _, p := range portRegistry.Records {
							found = false
							if p.Protocol == proto {
								pNum, err := strconv.Atoi(p.Number)
								if err == nil && uint64(pNum) == rawPort {
									tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, p.Service, p.Description))
									found = true
									break
								} else if strings.Contains(p.Number, "-") {
									pRange := helper.StrRangeToArray(p.Number)
									for _, n := range pRange {
										if uint64(n) == rawPort {
											tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, p.Service, p.Description))
											found = true
											break
										}
									}
								}
							}
						}
						if !found {
							tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, "N/A",
								"No description available"))
						}
					}
				}
			}
		}
		return tgtPorts
	}
	ports := strings.Split(portArgs, ",")
	for _, portArg := range ports {
		if strings.Contains(portArg, "-") {
			rawPorts := helper.StrRangeToArray(portArg)
			for _, rawPort := range rawPorts {
				if helper.ValidatePort(strconv.Itoa(rawPort)) {
					tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(rawPort), proto, "N/A",
						"Make sure you have provided the port-numbers.xml file"))
				}
			}
		} else {
			if helper.ValidatePort(portArg) {
				p, _ := strconv.ParseUint(portArg, 10, 16)
				tgtPorts = append(tgtPorts, netUtil.NewPort(uint16(p), proto, "N/A",
					"Make sure you have provided the port-numbers.xml file"))
			}
		}
	}
	return tgtPorts
}
