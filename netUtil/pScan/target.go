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

package pScan

import (
	"context"
	"fmt"
	"github.com/ElCap1tan/gort/internal/colorFmt"
	"github.com/ElCap1tan/gort/internal/helper"
	"github.com/ElCap1tan/gort/internal/helper/ulimit"
	"github.com/ElCap1tan/gort/internal/symbols"
	"github.com/ElCap1tan/gort/netUtil"
	"github.com/ElCap1tan/gort/netUtil/macLookup"
	"github.com/mdlayher/arp"
	quickArp "github.com/mostlygeek/arp"
	"github.com/sparrc/go-ping"
	"golang.org/x/sync/semaphore"
	"net"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Targets is an array of Target pointers.
type Targets []*Target

// HostName is a string representing the hostname of the Target.
type HostName string

// TargetStatus is an integer representing the status of the Target.
// The values can be Online, OfflineFiltered or Unknown.
type TargetStatus int

// NetworkLocation is an integer representing the Target location inside the network.
// The values can be Local, Global or UnknownLoc.
type NetworkLocation int

const (
	Online TargetStatus = iota
	OfflineFiltered
	Unknown
	Local NetworkLocation = iota
	Global
	UnknownLoc
)

// Target represent a network host with all the necessary fields to conduct a port scan.
// A target always should be initialized with the NewTarget, AsyncNewTarget or ParseHostString methods.
type Target struct {
	// HostName is string containing the host name of the Target.
	HostName HostName

	//Vendor is a string containing the name of the Target vendor if the lookup was successful.
	Vendor string

	// IPAddr is the ip of the Target.
	IPAddr net.IP

	// MACAddr is the MAC-address of the target if it could be found.
	MACAddr net.HardwareAddr

	// InitialTarget is a string containing the initial IP-address or host name the target was created with.
	InitialTarget string

	// Status contains the TargetStatus
	Status TargetStatus

	// Location contains the NetworkLocation of the Target.
	Location NetworkLocation

	// Ports is a list of type netUtil.ports to be scanned.
	Ports netUtil.Ports

	// RTTs contains the round trip times of the ping requests if they could be send successfully.
	RTTs []time.Duration
}

// NewTarget returns a pointer to an initialized instance of Target as defined
// by the targetAddress and ports. Before returning the Target, it is resolved by calling Target.Resolve.
// If the resolve was successful, AsyncNewTarget will try to send a ping request by calling Target.Ping and to query
// the MAC-address and vendor name  by calling Target.QueryMac and Target.LookUpVendor. scanLock is used to controls
// how many targets may be resolved simultaneously and privileged controls if the scan should be run either
// in a (more detailed) mode that require root privileges, or (in the less detailed) 'user' mode.
func NewTarget(targetAddress string, ports netUtil.Ports, privileged bool) *Target {
	h := &Target{InitialTarget: targetAddress, Ports: ports, Status: Unknown}
	h.Resolve()
	if h.IPAddr != nil {
		stats, _ := h.Ping(3, privileged)
		if stats.PacketsRecv > 0 {
			h.Status = Online
		}
		h.QueryMac()
		h.LookUpVendor()
	} else {
		h.MACAddr = nil
		h.Location = UnknownLoc
	}
	return h
}

// AsyncNewTarget asynchronously creates a pointer to an initialized instance of Target as defined
// by the targetAddress and ports. Before returning the Target over ch, it is resolved by calling Target.Resolve.
// If the resolve was successful, AsyncNewTarget will try to send a ping request by calling Target.Ping and to query
// the MAC-address and vendor name  by calling Target.QueryMac and Target.LookUpVendor. scanLock is used to controls
// how many targets may be resolved simultaneously and privileged controls if the scan should be run either
// in a (more detailed) mode that require root privileges, or (in the less detailed) 'user' mode.
func AsyncNewTarget(targetAddress string, ports netUtil.Ports, ch chan *Target, scanLock *semaphore.Weighted, privileged bool) {
	// TODO Add writeMutex
	h := &Target{InitialTarget: targetAddress, Ports: ports, Status: Unknown}
	scanLock.Acquire(context.TODO(), 1)
	h.Resolve()
	scanLock.Release(1)
	if h.IPAddr != nil {
		scanLock.Acquire(context.TODO(), 1)
		stats, _ := h.Ping(3, privileged)
		scanLock.Release(1)
		if stats.PacketsRecv > 0 {
			h.Status = Online
		}
		scanLock.Acquire(context.TODO(), 1)
		h.QueryMac()
		scanLock.Release(1)
		scanLock.Acquire(context.TODO(), 1)
		h.LookUpVendor()
		scanLock.Release(1)
	} else {
		h.MACAddr = nil
		h.Location = UnknownLoc
	}
	ch <- h
}

// ParseHostString parses hosts and returns the initialized Targets.
//
// hosts is comma separated list of values that can be in either of the following formats:
// - A single IP address: 192.88.99.1
// - A range of IP addresses: 192.88.99-100.1-100
// - A CIDR formatted IP address range: 192.88.99.1/24
// - A host name: example.com
//
// ports is a list of type ports that should be scanned for every host in hosts.
//
// privileged controls if the targets should be resolved either in
// a (more detailed) mode that require root privileges or (in the less detailed) 'user' mode.
func ParseHostString(hosts string, ports netUtil.Ports, privileged bool) Targets {
	var tgtHosts Targets
	hostCount := 0
	out := make(chan *Target)

	var limit int64
	l, err := ulimit.GetUlimit()
	if err != nil {
		limit = 1024
	} else {
		limit = int64(l)
	}

	lock := semaphore.NewWeighted(limit)

	hostList := strings.Split(hosts, ",")
	for _, hostArg := range hostList {
		if ip, ipNet, err := net.ParseCIDR(hostArg); err == nil {
			for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); helper.IncIp(ip) {
				go AsyncNewTarget(ip.String(), ports, out, lock, privileged)
				hostCount++
			}
		} else if helper.ValidateIPOrRange(hostArg) {
			if strings.Contains(hostArg, "-") {
				addrParts := strings.Split(hostArg, ".")
				var octets [4][]int
				for i, addrPart := range addrParts {
					if strings.Contains(addrPart, "-") {
						octets[i] = append(octets[i], helper.StrRangeToArray(addrPart)...)
					} else {
						p, _ := strconv.Atoi(addrPart)
						octets[i] = append(octets[i], p)
					}
				}
				for _, t := range octetsToTargets(octets) {
					go AsyncNewTarget(t, ports, out, lock, privileged)
					hostCount++
				}
			} else {
				go AsyncNewTarget(hostArg, ports, out, lock, privileged)
				hostCount++
			}
		} else {
			go AsyncNewTarget(hostArg, ports, out, lock, privileged)
			hostCount++
		}
	}
	for i := 1; i <= hostCount; i++ {
		tgtHosts = append(tgtHosts, <-out)
	}
	return tgtHosts
}

// Resolve tries to resolve the IP address and the host name of the Target pointer.
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

// QueryMac tries to query the MAC address of the Target pointer either by ARP cache lookup or alternatively if
// not found in cache by sending an ARP-request.
func (t *Target) QueryMac() {
	if b, err := t.IsHost(); b == true && err == nil {
		return
	}
	if macAddr, err := net.ParseMAC(quickArp.Search(t.IPAddr.String())); err == nil && macAddr.String() != "00:00:00:00:00:00" {
		colorFmt.Infof("%s Found MAC address for target '%s' via arp cache lookup: %s\n",
			symbols.INFO, t.InitialTarget, macAddr.String())
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
						err = arpCli.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
						if err != nil {
							colorFmt.Warnf("%s %s: Error setting read timeout for arp request. Skipping mac lookup...\n",
								symbols.INFO, t.IPAddr.String())
							t.MACAddr = nil
							_ = arpCli.Close()
							return
						}
						hwAddr, err := arpCli.Resolve(t.IPAddr)
						if err != nil || hwAddr.String() == "00:00:00:00:00:00" {
							t.MACAddr = nil
							_ = arpCli.Close()
							return
						}
						colorFmt.Infof("%s Found MAC address for target '%s' via arp request: %s\n",
							symbols.INFO, t.InitialTarget, hwAddr.String())
						t.MACAddr = hwAddr
						t.Status = Online
						_ = arpCli.Close()
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

// LookUpVendor tries to perform a vendor lookup based on the MAC address of the Target pointer by sending
// a HTTP request to the vendor lookup API of 'macvendors.co'.
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

// Ping sends a ping request to the IP address of the Target pointer. count specifies how many requests should be send
// and privileged
func (t *Target) Ping(count int, privileged bool) (*ping.Statistics, error) {
	pinger, err := ping.NewPinger(t.IPAddr.String())
	if err != nil {
		t.RTTs = nil
		return nil, err
	}
	pinger.Timeout = time.Millisecond * 3000
	pinger.Count = count

	pinger.SetPrivileged(privileged || runtime.GOOS == "windows")

	pinger.Run()

	stats := pinger.Statistics()
	t.RTTs = stats.Rtts
	return stats, nil
}

// AvgRTT calculates the average RTT of the last Target.Ping call.
func (t Target) AvgRTT() time.Duration {
	if len(t.RTTs) == 0 {
		return -1
	}
	var avgNS int64 = 0
	for _, rtt := range t.RTTs {
		avgNS += rtt.Nanoseconds()
	}
	avgNS /= int64(len(t.RTTs))
	return time.Duration(avgNS)
}

// IsHost returns true when the Target pointer is the host and false otherwise.
func (t *Target) IsHost() (bool, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	for _, inf := range interfaces {
		infAddresses, err := inf.Addrs()
		if err != nil {
			continue
		}

		// Check if target ip is the host
		for _, addr := range infAddresses {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
				if t.IPAddr.Equal(ip) {
					t.Location = Local
					t.MACAddr = inf.HardwareAddr
					return true, nil
				}
			case *net.IPNet:
				ip = net.ParseIP(strings.Split(v.IP.String(), ":")[0])
				if t.IPAddr.Equal(ip) {
					t.Location = Local
					t.MACAddr = inf.HardwareAddr
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// String returns a string representation of the Target pointer.
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
		"Avg Ping [%d send]: %v\n"+
		"Vendor: %s\n"+
		"MacAddress: %s\n"+
		"Network Location: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		len(t.RTTs), t.AvgRTT(),
		t.Vendor,
		mac,
		t.Location,
		t.Status,
		t.Ports.Preview(30))
}

// ColorString returns a colored string representation of the Target pointer.
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

	var vendor string
	if t.Location == Local && t.Vendor != "" && t.Vendor != "N/A" {
		vendor = colorFmt.Ssuccessf(t.Vendor)
	} else if t.Location == Local {
		vendor = colorFmt.Sfatalf("N/A")
	} else if t.Location == Global {
		vendor = colorFmt.Ssuccessf("N/A")
	}

	var rtt string
	avg := t.AvgRTT()
	if avg > 0 {
		rtt = colorFmt.Ssuccessf("%v", avg)
	} else {
		rtt = colorFmt.Sfatalf("%v", avg)
	}

	return fmt.Sprintf(""+
		"~~~~~~~~~~~~~~~ TARGET INFO ~~~~~~~~~~~~~~~~~~~~~~\n"+
		"Target: %s | IP: %s | Hostname: %s\n"+
		"Avg Ping [%d send]: %v\n"+
		"Vendor: %s\n"+
		"MacAddress: %s\n"+
		"Network Location: %s\n"+
		"Status: %s\n"+
		"Ports:\n%s\n"+
		"~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~",
		t.InitialTarget, t.IPAddr, t.HostName,
		len(t.RTTs), rtt,
		vendor,
		mac,
		t.Location.ColorString(),
		t.Status.ColorString(),
		t.Ports.Preview(30))
}

// String returns a string representation of TargetStatus.
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

// String returns a string representation of NetworkLocation.
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

// ColorString returns a colored string representation of TargetStatus.
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

// ColorString returns a colored string representation of NetworkLocation.
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

// octetsToTargets takes a two-dimensional array containing a list of values for each of the four octets as values
// and returns an string array containing all the possible IP combinations you can build up with it.
func octetsToTargets(octets [4][]int) []string {
	var targets []string
	for _, oc0 := range octets[0] {
		for _, oc1 := range octets[1] {
			for _, oc2 := range octets[2] {
				for _, oc3 := range octets[3] {
					targets = append(targets, fmt.Sprintf("%d.%d.%d.%d", oc0, oc1, oc2, oc3))
				}
			}
		}
	}
	return targets
}
