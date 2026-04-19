## Context

ubox-crosser is a SOCKS5-based reverse proxy written in Go. All existing communication uses encrypted TCP connections with custom JSON message framing. There is no HTTP server. The `cmd/server/server.go` entry point uses Cobra and starts a `ProxyServer` that listens on TCP.

## Goals / Non-Goals

**Goals:**
- Provide a `GET /version` HTTP endpoint returning `{"version": "v3"}`.
- Keep the HTTP server minimal — stdlib `net/http` only.
- Make the HTTP listen address configurable via CLI flag and config file.

**Non-Goals:**
- Full REST API or management plane (future work).
- Authentication or TLS on the version endpoint.
- Version sourced from build-time injection (hardcoded `v3` is sufficient for now).

## Decisions

### 1. Stdlib `net/http` over framework
Use Go's `net/http` package directly. The endpoint is trivial; adding a framework (Chi, Gin) would be over-engineering. If more endpoints are added later, a framework can be introduced then.

### 2. Separate HTTP listener
Run `http.ListenAndServe` on its own goroutine alongside the existing TCP listeners. This avoids mixing HTTP framing into the custom TCP protocol. A dedicated `--http-addr` flag (default `:8080`) controls the bind address.

### 3. Handler location
Add a thin `api/` package containing the version handler. This keeps HTTP concerns separated from the TCP proxy logic in `server/`.

### 4. Response format
Return `Content-Type: application/json` with body `{"version":"v3"}`. JSON is the conventional format for programmatic consumption and aligns with future API expansion.

## Risks / Trade-offs

- **New listening port**: Operators must expose an additional port. Mitigated by clear default (`:8080`) and CLI flag.
- **No graceful shutdown**: The HTTP server uses a simple goroutine without `context` cancellation. Acceptable for a single diagnostic endpoint; can be improved when more endpoints are added.
