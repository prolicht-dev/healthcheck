package main

import (
	"net/http"
	"os"
	"time"
)

// main checks the given URL, if the response is not 200, it will return with exit code 1
// on success, exit code 0 will be returned
func main() {
	if len(os.Args) < 2 {
		os.Exit(1)
	}
	url := os.Args[1]
	var client = &http.Client{
		Timeout: time.Second * 2,
	}
	if resp, err := client.Get(url); err != nil || resp.StatusCode < 200 || resp.StatusCode > 299 {
		os.Exit(1)
	}
}
