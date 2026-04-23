# REQ-final6-1776957940: /buildinfo Endpoint

## Summary

Expose a `/buildinfo` HTTP endpoint on `cmd/server` that returns build metadata as JSON.

## Motivation

Operators need a machine-readable way to verify which git revision and CI build is running in production without reading container labels or restarting the process.

## Design

- Reuse the existing HTTP server already listening on `:8080` (introduced for `/healthz`).
- Add a `/buildinfo` handler that returns `{"git_sha","build_id","go_version"}`.
- `git_sha` is injected at build time via `-ldflags "-X main.GitSHA=<sha>"`.
- `build_id` is read from the `BUILD_ID` environment variable at runtime; defaults to `"dev"`.
- `go_version` is hard-coded `"go1.23"` to match the toolchain used.
- No authentication is required.

## Changes

| File | Change |
|------|--------|
| `server/http_server.go` | New – HTTP server with `/healthz` + `/buildinfo` handlers |
| `cmd/server/server.go` | Add `GitSHA` ldflag var; start HTTP server goroutine |
| `Dockerfile` | Accept `GIT_SHA` build arg; pass to ldflag |
| `Makefile` | `GIT_SHA` variable; updated `build` + `ci-build`; `ci-test` target |
| `server/http_server_test.go` | Unit tests for handlers |
| `tests/acceptance/buildinfo_test.go` | Acceptance tests (docker-compose) |
| `tests/acceptance/docker-compose.yml` | Pass `GIT_SHA` build arg |
