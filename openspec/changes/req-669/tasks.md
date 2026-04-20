## 1. Version Package

- [ ] 1.1 Create `version/version.go` with `Version` constant and `Commit`, `BuildTime` vars
- [ ] 1.2 Create `version/handler.go` with JSON handler returning `Info` struct
- [ ] 1.3 Create `version/handler_test.go` with unit tests for GET response and JSON schema

## 2. Server Integration

- [ ] 2.1 Add `--admin-addr` flag to `cmd/server` (default `:8080`)
- [ ] 2.2 Wire `http.ServeMux` with `"GET /version"` pattern on admin listener goroutine

## 3. Build Pipeline

- [ ] 3.1 Update `Makefile` with `-ldflags -X` for `Commit` and `BuildTime` injection
- [ ] 3.2 Update `Dockerfile` with build args for commit and build_time

## 4. Acceptance Verification

- [ ] 4.1 Verify `GET /version` returns 200 with correct JSON schema
- [ ] 4.2 Verify non-GET methods return 405
- [ ] 4.3 Verify default values (`"unknown"`) when built without ldflags
