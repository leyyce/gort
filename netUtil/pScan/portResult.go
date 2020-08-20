package pScan

import (
	"fmt"
	"github.com/ElCap1tan/gort/netUtil"
)

type PortResults struct {
	Open     netUtil.Ports
	Closed   netUtil.Ports
	Filtered netUtil.Ports
}

func NewPortResults() *PortResults {
	return &PortResults{}
}

func (p *PortResults) String() string {
	ret := "" +
		"*************** PORT RESULT **********************\n" +
		"Open TCP Ports:\n"
	if len(p.Open) == 0 {
		ret += "\n"
	} else {
		for _, oP := range p.Open {
			ret += fmt.Sprintf("[+] %s\n", oP)
		}
	}
	ret += "Closed TCP Ports:\n"
	if len(p.Closed) == 0 {
		ret += "\n"
	} else {
		for _, cP := range p.Closed {
			ret += fmt.Sprintf("[-] %s\n", cP)
		}
	}
	ret += "**************************************************"
	return ret
}
