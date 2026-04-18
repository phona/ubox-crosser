# Dev Spec: /api/metrics Endpoint

## Overview

Add an HTTP `/api/metrics` endpoint to the proxy server that exposes runtime metrics in JSON format. This provides observability into the running server without requiring external monitoring agents.

## Implementation Details

### 1. Metrics Data Model (`server/metrics.go`)

Create a thread-safe metrics collector using `sync/atomic` counters:

- `uptime_seconds` (float64): seconds since server start
- `active_connections` (int64): currently active tunnel connections
- `total_connections` (int64): cumulative tunnel connections since start
- `active_controllers` (int64): currently logged-in client controllers
- `total_logins` (int64): cumulative login attempts
- `total_errors` (int64): cumulative error count

### 2. HTTP Handler (`server/metrics_handler.go`)

- Endpoint: `GET /api/metrics`
- Response: `application/json; charset=utf-8`
- HTTP 200 with JSON body containing all metrics fields
- HTTP 405 for non-GET methods

### 3. Server Integration

- Add optional `metrics_address` field to `ServerConfig` (e.g., `":9100"`)
- `ProxyServer` starts an HTTP server on `metrics_address` if configured
- Instrument connection lifecycle events to update counters

### 4. Configuration

Add `metrics_address` to server config JSON:

```json
{
  "common": {
    "metrics_address": ":9100"
  }
}
```

If `metrics_address` is empty or unset, the metrics HTTP server is not started (backward compatible).

## Testing

- Unit tests for metrics counter operations (increment, decrement, snapshot)
- Unit tests for HTTP handler (correct JSON response, correct content-type, method filtering)
