package netUtil

import "fmt"

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
