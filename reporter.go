package main

import (
	"fmt"
	"strings"
	"time"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// PrintResults displays the health check results with colored output.
func PrintResults(results []CheckResult) {
	fmt.Printf("\n%s%s=== Health Check Results ===%s\n", colorBold, colorCyan, colorReset)
	fmt.Printf("%-25s %-8s %-6s %-12s %s\n", "ENDPOINT", "STATUS", "CODE", "RESPONSE", "ERROR")
	fmt.Println(strings.Repeat("-", 80))

	for _, r := range results {
		status := fmt.Sprintf("%s%sUP%s", colorBold, colorGreen, colorReset)
		if !r.Up {
			status = fmt.Sprintf("%s%sDOWN%s", colorBold, colorRed, colorReset)
		}

		code := fmt.Sprintf("%d", r.StatusCode)
		if r.StatusCode == 0 {
			code = "-"
		}

		responseTime := r.ResponseTime.Round(time.Millisecond).String()

		errMsg := ""
		if r.Error != "" {
			errMsg = truncate(r.Error, 40)
		}

		name := truncate(r.Endpoint.Name, 24)
		fmt.Printf("%-25s %-18s %-6s %-12s %s\n", name, status, code, responseTime, errMsg)
	}

	fmt.Println()
}

// PrintSummary prints a summary line with counts.
func PrintSummary(results []CheckResult) bool {
	up := 0
	down := 0
	for _, r := range results {
		if r.Up {
			up++
		} else {
			down++
		}
	}

	total := len(results)
	if down > 0 {
		fmt.Printf("%s%sSummary: %d/%d endpoints up, %d down%s\n", colorBold, colorRed, up, total, down, colorReset)
		return false
	}

	fmt.Printf("%s%sSummary: %d/%d endpoints up, all healthy%s\n", colorBold, colorGreen, up, total, colorReset)
	return true
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
