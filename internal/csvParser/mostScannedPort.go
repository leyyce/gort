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

type MostScannedPorts []*MostScannedPort

type MostScannedPort struct {
	Service     string
	Number      uint16
	Protocol    string
	Frequency   float64
	Description string
}

func NewMostScannedPort(service string, number uint16, proto string, freq float64, desc string) *MostScannedPort {
	return &MostScannedPort{
		Service:     service,
		Number:      number,
		Protocol:    proto,
		Frequency:   freq,
		Description: desc,
	}
}

func NewMostScannedPorts() *MostScannedPorts {
	csvFile, _ := os.Open("data/port_open_freq.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var mostScanned MostScannedPorts
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

		mostScanned = append(
			mostScanned,
			NewMostScannedPort(
				service,
				num,
				proto,
				freq,
				desc,
			),
		)
	}

	return &mostScanned
}
