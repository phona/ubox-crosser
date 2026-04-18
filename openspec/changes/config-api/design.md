## Context

ubox-crosser is a Go TCP proxy tunnel with three binaries (client, server, auth_server). Configuration is loaded from JSON files at startup via `utils/conf.ParseServerConfigFile()` and stored as `map[string]config.ServerConfig` inside `ProxyServer`. There is currently no HTTP server or runtime management interface. All communication uses a custom newline-delimited JSON protocol over TCP.

The proxy server manages multiple named services, each with its own `ServerConfig`. Configuration fields include encryption settings (`key`, `method`), network addresses (`address`), authentication passwords (`login_password`, `auth_password`), and logging settings (`log_file`, `log_level`).

## Goals / Non-Goals

**Goals:**
- Provide a read-only HTTP endpoint to inspect the full running configuration
- Provide a write endpoint to update mutable configuration fields at runtime
- Apply configuration changes immediately without process restart (hot-reload)
- Ensure thread-safe access to configuration shared between the TCP proxy and HTTP API
- Keep the management server isolated from the proxy protocol (separate port)

**Non-Goals:**
- Adding/removing entire services at runtime (only updating fields of existing services)
- Authentication or authorization on the management API (rely on network-level access control)
- WebSocket or streaming config change notifications
- Persisting runtime config changes back to the JSON file on disk
- Management endpoints beyond config (metrics, health checks, etc.)

## Decisions

### 1. Embedded HTTP server using Go stdlib `net/http`

**Decision**: Use `net/http` from the standard library with a simple `http.ServeMux`.

**Alternatives considered**:
- Chi/Gin router: Unnecessary for two endpoints. Adds an external dependency.
- Unix socket: Less portable and harder to test with standard HTTP tools.

**Rationale**: The project currently has no HTTP dependencies. Two endpoints don't justify a framework. `net/http` is well-tested and sufficient.

### 2. Concurrency-safe config access via `sync.RWMutex`

**Decision**: Wrap the `map[string]config.ServerConfig` in a struct with `sync.RWMutex`. `GET` acquires a read lock; `PUT` acquires a write lock.

**Alternatives considered**:
- `atomic.Value`: Good for full replacement, but awkward for partial field updates across a map.
- Channel-based serialization: Over-engineered for simple read/write patterns.

**Rationale**: `RWMutex` allows concurrent reads (the common case) while serializing writes. The config map is small enough that lock contention is not a concern.

### 3. Mutable vs immutable field classification

**Decision**: Only `log_level` and `log_file` are mutable at runtime. All other fields (`key`, `method`, `address`, `login_password`, `auth_password`) are immutable because changing them would require re-creating listeners, ciphers, or re-authenticating clients.

**Rationale**: Changing a listener address or encryption key mid-session would break active connections. The safe set of mutable fields is deliberately small. It can be expanded later if hot-reload support is added for specific fields.

### 4. Management server address configuration

**Decision**: Add a `management_address` field to the server config's `common` section (default: empty, meaning management API is disabled). Example: `"management_address": "127.0.0.1:8080"`.

**Rationale**: Disabled by default avoids exposing an unsecured API. Binding to localhost by convention limits attack surface.

### 5. Sensitive field masking in GET response

**Decision**: `GET /api/config` SHALL mask sensitive fields (`key`, `login_password`, `auth_password`) in the response, replacing their values with `"***"`.

**Rationale**: The management API may be accessed by operators who don't need to see plaintext secrets. Defense-in-depth.

## Risks / Trade-offs

- **[Risk] Unsecured management API** → Mitigation: disabled by default; documentation recommends localhost binding; future change can add auth.
- **[Risk] Config drift between runtime and file** → Mitigation: explicitly a non-goal to persist changes. Operators must update the file manually for persistence across restarts.
- **[Risk] Log level change not taking effect** → Mitigation: the PUT handler must call `logrus.SetLevel()` directly after updating the config value.
- **[Trade-off] Small mutable field set** → Limits utility but prevents dangerous runtime changes. Can be expanded incrementally.
