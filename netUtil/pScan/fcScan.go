package pScan

import (
	"context"
	"github.com/ElCap1tan/gort/internal/helper/ulimit"
	"github.com/ElCap1tan/gort/netUtil"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
)

func (t Targets) Scan() MultiScanResult {
	var multiScanRes MultiScanResult
	out := make(chan *ScanResult)

	var limit int64
	l, err := ulimit.GetUlimit()
	if err != nil {
		limit = 1024
	} else {
		limit = int64(l)
	}

	lock := semaphore.NewWeighted(limit)
	for _, t := range t {
		if t.IPAddr == nil {
			multiScanRes.Unresolved = append(multiScanRes.Unresolved, t)
		} else {
			go t.scan(out, lock)
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

	var limit int64
	l, err := ulimit.GetUlimit()
	if err != nil {
		limit = 1024
	} else {
		limit = int64(l)
	}

	lock := semaphore.NewWeighted(limit)
	for _, p := range t.Ports {
		go t.scanPort(p, ch, lock)
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

func (t *Target) scan(out chan *ScanResult, lock *semaphore.Weighted) {
	r := NewScanResult(t, time.Now())
	ch := make(chan *PortResults)
	for _, p := range t.Ports {
		go t.scanPort(p, ch, lock)
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

func (t *Target) scanPort(p *netUtil.Port, ch chan *PortResults, lock *semaphore.Weighted) {
	res := NewPortResults()
	milli := 3000
	timeOut := time.Duration(milli) * time.Millisecond
	lock.Acquire(context.TODO(), 1)
	conn, err := net.DialTimeout("tcp", t.IPAddr.String()+":"+strconv.Itoa(int(p.PortNo)), timeOut)
	if err == nil {
		defer conn.Close()
		t.Status = Online
		res.Open = append(res.Open, p)
		ch <- res
		lock.Release(1)
		return
	} else if _, ok := err.(*net.OpError); ok {
		if t.Status == Unknown || t.Status == OfflineFiltered {
			if strings.HasSuffix(err.Error(), "No connection could be made because the target machine actively refused it.") ||
				strings.HasSuffix(err.Error(), "connect: connection refused") {
				t.Status = Online
			} else if strings.HasSuffix(err.Error(), "i/o timeout") && t.Status == Unknown {
				t.Status = OfflineFiltered
			}
		}
		if strings.HasSuffix(err.Error(), "No connection could be made because the target machine actively refused it.") ||
			strings.HasSuffix(err.Error(), "connect: connection refused") {
			res.Closed = append(res.Closed, p)
		}
		if strings.HasSuffix(err.Error(), "i/o timeout") {
			res.Filtered = append(res.Filtered, p)
		}
		if strings.HasSuffix(err.Error(), "too many open files") {
			time.Sleep(timeOut)
			go t.scanPort(p, ch, lock)
			lock.Release(1)
			return
		}
	}
	lock.Release(1)
	ch <- res
}
