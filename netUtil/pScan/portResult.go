package pScan

import (
	"fmt"
	"github.com/ElCap1tan/gort/netUtil"
	"strings"
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
	ret := "*************** PORT RESULT **********************\n"
	if len(p.Open) > 0 {
		ret += "Open TCP Ports:\n"
		for _, oP := range p.Open {
			ret += fmt.Sprintf("[+] %s - %s\n", oP, strings.Replace(oP.Description, "\n", " ", -1))
		}
	}
	if len(p.Closed) > 0 {
		ret += "Closed TCP Ports:\n"
		for _, cP := range p.Closed {
			ret += fmt.Sprintf("[-] %s - %s\n", cP, strings.Replace(cP.Description, "\n", " ", -1))
		}
	}
	if len(p.Filtered) > 0 {
		ret += "Offline or filtered TCP Ports:\n"
		for _, fP := range p.Filtered {
			ret += fmt.Sprintf("[?] %s - %s\n", fP, strings.Replace(fP.Description, "\n", " ", -1))
		}
	}
	ret += "**************************************************"
	return ret
}

/*
func (p *PortResults) String() string {
	var c int
	maxPerLine := 5

	ret := "*************** PORT RESULT **********************\n"
	if len(p.Open) > 0 {
		c = 0
		ret += "Open TCP Ports:\n\t"
		for _, oP := range p.Open {
			c++
			if c == len(p.Open) {
				ret += fmt.Sprintf("[+] %s\n", oP)
			} else if c%maxPerLine == 0 {
				ret += fmt.Sprintf("[+] %s\n\t", oP)
			} else {
				ret += fmt.Sprintf("[+] %-16s", oP)
			}
		}
	}
	if len(p.Closed) > 0 {
		c = 0
		ret += "Closed TCP Ports:\n\t"
		for _, cP := range p.Closed {
			c++
			if c == len(p.Closed) {
				ret += fmt.Sprintf("[-] %s\n", cP)
			} else if c%maxPerLine == 0 {
				ret += fmt.Sprintf("[-] %s\n\t", cP)
			} else {
				ret += fmt.Sprintf("[-] %-16s", cP)
			}
		}
	}
	if len(p.Filtered) > 0 {
		c = 0
		ret += "Offline or filtered TCP Ports:\n\t"
		for _, fP := range p.Filtered {
			c++
			if c == len(p.Filtered) {
				ret += fmt.Sprintf("[?] %s\n", fP)
			} else if c%maxPerLine == 0 {
				ret += fmt.Sprintf("[?] %s\n\t", fP)
			} else {
				ret += fmt.Sprintf("[?] %-16s", fP)
			}
		}
	}
	ret += "**************************************************"
	return ret
}
*/
