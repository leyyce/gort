package pScan

import (
	"fmt"
	"github.com/ElCap1tan/gort/internal/colorFmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/netUtil"
	"github.com/ElCap1tan/gort/netUtil/macLookup"
	"github.com/mdlayher/arp"
	quickArp "github.com/mostlygeek/arp"
	"net"
	"runtime"
)

type Targets []*Target

type HostName string

type TargetStatus int
type NetworkLocation int

const (
	Online TargetStatus = iota + 1
	OfflineFiltered
	Unknown
	Local NetworkLocation = iota + 1
	Global
	UnknownLoc
)

type Target struct {
	HostName      HostName
	Vendor        string
	IPAddr        net.IP
	MACAddr       net.HardwareAddr
	InitialTarget string
	Status        TargetStatus
	Location      NetworkLocation
	Ports         netUtil.Ports
}

func NewTarget(t string, ports netUtil.Ports) *Target {
	h := &Target{InitialTarget: t, Ports: ports, Status: Unknown}
	h.Resolve()
	if h.IPAddr != nil {
		h.QueryMac()
		h.LookUpVendor()
	} else {
		h.MACAddr = nil
		h.Location = UnknownLoc
	}
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

func (t *Target) QueryMac() {
	if macAddr, err := net.ParseMAC(quickArp.Search(t.IPAddr.String())); err == nil {
		t.Location = Local
		t.MACAddr = macAddr
		return
	}
	// Fallback if not found in arp cache
	interfaces, err := net.Interfaces()
	if err != nil {
		t.Location = UnknownLoc
		t.MACAddr = nil
		return
	}
	for _, inf := range interfaces {
		infAddresses, err := inf.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range infAddresses {
			if _, ipNet, err := net.ParseCIDR(addr.String()); err == nil {
				if ipNet.Contains(t.IPAddr) {
					t.Location = Local
					if arpCli, err := arp.Dial(&inf); err == nil {
						hwAddr, err := arpCli.Resolve(t.IPAddr)
						if err != nil {
							t.MACAddr = nil
							return
						}
						fmt.Println(hwAddr.String())
						t.MACAddr = hwAddr
						t.Status = Online
						return
					} else {
						t.MACAddr = nil
						return
					}
				}
			}
		}
	}
	t.Location = Global
	t.MACAddr = nil
	return
}

func (t *Target) LookUpVendor() {
	if t.MACAddr != nil {
		vendorRes, err := macLookup.LookupVendor(t.MACAddr)
		if err != nil {
			t.Vendor = "N/A"
			return
		}
		t.Vendor = vendorRes.Company
		return
	}
	t.Vendor = "N/A"
}

func (t *Target) String() string {
	var mac string
	if t.MACAddr == nil {
		mac = "N/A"
	} else {
		mac = t.MACAddr.String()
	}

	return fmt.Sprintf(""+
		"~~~~~~~~~~~~~~~ TARGET INFO ~~~~~~~~~~~~~~~~~~~~~~\n"+
		"Target: %s | IP: %s | Hostname: %s\n"+
		"Vendor: %s\n"+
		"MacAddress: %s\n"+
		"Network Location: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		t.Vendor,
		mac,
		t.Location,
		t.Status,
		t.Ports.Preview())
}

func (t *Target) ColorString() string {
	var mac string
	if t.MACAddr == nil && t.Location == Local && runtime.GOOS == "windows" {
		mac = colorFmt.Swarnf("N/A")
	} else if t.MACAddr == nil && t.Location == Local {
		mac = colorFmt.Sfatalf("N/A")
	} else if t.MACAddr == nil {
		mac = colorFmt.Ssuccessf("N/A")
	} else {
		mac = colorFmt.Ssuccessf(t.MACAddr.String())
	}

	return fmt.Sprintf(""+
		"~~~~~~~~~~~~~~~ TARGET INFO ~~~~~~~~~~~~~~~~~~~~~~\n"+
		"Target: %s | IP: %s | Hostname: %s\n"+
		"Vendor: %s\n"+
		"MacAddress: %s\n"+
		"Network Location: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		t.Vendor,
		mac,
		t.Location.ColorString(),
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

func (n NetworkLocation) String() string {
	if n == Local {
		return "LOCAL"
	} else if n == Global {
		return "GLOBAL"
	} else if n == UnknownLoc {
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

func (n NetworkLocation) ColorString() string {
	if n == Local {
		return colorFmt.Sinfof("LOCAL")
	} else if n == Global {
		return colorFmt.Sinfof("GLOBAL")
	} else if n == UnknownLoc {
		return colorFmt.Swarnf("UNKNOWN")
	}
	return "N/A"
}
