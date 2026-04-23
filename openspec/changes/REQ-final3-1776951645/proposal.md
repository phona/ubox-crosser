# REQ-final3-1776951645: /buildinfo HTTP Endpoint

## Why

Operators and CI need a cheap, unauthenticated way to identify which exact
build of `ubox-crosser server` is running in a cluster or container. Today
there is no endpoint that exposes the git commit, build identifier, or Go
toolchain version — debugging a bad deploy requires SSH, filesystem
inspection, or correlating container image tags by hand.

This REQ adds a `GET /buildinfo` endpoint to `cmd/server` that returns a
JSON document with three fields:

```json
{"git_sha": "<7-char>", "build_id": "<env BUILD_ID or 'dev'>", "go_version": "go1.23"}
```

It is the first endpoint to introduce the `ldflags -X main.GitSHA=…`
injection pattern into this repo and ships the minimal admin HTTP listener
that subsequent REQs (notably REQ-e2e-1776916220 for `/healthz`) can
register additional handlers on.

## What changes

- **New HTTP admin listener** on port `8080` inside `cmd/server`. A single
  `*http.ServeMux` is wired up in `main()` and served in a goroutine
  alongside the existing TCP proxy loop.
- **New handler** `GET /buildinfo` returning a fixed JSON shape with
  `Content-Type: application/json`.
- **Build-time git SHA injection** via
  `-X main.GitSHA=$(git rev-parse --short HEAD)` threaded through the
  `Makefile` `build` / `ci-build` targets and the root `Dockerfile` (as a
  `GIT_SHA` build arg).
- **Unit tests** for the handler in `cmd/server/buildinfo_test.go` using
  `httptest`, covering status code, JSON shape, `BUILD_ID` env override,
  `GitSHA` default, and `go_version` literal.
- **Integration test** in `tests/integration/buildinfo_test.go`
  (`//go:build integration`) hitting `http://proxy-server:8080/buildinfo`
  through the existing docker-compose stack (`tests/docker-compose.yml`).
- **docker-compose** edit: expose container port `8080` inside the compose
  network and pass `BUILD_ID` / `GIT_SHA` down to the `proxy-server`
  service so the integration test can assert real values.

## Impact

- **Affected capabilities**: `buildinfo-endpoint` (new).
- **Affected code**: `cmd/server/`, `server/` (admin HTTP listener),
  `Dockerfile`, `Makefile`, `tests/docker-compose.yml`,
  `tests/integration/`.
- **No breaking changes** to the existing proxy protocol, control
  channel, or the `client` / `auth_server` binaries. Only `server`
  binary gains an extra listener.
- **Composes cleanly** with REQ-e2e-1776916220 (`/healthz`): both share
  the same admin mux helper. Whichever REQ lands first brings the mux
  up; the second just registers its handler.
