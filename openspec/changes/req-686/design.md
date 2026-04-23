---
change_id: req-686
title: "GET /ping endpoint (dev-spec) — Design"
---

## Context

ubox-crosser is a TCP-based proxy tunnel using a custom binary protocol. The admin HTTP listener on `--admin-addr` (default `:8080`) already serves `GET /version` returning JSON build metadata. REQ-685 adds `GET /healthz` for liveness probes.

Operators and load-balancers need the simplest possible endpoint to confirm the server process is reachable at the network level — a classic ping/pong pattern. Unlike `/healthz` (liveness semantics) or `/version` (build metadata), `/ping` is a zero-overhead echo that confirms TCP + HTTP connectivity without implying any health contract.

The `version` package (`version/handler.go`) provides the existing pattern: an exported `Handler` func with signature `func(http.ResponseWriter, *http.Request)`, registered on the admin mux in `cmd/server/server.go` via `mux.HandleFunc("GET /version", version.Handler)`.

## Goals / Non-Goals

**Goals:**
- Expose `GET /ping` on the existing admin HTTP listener returning plain-text `pong`
- Minimal overhead: no JSON encoding, no Content-Type negotiation, no struct allocation
- Accept only GET; reject other methods with 405 (Go 1.22+ method-based routing)
- Follow the established package-per-handler pattern (`version/`, `health/`, `ping/`)

**Non-Goals:**
- Health semantics (that is `/healthz`)
- Build metadata (that is `/version`)
- Round-trip latency measurement or timestamped responses
- Response body other than the exact string `pong`

## Decisions

### Decision 1: Plain-text response (not JSON)

A ping/pong endpoint is an echo check. Plain text avoids unnecessary JSON encoding overhead and is the idiomatic format for this pattern.

**Chosen:** Return `pong` as `text/plain; charset=utf-8` using `io.WriteString`.

| Option | Pros | Cons |
|--------|------|------|
| Plain text `pong` via `io.WriteString` | Minimal, idiomatic, zero-alloc | Not JSON like `/version` |
| JSON `{"ping":"pong"}` | Consistent with other endpoints | Unnecessary `json.Encoder` overhead for an echo |

### Decision 2: Separate `ping` package

Follow the same pattern as the `version` package — a small, self-contained package with `handler.go` (exported `Handler` func) and `handler_test.go`. Keeps admin endpoint handlers decoupled and independently testable.

**Implementation mapping:**
- `ping/handler.go` — `func Handler(w http.ResponseWriter, _ *http.Request)` sets `Content-Type: text/plain; charset=utf-8`, writes `pong`
- `ping/handler_test.go` — unit tests mirroring `version/handler_test.go` structure: status code, content-type, body content, method enforcement via mux, 405 for non-GET

### Decision 3: Registration on existing admin mux

Add `mux.HandleFunc("GET /ping", ping.Handler)` in `cmd/server/server.go` alongside the existing `GET /version` registration. Go 1.22+ `http.ServeMux` method-based routing ensures 405 for non-GET automatically.

## Risks / Trade-offs

- **[Plain text vs JSON]** → Breaks JSON consistency with `/version` and `/healthz`. Acceptable because ping/pong is universally understood as plain text and the simplicity is the point.
- **[No auth]** → Same as `/version` and `/healthz` — relies on network-level access control via `--admin-addr` binding.
- **[No response body variation]** → Always returns `pong` regardless of server state. This is intentional — state-aware responses belong in `/healthz`.
