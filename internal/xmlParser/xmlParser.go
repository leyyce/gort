package xmlParser

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

func NewPortRegistry() (*PortRegistry, error) {
	xmlFile, err := os.Open("data/port-numbers.xml")
	if err != nil {
		return nil, err
	}

	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)

	var portRegistry PortRegistry

	err = xml.Unmarshal(byteValue, &portRegistry)

	if err != nil {
		println(err.Error())
		return nil, err
	}

	return &portRegistry, nil
}
