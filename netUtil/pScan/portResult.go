package pScan

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/colorFmt"
	"github.com/ElCap1tan/gort/internal/symbols"
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
			if oP.Description == "" {
				ret += fmt.Sprintf("\t%s %s\n", symbols.OPEN, oP)
			} else {
				ret += fmt.Sprintf("\t%s %s - %s\n", symbols.OPEN, oP, strings.Replace(oP.Description, "\n", " ", -1))
			}
		}
	}
	if len(p.Closed) > 0 {
		ret += "Closed TCP Ports:\n"
		for _, cP := range p.Closed {
			if cP.Description == "" {
				ret += fmt.Sprintf("\t%s %s\n", symbols.CLOSED, cP)
			} else {
				ret += fmt.Sprintf("\t%s %s - %s\n", symbols.CLOSED, cP, strings.Replace(cP.Description, "\n", " ", -1))
			}
		}
	}
	if len(p.Filtered) > 0 {
		ret += "Offline or filtered TCP Ports:\n"
		for _, fP := range p.Filtered {
			if fP.Description == "" {
				ret += fmt.Sprintf("\t%s %s\n", symbols.UNKNOWN, fP)
			} else {
				ret += fmt.Sprintf("\t%s %s - %s\n", symbols.UNKNOWN, fP, strings.Replace(fP.Description, "\n", " ", -1))
			}
		}
	}
	ret += "**************************************************"
	return ret
}

func (p *PortResults) ColorString() string {
	ret := "*************** PORT RESULT **********************\n"
	if len(p.Open) > 0 {
		ret += "Open TCP Ports:\n"
		for _, oP := range p.Open {
			if oP.Description == "" {
				ret += colorFmt.Sopenf("\t%s %s\n", symbols.OPEN, oP)
			} else {
				ret += colorFmt.Sopenf("\t%s %s - %s\n", symbols.OPEN, oP, strings.Replace(oP.Description, "\n", " ", -1))
			}
		}
	}
	if len(p.Closed) > 0 {
		ret += "Closed TCP Ports:\n"
		for _, cP := range p.Closed {
			if cP.Description == "" {
				ret += colorFmt.Sclosedf("\t%s %s\n", symbols.CLOSED, cP)
			} else {
				ret += colorFmt.Sclosedf("\t%s %s - %s\n", symbols.CLOSED, cP, strings.Replace(cP.Description, "\n", " ", -1))
			}
		}
	}
	if len(p.Filtered) > 0 {
		ret += "Offline or filtered TCP Ports:\n"
		for _, fP := range p.Filtered {
			if fP.Description == "" {
				ret += colorFmt.Sfilteredf("\t%s %s\n", symbols.UNKNOWN, fP)
			} else {
				ret += colorFmt.Sfilteredf("\t%s %s - %s\n", symbols.UNKNOWN, fP, strings.Replace(fP.Description, "\n", " ", -1))
			}
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
