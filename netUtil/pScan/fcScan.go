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
	"github.com/ElCap1tan/gort/internal/helper/ulimit"
	"github.com/ElCap1tan/gort/netUtil"
	"net"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sync/semaphore"
)

// Scan performs a concurrent full connection scan for every Target in Targets and once finished,
// returns the scan result as an MultiScanResult.
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

// Scan performs a concurrent full connection scan for all ports of a singe Target and
// returns a pointer to the ScanResult when finished.
func (t *Target) Scan() *ScanResult {
	r := NewScanResult(t, time.Now())
	ch := make(chan *PortResult)

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

// scan starts a concurrent full connection port scan for every port of the Target and writes the results to out.
// The lock can be used to control how many concurrent scans are allowed to run.
func (t *Target) scan(out chan *ScanResult, lock *semaphore.Weighted) {
	r := NewScanResult(t, time.Now())
	ch := make(chan *PortResult)
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

// scanPort scans a single port of the Target as specified by p. When finished the result is written to ch.
// The parameter lock can be used to control how many concurrent scans are allowed to run.
func (t *Target) scanPort(p *netUtil.Port, ch chan *PortResult, lock *semaphore.Weighted) {
	res := NewPortResult()
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
			// TODO Check if it makes sense to not release the lock if to many files are open already
			lock.Release(1)
			return
		}
	}
	lock.Release(1)
	ch <- res
}
