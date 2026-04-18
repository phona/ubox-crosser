## Context

The ubox-crosser proxy server (`server/server.go`) manages TCP connections via a custom binary protocol. It has no HTTP surface. The server binary (`cmd/server/server.go`) initializes a `ProxyServer`, calls `Process()` in a goroutine, and blocks on an error-reporting loop. Configuration is loaded from JSON files or CLI flags into `models/config.ServerConfig`.

Health checks are needed for container orchestration and load-balancer integration.

## Goals / Non-Goals

**Goals:**

- Provide a `GET /health` HTTP endpoint returning `{"status":"ok"}` with HTTP 200.
- Make the health listener address configurable via config file and CLI flag.
- Integrate the health server lifecycle into the existing `ProxyServer` without disrupting the TCP proxy flow.

**Non-Goals:**

- Readiness checks (checking TCP listener or client-controller state).
- Metrics, tracing, or observability endpoints.
- Health endpoints on the `auth_server` or `client` binaries.
- TLS or authentication on the health endpoint.

## Decisions

### 1. Stdlib `net/http` with a single handler

**Decision**: Use Go's `net/http` package directly — `http.ServeMux` with one route.

**Rationale**: The endpoint is trivial (one route, static JSON response). Adding a framework (Chi, Gin) would introduce a dependency for no benefit. The stdlib mux supports exact path matching and method checks.

**Alternatives considered**:
- *Chi/Gin router*: Rejected — unnecessary dependency for a single route.
- *Embedding in the TCP dispatcher*: Rejected — mixing HTTP into the custom protocol layer adds complexity and couples concerns.

### 2. Separate HTTP listener on a dedicated port

**Decision**: The health HTTP server listens on its own address, independent of the TCP proxy listeners.

**Rationale**: The TCP listeners use a custom binary protocol with optional Shadowsocks encryption. An HTTP handler cannot share those listeners. A separate port also allows firewall rules to restrict health-check traffic to internal networks.

**Default address**: `:8080` — a conventional health/admin port that avoids conflicts with common proxy ports.

### 3. Configuration via `health_address` field

**Decision**: Add `health_address` (JSON: `"health_address"`) to `ServerConfig` and a `--health-address` CLI flag. When empty, the health server does not start.

**Rationale**: Mirrors the existing `address` field pattern. Making it optional (empty = disabled) avoids breaking existing deployments that don't need health checks. The field lives in the common config section so it applies once per server instance, not per proxy entry.

### 4. Lifecycle integration

**Decision**: Start the HTTP health server in a goroutine from `ProxyServer`. Report startup errors through the existing `errs` channel. Use `http.Server` with `ListenAndServe` so the server can be cleanly shut down in the future if needed.

**Rationale**: Keeps the health server co-located with the proxy lifecycle. Errors surface through the same channel operators already monitor.

### 5. Strict method and path handling

**Decision**: Only `GET /health` returns 200. Other methods on `/health` return 405 Method Not Allowed. Other paths return 404 Not Found.

**Rationale**: Prevents misuse and clearly signals that this is a single-purpose endpoint. Returning 405 instead of silently ignoring non-GET requests follows HTTP semantics.

## Risks / Trade-offs

- **[Health port conflict]** → The default `:8080` may conflict with other services on the host. Mitigation: the address is configurable; documentation should note the default.
- **[No readiness signal]** → The endpoint reports liveness only — a running process that lost all TCP listeners still returns OK. Mitigation: this is explicitly a non-goal; a future readiness endpoint can be added if needed.
- **[Unencrypted HTTP]** → The health endpoint has no TLS. Mitigation: health checks are typically internal-only; TLS can be layered via sidecar proxy if required.
