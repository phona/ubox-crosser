# Capability Spec: Health Endpoint

## Overview

The proxy server exposes an HTTP health-check endpoint for external orchestrators and monitoring systems.

## Behavior

### CAP-01: GET /health returns 200 with status ok

- **Given** the proxy server is running and the health HTTP listener has started successfully
- **When** a client sends `GET /health`
- **Then** the server responds with:
  - HTTP status code: `200`
  - Header: `Content-Type: application/json; charset=utf-8`
  - Body: `{"status":"ok"}`

### CAP-02: Non-health paths return 404

- **Given** the health HTTP listener is running
- **When** a client sends a request to any path other than `/health` (e.g., `GET /`, `GET /ready`, `POST /health`)
- **Then** the server responds with HTTP status code `404`

### CAP-03: Health address is configurable via CLI flag

- **Given** the server binary is invoked with `--health-address <addr>`
- **Then** the health HTTP listener binds to `<addr>` instead of the default `:8080`

### CAP-04: Health address is configurable via config file

- **Given** the server config file contains `"health_address": "<addr>"`
- **Then** the health HTTP listener binds to `<addr>` instead of the default `:8080`
- CLI flag takes precedence over config file value.

### CAP-05: Health listener failure does not crash the server

- **Given** the health address port is already in use or otherwise unavailable
- **When** the proxy server starts
- **Then** the main TCP listeners start normally, and an error is logged indicating the health listener failed to start

### CAP-06: Health endpoint responds independently of TCP state

- **Given** the health HTTP listener is running
- **When** no TCP clients are connected, or TCP listeners have errors
- **Then** `GET /health` still returns `200 {"status":"ok"}`

## Non-Requirements

- No authentication on the health endpoint.
- No request body parsing.
- No rate limiting.
- No graceful shutdown coordination (health endpoint stops when process exits).
- No metrics or Prometheus-format output.
