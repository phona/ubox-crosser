## Context

The ubox-crosser proxy server (`server/server.go`) manages TCP connections via a custom binary protocol. It has no HTTP surface. The server binary (`cmd/server/server.go`) initializes a `ProxyServer`, calls `Process()` in a goroutine, and blocks on an error-reporting loop. Configuration is loaded from JSON files or CLI flags into `models/config.ServerConfig`.

Health checks are needed for Kubernetes liveness probes and container orchestration.

## Goals / Non-Goals

**Goals:**

- Provide a `GET /healthz` HTTP endpoint returning `{"status":"ok"}` with HTTP 200.
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

**Rationale**: The endpoint is trivial (one route, static JSON response). Adding a framework (Chi, Gin) would introduce a dependency for no benefit.

**Alternatives considered**:
- *Chi/Gin router*: Rejected — unnecessary dependency for a single route.
- *Embedding in the TCP dispatcher*: Rejected — mixing HTTP into the custom protocol layer adds complexity and couples concerns.

### 2. Separate HTTP listener on a dedicated port

**Decision**: The health HTTP server listens on its own address, independent of the TCP proxy listeners.

**Rationale**: The TCP listeners use a custom binary protocol with optional Shadowsocks encryption. An HTTP handler cannot share those listeners. A separate port also allows firewall rules to restrict health-check traffic to internal networks.

**Default address**: `:8080` — a conventional health/admin port that avoids conflicts with common proxy ports.

### 3. Configuration via `health_address` field

**Decision**: Add `health_address` (JSON: `"health_address"`) to `ServerConfig` and a `--health-address` CLI flag. When empty, the health server does not start.

**Rationale**: Mirrors the existing `address` field pattern. Making it optional (empty = disabled) avoids breaking existing deployments that don't need health checks.

### 4. Lifecycle integration

**Decision**: Start the HTTP health server in a goroutine from `ProxyServer`. Report startup errors through the existing `errs` channel.

**Rationale**: Keeps the health server co-located with the proxy lifecycle. Errors surface through the same channel operators already monitor.

### 5. Path: `/healthz` (Kubernetes convention)

**Decision**: Use `/healthz` instead of `/health`.

**Rationale**: `/healthz` is the de facto standard for Kubernetes liveness probes. Using this convention means Kubernetes manifests can use the default path without extra configuration.

### 6. Strict method and path handling

**Decision**: Only `GET /healthz` returns 200. Other methods on `/healthz` return 405 with `Allow: GET` header. Other paths return 404.

**Rationale**: Follows HTTP semantics and prevents misuse.

## Risks / Trade-offs

- **[Health port conflict]** → The default `:8080` may conflict with other services. Mitigation: the address is configurable.
- **[No readiness signal]** → The endpoint reports liveness only. Mitigation: explicitly a non-goal; readiness can be added later.
- **[Unencrypted HTTP]** → No TLS. Mitigation: health checks are typically internal-only.
