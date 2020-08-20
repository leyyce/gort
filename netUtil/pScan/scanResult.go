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
		"=============== SCAN RESULT ================================\n"+
		"Scan started  @ %s\n"+
		"Scan finished @ %s\n"+
		"%s\n"+
		"%s\n"+
		"============================================================",
		s.StartTime.Format(time.RFC1123), s.EndTime.Format(time.RFC1123),
		s.Target,
		s.Ports)
}

func (m *MultiScanResult) String() string {
	ret := "" +
		"#################################################################\n" +
		"############### SCAN RESULTS ####################################\n" +
		"#################################################################\n"
	if len(m.Resolved) == 0 {
		ret += "NONE\n"
	}
	for _, scanResult := range m.Resolved {
		ret += scanResult.String() + "\n"
	}
	ret += "" +
		"#################################################################\n" +
		"############### UNRESOLVED ######################################\n" +
		"#################################################################\n"
	if len(m.Unresolved) == 0 {
		ret += "NONE\n"
	}
	for _, target := range m.Unresolved {
		ret += target.String() + "\n"
	}
	ret += "#################################################################"
	return ret
}
