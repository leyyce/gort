package netUtil

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/internal/xmlParser"
	"strconv"
	"strings"
)

type Ports []*Port

type Port struct {
	PortNo      uint16
	Protocol    string
	Service     string
	Description string
}

func NewPort(portNo uint16, proto string, service string, desc string) *Port {
	return &Port{PortNo: portNo, Protocol: proto, Service: service, Description: desc}
}

func ParsePortString(portArgs string, proto string) Ports {
	var tgtPorts Ports
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
									tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, p.Service, p.Description))
									found = true
									break
								} else if strings.Contains(p.Number, "-") {
									pRange := helper.StrRangeToArray(p.Number)
									for _, n := range pRange {
										if n == rawPort {
											tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, p.Service, p.Description))
											found = true
											break
										}
									}
								}
							}
						}
						if !found {
							tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, "N/A",
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
									tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, p.Service, p.Description))
									found = true
									break
								} else if strings.Contains(p.Number, "-") {
									pRange := helper.StrRangeToArray(p.Number)
									for _, n := range pRange {
										if uint64(n) == rawPort {
											tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, p.Service, p.Description))
											found = true
											break
										}
									}
								}
							}
						}
						if !found {
							tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, "N/A",
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
					tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, "N/A",
						"Make sure you have provided the port-numbers.xml file"))
				}
			}
		} else {
			if helper.ValidatePort(portArg) {
				p, _ := strconv.ParseUint(portArg, 10, 16)
				tgtPorts = append(tgtPorts, NewPort(uint16(p), proto, "N/A",
					"Make sure you have provided the port-numbers.xml file"))
			}
		}
	}
	return tgtPorts
}

func (p *Port) String() string {
	if p.Service == "" || p.Service == "N/A" {
		return fmt.Sprintf("%5d/%s", p.PortNo, p.Protocol)
	}
	return fmt.Sprintf("%5d/%s [%s]", p.PortNo, p.Protocol, p.Service)
}

func (ps Ports) String() string {
	ret := ""
	for i, p := range ps {
		ret += p.String() + ", "
		if i != 0 && i%10 == 0 && i+1 != len(ps) {
			ret += "\n"
		}
	}
	return ret[:len(ret)-2]
}

func (ps Ports) Preview() string {
	maxPerLine := 30
	ret := ""
	if len(ps) < maxPerLine {
		for i, p := range ps {
			ret += p.String() + ", "
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
		}
		ret = ret[:len(ret)-2]
	} else {
		for i, p := range ps[:maxPerLine] {
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
			ret += p.String() + ", "
		}
		ret = ret[:len(ret)-2] + "..."
	}
	return ret
}
