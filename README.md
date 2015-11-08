# #Metrofail PTV Performance Data Scraper

Scrapes performance data from https://www.ptv.vic.gov.au/about-ptv/ptv-data-and-reports/daily-operational-performance-reports/.

Used to provide updated performance information in [Metrofail](http://www.metrofail.org).

## Usage

```
go get ./...
go build -o ptvperf .
./ptvperf --format csv --url http://www.ptv.vic.gov.au/about-ptv/ptv-data-and-reports/daily-operational-performance-reports/
```

## Running in Docker

```
gox -osarch="linux/amd64" -output="ptvperf" && docker build -t mfellows/ptvperf
docker run  mfellows/ptvperf
```
