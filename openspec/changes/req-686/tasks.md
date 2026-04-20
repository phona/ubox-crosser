## 1. Ping Package

- [ ] 1.1 Create `ping/handler.go` with plain-text handler returning `pong`
- [ ] 1.2 Create `ping/handler_test.go` with unit tests for GET response, body content, and method enforcement

## 2. Server Integration

- [ ] 2.1 Register `"GET /ping"` on the admin `http.ServeMux` in `cmd/server/server.go`

## 3. Acceptance Verification

- [ ] 3.1 Verify `GET /ping` returns 200 with body `pong`
- [ ] 3.2 Verify non-GET methods return 405
- [ ] 3.3 Verify Content-Type is `text/plain; charset=utf-8`
