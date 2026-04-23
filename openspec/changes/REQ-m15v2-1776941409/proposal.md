# REQ-m15v2-1776941409: /buildinfo Endpoint

## Why
End-to-end validation of the sisyphus M15 verifier-stage rename (`dev → dev_cross_check`, `spec → spec_lint`). We need a minimal, self-contained, mechanically-verifiable feature so the rename can be exercised through the full analyze → spec → dev → staging-test → pr-ci → accept pipeline.

## What
Add an unauthenticated `GET /buildinfo` HTTP endpoint to the `server` subcommand (`cmd/server`) of ubox-crosser that returns:

```json
{
  "git_sha":    "<7-char short commit SHA>",
  "build_id":   "<env BUILD_ID, default 'dev'>",
  "go_version": "go1.23"
}
```

## Scope

### In scope
- One new HTTP route, `GET /buildinfo`, served on the **same listener** the server already exposes for `/healthz` (no new port, no new listener).
- Build-time injection of `git_sha` via Go ldflags `-X main.GitSHA=$(git rev-parse --short HEAD)`.
- Runtime read of `BUILD_ID` env var, default `"dev"`.
- Hardcoded `go_version: "go1.23"` (no `runtime.Version()`; spec-level requirement).
- Unit test covering the handler.
- Integration / acceptance test against the real docker stack: `GET /buildinfo` → 200 with the three required JSON fields.

### Out of scope
- No changes to client / auth_server / proxy logic.
- No auth, rate limiting, or caching on the new endpoint.
- No changes to `/healthz`.
- No new config options (the endpoint sits on the existing health-check port).

## Affected components
- **leader source repo:** `phona/ubox-crosser` (single-repo task; no integration repo needed)
- **Go packages:** `cmd/server` (handler registration, ldflags var), `tests/acceptance` (new acceptance test)
- **Build tooling:** `Makefile` / `Dockerfile` if they currently build the server binary — they need the `-X main.GitSHA=...` ldflag wired in for the integration test to assert a non-empty 7-char SHA.

## Supportability
**Supported.** The repo already exposes an HTTP server inside `cmd/server` (REQ-e2e-1776916220 added `/healthz` on a dedicated health-check port with a docker-compose-based acceptance test). `/buildinfo` is a strictly additive route on the same listener and reuses the same acceptance-test scaffolding. No architectural change required.

## Note on the prompt's "/version" reference
The intent prompt says "参考已有 /version 端点的实现路径". There is **no `/version` endpoint** in the current tree — the closest analog is `/healthz` (REQ-e2e-1776916220). The `/healthz` path is the de-facto reference and what `dev-agent` should follow.
