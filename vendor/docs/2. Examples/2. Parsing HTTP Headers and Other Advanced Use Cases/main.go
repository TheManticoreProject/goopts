package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/TheManticoreProject/goopts/parser"
)

var (
	// Positional arguments
	url string

	// Request flags
	method  string
	headers map[string]string

	// Feature flags
	verbose bool
)

func parseArgs() {
	// Create a new arguments parser with a custom banner
	ap := parser.ArgumentsParser{Banner: "PoC of goopts parsing - HTTP request example v1.0 - by Remi GASCOU (Podalirius)"}

	// Define positional argument
	ap.NewStringPositionalArgument(&url, "url", "URL to send the HTTP request to.")

	// Define global flags
	ap.NewStringArgument(&method, "-X", "--method", "GET", false, "HTTP request method (e.g., GET, POST, PUT, DELETE).")
	ap.NewMapOfHttpHeadersArgument(&headers, "-H", "--header", map[string]string{}, false, "Header for the request (can be specified multiple times). Example: -H 'Authorization: Bearer token'")
	ap.NewBoolArgument(&verbose, "-v", "--verbose", false, "Enable verbose mode to see detailed request and response information.")

	// Parse the flags
	ap.Parse()
}

func main() {
	parseArgs()

	// Create HTTP request
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Fatalf("Error creating request: %s", err)
	}

	// Add headers to the request
	for headerName, headerValue := range headers {
		req.Header.Add(headerName, headerValue)
	}

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error performing HTTP request: %s", err)
	}
	defer resp.Body.Close()

	// Read response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %s", err)
	}

	// Display response
	if verbose {
		fmt.Printf("Response Status: %s\n", resp.Status)
		fmt.Printf("Response Headers: %v\n", resp.Header)
	}
	fmt.Printf("Response Body: %s\n", string(bodyBytes))
}
