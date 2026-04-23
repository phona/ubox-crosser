# /healthz Endpoint Specification

## Overview

The ProxyServer must expose an HTTP `/healthz` endpoint on a dedicated HTTP listener port that returns the service uptime.

## API Contract

### Endpoint

```
GET /healthz
```

### Request

No body required.

### Response

**Status:** `200 OK`

**Content-Type:** `application/json`

**Body:**
```json
{
  "status": "healthy",
  "uptime": {
    "seconds": 3661,
    "duration": "1h1m1s"
  }
}
```

### Response Fields

- `status` (string): Always `"healthy"` (reserved for future expansion)
- `uptime.seconds` (integer): Service uptime in seconds since ProxyServer start
- `uptime.duration` (string): Human-readable duration format (Go `time.Duration` string)

## Implementation Requirements

1. **HTTP Listener:** ProxyServer must bind to a configurable HTTP port (env: `PROXY_HTTP_PORT`, default: `8080`)
2. **Uptime Tracking:** ProxyServer tracks start time when initialized; uptime is calculated as `time.Now() - startTime`
3. **Non-Blocking:** HTTP listener runs in a separate goroutine; does not block message processing
4. **Graceful Degradation:** If HTTP listener fails to start, ProxyServer continues processing messages

## Scenarios

### REQ-e2e-1776916220-S1: Healthz Returns Valid Response

**Given** ProxyServer is running
**When** client sends `GET /healthz`
**Then** response status is 200 and body contains valid JSON with `status="healthy"` and `uptime.seconds >= 0`

### REQ-e2e-1776916220-S2: Uptime Increases Over Time

**Given** ProxyServer is running and has processed messages for 2+ seconds
**When** client sends two consecutive `GET /healthz` requests with 1 second delay
**Then** second request's `uptime.seconds` > first request's `uptime.seconds`

### REQ-e2e-1776916220-S3: Duration Format Valid

**Given** ProxyServer uptime is 3661 seconds (1h1m1s)
**When** client sends `GET /healthz`
**Then** response contains `uptime.duration = "1h1m1s"`

## Integration Test Environment

- HTTP listener listens on `PROXY_HTTP_ADDR=localhost:8080` (via env variable)
- Tests run in Docker Compose with ProxyServer service
- Test validates JSON schema, numeric fields, and string formatting
