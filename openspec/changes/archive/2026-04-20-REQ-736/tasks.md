## Stage: contract-tests (owner: contract-test-agent)
- [x] [REQ-736-S1] GET /uptime returns 200, Content-Type: application/json, body `{"uptime_seconds": <int>}`
- [x] [REQ-736-S2] uptime_seconds reflects elapsed time (>= 1 after 1s)
- [x] [REQ-736-S3] POST/PUT/DELETE /uptime return 405 Method Not Allowed
- [x] [REQ-736-S4] Response body has exactly one key `uptime_seconds`, no extra fields
- [x] OpenAPI contract spec: `contract.spec.yaml`
- [x] Contract test suite: `tests/contract/uptime_test.go`

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] [REQ-736-A1] Normal uptime response: GET /uptime returns 200 with JSON {"uptime_seconds": <int>}
- [x] [REQ-736-A2] Uptime value increases over time (monotonic)
- [x] [REQ-736-A3] Uptime resets to near-zero after service restart
- [x] [REQ-736-A4] Non-GET methods rejected with 405
- [x] [REQ-736-A5] Endpoint accessible without authentication

## Stage: implementation (owner: dev-agent)
- [ ] [IMP-1] 创建 `uptime/handler.go`：package-level `startTime`，`Init()` 记录启动时间，`Handler` 返回 JSON `{"uptime_seconds": N}`，非 GET 返回 405
- [ ] [IMP-2] 创建 `uptime/handler_test.go`：单元测试覆盖 S1–S4 场景（200+JSON、elapsed time、405、no extra fields）
- [ ] [IMP-3] `models/config/config.go` ServerConfig 增加 `AdminAddress` 字段（json: `admin_address`）
- [ ] [IMP-4] `cmd/server/server.go`：增加 `--admin-address` flag，启动时调用 `uptime.Init()`，创建 admin `http.ServeMux` 注册 `GET /uptime`，启动 HTTP listener
- [ ] [IMP-5] `go vet` / `go build` / `make unit-test` 全部通过
