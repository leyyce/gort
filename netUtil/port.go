// Copyright (c) 2020 Yannic Wehner
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// Package netUtil contains all the publicly exposed networking functionality of gort like the port scanning types and
// functionality provided in 'netUtil.pScan'.
package netUtil

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/internal/xmlParser"
	"strconv"
	"strings"
)

// Ports is a list of Port pointers.
type Ports []*Port

// Port represents a single target port.
type Port struct {
	// PortNo is the number of the target port.
	// A port number is a 16-bit unsigned integer and thus can range from 0 to 65535.
	PortNo uint16

	// Protocol is the transport protocol of the port (either TCP or UDP).
	Protocol string

	// Service is the default service running on a given port as specified by IANA.
	Service string

	// Description is a short description of the Service.
	Description string
}

// NewPort returns a pointer to a new instance of Port.
// portNo is the port number of the new Port and proto is the transport protocol used by it.
// service is the default service running on a given port as specified by IANA and
// desc is a short description of that service.
// To initialize a new port you usually should use ParsePortString as this method will try to lookup the
// service and description of the given ports automatically.
func NewPort(portNo uint16, proto string, service string, desc string) *Port {
	return &Port{PortNo: portNo, Protocol: proto, Service: service, Description: desc}
}

// ParsePortString parses ports and returns the newly initialized Ports.
//
// ports is comma separated list of values that can be in either of the following formats:
// - A single port: 23
// - A range of ports: 23-100
//
// The parameter proto specifies the transport protocol of the given ports and dataFolder is the path to the folder
// which contains the 'service-names-port-numbers.xml' file by IANA that is used for port service lookups. (The file can
// be found under https://www.iana.org/assignments/service-names-port-numbers/service-names-port-numbers.xml).
func ParsePortString(ports string, proto string, dataFolder string) Ports {
	var tgtPorts Ports
	portRegistry, err := xmlParser.NewPortRegistry(dataFolder)
	if err == nil {
		portList := strings.Split(ports, ",")
		for _, portArg := range portList {
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
	portList := strings.Split(ports, ",")
	for _, portArg := range portList {
		if strings.Contains(portArg, "-") {
			rawPorts := helper.StrRangeToArray(portArg)
			for _, rawPort := range rawPorts {
				if helper.ValidatePort(strconv.Itoa(rawPort)) {
					tgtPorts = append(tgtPorts, NewPort(uint16(rawPort), proto, "N/A",
						"Make sure you have provided the service-names-port-numbers.xml file"))
				}
			}
		} else {
			if helper.ValidatePort(portArg) {
				p, _ := strconv.ParseUint(portArg, 10, 16)
				tgtPorts = append(tgtPorts, NewPort(uint16(p), proto, "N/A",
					"Make sure you have provided the service-names-port-numbers.xml file"))
			}
		}
	}
	return tgtPorts
}

// String returns a string representation of the Port pointer.
func (p *Port) String() string {
	if p.Service == "" || p.Service == "N/A" {
		return fmt.Sprintf("%5d/%s", p.PortNo, p.Protocol)
	}
	return fmt.Sprintf("%5d/%s [%s]", p.PortNo, p.Protocol, p.Service)
}

// String returns a string representation of Ports.
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

// Preview returns a preview string representation of Ports where maxToDisplay is the maximum number of ports that
// should be shown.
func (ps Ports) Preview(maxToDisplay int) string {
	ret := ""
	if len(ps) < maxToDisplay {
		for i, p := range ps {
			ret += p.String() + ", "
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
		}
		ret = ret[:len(ret)-2]
	} else {
		for i, p := range ps[:maxToDisplay] {
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
			ret += p.String() + ", "
		}
		ret = ret[:len(ret)-2] + "..."
	}
	return ret
}
