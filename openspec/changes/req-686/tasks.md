## 1. Ping Package

- [ ] 1.1 Create `ping/handler.go` — exported `Handler(http.ResponseWriter, *http.Request)` that sets `Content-Type: text/plain; charset=utf-8` and writes `pong` via `io.WriteString`
- [ ] 1.2 Create `ping/handler_test.go` — unit tests mirroring `version/handler_test.go`: status 200, Content-Type header, body equals `pong`, method enforcement (405 for POST/PUT/DELETE via mux)

## 2. Server Integration

- [ ] 2.1 Add `import "github.com/phona/ubox-crosser/ping"` in `cmd/server/server.go`
- [ ] 2.2 Register `mux.HandleFunc("GET /ping", ping.Handler)` on admin mux alongside existing `GET /version` route

## 3. Acceptance Verification

- [ ] 3.1 `go test ./ping/...` passes — all handler tests green
- [ ] 3.2 `go vet ./...` and `make lint` pass with no new warnings
- [ ] 3.3 Verify `GET /ping` returns 200 with body `pong` and `Content-Type: text/plain; charset=utf-8`
- [ ] 3.4 Verify POST/PUT/DELETE to `/ping` return 405
- [ ] 3.5 Verify existing `GET /version` endpoint is unaffected
