## Stage: contract-tests (owner: contract-test-agent)
- [x] [REQ-722-S1] GET /whoami returns 200 with Content-Type text/plain; charset=utf-8
- [x] [REQ-722-S2] GET /whoami body is a non-empty hostname string
- [x] [REQ-722-S3] os.Hostname failure returns fallback "unknown"
- [x] [REQ-722-S4] POST /whoami returns 405 Method Not Allowed
- [x] [REQ-722-S5] PUT /whoami returns 405 Method Not Allowed
- [x] [REQ-722-S6] DELETE /whoami returns 405 Method Not Allowed
- [x] contract.spec.yaml — OpenAPI 3.1 spec for /whoami endpoint
- [x] tests/contract/whoami_test.go — contract test suite

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] [ACCEPT-A1] GET /whoami returns 200 with text/plain hostname
- [x] [ACCEPT-A2] Response body is a valid, non-empty hostname string
- [x] [ACCEPT-A3] POST /whoami returns 405
- [x] [ACCEPT-A4] PUT /whoami returns 405
- [x] [ACCEPT-A5] DELETE /whoami returns 405
- [x] [ACCEPT-A6] os.Hostname failure returns fallback "unknown"
- [x] [ACCEPT-A7] GET /ping still returns 200 (regression check)
- [x] [ACCEPT-A8] GET /healthz still returns 200 (regression check)

## Stage: implementation (owner: dev-agent)
- [x] whoami/handler.go — Handler 函数：调用 os.Hostname()，失败回退 "unknown"，返回 text/plain 200
- [x] whoami/handler_test.go — 单元测试：状态码+Content-Type、body 匹配 hostname、非空 body、POST/PUT/DELETE 返回 405
- [x] cmd/server/server.go — 在 admin mux 注册 `GET /whoami` 路由
- [x] go vet / go build / unit test 全部通过
