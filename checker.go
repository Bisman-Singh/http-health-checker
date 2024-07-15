package main

import (
	"fmt"
	"net/http"
	"time"
)

// CheckResult holds the outcome of a single health check.
type CheckResult struct {
	Endpoint     Endpoint
	StatusCode   int
	ResponseTime time.Duration
	Up           bool
	Error        string
	CheckedAt    time.Time
}

// CheckEndpoint performs an HTTP request to the given endpoint and returns the result.
func CheckEndpoint(ep Endpoint) CheckResult {
	result := CheckResult{
		Endpoint:  ep,
		CheckedAt: time.Now(),
	}

	client := &http.Client{
		Timeout: time.Duration(ep.Timeout) * time.Second,
	}

	req, err := http.NewRequest(ep.Method, ep.URL, nil)
	if err != nil {
		result.Up = false
		result.Error = fmt.Sprintf("invalid request: %v", err)
		return result
	}

	req.Header.Set("User-Agent", "http-health-checker/1.0")

	start := time.Now()
	resp, err := client.Do(req)
	result.ResponseTime = time.Since(start)

	if err != nil {
		result.Up = false
		result.Error = fmt.Sprintf("request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	result.Up = resp.StatusCode >= 200 && resp.StatusCode < 400

	if !result.Up {
		result.Error = fmt.Sprintf("unhealthy status code: %d", resp.StatusCode)
	}

	return result
}

// CheckAll runs health checks on all endpoints and returns the results.
func CheckAll(endpoints []Endpoint) []CheckResult {
	results := make([]CheckResult, len(endpoints))
	done := make(chan struct{}, len(endpoints))

	for i, ep := range endpoints {
		go func(idx int, endpoint Endpoint) {
			results[idx] = CheckEndpoint(endpoint)
			done <- struct{}{}
		}(i, ep)
	}

	for range endpoints {
		<-done
	}

	return results
}
