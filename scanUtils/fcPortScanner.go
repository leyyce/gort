package scanUtils

import (
	"net"
	"strconv"
	"time"
)

type FCPortScanner struct {
	Host  *TargetInfo
	Ports []Port
}

func NewFCPortScanner(host *TargetInfo, ports []Port) *FCPortScanner {
	return &FCPortScanner{Host: host, Ports: ports}
}

func (f *FCPortScanner) Scan(out chan *ScanResult) {
	r := &ScanResult{Target: f.Host, StartTime: time.Now()}
	ch := make(chan PortResults)
	for _, p := range f.Ports {
		go f.scanPort(p, ch)
	}
	for range f.Ports {
		pI := <-ch
		r.Open = append(r.Open, pI.Open...)
		r.Closed = append(r.Closed, pI.Closed...)
	}
	r.EndTime = time.Now()
	out <- r
}

func (f *FCPortScanner) scanPort(p Port, ch chan PortResults) {
	var res PortResults
	seconds := 5
	timeOut := time.Duration(seconds) * time.Second
	_, err := net.DialTimeout("tcp", f.Host.IPAddr.String()+":"+strconv.Itoa(int(p)), timeOut)
	if err == nil {
		res.Open = append(res.Open, p)
		ch <- res
		return
	}
	res.Closed = append(res.Closed, p)
	ch <- res
}
