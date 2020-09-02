package macLookup

type Vendor struct {
	Address   string `json:"address"`
	Company   string `json:"company"`
	Country   string `json:"country"`
	Type      string `json:"type"`
	MacPrefix string `json:"mac_prefix"`
	StartHex  string `json:"start_hex"`
	EndHex    string `json:"end_hex"`
}

type VendorResult struct {
	*Vendor `json:"result"`
}
