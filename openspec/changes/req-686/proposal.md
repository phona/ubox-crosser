---
change_id: req-686
title: "Add GET /ping endpoint"
layers:
  - backend
status: draft
---

## Why

A minimal connectivity check endpoint is needed so that clients and monitoring tools can verify the admin HTTP server is reachable with the lowest possible overhead. Unlike `/healthz` (liveness semantics) or `/version` (build metadata), `/ping` is a simple echo that confirms TCP + HTTP connectivity without implying any health contract.

## What Changes

- Add a `ping` package with an HTTP handler returning plain-text `pong` on `GET /ping`
- Wire the handler into the admin `http.ServeMux` alongside the existing `/version` and `/healthz` routes
- Go 1.22+ method-based routing ensures only GET is accepted; other methods receive 405

## Capabilities

### New Capabilities
- `ping-endpoint`: HTTP GET /ping returning plain-text `pong` with HTTP 200

### Modified Capabilities

## Impact

- `ping/` — new package (handler.go, handler_test.go)
- `cmd/server/server.go` — register `GET /ping` on existing admin mux
- No new external dependencies (stdlib only)
