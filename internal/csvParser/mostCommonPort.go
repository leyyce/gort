package csvParser

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
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

func NewMostCommonPorts() *MostCommonPorts {
	csvFile, _ := os.Open("data/port_open_freq.csv")
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
