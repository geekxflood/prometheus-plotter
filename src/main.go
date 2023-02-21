package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// Set the Prometheus URL, query, and time range
	url := "http://localhost:9090/api/v1/query_range"
	query := "sum(rate(http_requests_total[5m]))"
	startTime := time.Now().Add(-24 * time.Hour).Unix()
	endTime := time.Now().Unix()

	// Retrieve data from Prometheus
	data, err := getData(url, query, startTime, endTime)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Extract the relevant data points
	values, err := extractData(data)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// Create a plot of the data
	err = createPlot(values, "Prometheus Query", "Time", "Requests", "output.png")
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println("Done!")
}
