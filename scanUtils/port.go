package scanUtils

type Port uint16

type PortResults struct {
	Open   []Port
	Closed []Port
}
