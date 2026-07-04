# HTTP Health Checker

A lightweight HTTP endpoint health monitoring tool written in Go. Reads endpoint definitions from a JSON config file and continuously checks their availability, reporting status, response times, and errors with colored console output.

## Features

- Monitor multiple HTTP endpoints from a JSON config file
- Concurrent health checks across all endpoints
- Configurable check interval, timeout, and HTTP method per endpoint
- Colored console output (green for UP, red for DOWN)
- Single-run mode with exit code 1 if any endpoint is down
- Graceful shutdown on SIGINT/SIGTERM

## Build

```bash
go build -o http-health-checker .
```

## Usage

```bash
# Run with default config file (endpoints.json)
./http-health-checker

# Specify a custom config file
./http-health-checker -config /path/to/config.json

# Override check interval to 10 seconds
./http-health-checker -interval 10

# Run once and exit (useful for CI/CD)
./http-health-checker -once
```

### Flags

| Flag        | Default           | Description                              |
|-------------|-------------------|------------------------------------------|
| `-config`   | `endpoints.json`  | Path to endpoints config file            |
| `-interval` | from config (30s) | Check interval in seconds                |
| `-once`     | `false`           | Run checks once and exit                 |

### Config File Format

```json
{
  "endpoints": [
    {
      "name": "My API",
      "url": "https://api.example.com/health",
      "method": "GET",
      "timeout_seconds": 5
    },
    {
      "name": "Dashboard",
      "url": "https://dashboard.example.com",
      "method": "HEAD",
      "timeout_seconds": 10
    }
  ],
  "interval_seconds": 30
}
```

### Example Output

```
Health Checker started. Monitoring 3 endpoint(s) every 30s

=== Health Check Results ===
ENDPOINT                  STATUS   CODE   RESPONSE     ERROR
--------------------------------------------------------------------------------
Example Site              UP       200    123ms
Example API               UP       200    456ms
Nonexistent Host          DOWN     -      0s           request failed: dial tcp...

Summary: 2/3 endpoints up, 1 down
```



<sub><sup>Originally developed and tested locally during learning. Later organized and pushed to GitHub for portfolio visibility.</sup></sub>
