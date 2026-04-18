# Dev Spec: GET /api/routes Management API (REQ-70)

## Overview

Add a read-only HTTP management endpoint `GET /api/routes` to the proxy server that returns the list of configured proxy routes with their connection status. The endpoint follows the contract defined in `contract.spec.yaml` (OpenAPI 3.0.3).

## File Structure

```
api/
  routes.go               # RouteInfo type, RouteProvider interface, HTTP handler
  routes_contract_test.go  # External tests validating OpenAPI contract compliance
  routes_unit_test.go      # Internal tests for handler logic
contract.spec.yaml         # OpenAPI 3.0.3 spec (source of truth)
server/
  server.go                # ProxyServer — add RouteProvider implementation
cmd/
  server/server.go         # Wire management HTTP server startup
```

No new packages or directories beyond `api/`.

## Types and Interfaces

### `api.RouteInfo`

```go
// api/routes.go
type RouteInfo struct {
    Name          string `json:"name"`
    ListenAddress string `json:"listen_address"`
    Method        string `json:"method"`
    Active        bool   `json:"active"`
}
```

Maps 1:1 to `components/schemas/Route` in `contract.spec.yaml`. Four required fields, no optional fields, no extra fields.

### `api.RouteProvider`

```go
// api/routes.go
type RouteProvider interface {
    ListRoutes() []RouteInfo
}
```

Single-method interface. Implemented by `ProxyServer` to decouple the HTTP layer from the server internals.

### `api.NewRoutesHandler`

```go
// api/routes.go
func NewRoutesHandler(provider RouteProvider) http.Handler
```

Returns an `http.Handler` for mounting at `/api/routes`.

## Handler Behavior

| Condition | Response |
|---|---|
| `GET /api/routes` | 200 + `application/json` + JSON array of `RouteInfo` |
| `GET` with nil/empty provider | 200 + `[]` (empty JSON array, never `null`) |
| `POST`, `PUT`, `DELETE`, etc. | 405 Method Not Allowed, empty body |

### Response rules

- Content-Type header: `application/json` (only set for 200 responses).
- Response body for 200: always a JSON array (`[]` when empty, never `null`).
- Response body for 405: empty (no body written).
- No sensitive fields (`key`, `login_password`, `auth_password`) in the response.
- Only the four contract-defined fields appear in each route object.

## ProxyServer Integration

### `ListRoutes()` on `server.ProxyServer`

```go
// server/server.go
func (p *ProxyServer) ListRoutes() []api.RouteInfo
```

Implementation strategy:

1. Iterate over `p.context` (`map[string]config.ServerConfig`).
2. For each entry, build a `RouteInfo`:
   - `Name`: the map key (service name, e.g. `"UBox_cytm"`)
   - `ListenAddress`: `config.Address`
   - `Method`: `config.Method` (empty string if unencrypted)
   - `Active`: `true` if `p.controllers[name]` exists and is non-nil, `false` otherwise
3. Return the slice.

### Concurrency consideration

- `p.context` is written once at construction and never mutated — safe to read without a lock.
- `p.controllers` is mutated by `handleLoginRequest` (adds entries). Reading it from the HTTP goroutine while the TCP goroutine writes requires synchronization. Options:
  - Add a `sync.RWMutex` to `ProxyServer` protecting `controllers`.
  - Or accept the race for the initial implementation since the map is only added-to (never deleted) and `active` is advisory. Document the decision.
- Recommended: add `sync.RWMutex` — the cost is negligible and it prevents data-race detector failures.

## Management HTTP Server Startup

### In `cmd/server/server.go`

After creating the `ProxyServer`, start an HTTP server:

```go
mux := http.NewServeMux()
mux.Handle("/api/routes", api.NewRoutesHandler(proxy))
go http.ListenAndServe("127.0.0.1:8080", mux)
```

### Configuration

The management address is currently hardcoded to `127.0.0.1:8080` per `contract.spec.yaml`. A future change (config-api) will make this configurable via `management_address` in the server config. For REQ-70, hardcoding is acceptable.

### Binding

- Localhost only (`127.0.0.1`) — never bind to `0.0.0.0`.
- This is a read-only endpoint with no authentication; localhost binding is the security boundary.

## Dependencies

### New dependencies: none

The implementation uses only Go stdlib (`net/http`, `encoding/json`).

### Test dependencies

- `github.com/stretchr/testify` — already in `go.mod` (indirect). Promote to direct dependency.
- `net/http/httptest` — stdlib.

## Edge Cases

| Case | Expected behavior |
|---|---|
| No services configured (empty config map) | Return `[]` |
| Provider is `nil` (handler created with nil) | Return `[]` |
| Provider's `ListRoutes()` returns `nil` | Return `[]` (handler normalizes nil to empty slice) |
| Service configured but no client connected | `active: false` |
| Service configured and client connected | `active: true` |
| Client disconnects after connecting | `active` should reflect current state (requires controller cleanup, out of REQ-70 scope) |
| Multiple services sharing same listen address | Each appears as a separate route entry |
| Method is empty string (no encryption) | `method: ""` in JSON response |
| Concurrent GET requests | Must be safe (RWMutex read lock) |
| Very large number of services | No pagination required per contract; single array response |

## What NOT to implement

- No authentication/authorization on the management endpoint.
- No CORS headers.
- No pagination, filtering, or query parameters.
- No PUT/POST/DELETE handlers beyond returning 405.
- No management address configuration (hardcoded for REQ-70).
- No graceful shutdown of the management HTTP server.
- No controller cleanup or lifecycle tracking beyond checking map presence.

## Task Breakdown

1. Create `api/routes.go` with `RouteInfo`, `RouteProvider`, `NewRoutesHandler`
2. Add `ListRoutes()` method to `server.ProxyServer`
3. Wire management HTTP server in `cmd/server/server.go`
4. Write contract tests (`api/routes_contract_test.go`)
5. Write unit tests (`api/routes_unit_test.go`)

## Contract Reference

Source of truth: `contract.spec.yaml` at repository root.

- Endpoint: `GET /api/routes`
- Base URL: `http://127.0.0.1:8080`
- Response schema: array of `Route` objects
- Route fields: `name` (string), `listen_address` (string), `method` (string), `active` (boolean)
- All four fields are required
- Only responses: 200 (success) and 405 (method not allowed)
