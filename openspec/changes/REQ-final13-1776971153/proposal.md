# REQ-final13-1776971153: /buildinfo Endpoint

## Summary

Add a `/buildinfo` HTTP endpoint to the ubox-crosser management server that returns build
metadata as JSON: `git_sha`, `build_id`, and `go_version`.

## Approach

The server exposes a management HTTP server at `:8080` (configurable via `--management-addr`).
`/buildinfo` is served on this same server alongside `/healthz` — no new port is opened.

- `git_sha`: injected at build time via `-ldflags "-X main.GitSHA=<7-char>"`.
- `build_id`: read from `BUILD_ID` environment variable at startup; defaults to `"dev"`.
- `go_version`: hardcoded `"go1.23"` (matches the toolchain used to build the project).

## Files Changed

- `server/management.go` — ManagementServer with `/healthz` and `/buildinfo` handlers
- `server/management_test.go` — unit tests (internal, pure, no external services)
- `cmd/server/server.go` — add `var GitSHA`, `--management-addr` flag, start management goroutine
- `Makefile` — add `ci-test` and `dev-cross-check` targets; inject `GIT_SHA` via ldflags in `build`

## spec home repo
phona/ubox-crosser
