## Context

ubox-crosser is a TCP-only reverse-proxy server. It currently has no HTTP surface. The proxy server (`server.ProxyServer`) accepts TCP connections through a `Dispatcher`, manages client controllers, and bridges tunnelled traffic. Configuration is loaded via `models/config.ServerConfig` and bound through Cobra flags in `cmd/server/server.go`.

Container orchestrators and load balancers need a simple HTTP probe to determine liveness. Today, the only option is a raw TCP connect check on the control port, which cannot distinguish "process is listening" from "process is healthy."

## Goals / Non-Goals

**Goals:**
- Provide a `GET /health` HTTP endpoint that returns `200 {"status":"ok"}` when the server process is running.
- Make the health-listen address configurable with a sensible default (`:8080`).
- Keep the implementation minimal — stdlib `net/http` only, no new dependencies.

**Non-Goals:**
- Deep readiness checks (database connectivity, upstream reachability) — this is a liveness probe only.
- Health endpoints on the client or auth-server binaries.
- Metrics, Prometheus scraping, or any observability beyond health.
- TLS on the health endpoint.

## Decisions

### 1. Stdlib `net/http` server — no framework

**Choice:** Use `net/http.ServeMux` + `http.ListenAndServe`.

**Rationale:** The endpoint is trivial (one route, static JSON). Adding a router framework (Chi, Gin) would introduce a dependency with no benefit. Stdlib keeps the binary lean and the dependency graph unchanged.

**Alternatives considered:**
- *Chi router* — unnecessary for a single route; adds a dependency.
- *Embed health in the existing TCP dispatcher* — would require custom protocol support on the probe side; defeats the purpose of HTTP compatibility.

### 2. Separate goroutine, separate port

**Choice:** Start the HTTP server in its own goroutine on a dedicated port (default `:8080`), independent of the TCP control/data listeners.

**Rationale:** Health probes must not compete with proxy traffic. A separate port makes firewall and security-group rules straightforward (expose only `:8080` to the orchestrator, keep proxy ports internal).

**Alternatives considered:**
- *Multiplex HTTP and TCP on the same port* (e.g., `cmux`) — adds complexity and a new dependency for no user benefit.

### 3. Configuration via `ServerConfig.HealthAddr`

**Choice:** Add a `HealthAddr string` field to `models/config.ServerConfig`, exposed as `--health-addr` CLI flag with default `:8080`.

**Rationale:** Follows the existing pattern (`ControlAddr`, `TunnelAddr`) and lets operators pick a port that avoids collisions.

### 4. Lifecycle: start with server, log errors, non-fatal

**Choice:** Launch the health HTTP server inside `ProxyServer.Run()`. If the health listener fails to bind, log the error but do not crash the proxy — the core proxy function is more important than the health endpoint.

**Rationale:** A health-check failure should degrade observability, not availability.

## Risks / Trade-offs

- **[Port collision]** Default `:8080` may conflict with other services on the host. → Mitigation: configurable via `--health-addr`; document the default.
- **[No graceful shutdown]** The health HTTP server will not have explicit `Shutdown()` handling in this iteration. → Acceptable: the proxy itself also lacks graceful shutdown; can be addressed holistically later.
