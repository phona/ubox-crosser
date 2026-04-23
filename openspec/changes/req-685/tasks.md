## 1. Health Package

- [ ] 1.1 Create `health/handler.go` with JSON handler returning `{"status":"ok"}`
- [ ] 1.2 Create `health/handler_test.go` with unit tests for GET response, JSON schema, and method enforcement

## 2. Server Integration

- [ ] 2.1 Register `"GET /healthz"` on the admin `http.ServeMux` in `cmd/server/server.go`

## 3. Acceptance Verification

- [ ] 3.1 Verify `GET /healthz` returns 200 with `{"status":"ok"}`
- [ ] 3.2 Verify non-GET methods return 405
- [ ] 3.3 Verify response matches `HealthStatus` schema in contract
