# DEV SPEC — REQ-03 Health Endpoint

## 1. File Structure

### Modified Files

| File | Change |
|------|--------|
| `models/config/config.go` | Add `HealthAddr` field to `ServerConfig` |
| `cmd/server/server.go` | Add `--health-addr` Cobra flag binding |
| `server/server.go` | Start health HTTP server in `NewProxyServer` or alongside `Process()` |

### New Files

| File | Purpose |
|------|---------|
| `server/health.go` | Health HTTP handler and server startup logic |
| `server/health_test.go` | Unit tests for health handler |

## 2. Function Signatures & Responsibilities

### `models/config/config.go`

```go
type ServerConfig struct {
    // ... existing fields ...
    HealthAddr string `json:"health_addr"`
}
```

- `HealthAddr`: listen address for the health HTTP server. Default `:8080`.

### `server/health.go`

```go
func newHealthHandler() http.Handler
```
- Creates an `http.ServeMux` with a single `GET /health` route.
- `GET /health` → 200, `Content-Type: application/json`, body `{"status":"ok"}`.
- Non-GET methods on `/health` → 405 Method Not Allowed.
- Unknown paths → 404 (default `ServeMux` behavior).

```go
func startHealthServer(addr string, errs chan<- error)
```
- Calls `http.ListenAndServe(addr, handler)` in a goroutine.
- If `ListenAndServe` returns an error (e.g. port bind failure), sends the error to `errs` and logs via `logrus.Errorf`. Does NOT crash or propagate fatally.

### `cmd/server/server.go`

Add one flag:

```go
cmd.Flags().StringVar(&cmdConfig.HealthAddr, "health-addr", ":8080", "health check listen address")
```

### `server/server.go`

Call `startHealthServer` once during server initialization. The call site is inside the `Run` function in `cmd/server/server.go` (after `NewProxyServer`) or inside `NewProxyServer` — whichever keeps health lifecycle co-located with the proxy lifecycle. The recommended approach is to call it in `cmd/server/server.go` right after creating the proxy, using the first config's `HealthAddr` (or a dedicated top-level field). This avoids pushing HTTP concerns into the TCP proxy struct.

Recommended call site in `cmd/server/server.go`:

```go
proxy := server.NewProxyServer(configs)
go server.StartHealthServer(configs, proxy.Errs()) // or extract HealthAddr
go proxy.Process()
```

If `HealthAddr` is empty string, skip starting the health server (allows explicit opt-out).

## 3. Dependencies

- **No new third-party dependencies.** Uses Go stdlib `net/http` only.
- Existing: `github.com/sirupsen/logrus` (for logging errors).

## 4. Error Handling Strategy

| Scenario | Behavior |
|----------|----------|
| Health port bind failure (EADDRINUSE) | Log error via `logrus.Errorf`, send to `errs` channel. Proxy continues normally. |
| Health handler panic | Standard `net/http` recovery — logs stack trace, connection closed. No proxy impact. |
| Invalid HealthAddr format | `ListenAndServe` returns error immediately → logged, non-fatal. |

Key principle from design.md: **health-check failure degrades observability, not availability.** The health server must never crash or block the proxy server.

## 5. Boundary Conditions

| Condition | Expected Behavior |
|-----------|-------------------|
| `GET /health` | 200, `{"status":"ok"}`, `Content-Type: application/json` |
| `POST /health` | 405 Method Not Allowed |
| `PUT /health` | 405 Method Not Allowed |
| `DELETE /health` | 405 Method Not Allowed |
| `GET /unknown` | 404 Not Found |
| `GET /` | 404 Not Found |
| `--health-addr` not specified | Defaults to `:8080` |
| `--health-addr :9090` | Listens on `:9090` |
| `--health-addr ""` (empty) | Skip health server startup |
| Port already in use | Log error, proxy runs without health endpoint |
| Multiple server configs | Single health server instance (one `HealthAddr` for the whole process) |

## 6. Storage / DB Requirements

None. This feature is stateless — no persistence required.

## 7. Config File Support

When `health_addr` is present in the JSON config file, it should be picked up via the existing `conf.ParseServerConfigFile` flow. Since `HealthAddr` is on `ServerConfig`, it will be parsed automatically by `json.Unmarshal`. The `--health-addr` CLI flag takes precedence over config file values (following existing Cobra/config merge pattern).

## 8. Implementation Notes for Dev Agent

1. **Response body must be exactly `{"status":"ok"}`** — use `json.Marshal` or a hardcoded constant. Prefer hardcoded `[]byte(`{"status":"ok"}`)` since the response is static and avoids allocation.
2. **Set `Content-Type: application/json`** header explicitly before writing the response body.
3. **Method check**: `http.ServeMux` does not enforce method restrictions. The handler must check `r.Method == http.MethodGet` and return 405 for anything else.
4. **Exported vs unexported**: `StartHealthServer` should be exported so `cmd/server/server.go` can call it. `newHealthHandler` can remain unexported.
5. **No graceful shutdown** in this iteration (per design.md decision #4).
6. **Logging**: use `logrus` (aliased as `log` per project convention in `server/` package).
7. **Test approach**: unit test the handler via `httptest.NewServer` / `httptest.NewRecorder`. Test all method/path combinations from the boundary conditions table.
