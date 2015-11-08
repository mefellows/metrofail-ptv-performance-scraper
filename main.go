package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gocarina/gocsv"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type PerformanceData struct {
	Delivery      float64 `csv:"delivery" json:"delivery"`
	Punctuality   float64 `csv:"punctuality" json:"punctuality"`
	Date          string  `csv:"date" json:"date"`
	TransportType string  `csv:"transport_type" json:"transport_type"`
}

func main() {
	//var url = "http://localhost:8000/perf.html" // python -m SimpleHTTPServer
	var DEFAULT_URL = "http://www.ptv.vic.gov.au/about-ptv/ptv-data-and-reports/daily-operational-performance-reports/"
	var format string
	var url string

	flag.StringVar(&format, "format", "json", "--format json,csv")
	flag.StringVar(&url, "url", DEFAULT_URL, "--url <path to PTV perf data>")
	flag.Parse()

	resp, err := http.Get(url)

	if err != nil {
		log.Fatalf("Error reading performance data: %s", err.Error())
	}

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	types := map[int]string{1: "train", 2: "tram", 3: "v/line"}
	values := make(map[int][][]string, 0)

	i := 0    // How many tables we've gone through
	rows := 0 // how many rows in table
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "table" {
			if values[i] == nil {
				values[i] = make([][]string, 1)
				values[i][0] = make([]string, 3)
			}
			i++
			rows = 0
		}
		if n.Type == html.ElementNode && n.Data == "tr" {
			values[i] = append(values[i], make([]string, 0))

			rows++
		}
		if n.Type == html.ElementNode && (n.Data == "th" || n.Data == "td") && rows > 2 {
			values[i][rows-2] = append(values[i][rows-2], n.FirstChild.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if len(values) != 4 {
		log.Fatalf("Error reading performance data, only have %d values: %s", len(values), values)
	}

	data := make([]PerformanceData, 0)

	for i, transportValues := range values {
		if len(transportValues) > 0 {
			for _, dataSet := range transportValues {
				if len(dataSet) >= 3 && dataSet[0] != "" {
					date, e := time.Parse("Monday, 2 January 2006", dataSet[0])
					dateFmt := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
					failOnErr(e)
					delivery, e := strconv.ParseFloat(strings.Replace(dataSet[1], "%", "", -1), 64)
					failOnErr(e)
					punctuality, e := strconv.ParseFloat(strings.Replace(dataSet[2], "%", "", -1), 64)
					failOnErr(e)
					data = append(data, PerformanceData{Delivery: delivery, Punctuality: punctuality, Date: dateFmt, TransportType: types[i]})
				}

			}
		}

	}

	switch format {
	case "json":
		d, err := json.Marshal(data)
		failOnErr(err)
		fmt.Printf("%s", string(d))
	case "csv":
		csvContent, err := gocsv.MarshalString(&data) // Get all clients as CSV string
		failOnErr(err)
		fmt.Println(csvContent)
	}
}

func failOnErr(e error) {
	if e != nil {
		log.Fatalf("Error: %v", e)
	}
}
