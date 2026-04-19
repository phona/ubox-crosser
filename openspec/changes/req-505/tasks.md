## 1. Version Package

- [ ] 1.1 Create `version/version.go` with a `Version` variable (default `"dev"`) for build-time injection via ldflags

## 2. HTTP Management Server

- [ ] 2.1 Implement `GET /version` handler that returns `{"version":"<value>"}` with `Content-Type: application/json`
- [ ] 2.2 Create management HTTP server startup function that accepts a listen address and registers the version route

## 3. Configuration & Wiring

- [ ] 3.1 Add `ManagementAddress` field to `ServerConfig` and `--management-address` CLI flag in `cmd/server/server.go`
- [ ] 3.2 Wire management server startup in `cmd/server/server.go` — start in a goroutine when address is configured, skip when empty

## 4. Build

- [ ] 4.1 Update Makefile `build` target to accept `VERSION` variable and pass `-X github.com/phona/ubox-crosser/version.Version=$(VERSION)` via ldflags

## 5. Contract

- [ ] 5.1 Create `openspec/changes/req-505/contract.spec.yaml` defining the `GET /version` HTTP contract

## 6. Tests

- [ ] 6.1 Add unit test for the version handler (correct status code, content-type, body)
- [ ] 6.2 Add integration test verifying `GET /version` against a running server container
