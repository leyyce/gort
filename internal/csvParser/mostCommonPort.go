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

package csvParser

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
)

type MostCommonPorts []*MostCommonPort

type MostCommonPort struct {
	Service     string
	Number      uint16
	Protocol    string
	Frequency   float64
	Description string
}

func NewMostCommonPort(service string, number uint16, proto string, freq float64, desc string) *MostCommonPort {
	return &MostCommonPort{
		Service:     service,
		Number:      number,
		Protocol:    proto,
		Frequency:   freq,
		Description: desc,
	}
}

func NewMostCommonPorts(dataDir string) *MostCommonPorts {
	csvFile, _ := os.Open(path.Join(dataDir, "port_open_freq.csv"))
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var commonPorts MostCommonPorts
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		service, protoNum, f, desc := line[0], line[1], line[2], line[3]
		x := strings.Split(protoNum, "/")
		n, _ := strconv.ParseInt(x[0], 10, 16)
		num := uint16(n)
		proto := x[1]
		freq, _ := strconv.ParseFloat(f, 64)

		commonPorts = append(
			commonPorts,
			NewMostCommonPort(
				service,
				num,
				proto,
				freq,
				desc,
			),
		)
	}

	return &commonPorts
}

func (ports *MostCommonPorts) GetMostCommonString(n int) string {
	var i int
	var str string
	for _, mS := range *ports {
		if i == n-1 {
			break
		}
		if mS.Protocol == "tcp" {
			str += strconv.Itoa(int(mS.Number)) + ","
			i++
		}
	}
	return str[:len(str)-2]
}
