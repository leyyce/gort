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

package helper

import (
	"net"
	"strconv"
	"strings"
)

func ValidateIPOrRange(hosts string) bool {
	addrParts := strings.Split(hosts, ".")
	if len(addrParts) != 4 {
		return false
	}
	for _, addrPart := range addrParts {
		if strings.Contains(addrPart, "-") && strings.Count(addrPart, "-") == 1 {
			hRange := strings.Split(addrPart, "-")
			if len(hRange[0]) > 3 || len(hRange[1]) > 3 {
				return false
			}
			rStart, _ := strconv.Atoi(hRange[0])
			rEnd, _ := strconv.Atoi(hRange[0])
			if rStart > rEnd || !ValidateIPSegment(hRange[0]) || !ValidateIPSegment(hRange[1]) {
				return false
			}
		} else if !ValidateIPSegment(addrPart) {
			return false
		}
	}
	return true
}

func ValidateIPSegment(ipSeg string) bool {
	n, err := strconv.Atoi(ipSeg)
	if err != nil {
		return false
	}
	return n <= 255 && n >= 0
}

func StrRangeToArray(r string) []int {
	var ret []int
	temp := strings.Split(r, "-")
	// TODO Size check
	rStart, _ := strconv.Atoi(temp[0])
	rEnd, _ := strconv.Atoi(temp[1])
	for ; rStart <= rEnd; rStart++ {
		ret = append(ret, rStart)
	}
	return ret
}

func ValidatePort(port string) bool {
	_, err := strconv.ParseUint(port, 10, 16)
	if err != nil {
		return false
	}
	return true
}

func IncIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
