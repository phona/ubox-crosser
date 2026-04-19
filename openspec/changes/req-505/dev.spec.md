# DEV SPEC — REQ-505: GET /version Endpoint

## File Changes

### New Files

| File | Purpose |
|------|---------|
| `version/version.go` | Package-level `Version` variable for ldflags injection |
| `server/management.go` | Management HTTP server + version handler |

### Modified Files

| File | Change |
|------|--------|
| `models/config/config.go` | Add `ManagementAddress` field to `ServerConfig` |
| `cmd/server/server.go` | Add `--management-address` CLI flag; start management server |
| `Makefile` | Add `VERSION` variable; pass `-X` ldflags to build targets |

---

## Detailed Design

### 1. `version/version.go` (new)

```go
package version

var Version = "dev"
```

- Single exported variable, zero dependencies.
- Build-time injection target: `-X github.com/phona/ubox-crosser/version.Version=$(VERSION)`

### 2. `server/management.go` (new)

#### Function: `NewManagementServer`

```go
func NewManagementServer(addr string) *http.Server
```

- Creates an `http.ServeMux`, registers `GET /version` handler.
- Returns `*http.Server{Addr: addr, Handler: mux}`.

#### Function: `handleVersion`

```go
func handleVersion(w http.ResponseWriter, r *http.Request)
```

- **Method check**: If `r.Method != http.MethodGet`, respond `405 Method Not Allowed`.
- Set `Content-Type: application/json`.
- Write `200` status.
- Encode `{"version":"<version.Version>"}` via `json.NewEncoder(w).Encode(...)`.
- Response struct: anonymous `struct{ Version string \`json:"version"\` }`.

#### Dependencies

- `net/http` (stdlib)
- `encoding/json` (stdlib)
- `github.com/phona/ubox-crosser/version`

### 3. `models/config/config.go` (modify)

Add one field to `ServerConfig`:

```go
type ServerConfig struct {
    LoginPass         string `json:"login_password"`
    AuthPass          string `json:"auth_password"`
    Address           string `json:"address"`
    ManagementAddress string `json:"management_address"`
    Config
}
```

- JSON key: `management_address`
- Zero value (`""`) means management server is disabled.
- The existing `Config.Update()` reflection-based merge handles `string` fields, so `ManagementAddress` inherits from `"common"` automatically if present.

### 4. `cmd/server/server.go` (modify)

#### New CLI flag

Add after existing flags:

```go
cmd.Flags().StringVar(&cmdConfig.ManagementAddress, "management-address", "", "HTTP management API listen address (e.g. 127.0.0.1:8080)")
```

#### Start management server

After `proxy.Process()` goroutine launch, before the error-reading loop:

```go
if mgmtAddr := getManagementAddress(configs); mgmtAddr != "" {
    mgmtServer := server.NewManagementServer(mgmtAddr)
    go func() {
        if err := mgmtServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Errorf("management server error: %v", err)
        }
    }()
}
```

Helper to resolve management address — use the first non-empty `ManagementAddress` from the configs map (CLI mode has a single entry; file mode may have it in `"common"` or any service key):

```go
func getManagementAddress(configs map[string]config.ServerConfig) string {
    for _, cfg := range configs {
        if cfg.ManagementAddress != "" {
            return cfg.ManagementAddress
        }
    }
    return ""
}
```

### 5. `Makefile` (modify)

```makefile
VERSION ?= dev

build: $(SOURCES)
	@echo "=== Building binaries ==="
	CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/phona/ubox-crosser/version.Version=$(VERSION)" -o bin/client ./cmd/client
	CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/phona/ubox-crosser/version.Version=$(VERSION)" -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/phona/ubox-crosser/version.Version=$(VERSION)" -o bin/auth_server ./cmd/auth_server
```

- `VERSION` defaults to `dev`; overridden via `make build VERSION=v2`.
- All three binaries get the same version string.

---

## Error Handling

| Scenario | Behavior |
|----------|----------|
| `ManagementAddress` empty | No HTTP server started; no error |
| `ManagementAddress` port in use | `ListenAndServe` returns error; logged via logrus; proxy continues |
| Non-GET request to `/version` | 405 Method Not Allowed |
| Management server crash | Logged; does not affect TCP proxy goroutines |

---

## Boundary Conditions

| Condition | Expected |
|-----------|----------|
| `Version` not injected at build time | Returns `{"version":"dev"}` |
| Multiple configs with different `ManagementAddress` | First non-empty wins |
| Concurrent requests to `/version` | Safe — handler is stateless, reads a package-level string |
| Very long version string | No artificial limit; JSON-encoded as-is |

---

## Dependencies

- **No new external dependencies.** Only Go stdlib (`net/http`, `encoding/json`).
- Internal dependency: `server/management.go` imports `github.com/phona/ubox-crosser/version`.

---

## Storage / DB

None. Version is a compile-time constant.

---

## Config File Example

```json
{
  "common": {
    "key": "my-secret",
    "method": "aes-256-cfb",
    "management_address": "0.0.0.0:8080"
  },
  "service1": {
    "address": "0.0.0.0:7001",
    "login_password": "pass1",
    "auth_password": "auth1"
  }
}
```

---

## Test Expectations (for test authors)

### Unit Tests (`server/management_test.go`)

1. `GET /version` → 200, `Content-Type: application/json`, body `{"version":"dev"}\n`
2. `POST /version` → 405
3. Set `version.Version = "v2"` before request → body contains `"v2"`

### Integration Tests

1. Build server with `-X ...version.Version=v2`, start with `--management-address`, curl `GET /version`, assert `{"version":"v2"}`.
2. Start server without `--management-address`, confirm no HTTP listener on default port.
