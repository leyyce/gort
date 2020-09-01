package pScan

import (
	"fmt"
	"time"
)

type ScanResult struct {
	StartTime time.Time
	EndTime   time.Time
	Target    *Target
	Ports     *PortResults
}

type MultiScanResult struct {
	Resolved   []*ScanResult
	Unresolved Targets
}

func NewScanResult(t *Target, startTime time.Time) *ScanResult {
	return &ScanResult{Target: t, StartTime: startTime, Ports: NewPortResults()}
}

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
