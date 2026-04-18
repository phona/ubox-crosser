## 1. Configuration Model

- [ ] 1.1 Add `ManagementAddress` field to `ServerConfig` in `models/config/config.go`
- [ ] 1.2 Create `ConfigStore` struct with `sync.RWMutex` wrapping `map[string]config.ServerConfig` in a new `api/` package
- [ ] 1.3 Define mutable field whitelist (`log_level`, `log_file`) and field validation logic in `ConfigStore`

## 2. API Handlers

- [ ] 2.1 Implement `GET /api/config` handler — read-lock config, mask sensitive fields (`key`, `login_password`, `auth_password`), return JSON
- [ ] 2.2 Implement `PUT /api/config` handler — parse JSON body, validate all fields are mutable and services exist, apply updates atomically
- [ ] 2.3 Implement method-not-allowed fallback for unsupported HTTP methods on `/api/config`
- [ ] 2.4 Apply hot-reload side effects: call `logrus.SetLevel()` when `log_level` is updated

## 3. Management Server Lifecycle

- [ ] 3.1 Create management HTTP server startup function that reads `management_address` from common config
- [ ] 3.2 Wire management server into `cmd/server/server.go` — start alongside proxy when configured, skip when unconfigured
- [ ] 3.3 Register `/api/config` route on the management `http.ServeMux`

## 4. Testing

- [ ] 4.1 Write unit tests for `ConfigStore` (concurrent read/write, mutable/immutable field enforcement)
- [ ] 4.2 Write unit tests for GET handler (response format, sensitive field masking, method not allowed)
- [ ] 4.3 Write unit tests for PUT handler (valid update, immutable field rejection, invalid JSON, unknown service, atomicity)
- [ ] 4.4 Write integration test for management server lifecycle (start when configured, skip when unconfigured)

## 5. Configuration & Documentation

- [ ] 5.1 Add `management_address` to example config files (`cmd/server/server.json`, `tests/configs/server.json`)
- [ ] 5.2 Update README with management API usage examples
