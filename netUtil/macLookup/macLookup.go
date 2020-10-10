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

// MACFormatError is returned if the provided MAC-address is an unsupported format.
var MACFormatError = errors.New("invalid mac format. Supported formats hex ':' bit '-' dot '.'")

const lookUpURL, format string = "http://macvendors.co/api", "JSON"

// LookupVendor tries to look up the vendor of hardwareAddr by sending a HTTP-request to
// the API of 'http://macvendors.co/' and if successful returns a pointer to the VendorResult.
func LookupVendor(hardwareAddr net.HardwareAddr) (*VendorResult, error) {
	url := fmt.Sprintf("%s/%s/%s", lookUpURL, hardwareAddr.String(), format)
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
