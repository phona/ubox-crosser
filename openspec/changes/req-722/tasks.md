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
- [ ] TODO: 实现 whoami 包（handler.go + handler_test.go）
- [ ] TODO: 注册路由到 admin mux（cmd/server/server.go）
