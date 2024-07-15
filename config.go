package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Endpoint represents a single HTTP endpoint to monitor.
type Endpoint struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	Method  string `json:"method,omitempty"`
	Timeout int    `json:"timeout_seconds,omitempty"`
}

// Config holds the full configuration for the health checker.
type Config struct {
	Endpoints []Endpoint `json:"endpoints"`
	Interval  int        `json:"interval_seconds,omitempty"`
}

// LoadConfig reads and parses a JSON config file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file: %w", err)
	}

	if len(cfg.Endpoints) == 0 {
		return nil, fmt.Errorf("no endpoints defined in config")
	}

	// Apply defaults
	for i := range cfg.Endpoints {
		if cfg.Endpoints[i].Method == "" {
			cfg.Endpoints[i].Method = "GET"
		}
		if cfg.Endpoints[i].Timeout == 0 {
			cfg.Endpoints[i].Timeout = 10
		}
	}

	if cfg.Interval == 0 {
		cfg.Interval = 30
	}

	return &cfg, nil
}
