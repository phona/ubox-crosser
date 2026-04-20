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

### 3. HTTP listener lifecycle

**Choice:** Start the HTTP server in a separate goroutine from `cmd/server/main.go`, on a configurable `--http-addr` flag (default `:8080`). Log fatal if bind fails.

**Alternatives considered:**
- Embedding in the existing TCP accept loop: Mixing protocols on the same listener adds complexity.
- Separate binary: Over-engineering; version info belongs to the binary it describes.

**Rationale:** Independent goroutine keeps HTTP decoupled from the TCP proxy path. Fail-fast on bind error prevents silent misconfiguration.

### 4. Contract test approach

**Choice:** Go test that loads `contract.spec.yaml`, makes an HTTP request to the handler, and validates the response against the OpenAPI 3.0 schema using `kin-openapi` (or equivalent).

**Alternatives considered:**
- Shell-based `curl` + `yq` validation: Fragile, not portable.
- Schemathesis / Dredd: External tool dependency; heavy for one endpoint.

**Rationale:** Keeps validation in Go test toolchain. `kin-openapi` is a well-maintained OpenAPI 3 validator.

## Risks / Trade-offs

- **New listening port** → Operators must open/configure a second port. Mitigation: default to `:8080`, document clearly.
- **`kin-openapi` dependency for tests only** → Adds to `go.sum`. Mitigation: test-only import; does not affect binary size.
- **Version defaults to `dev`** → If built without ldflags (e.g., `go run`), version fields show defaults. Mitigation: acceptable for development; CI always injects real values.
