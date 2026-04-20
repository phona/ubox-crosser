---
change_id: req-685
title: "Add GET /healthz endpoint"
layers:
  - backend
status: draft
---

## Why

Operators and orchestration platforms (Kubernetes, Docker health checks, load balancers) need a lightweight endpoint to verify that the ubox-crosser admin HTTP server is alive and accepting requests. The existing `/version` endpoint serves a different purpose (build metadata) and is not idiomatic for health probes.

## What Changes

- Add a `health` package with an HTTP handler returning `{"status":"ok"}` on `GET /healthz`
- Wire the handler into the admin `http.ServeMux` alongside the existing `/version` route
- Go 1.22+ method-based routing ensures only GET is accepted; other methods receive 405

## Capabilities

### New Capabilities
- `healthz-endpoint`: HTTP GET /healthz returning JSON `{"status":"ok"}` with HTTP 200

### Modified Capabilities

## Impact

- `health/` — new package (handler.go, handler_test.go)
- `cmd/server/server.go` — register `GET /healthz` on existing admin mux
- No new external dependencies (stdlib only)
