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
	Delivery      float64 `csv:"delivery"`
	Punctuality   float64 `csv:"punctuality"`
	Date          string  `csv:"date"`
	TransportType string  `csv:"transport_type"`
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
	values := make([]string, 0)
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "span" {
			for _, a := range n.Parent.Parent.Attr {
				if a.Key == "class" && a.Val == "results" {
					values = append(values, n.FirstChild.Data)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if len(values) < 2 || len(values)%2 != 0 {
		log.Fatalf("Error reading performance data, only have %d values: %s", len(values), values)
	}

	data := make([]PerformanceData, 0)
	_date := time.Now()
	date := fmt.Sprintf("%d-%d-%d", _date.Year(), _date.Month(), _date.Day())
	types := map[int]string{0: "train", 2: "tram", 4: "v/line"}

	for i := 0; i < len(values); i += 2 {
		delivery, e := strconv.ParseFloat(strings.Replace(values[i], "%", "", -1), 64)
		failOnErr(e)
		punctuality, e := strconv.ParseFloat(strings.Replace(values[i+1], "%", "", -1), 64)
		failOnErr(e)
		data = append(data, PerformanceData{Delivery: delivery, Punctuality: punctuality, Date: date, TransportType: types[i]})
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
