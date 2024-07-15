package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	configPath := flag.String("config", "endpoints.json", "path to endpoints config file")
	interval := flag.Int("interval", 0, "check interval in seconds (overrides config)")
	once := flag.Bool("once", false, "run checks once and exit")
	flag.Parse()

	cfg, err := LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	if *interval > 0 {
		cfg.Interval = *interval
	}

	fmt.Printf("Health Checker started. Monitoring %d endpoint(s) every %ds\n", len(cfg.Endpoints), cfg.Interval)

	if *once {
		results := CheckAll(cfg.Endpoints)
		PrintResults(results)
		allHealthy := PrintSummary(results)
		if !allHealthy {
			os.Exit(1)
		}
		return
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Duration(cfg.Interval) * time.Second)
	defer ticker.Stop()

	// Run immediately on start
	results := CheckAll(cfg.Endpoints)
	PrintResults(results)
	PrintSummary(results)

	for {
		select {
		case <-ticker.C:
			results := CheckAll(cfg.Endpoints)
			PrintResults(results)
			PrintSummary(results)
		case sig := <-sigChan:
			fmt.Printf("\nReceived %v, shutting down...\n", sig)
			os.Exit(0)
		}
	}
}
