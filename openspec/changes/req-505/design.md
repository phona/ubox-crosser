## Context

ubox-crosser is a Go SOCKS5 reverse proxy that communicates entirely over raw TCP with JSON-encoded messages. There are no HTTP endpoints today. The server binary (`cmd/server`) creates a `ProxyServer` that listens for TCP connections and dispatches them by message type (LOGIN, GEN_WORKER, AUTHENTICATION).

The requirement is to add `GET /version` returning `v2`. This is the first HTTP endpoint in the project.

## Goals / Non-Goals

**Goals:**
- Expose `GET /version` returning the application version as JSON on an HTTP management port.
- Inject the version string at build time so it can be changed without code edits.
- Keep the management server optional — zero-config deployments continue to work as before.

**Non-Goals:**
- Full management API (health checks, metrics, admin operations) — out of scope for this change.
- HTTPS/TLS on the management port — can be added later if needed.
- Version endpoint on the auth_server or client binaries.

## Decisions

### 1. Stdlib `net/http` for the management server

Use Go's `net/http` with a simple `http.ServeMux`. No framework needed for a single endpoint. This avoids adding external dependencies.

**Alternative considered:** Embedding version info in the existing TCP message protocol. Rejected because operational tooling (load balancers, curl scripts) expects HTTP, and mixing version queries into the proxy protocol adds complexity.

### 2. Build-time version injection via `ldflags`

Define a `var Version string` in a `version` package (or top-level variable). Set it at build time with `-ldflags "-X github.com/phona/ubox-crosser/version.Version=v2"`. The Makefile `build` target gets a `VERSION` variable defaulting to `dev`.

**Alternative considered:** Reading version from a file at runtime. Rejected because it adds a runtime dependency on a file path and complicates container builds.

### 3. Optional management listener via config

Add `ManagementAddress` field to `ServerConfig` and a `--management-address` CLI flag. When empty, no HTTP server starts. This preserves backward compatibility.

### 4. Response format

```json
{"version":"v2"}
```

Simple JSON object. `Content-Type: application/json`. HTTP 200.

## Risks / Trade-offs

- **[Additional port exposure]** → Operators must open/firewall the management port. Mitigated by making it opt-in (no default address).
- **[No auth on management endpoint]** → The version endpoint leaks no sensitive data. Acceptable for now; future management endpoints may need auth.
