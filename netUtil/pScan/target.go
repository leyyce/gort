package pScan

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/colorFmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/netUtil"
	"net"
)

type Targets []*Target

type HostName string

type TargetStatus int

const (
	Online TargetStatus = iota + 1
	OfflineFiltered
	Unknown
)

type Target struct {
	HostName      HostName
	IPAddr        net.IP
	InitialTarget string
	Status        TargetStatus
	Ports         netUtil.Ports
}

func NewTarget(t string, ports netUtil.Ports) *Target {
	h := &Target{InitialTarget: t, Ports: ports, Status: Unknown}
	h.Resolve()
	return h
}

func (t *Target) Resolve() {
	if helper.ValidateIPOrRange(t.InitialTarget) {
		hostNames, err := net.LookupAddr(t.InitialTarget)
		if err != nil {
			t.HostName = "N/A"
		} else {
			t.HostName = HostName(hostNames[0])
		}
		t.IPAddr = net.ParseIP(t.InitialTarget)
	} else {
		ips, err := net.LookupIP(t.InitialTarget)
		if err == nil {
			t.IPAddr = ips[0]

			hostNames, err := net.LookupAddr(string(t.IPAddr))
			if err == nil {
				t.HostName = HostName(hostNames[0])
			} else {
				t.HostName = HostName(t.InitialTarget)
			}
		} else {
			t.IPAddr = nil
			t.Status = OfflineFiltered
		}
	}
}

func (t *Target) String() string {
	return fmt.Sprintf(""+
		"~~~~~~~~~~~~~~~ TARGET INFO ~~~~~~~~~~~~~~~~~~~~~~\n"+
		"Target: %s | IP: %s | Hostname: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		t.Status,
		t.Ports.Preview())
}

func (t *Target) ColorString() string {
	return fmt.Sprintf(""+
		"~~~~~~~~~~~~~~~ TARGET INFO ~~~~~~~~~~~~~~~~~~~~~~\n"+
		"Target: %s | IP: %s | Hostname: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		t.Status.ColorString(),
		t.Ports.Preview())
}

func (ts TargetStatus) String() string {
	if ts == Online {
		return "ONLINE"
	} else if ts == OfflineFiltered {
		return "OFFLINE / FILTERED"
	} else if ts == Unknown {
		return "UNKNOWN"
	}
	return "N/A"
}

func (ts TargetStatus) ColorString() string {
	if ts == Online {
		return colorFmt.Sopenf("ONLINE")
	} else if ts == OfflineFiltered {
		return colorFmt.Sclosedf("OFFLINE / FILTERED")
	} else if ts == Unknown {
		return colorFmt.Sfilteredf("UNKNOWN")
	}
	return "N/A"
}
