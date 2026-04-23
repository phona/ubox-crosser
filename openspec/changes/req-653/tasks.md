---
change_id: req-653
title: "Tasks: GET /version endpoint (v3)"
---

## 1. Version Package

- [ ] 1.1 Create `version/version.go`: declare `const Version = "0.1.0"` and `var (Commit = "unknown"; BuildTime = "unknown")`
- [ ] 1.2 Create `version/handler.go`: implement `Handler(w, r)` — GET returns 200 + JSON with version/commit/build_time; non-GET returns 405
- [ ] 1.3 Create `version/handler_test.go`: test GET 200 + JSON body, version field matches constant, POST/PUT/DELETE return 405

## 2. Server Integration

- [ ] 2.1 Add `--http-addr` flag (default `:8080`) to `cmd/server/server.go`
- [ ] 2.2 Start HTTP server goroutine with `http.ServeMux` registering `/version` handler
- [ ] 2.3 Log HTTP server listen address on startup

## 3. Build Pipeline

- [ ] 3.1 Update `Makefile`: add `GIT_COMMIT`, `BUILD_TIME`, `VERSION_PKG`, `LDFLAGS` variables; pass `-ldflags` in build target
- [ ] 3.2 Update `Dockerfile`: add `GIT_COMMIT` and `BUILD_TIME` build args, pass to `go build -ldflags`

## 4. Verification

- [ ] 4.1 `go build ./...` compiles without errors
- [ ] 4.2 `go test ./version/...` all tests pass
- [ ] 4.3 `go vet ./...` no warnings
- [ ] 4.4 `make build` produces binaries with injected commit and build_time
