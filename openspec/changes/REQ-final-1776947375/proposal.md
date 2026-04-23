# REQ-final-1776947375: /buildinfo HTTP Endpoint

## Summary

Add a `/buildinfo` HTTP endpoint to the ubox-crosser **server** binary
(`cmd/server`) that returns a small JSON document identifying the running
build: git SHA (short), build identifier, and Go toolchain version. The
endpoint is unauthenticated and reuses the server's existing HTTP listener
(no new port).

## Motivation

This change is the final end-to-end (E2E) acceptance run for the sisyphus
multi-repo refactor (PR #31). It is intentionally a minimal, closed,
mechanically verifiable requirement covering the full pipeline
`intent:analyze → … → DONE` in one small surface.

Operationally, the endpoint also gives downstream ops/monitoring a
zero-auth way to confirm which artifact is actually running on a given
host — useful for bisecting bad rollouts without shelling into the pod.

## Scope

**In scope**
- Single source repo: `phona/ubox-crosser`
- New HTTP handler `GET /buildinfo` on the existing server listener
- `ldflags` injection of `main.GitSHA` at build time (Makefile + Dockerfile)
- Runtime read of `BUILD_ID` env var (default `"dev"`)
- Hardcoded `go_version` string `"go1.23"` in the response
- Unit test for the handler
- Integration / acceptance test via `tests/acceptance/docker-compose.yml`
  asserting `GET /buildinfo` returns 200 + required JSON fields

**Out of scope**
- Authentication / authorization (endpoint is explicitly public)
- New listeners or port configuration
- Structured build metadata beyond the three agreed fields
- Arbitrary version negotiation / API versioning

## Response Contract

```json
{
  "git_sha":    "<7-char short SHA, injected via ldflags>",
  "build_id":   "<env BUILD_ID, or 'dev' if unset>",
  "go_version": "go1.23"
}
```

- HTTP 200, `Content-Type: application/json`
- Unauthenticated
- Same listener as `/healthz` (port 8080 per current docker-compose)

## Support Judgment

**Supportable as specified.** The change is small, contained, and does not
interact with the tunnel / proxy data plane. One non-blocking nuance worth
calling out in `design.md`: the prompt references "existing `/version`
endpoint" — there is currently no `/version` handler in-tree. The closest
precedent is the in-flight `/healthz` work (REQ-e2e-1776916220). Resolution
documented in `design.md` §"Nuance on the 'existing /version' reference".

## Repo Layout Impact

**spec home repo:** `phona/ubox-crosser` (single-repo REQ — it is also the
only source repo).

Files touched by this REQ's dev stage (illustrative, not a hard contract):
- `server/buildinfo.go` (new) — handler + `GitSHA` var declaration
- `server/buildinfo_test.go` (new) — unit test for the handler
- `cmd/server/server.go` — wire handler onto the HTTP mux
- `Makefile` — add `-X github.com/phona/ubox-crosser/server.GitSHA=$(shell git rev-parse --short HEAD)` to the server build ldflags
- `Dockerfile` — pass the same `-X ...=<build-time SHA>` to the `go build` step (`ARG GIT_SHA`)
- `tests/acceptance/buildinfo_test.go` (new) — integration test
- `tests/acceptance/docker-compose.yml` — extend `test-client` to also curl `/buildinfo`
- `openspec/changes/REQ-final-1776947375/**` — this change package
