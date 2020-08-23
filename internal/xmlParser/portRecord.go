package xmlParser

import "encoding/xml"

type PortRecord struct {
	XMLName     xml.Name `xml:"record"`
	Service     string   `xml:"name"`
	Protocol    string   `xml:"protocol"`
	Description string   `xml:"description"`
	Number      string   `xml:"number"`
}
