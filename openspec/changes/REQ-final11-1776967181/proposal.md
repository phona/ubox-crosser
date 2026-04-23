# REQ-final11-1776967181: /buildinfo Endpoint

## Overview

Expose a `/buildinfo` HTTP endpoint on `cmd/server` that returns build metadata as JSON: git commit SHA, build pipeline identifier, and Go runtime version.

## Motivation

Operators need a lightweight, unauthenticated endpoint to verify which build artifact is running in any environment without shell access. The three fields cover the three most common "is this the right binary?" questions.

## Design

- Reuse the existing HTTP management server (same port as `/healthz`, default `:8080`)
- `git_sha` is injected at compile time via ldflags: `-X main.GitSHA=$(git rev-parse --short HEAD)`
- `build_id` is read from `$BUILD_ID` env var at request time; defaults to `"dev"` when unset
- `go_version` is hardcoded `"go1.23"` (matches the repo's `go.mod` language requirement)
- Endpoint requires no authentication

## Files Changed

| File | Change |
|------|--------|
| `server/http.go` | New: HTTP management server with `/buildinfo` + `/healthz` handlers |
| `server/http_test.go` | New: unit tests for buildinfo and healthz handlers |
| `cmd/server/server.go` | Add `GitSHA` ldflags var; start HTTP server on `:8080` |
| `models/config/config.go` | Add `HTTPAddr` field to `ServerConfig` |
| `Makefile` | Inject `GIT_SHA` ldflags in `build` target |
| `Dockerfile` | Pass `GIT_SHA` build-arg through ldflags |
