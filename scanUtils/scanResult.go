package scanUtils

import (
	"fmt"
	"time"
)

type ScanResult struct {
	StartTime time.Time
	EndTime   time.Time
	Target    *TargetInfo
	Open      []Port
	Closed    []Port
}

func (s ScanResult) String() string {
	return fmt.Sprintf(""+
		"--------------------------------------------\n"+
		"Target: %s | Hostname: %s | IP: %s\n"+
		"Scan started @: %s | Scan finished @: %s\n"+
		"Open TCP Ports:\n"+
		"%d\n"+
		"Closed TCP Ports:\n"+
		"%d\n"+
		"--------------------------------------------",
		s.Target.InitialTarget, string(s.Target.HostName), s.Target.IPAddr.String(),
		s.StartTime.Format(time.RFC1123), s.EndTime.Format(time.RFC1123),
		s.Open,
		s.Closed)
}
