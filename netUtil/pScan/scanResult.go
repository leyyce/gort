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
	"fmt"
	"time"
)

// ScanResult represents the result of a single port scan.
// It contains the scans StartTime and EndTime, the scan Target and the PortResult.
type ScanResult struct {
	// StartTime is the time the scan was started at.
	StartTime time.Time

	// EndTime is the time the scan was finished at.
	EndTime time.Time

	// Target is the scan Target.
	Target *Target

	// Ports is the PortResult of the scan.
	Ports *PortResult
}

// ScanResults is an array of ScanResult pointers.
type ScanResults []*ScanResult

// MultiScanResult represents the scan result of multiple Targets.
type MultiScanResult struct {
	// Resolved contains the ScanResults of resolved hosts
	Resolved ScanResults

	// Unresolved contains the unresolved Targets
	Unresolved Targets
}

// NewScanResult returns the pointer to a new ScanResult instance.
// The parameter t is the Target of the scan and startTime is the start time.Time of the scan.
func NewScanResult(t *Target, startTime time.Time) *ScanResult {
	return &ScanResult{Target: t, StartTime: startTime, Ports: NewPortResult()}
}

// String returns a string representation of the ScanResult pointer.
func (s *ScanResult) String() string {
	return fmt.Sprintf(""+
		"=============== SCAN RESULT ================================\n\n"+
		"Scan started  @ %s\n"+
		"Scan finished @ %s\n"+
		"%s\n"+
		"%s\n\n"+
		"============================================================",
		s.StartTime.Format(time.RFC1123), s.EndTime.Format(time.RFC1123),
		s.Target,
		s.Ports)
}

// ColorString returns a colored string representation of the ScanResult pointer.
func (s *ScanResult) ColorString() string {
	return fmt.Sprintf(""+
		"=============== SCAN RESULT ================================\n\n"+
		"Scan started  @ %s\n"+
		"Scan finished @ %s\n"+
		"%s\n"+
		"%s\n\n"+
		"============================================================",
		s.StartTime.Format(time.RFC1123), s.EndTime.Format(time.RFC1123),
		s.Target.ColorString(),
		s.Ports.ColorString())
}

// CustomColorString returns a colored string representation of the ScanResult pointer.
// The parameter showClosed controls if closed and filtered ports also will be incorporated into the string.
func (s *ScanResult) CustomColorString(showClosed bool) string {
	return fmt.Sprintf(""+
		"=============== SCAN RESULT ================================\n\n"+
		"Scan started  @ %s\n"+
		"Scan finished @ %s\n"+
		"%s\n"+
		"%s\n\n"+
		"============================================================",
		s.StartTime.Format(time.RFC1123), s.EndTime.Format(time.RFC1123),
		s.Target.ColorString(),
		s.Ports.CustomColorString(showClosed))
}

// String returns a string representation of the MultiScanResult pointer.
func (m *MultiScanResult) String() string {
	ret := "" +
		"#################################################################\n" +
		"############### MULTI SCAN RESULT ###############################\n" +
		"#################################################################\n\n"
	if len(m.Resolved) == 0 {
		ret += "\tNONE\n"
	}
	for _, scanResult := range m.Resolved {
		ret += scanResult.String() + "\n\n"
	}
	ret += "\n" +
		"#################################################################\n" +
		"############### UNRESOLVED ######################################\n" +
		"#################################################################\n\n"
	if len(m.Unresolved) == 0 {
		ret += "\tNONE\n"
	}
	for _, target := range m.Unresolved {
		ret += target.String() + "\n\n"
	}
	ret += "\n#################################################################"
	return ret
}

// ColorString returns a colored string representation of the ScanResult pointer.
func (m *MultiScanResult) ColorString() string {
	ret := "" +
		"#################################################################\n" +
		"############### MULTI SCAN RESULT ###############################\n" +
		"#################################################################\n\n"
	if len(m.Resolved) == 0 {
		ret += "\tNONE\n"
	}
	for _, scanResult := range m.Resolved {
		ret += scanResult.ColorString() + "\n\n"
	}
	ret += "\n" +
		"#################################################################\n" +
		"############### UNRESOLVED ######################################\n" +
		"#################################################################\n\n"
	if len(m.Unresolved) == 0 {
		ret += "\tNONE\n"
	}
	for _, target := range m.Unresolved {
		ret += target.ColorString() + "\n\n"
	}
	ret += "\n#################################################################"
	return ret
}

// CustomColorString returns a colored string representation of the ScanResult pointer.
// The parameter onlineOnly controls if targets not confirmed as online will be incorporated into the string.
// showClosed controls if closed and filtered ports also will be incorporated into the string.
func (m *MultiScanResult) CustomColorString(onlineOnly, showClosed bool) string {
	ret := "" +
		"#################################################################\n" +
		"############### MULTI SCAN RESULT ###############################\n" +
		"#################################################################\n\n"
	if len(m.Resolved) == 0 {
		ret += "\tNONE\n"
	}
	for _, scanResult := range m.Resolved {
		if onlineOnly {
			if scanResult.Target.Status == Online {
				ret += scanResult.CustomColorString(showClosed) + "\n\n"
			}
		} else {
			ret += scanResult.CustomColorString(showClosed) + "\n\n"
		}
	}
	if !onlineOnly {
		ret += "\n" +
			"#################################################################\n" +
			"############### UNRESOLVED ######################################\n" +
			"#################################################################\n\n"
		if len(m.Unresolved) == 0 {
			ret += "\tNONE\n"
		}
		for _, target := range m.Unresolved {
			ret += target.ColorString() + "\n\n"
		}
	}
	ret += "\n#################################################################"
	return ret
}
