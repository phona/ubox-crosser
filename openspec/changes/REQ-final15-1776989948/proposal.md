# REQ-final15-1776989948 Proposal

## Summary

Expose `GET /buildinfo` on the ubox-crosser server's management HTTP server, returning
build metadata as JSON: `git_sha` (7-char short SHA injected via ldflags), `build_id`
(from `BUILD_ID` env, defaulting to `"dev"`), and `go_version` (hardcoded `"go1.23"`).
The endpoint is unauthenticated.

## Approach

- Add `server/management.go` with a `ManagementServer` that serves `/healthz` and `/buildinfo`
  on a dedicated HTTP listener (default `:8080`, configurable via `--management-addr` flag).
- `GitSHA` is injected via `-ldflags "-X main.GitSHA=$(git rev-parse --short HEAD)"` at build time.
- `BUILD_ID` is read from the environment at server startup; defaults to `"dev"`.
- Unit tests use `net/http/httptest` — no running server required, safe for `go test ./...`.
- Acceptance tests (`tests/acceptance/buildinfo_test.go`) use `//go:build acceptance` to
  prevent execution during `go test ./...`.

## Spec home repo

phona/ubox-crosser
