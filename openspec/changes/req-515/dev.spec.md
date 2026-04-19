# DEV-SPEC: REQ-515 — GET /version endpoint

## File Structure

### New Files

| File | Purpose |
|------|---------|
| `api/version.go` | Version HTTP handler |
| `api/version_test.go` | Unit tests for version handler |
| `tests/integration/version_test.go` | Integration test for `/version` endpoint |

### Modified Files

| File | Change |
|------|--------|
| `models/config/config.go` | Add `HttpAddr` field to `ServerConfig` |
| `cmd/server/server.go` | Add `--http-addr` CLI flag; start HTTP server goroutine |
| `tests/docker-compose.yml` | Expose HTTP port on `proxy-server` service; add env var to test-runner |

## Function Signatures & Responsibilities

### `api/version.go`

```go
package api

// VersionHandler returns an http.HandlerFunc that responds to GET /version
// with {"version":"v3"} (200 OK, Content-Type: application/json).
// Non-GET methods receive 405 Method Not Allowed.
func VersionHandler() http.HandlerFunc
```

Implementation notes:
- Use `http.StatusMethodNotAllowed` for non-GET methods. Do **not** set an `Allow` header (contract does not require it).
- Hard-code version string `"v3"` as a constant within the package.
- Response body must be exactly `{"version":"v3"}` (no trailing newline from `json.Encoder` — use `json.Marshal` + `w.Write`).

### `models/config/config.go`

```go
type ServerConfig struct {
    // ... existing fields ...
    HttpAddr string `json:"http_addr"` // NEW — HTTP API listen address
    Config
}
```

### `cmd/server/server.go`

Add to Cobra flag registration:

```go
cmd.Flags().StringVar(&cmdConfig.HttpAddr, "http-addr", ":8080", "HTTP API listen address")
```

Start HTTP server after proxy initialization:

```go
// Inside Run func, after proxy creation:
mux := http.NewServeMux()
mux.HandleFunc("/version", api.VersionHandler())
go http.ListenAndServe(configs["common"].HttpAddr_OR_cmdConfig.HttpAddr, mux)
```

Logic for determining `HttpAddr`:
1. If running from CLI flags (no config file / empty configs): use `cmdConfig.HttpAddr`.
2. If running from config file: use the `common` section's `HttpAddr`, falling back to `":8080"` if empty.

The HTTP server runs in a fire-and-forget goroutine (design.md: "no graceful shutdown" is acceptable).

### `tests/integration/version_test.go`

```go
//go:build integration

package integration

// TestVersionEndpoint verifies GET /version returns {"version":"v3"} with status 200.
func TestVersionEndpoint(t *testing.T)

// TestVersionEndpointWrongMethod verifies POST /version returns 405.
func TestVersionEndpointWrongMethod(t *testing.T)
```

Uses `HTTP_API_ADDR` env var (e.g., `proxy-server:8080`).

## Dependencies

**No new third-party dependencies.** All functionality uses Go stdlib:
- `net/http` — HTTP server and handler
- `encoding/json` — JSON response marshaling

## Error Handling Strategy

| Scenario | Behavior |
|----------|----------|
| `GET /version` | 200 OK + `{"version":"v3"}` |
| Non-GET `/version` | 405 Method Not Allowed (empty body) |
| `json.Marshal` failure | Impossible (static struct), no handling needed |
| HTTP listen failure | `http.ListenAndServe` returns error; logged via `logrus.Fatal` to crash-fast — a bind failure is unrecoverable |
| Unknown path (e.g., `/foo`) | Default `http.ServeMux` 404 behavior (stdlib default) |

## Boundary Conditions

1. **Method filtering**: Only `GET` returns 200. `HEAD`, `POST`, `PUT`, `DELETE`, `PATCH`, `OPTIONS` all return 405.
2. **Content-Type**: Must be `application/json` on 200 response.
3. **Port conflict**: If `--http-addr` port is already in use, the server should fail fast with a log message. No retry logic.
4. **Config precedence**: CLI flag `--http-addr` takes effect when no config file is provided. Config file's `http_addr` field takes effect when config file is used. Default is `:8080` in both paths.
5. **Concurrent startup**: HTTP server starts in a goroutine alongside the TCP proxy. No ordering dependency between them.

## Storage / DB Requirements

None. Version is a hardcoded constant. No persistent state.

## Docker Compose Changes (Integration Tests)

In `tests/docker-compose.yml`, the `proxy-server` service needs:
- No additional port mapping needed (test-runner connects within the Docker network).
- The test-runner service needs a new env var:
  ```yaml
  - HTTP_API_ADDR=proxy-server:8080
  ```

The server config file `tests/configs/server.json` needs `"http_addr": ":8080"` added (or rely on default if the code defaults when the field is absent).

## Checklist for Dev Agent

- [ ] Create `api/version.go` with `VersionHandler()`
- [ ] Create `api/version_test.go` — test 200 response body/headers, test 405 for POST
- [ ] Add `HttpAddr` field to `ServerConfig` in `models/config/config.go`
- [ ] Add `--http-addr` flag and HTTP server startup in `cmd/server/server.go`
- [ ] Add integration test `tests/integration/version_test.go`
- [ ] Update `tests/docker-compose.yml` with `HTTP_API_ADDR` env var for test-runner
- [ ] Verify `go vet ./...` and `go build ./...` pass
