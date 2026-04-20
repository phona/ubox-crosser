---
change_id: REQ-642
title: "Design: GET /version endpoint"
---

# Design: GET /version endpoint

## Approach

Add a `version` package at the repo root level containing:
1. A `Version` constant (`0.1.0`)
2. `Commit` and `BuildTime` variables injected via `go build -ldflags -X`
3. An `http.HandlerFunc` that serializes version info as JSON

## Trade-offs Considered

### Where to host the HTTP server

**Option A (chosen):** Embed an HTTP mux in `cmd/server` alongside the existing TCP listener. Minimal code, no new dependencies.

**Option B:** Separate binary/service. Rejected — overkill for a single endpoint.

### Version source

**Option A (chosen):** Constant in Go source + ldflags for commit/build_time. Standard Go pattern, zero runtime cost.

**Option B:** Read from a file or env var at startup. More flexible but adds failure modes.

## Dependencies

None — uses `net/http` from the standard library.
