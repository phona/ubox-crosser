## Context

ubox-crosser is a TCP-only SOCKS5 reverse proxy. There is no HTTP server in the codebase today. The project uses Go 1.23, Cobra for CLI, and has no third-party HTTP framework. The build already uses `-ldflags="-s -w"` but does not inject version metadata.

## Goals / Non-Goals

**Goals:**
- Expose a `GET /version` JSON endpoint on the proxy-server binary.
- Inject `commit` (git short SHA) and `build_time` (ISO 8601) at compile time via `-ldflags -X`.
- Keep the HTTP surface minimal: one endpoint, stdlib only.
- Unit-test the handler; contract-test against an OpenAPI schema.

**Non-Goals:**
- Adding HTTP to the client or auth_server binaries.
- Building a general-purpose HTTP API framework or middleware pipeline.
- Health-check, readiness, or metrics endpoints (future work).
- TLS for the HTTP listener (operational concern, out of scope).

## Decisions

### 1. Stdlib `net/http` over Chi/Gin

**Choice:** Use `net/http.ServeMux` (Go 1.22+ pattern routing).

**Alternatives considered:**
- **Chi**: Adds a dependency for a single route; not justified.
- **Gin**: Heavy; brings reflect, validation, binding machinery we don't need.

**Rationale:** Zero dependencies. Go 1.22+ `ServeMux` supports method-based routing (`GET /version`), which is sufficient.

### 2. Shared `internal/version` package

**Choice:** Place `Version`, `Commit`, `BuildTime` as package-level `var`s in `internal/version/version.go`.

**Alternatives considered:**
- Vars in `cmd/server/main.go`: Ties version info to one binary; can't reuse if other binaries need it later.
- Build-time code generation: Over-engineered for three strings.

**Rationale:** `internal/` keeps it project-private. Any binary can import it. `-ldflags -X` targets the full package path.

### 3. HTTP listener lifecycle — shared with REQ-601

**Choice:** Register `/version` on REQ-601's health HTTP mux (`newHealthMux()`), configured via `--health-address` (default `:8080`). No new listener or flag.

**Alternatives considered:**
- Separate `--http-addr` listener: Adds a second port for operators to manage; both default to `:8080` causing conflicts.
- Embedding in the existing TCP accept loop: Mixing protocols on the same listener adds complexity.
- Separate binary: Over-engineering; version info belongs to the binary it describes.

**Rationale:** `/version` and `/healthz` are both operational endpoints. Sharing the health mux avoids extra port and configuration burden. REQ-601's health listener already runs in an independent goroutine with fail-fast on bind error.

> **Updated in dev-spec stage:** Original design proposed a standalone `--http-addr` listener. Stage 1 review confirmed sharing REQ-601's health mux is the better approach.

### 4. Contract test approach — stdlib only

**Choice:** Go test using `net/http/httptest` + hand-written JSON assertions. No OpenAPI validator library.

**Alternatives considered:**
- `kin-openapi` schema validator: Adds a test-only dependency to `go.sum` for a single endpoint with a trivial schema.
- Shell-based `curl` + `yq` validation: Fragile, not portable.
- Schemathesis / Dredd: External tool dependency; heavy for one endpoint.

**Rationale:** The contract has exactly one endpoint with three string fields. Hand-written assertions are sufficient and avoid any new dependency. `httptest.NewServer` provides a full HTTP round-trip without network overhead.

> **Updated in dev-spec stage:** Original design proposed `kin-openapi`. Stage 1 review confirmed hand-written assertions are adequate given the minimal schema.

## Risks / Trade-offs

- **REQ-601 dependency** → `/version` cannot be implemented until REQ-601's health listener (`newHealthMux()`) is merged. Mitigation: Stage 2 marks this as a prerequisite; the endpoint is additive and won't conflict.
- **Shared mux coupling** → Changes to the health mux affect `/version`. Mitigation: both are simple operational endpoints with no middleware; coupling risk is low.
- **Version defaults to `unknown`** → If built without ldflags (e.g., `go run`), `commit` and `build_time` show `"unknown"`. Mitigation: acceptable for development; CI always injects real values.
