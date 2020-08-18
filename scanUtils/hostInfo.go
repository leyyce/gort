package scanUtils

import (
	"github.com/ElCap1tan/gort/internal/helper"
	"net"
)

type HostName string

type TargetInfo struct {
	HostName      HostName
	IPAddr        net.IP
	InitialTarget string
}

func NewHostInfo(t string) *TargetInfo {
	h := &TargetInfo{InitialTarget: t}
	h.resolve()
	return h
}

func (h *TargetInfo) resolve() {
	if helper.ValidateIPOrRange(h.InitialTarget) {
		hostNames, err := net.LookupAddr(h.InitialTarget)
		if err != nil {
			h.HostName = "N/A"
		} else {
			h.HostName = HostName(hostNames[0])
		}
		h.IPAddr = net.ParseIP(h.InitialTarget)
	} else {
		ips, err := net.LookupIP(h.InitialTarget)
		if err == nil {
			h.IPAddr = ips[0]

			hostNames, err := net.LookupAddr(string(h.IPAddr))
			if err == nil {
				h.HostName = HostName(hostNames[0])
			} else {
				h.HostName = HostName(h.InitialTarget)
			}
		} else {
			h.IPAddr = nil
		}
	}
}
