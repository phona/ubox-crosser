---
change_id: req-653
title: "Add GET /version endpoint (v3)"
layers:
  - backend
status: draft
---

## Why

There is no way to query the running ubox-crosser server's version, git commit, or build timestamp at runtime. This makes it difficult to verify deployments and troubleshoot version mismatches across distributed proxy nodes.

## What Changes

- Add a `version` package exposing build-time metadata (`Version`, `Commit`, `BuildTime`)
- Add an HTTP `GET /version` handler returning JSON with version info
- Wire the handler into `cmd/server` via a new `--http-addr` flag (default `:8080`)
- Update `Makefile` and `Dockerfile` with `-ldflags` for compile-time metadata injection
- Add unit tests for the handler

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP GET /version returning JSON with version, commit, and build_time fields; non-GET returns 405

### Modified Capabilities

## Impact

- `version/` — new package (version.go, handler.go, handler_test.go)
- `cmd/server/server.go` — new `--http-addr` flag and HTTP server goroutine
- `Makefile` — ldflags variables and build target changes
- `Dockerfile` — build args for commit/build_time injection
- No new external dependencies (stdlib only)
