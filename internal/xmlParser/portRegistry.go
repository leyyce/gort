package xmlParser

import (
	"encoding/xml"
)

type PortRegistry struct {
	XMLName xml.Name     `xml:"registry"`
	Records []PortRecord `xml:"record"`
}
