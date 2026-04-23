# REQ-final14-1776975275: /buildinfo Endpoint

## Summary

Add a `GET /buildinfo` HTTP endpoint to the ubox-crosser server that returns build metadata
as JSON: `git_sha` (7-char SHA injected at build time), `build_id` (env `BUILD_ID` or `"dev"`),
and `go_version` (`"go1.23"`). The endpoint is unauthenticated and lives on the existing
management HTTP server (`:8080`).

## Design

- Add `server/management.go`: `ManagementServer` struct with HTTP mux handling `/healthz` and `/buildinfo`.
- Update `cmd/server/server.go`: add `var GitSHA = "dev"` (overridden via ldflags at build time),
  `--management-addr` flag (default `:8080`), start `ManagementServer` goroutine.
- Update `Makefile`: add `ci-test` and `dev-cross-check` targets.
- Unit tests via `net/http/httptest` — no running server needed.

## Trade-offs

- Single management HTTP server hosts both `/healthz` and `/buildinfo` to avoid opening a second port.
- `go_version` is hardcoded to `"go1.23"` per the spec — avoids `runtime.Version()` which returns the
  full toolchain string and would be brittle in tests.
