package netUtil

import "fmt"

type Ports []*Port

type Port struct {
	PortNo   uint16
	Protocol string
}

func NewPort(portNo uint16, proto string) *Port {
	return &Port{PortNo: portNo, Protocol: proto}
}

func (p *Port) String() string {
	return fmt.Sprintf("%d/%s", p.PortNo, p.Protocol)
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
	max := 30
	ret := ""
	if len(ps) < max {
		for i, p := range ps {
			ret += p.String() + ", "
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
		}
		ret = ret[:len(ret)-2]
	} else {
		for i, p := range ps[:max] {
			if i != 0 && i%10 == 0 && i+1 != len(ps) {
				ret += "\n"
			}
			ret += p.String() + ", "
		}
		ret = ret[:len(ret)-2] + "..."
	}
	return ret
}
