# Dev Spec: /api/connections - Connection Query API

## Overview

Add an HTTP management API to the ubox-crosser proxy server that exposes connection state via REST endpoints. This enables operators to inspect active connections without accessing logs.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/connections` | List connections with filtering & pagination |
| GET | `/api/connections/{id}` | Get a single connection by ID |

Full contract: `contract.spec.yaml` in repo root.

---

## File Structure

```
api/
  handler.go         # HTTP handlers for /api/connections
  handler_test.go    # (DO NOT CREATE - tests are locked)
  models.go          # Request/response structs
  server.go          # HTTP server setup and router
server/
  server.go          # MODIFY: expose connection registry
  controller.go      # MODIFY: track connection metadata
models/
  connection/
    connection.go    # Connection data model with ID, status, timestamps
```

---

## 1. Connection Model (`models/connection/connection.go`)

```go
package connection

import "time"

type Status string
type Type string

const (
    StatusActive     Status = "active"
    StatusIdle       Status = "idle"
    StatusTerminated Status = "terminated"
)

const (
    TypeControl Type = "control"
    TypeWorker  Type = "worker"
)

type Connection struct {
    ID            string     `json:"id"`
    ServeName     string     `json:"serve_name"`
    RemoteAddr    string     `json:"remote_addr"`
    Status        Status     `json:"status"`
    Type          Type       `json:"type"`
    ConnectedAt   time.Time  `json:"connected_at"`
    LastHeartbeat *time.Time `json:"last_heartbeat,omitempty"`
}
```

**Responsibilities:**
- Pure data struct, no business logic
- `ID` is generated using `fmt.Sprintf("conn-%s", uuid)` or a simple incrementing counter with prefix
- `LastHeartbeat` is a pointer (nullable) - only set for control connections

---

## 2. Connection Registry (`server/registry.go` - NEW FILE)

```go
package server

import "sync"

type ConnectionRegistry struct {
    mu    sync.RWMutex
    conns map[string]*connection.Connection
}

func NewConnectionRegistry() *ConnectionRegistry
func (r *ConnectionRegistry) Add(conn *connection.Connection)
func (r *ConnectionRegistry) Remove(id string)
func (r *ConnectionRegistry) Get(id string) (*connection.Connection, bool)
func (r *ConnectionRegistry) List(filter FilterOptions) ([]*connection.Connection, int)
func (r *ConnectionRegistry) UpdateHeartbeat(id string, t time.Time)
func (r *ConnectionRegistry) UpdateStatus(id string, status connection.Status)
```

**FilterOptions struct:**
```go
type FilterOptions struct {
    ServeName string
    Status    connection.Status
    Type      connection.Type
    Page      int
    PageSize  int
}
```

**Responsibilities:**
- Thread-safe (use `sync.RWMutex`) - multiple goroutines read/write connections
- `List()` applies filters, then paginates. Return (slice, totalCount)
- Pagination: skip `(page-1)*pageSize` items, take `pageSize` items
- If no filters, return all

---

## 3. Modify `server/server.go`

**Changes to `ProxyServer`:**
```go
type ProxyServer struct {
    // ... existing fields ...
    registry *ConnectionRegistry  // ADD THIS
}
```

**In `NewProxyServer()`:**
- Initialize `registry: NewConnectionRegistry()`

**In `handleLoginRequest()`:**
- After successful login, create a `connection.Connection` with:
  - `ID`: generate unique ID (e.g., `fmt.Sprintf("conn-%d", atomic.AddUint64(&idCounter, 1))`)
  - `ServeName`: from the login message
  - `RemoteAddr`: `coordinator.Conn.RemoteAddr().String()`
  - `Status`: `StatusActive`
  - `Type`: `TypeControl`
  - `ConnectedAt`: `time.Now()`
- Call `registry.Add(conn)`
- Store the connection ID on the controller so it can be referenced later

**In `handleConnection()` for `GEN_WORKER`:**
- Create a worker connection entry with `Type: TypeWorker`
- Call `registry.Add(conn)`

**In controller `daemonize()`:**
- On `HEART_BEAT`: call `registry.UpdateHeartbeat(id, time.Now())`
- On loop exit (connection closed): call `registry.UpdateStatus(id, StatusTerminated)` or `registry.Remove(id)`

**Expose registry getter:**
```go
func (p *ProxyServer) Registry() *ConnectionRegistry {
    return p.registry
}
```

---

## 4. API HTTP Handlers (`api/handler.go`)

```go
package api

import (
    "encoding/json"
    "net/http"
    "strconv"
)

type Handler struct {
    registry RegistryReader  // interface for testability
}

type RegistryReader interface {
    Get(id string) (*connection.Connection, bool)
    List(filter FilterOptions) ([]*connection.Connection, int)
}

func NewHandler(registry RegistryReader) *Handler

func (h *Handler) ListConnections(w http.ResponseWriter, r *http.Request)
func (h *Handler) GetConnection(w http.ResponseWriter, r *http.Request)
```

### `ListConnections` logic:
1. Parse query params: `serve_name`, `status`, `type`, `page`, `page_size`
2. Validate:
   - `status` must be one of: `active`, `idle`, `terminated` (if provided)
   - `type` must be one of: `control`, `worker` (if provided)
   - `page` must be >= 1 (default 1)
   - `page_size` must be 1-100 (default 20)
   - On invalid params: return 400 with `ErrorResponse{Error: "invalid_parameter", Message: "..."}`
3. Call `registry.List(filter)` to get filtered, paginated results
4. Build response:
   ```json
   {
     "connections": [...],
     "pagination": {
       "page": 1,
       "page_size": 20,
       "total": 42,
       "total_pages": 3
     }
   }
   ```
5. `total_pages` = `ceil(total / page_size)`
6. Set `Content-Type: application/json` header
7. Write 200 response

### `GetConnection` logic:
1. Extract `{id}` from URL path (parse manually or use a lightweight router)
2. Call `registry.Get(id)`
3. If not found: return 404 with `ErrorResponse{Error: "not_found", Message: "Connection not found"}`
4. If found: return 200 with `ConnectionResponse{Connection: conn}`
5. Set `Content-Type: application/json` header

---

## 5. API HTTP Server (`api/server.go`)

```go
package api

import (
    "net/http"
)

func NewServeMux(handler *Handler) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/connections", handler.ListConnections)
    mux.HandleFunc("/api/connections/", handler.GetConnection)
    return mux
}

func StartManagementServer(addr string, registry RegistryReader) error {
    handler := NewHandler(registry)
    mux := NewServeMux(handler)
    return http.ListenAndServe(addr, mux)
}
```

**Path routing note:** Using Go stdlib `http.ServeMux`:
- `/api/connections` (exact) -> `ListConnections`
- `/api/connections/` (prefix) -> `GetConnection` - extract ID by trimming prefix `/api/connections/`

---

## 6. Response Structs (`api/models.go`)

```go
package api

import "ubox-crosser/models/connection"

type ConnectionListResponse struct {
    Connections []*connection.Connection `json:"connections"`
    Pagination  Pagination              `json:"pagination"`
}

type ConnectionResponse struct {
    Connection *connection.Connection `json:"connection"`
}

type Pagination struct {
    Page       int `json:"page"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Message string `json:"message"`
}
```

---

## Dependencies

No new external dependencies required. Use only Go stdlib:
- `net/http` for HTTP server
- `encoding/json` for JSON serialization
- `sync` for thread-safe registry
- `strconv` for query param parsing
- `strings` for URL path parsing
- `math` for `Ceil` in pagination
- `time` for timestamps

---

## Edge Cases to Handle

1. **Empty connection list**: Return `{"connections": [], "pagination": {"page": 1, "page_size": 20, "total": 0, "total_pages": 0}}`
2. **Page beyond range**: Return empty connections array with correct total (not 404)
3. **Invalid status filter**: Return 400 `{"error": "invalid_parameter", "message": "Invalid status value. Must be one of: active, idle, terminated"}`
4. **Invalid type filter**: Return 400 `{"error": "invalid_parameter", "message": "Invalid type value. Must be one of: control, worker"}`
5. **Invalid page number** (0, negative, non-integer): Return 400
6. **Invalid page_size** (0, negative, >100, non-integer): Return 400
7. **Connection ID not found**: Return 404 `{"error": "not_found", "message": "Connection not found"}`
8. **Concurrent access**: Registry must be safe for concurrent reads/writes
9. **Empty ID in path** (`/api/connections/`): Return 400 `{"error": "invalid_parameter", "message": "Connection ID is required"}`
10. **Non-GET methods**: Return 405 Method Not Allowed
11. **Connections array must never be null**: Always return `[]` not `null` when empty
12. **Content-Type header**: Always set `Content-Type: application/json` on all responses

## Integration Point

Start the management HTTP server alongside the proxy server. In `cmd/server/main.go` (or equivalent entry point), after creating the `ProxyServer`, start the management API:

```go
go api.StartManagementServer(":8080", proxyServer.Registry())
```

The management server port should be configurable but defaults to `:8080`.
