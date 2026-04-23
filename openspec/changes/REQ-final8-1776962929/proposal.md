# REQ-final8-1776962929: /buildinfo Endpoint

## Summary

Expose a `/buildinfo` HTTP endpoint on the existing health-check HTTP server (`:8080`)
that returns build metadata: git commit SHA, CI build ID, and Go runtime version.

## Motivation

Operators need a lightweight, unauthenticated way to confirm which exact build is running
without grepping logs or exec-ing into the container.

## Design

- Reuse the HTTP server already planned for `/healthz` (port 8080, flag `--health-addr`).
- `git_sha` is injected at compile time via `-X main.GitSHA=$(git rev-parse --short HEAD)`.
- `build_id` is read from the `BUILD_ID` environment variable at startup; defaults to `"dev"`.
- `go_version` is hardcoded `"go1.23"` (matches the module's `go` directive).
- The endpoint requires no authentication.

## Scope

Single repo: `phona/ubox-crosser`. Changes touch `server/health.go` (new file),
`server/health_test.go` (new file), and `cmd/server/server.go` (wired up).
