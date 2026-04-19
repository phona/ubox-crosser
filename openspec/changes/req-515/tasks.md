## 1. API Package

- [ ] 1.1 Create `api/` package with version handler that returns `{"version":"v3"}` on `GET /version` and `405` on other methods
- [ ] 1.2 Add unit tests for the version handler

## 2. Server Integration

- [ ] 2.1 Add `HttpAddr` field to `ServerConfig` and wire up `--http-addr` CLI flag (default `:8080`)
- [ ] 2.2 Start HTTP server in a goroutine from `cmd/server/server.go`

## 3. Verification

- [ ] 3.1 Add integration test verifying `GET /version` returns `{"version":"v3"}` through Docker Compose
