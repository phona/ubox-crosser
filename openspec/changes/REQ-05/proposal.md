# REQ-05: Add Health Endpoint

## What

Add an HTTP `GET /health` endpoint to the **proxy server** binary (`cmd/server`) that returns `{"status":"ok"}` with HTTP 200. This provides a standard health-check target for orchestrators (Docker, Kubernetes, load balancers).

## Why

ubox-crosser currently has no way for external systems to probe whether the proxy server is alive. The only existing liveness signal is the internal TCP heartbeat between client and server, which is not accessible to infrastructure tooling. Container orchestrators and monitoring systems need an HTTP endpoint to determine service health.

## Scope

- **In scope:** A new HTTP listener on the proxy server, serving a single `GET /health` route. Configurable listen address via CLI flag and config file.
- **Out of scope:** Health endpoints on the client or auth_server binaries. Deep health checks (database connectivity, downstream service status). Readiness vs. liveness distinction. Metrics or Prometheus endpoints.

## Approach

The proxy server is a pure TCP application with no existing HTTP stack. We will add a minimal `net/http` server using Go's standard library (no new dependencies) that runs alongside the existing TCP listeners. The HTTP health server will:

1. Listen on a configurable address (default `:8080`).
2. Serve `GET /health` returning `{"status":"ok"}` with `Content-Type: application/json` and HTTP 200.
3. Return HTTP 404 for all other paths.
4. Start as a goroutine during server initialization; a failure to bind the health port logs an error but does not prevent the main TCP server from starting.
