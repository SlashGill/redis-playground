package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL      = "http://localhost:9999/nPsaI7"
	numRequests = 1000
)

func makeRequest(wg *sync.WaitGroup, client *http.Client, requestNumber int) {
	defer wg.Done()

	startTime := time.Now()
	resp, err := client.Get(apiURL)
	duration := time.Since(startTime)

	if err != nil {
		fmt.Printf("apiURL %s Request %d failed: %v (Duration: %v)\n", apiURL, requestNumber, err, duration)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request %d successful. Status Code: %d (Duration: %v)\n", requestNumber, resp.StatusCode, duration)
}

func main() {
	var wg sync.WaitGroup
	client := &http.Client{} // 建立一個共用的 HTTP Client

	fmt.Printf("Sending %d concurrent requests to %s...\n", numRequests, apiURL)
	startTime := time.Now()

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go makeRequest(&wg, client, i+1)
	}

	wg.Wait() // 等待所有 Goroutine 完成
	duration := time.Since(startTime)
	fmt.Printf("All requests completed. total duration: %v", duration)
}
