---
change_id: req-669
title: "Add GET /version endpoint (v4)"
layers:
  - backend
status: draft
---

## Why

There is no way to query the running ubox-crosser server's version, git commit, or build timestamp at runtime. Operators need a lightweight HTTP endpoint to verify which build is deployed on each proxy node without SSH access.

## What Changes

- Add a `version` package exposing build-time metadata (`Version`, `Commit`, `BuildTime`)
- Add an HTTP `GET /version` handler returning JSON with version info
- Wire the handler into `cmd/server` on a dedicated admin HTTP listener (`--admin-addr`, default `:8080`)
- Update `Makefile` and `Dockerfile` with `-ldflags` for compile-time metadata injection
- Go 1.22+ method-based routing ensures only GET is accepted; other methods receive 405

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP GET /version returning JSON with version, commit, and build_time fields

### Modified Capabilities

## Impact

- `version/` — new package (version.go, handler.go, handler_test.go)
- `cmd/server/server.go` — new admin HTTP listener with `--admin-addr` flag
- `Makefile` — ldflags variables and build target changes
- `Dockerfile` — build args for commit/build_time injection
- No new external dependencies (stdlib only)
