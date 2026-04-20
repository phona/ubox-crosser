## 1. Version Package

- [ ] 1.1 Create `internal/version/version.go` with `Version`, `Commit`, `BuildTime` vars (defaults: `"0.1.0"`, `"unknown"`, `"unknown"`)
- [ ] 1.2 Create `internal/version/handler.go` with `Handler()` returning an `http.HandlerFunc` that writes JSON `{"version","commit","build_time"}`

## 2. HTTP Server Integration

- [ ] 2.1 Add `--http-addr` flag (default `:8080`) to `cmd/server/server.go`
- [ ] 2.2 Register `GET /version` on `http.ServeMux` and start HTTP listener in a goroutine before the proxy loop

## 3. Build Pipeline

- [ ] 3.1 Update `Makefile` build target to inject `Version`, `Commit`, `BuildTime` via `-ldflags -X`
- [ ] 3.2 Verify `make build && ./bin/server --help` shows no regression

## 4. Contract Spec

- [ ] 4.1 Create `openspec/changes/req-630/contract.spec.yaml` — OpenAPI 3.0 schema for `GET /version`

## 5. Unit Tests

- [ ] 5.1 Create `internal/version/handler_test.go` — test HTTP 200, Content-Type, and JSON body fields (FEATURE-S1, S2, S5, S7)
- [ ] 5.2 Add test for POST → 405 (FEATURE-S6)

## Stage: Contract Test (owner: contract-test-agent)

- [x] [FEATURE-S1] TestVersion_S1_Returns200WithJSON — verify HTTP 200, Content-Type application/json, exactly 3 required string fields (version, commit, build_time), no extra fields
- [x] [FEATURE-S5] TestVersion_S5_NoAuthRequired — verify unauthenticated GET /version returns 200
- [x] [FEATURE-S6] TestVersion_S6_PostReturns405 — verify POST /version returns 405 Method Not Allowed
- [x] [FEATURE-S7] TestVersion_S7_DefaultFieldsPresent — verify all fields present and non-empty in default build
