package pScan

import (
	"github.com/ElCap1tan/gort/netUtil"
	"net"
	"strconv"
	"strings"
	"time"
)

func (t Targets) Scan() MultiScanResult {
	var multiScanRes MultiScanResult
	out := make(chan *ScanResult)
	for _, t := range t {
		if t.IPAddr == nil {
			multiScanRes.Unresolved = append(multiScanRes.Unresolved, t)
		} else {
			go t.scan(out)
		}
	}
	for i := 0; i < len(t)-len(multiScanRes.Unresolved); i++ {
		multiScanRes.Resolved = append(multiScanRes.Resolved, <-out)
	}
	return multiScanRes
}

func (t *Target) Scan() *ScanResult {
	r := NewScanResult(t, time.Now())
	ch := make(chan *PortResults)
	for _, p := range t.Ports {
		go t.scanPort(p, ch)
	}
	for range t.Ports {
		pI := <-ch
		r.Ports.Open = append(r.Ports.Open, pI.Open...)
		r.Ports.Closed = append(r.Ports.Closed, pI.Closed...)
		r.Ports.Filtered = append(r.Ports.Filtered, pI.Filtered...)
	}
	r.EndTime = time.Now()
	return r
}

func (t *Target) scan(out chan *ScanResult) {
	r := NewScanResult(t, time.Now())
	ch := make(chan *PortResults)
	for _, p := range t.Ports {
		go t.scanPort(p, ch)
	}
	for range t.Ports {
		pI := <-ch
		r.Ports.Open = append(r.Ports.Open, pI.Open...)
		r.Ports.Closed = append(r.Ports.Closed, pI.Closed...)
		r.Ports.Filtered = append(r.Ports.Filtered, pI.Filtered...)
	}
	r.EndTime = time.Now()
	out <- r
}

func (t *Target) scanPort(p *netUtil.Port, ch chan *PortResults) {
	res := NewPortResults()
	seconds := 5
	timeOut := time.Duration(seconds) * time.Second
	conn, err := net.DialTimeout("tcp", t.IPAddr.String()+":"+strconv.Itoa(int(p.PortNo)), timeOut)
	if err == nil {
		t.Status = Online
		res.Open = append(res.Open, p)
		ch <- res
		_ = conn.Close()
		return
	} else if _, ok := err.(*net.OpError); ok {
		if t.Status == Unknown || t.Status == OfflineFiltered {
			if strings.HasSuffix(err.Error(), "No connection could be made because the target machine actively refused it.") {
				t.Status = Online
			} else if strings.HasSuffix(err.Error(), "i/o timeout") && t.Status == Unknown {
				t.Status = OfflineFiltered
			}
		}
		if strings.HasSuffix(err.Error(), "No connection could be made because the target machine actively refused it.") {
			res.Closed = append(res.Closed, p)
		}
		if strings.HasSuffix(err.Error(), "i/o timeout") {
			res.Filtered = append(res.Filtered, p)
		}
	}
	ch <- res
}
