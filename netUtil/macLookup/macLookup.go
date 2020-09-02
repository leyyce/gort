package macLookup

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

var MACFormatError = errors.New("invalid mac format. Supported formats hex ':' bit '-' dot '.'")

const lookUpURL, format string = "http://macvendors.co/api/", "JSON"

func LookupVendor(hardwareAddr net.HardwareAddr) (*VendorResult, error) {
	url := fmt.Sprintf("%s%s/%s", lookUpURL, hardwareAddr.String(), format)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "API Browser")
	rsp, err := do(req)
	if rsp != nil {
		defer rsp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}
	if bytes.Contains(body, []byte("error")) {
		return nil, MACFormatError
	}
	vr := &VendorResult{}
	err = json.Unmarshal(body, &vr)
	if err != nil {
		return nil, err
	}
	return vr, nil
}

func do(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
