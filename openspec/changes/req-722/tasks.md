## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 列出要覆盖的 API 契约点（GET /whoami 路径、状态码、Content-Type、响应体格式）

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
- [ ] TODO: 实现 whoami 包（handler.go + handler_test.go）
- [ ] TODO: 注册路由到 admin mux（cmd/server/server.go）
