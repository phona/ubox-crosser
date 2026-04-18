# Dev Spec: Add Echo API

## File Structure

```
server/api/
├── echo.go          # Echo HTTP handler
├── echo_test.go     # Unit tests for the echo handler
└── server.go        # HTTP server setup and route registration

models/config/
└── config.go        # Modified: add HTTPAddress field to ServerConfig

cmd/server/
└── server.go        # Modified: wire HTTP server startup alongside ProxyServer
```

## Dependencies

- **Go stdlib only**: `net/http`, `encoding/json`, `context`, `net/http/httptest` (tests)
- No new external dependencies are introduced

## Function Signatures

### `server/api/echo.go`

```go
package api

// EchoHandler returns an http.HandlerFunc that reads the "msg" query parameter
// and responds with JSON.
func EchoHandler() http.HandlerFunc
```

**Internal types** (unexported, file-scoped):

```go
type echoResponse struct {
    Message string `json:"message"`
}

type errorResponse struct {
    Error string `json:"error"`
}
```

### `server/api/server.go`

```go
package api

// NewHTTPServer creates an *http.Server with all routes registered.
// addr is the listen address (e.g. ":8080").
func NewHTTPServer(addr string) *http.Server

// ListenAndServe starts the HTTP server. It blocks until the server stops.
// Returns http.ErrServerClosed on graceful shutdown.
func ListenAndServe(srv *http.Server) error

// Shutdown gracefully stops the HTTP server using the provided context.
func Shutdown(ctx context.Context, srv *http.Server) error
```

### `models/config/config.go` (modification)

```go
type ServerConfig struct {
    LoginPass   string `json:"login_password"`
    AuthPass    string `json:"auth_password"`
    Address     string `json:"address"`
    HTTPAddress string `json:"http_address"` // NEW — default ":8080"
    Config
}
```

### `cmd/server/server.go` (modification)

Wire HTTP server startup inside the cobra `Run` function:

```go
// After proxy.Process() goroutine launch:
httpAddr := resolveHTTPAddr(configs) // extract http_address from config, default ":8080"
httpSrv := api.NewHTTPServer(httpAddr)
go func() {
    if err := api.ListenAndServe(httpSrv); err != nil && err != http.ErrServerClosed {
        logrus.Errorf("HTTP server error: %s", err)
    }
}()
```

## Handler Logic: `EchoHandler`

```
1. Check request method
   - If method != GET → respond 405 with {"error": "method not allowed"}
   - Return immediately

2. Read query parameter "msg"
   - If "msg" key is absent from query → respond 400 with {"error": "msg parameter is required"}
   - Note: ?msg= (present but empty) is valid — respond 200

3. Respond 200 with {"message": "<value>"}
```

Response headers for all paths: `Content-Type: application/json`.

## Route Registration

`NewHTTPServer` uses `http.NewServeMux`:

```go
mux := http.NewServeMux()
mux.HandleFunc("/api/echo", EchoHandler())
```

No middleware. No catch-all. Unregistered paths get the default 404 from `ServeMux`.

## Error Handling

| Condition | Status | Body |
|---|---|---|
| Non-GET method | 405 | `{"error":"method not allowed"}` |
| Missing `msg` param | 400 | `{"error":"msg parameter is required"}` |
| JSON marshal failure | 500 | `Internal Server Error` (plain text, stdlib default) |

- `json.Marshal` on these small structs will not fail in practice; the 500 path is a defensive stdlib fallback, not an explicit handler branch.
- `ListenAndServe` returning a non-`ErrServerClosed` error is logged by the caller in `cmd/server/server.go`.

## Edge Cases

1. **`?msg=` (empty string)**: Valid. Returns `{"message":""}` with 200. The check is for parameter *presence*, not *non-empty*.
2. **`?msg=hello&msg=world` (duplicate keys)**: `r.URL.Query().Get("msg")` returns the first value (`"hello"`). This is acceptable per the spec (single value expected).
3. **Request body on GET**: Ignored. Only the query string is read.
4. **Very long `msg` value**: No explicit length limit in the spec. The Go stdlib enforces `http.DefaultMaxHeaderBytes` (1 MB) on the request line + headers, which implicitly caps query string length. No additional validation needed.
5. **Unicode / special characters in `msg`**: Handled transparently — `Query().Get()` returns the URL-decoded value, `json.Marshal` encodes it as a valid JSON string (with proper escaping).
6. **Concurrent requests**: `EchoHandler` is stateless — safe for concurrent use with no synchronization needed.
7. **Port conflict on `:8080`**: `ListenAndServe` returns an error; logged by the caller. Operator resolves via `http_address` config.
8. **Graceful shutdown**: `Shutdown` is exposed but not wired in the initial implementation (the existing server uses an infinite error loop with no signal handling). This is intentional — shutdown plumbing is out of scope for this change.

## Config Resolution

The `http_address` field is read from the server config JSON. Resolution order:

1. If `http_address` is set in the config file → use it
2. Otherwise → default to `":8080"`

A new CLI flag `--http-address` is added to `cmd/server/server.go` for command-line override, following the existing flag pattern.

## Testing Strategy (structure only — no test code)

`server/api/echo_test.go` uses `httptest.NewRequest` + `httptest.NewRecorder`:

| Test Case | Method | URL | Expected Status | Expected Body |
|---|---|---|---|---|
| echo with message | GET | `/api/echo?msg=hello` | 200 | `{"message":"hello"}` |
| echo with empty msg | GET | `/api/echo?msg=` | 200 | `{"message":""}` |
| missing msg param | GET | `/api/echo` | 400 | `{"error":"msg parameter is required"}` |
| wrong method | POST | `/api/echo?msg=hello` | 405 | `{"error":"method not allowed"}` |
