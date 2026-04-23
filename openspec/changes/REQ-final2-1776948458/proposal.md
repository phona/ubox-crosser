# REQ-final2-1776948458: /buildinfo HTTP Endpoint

## Summary

Expose a `/buildinfo` HTTP endpoint on the ubox-crosser `server` binary's admin HTTP
listener. It returns build provenance as JSON: the 7-char git SHA, a build
identifier (env `BUILD_ID` or `"dev"`), and the Go toolchain version (`"go1.23"`).
The endpoint is unauthenticated and intended for operators and CI smoke tests to
confirm exactly which artefact is running.

## Motivation

Final-phase e2e verification of the sisyphus multi-repo refactor needs a
machine-readable way to pin down which build a running server is. Log lines are
ad-hoc; a stable `GET /buildinfo` lets acceptance tests assert on the identity
of the binary without parsing logs or shelling into containers.

## In scope

- `cmd/server` binary — adds a `/buildinfo` handler mounted on the existing
  admin HTTP listener (the one the `/healthz` REQ introduces on port `8080`).
- `Dockerfile` / `Makefile` build — inject `main.GitSHA` through
  `-ldflags "-X main.GitSHA=$(git rev-parse --short HEAD)"`.
- Unit test for the handler (pure `httptest` — no listener).
- Integration test in `tests/acceptance/` (or equivalent) that stands up the
  docker stack and asserts `GET /buildinfo` returns 200 with the three fields.

## Out of scope

- `cmd/client` and `cmd/auth_server` binaries — they don't expose an admin HTTP
  surface today; keep blast radius minimal.
- New HTTP port or new listener. `/buildinfo` rides on the same admin listener
  `/healthz` uses (`:8080` per `tests/acceptance/docker-compose.yml`).
- Authentication / TLS on the admin surface — both `/healthz` and `/buildinfo`
  stay open.
- Richer build metadata (build time, dirty flag, Go runtime info). If that's
  wanted later, extend the JSON — don't add a second endpoint.

## Dependency on REQ-e2e-1776916220 (/healthz)

This REQ assumes the admin HTTP listener and its mux (introduced by the
`/healthz` REQ) already exist at implementation time. If `/healthz`'s code has
not yet landed when dev-agent starts, dev-agent bootstraps the listener here and
registers both handlers; otherwise it only adds `/buildinfo` to the existing mux.
See `design.md` for the branching plan.

## Feasibility

All three data sources are trivially available:

- `git_sha`: compile-time `-ldflags -X` — standard Go pattern, already used in
  many Go admin endpoints.
- `build_id`: `os.Getenv("BUILD_ID")` at request time, default `"dev"`.
- `go_version`: string literal `"go1.23"` (matches `Makefile` `ci-env` line 115).

No external deps. Handler is a few dozen lines. Supportable.
