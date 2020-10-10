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

// Vendor represents a single device manufacturer.
type Vendor struct {
	// Address is the postal address of the vendor
	Address string `json:"address"`

	// Company is the company name of the vendor.
	Company string `json:"company"`

	// Country is the country in which the vendor operates.
	Country string `json:"country"`

	// Type is the address block type of the vendor.
	Type string `json:"type"`

	// MacPrefix is the vendor specific MAC-prefix.
	MacPrefix string `json:"mac_prefix"`

	// StartHex is the first address of the vendor prefix.
	StartHex string `json:"start_hex"`

	//EndHex is the last address of the vendor prefix.
	EndHex string `json:"end_hex"`
}

// VendorResult represents a result returned from the API.
type VendorResult struct {
	*Vendor `json:"result"`
}
