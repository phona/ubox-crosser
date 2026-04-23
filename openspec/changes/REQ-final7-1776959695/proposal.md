# REQ-final7-1776959695: /buildinfo HTTP Endpoint

## Summary

Expose a `/buildinfo` endpoint on the existing `cmd/server` HTTP server (port 8080)
that returns a JSON object with three build-time fields:

```json
{"git_sha": "<7-char>", "build_id": "<BUILD_ID env or 'dev'>", "go_version": "go1.23"}
```

## Motivation

Operators need a machine-readable way to verify which binary revision is running in
production without shelling into the container.  The three fields cover the three
identifiers that uniquely describe a deployed artifact: source revision, CI pipeline run,
and runtime version.

## Design

- Reuse the HTTP server (`:8080`) already started for `/healthz`; no new port or listener.
- `git_sha` injected at link time via `-X main.GitSHA=$(git rev-parse --short HEAD)`.
- `build_id` read from `BUILD_ID` env var at request time; defaults to `"dev"`.
- `go_version` hardcoded `"go1.23"` matching the module and toolchain constraint.
- No authentication required (same as `/healthz`).

## Files Changed

| File | Change |
|------|--------|
| `server/http.go` | New: `HealthzHandler`, `BuildInfoHandler`, `StartHTTPServer` |
| `server/http_test.go` | New: unit tests for both handlers |
| `cmd/server/server.go` | Add `var GitSHA string`; call `StartHTTPServer(":8080", GitSHA)` |
| `Makefile` | Add ldflags, `ci-test`, `ci-accept-env-up/down` targets |
| `Dockerfile` | Add `GIT_SHA` build arg; pass via ldflags for server binary |
| `tests/Dockerfile.test` | Add `GIT_SHA` build arg; pass via ldflags |
| `tests/acceptance/buildinfo_test.go` | New: acceptance tests |
| `tests/acceptance/docker-compose.yml` | Add `BUILD_ID`, test-runner service |
