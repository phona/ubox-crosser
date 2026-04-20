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
- [ ] TODO: 列出要实现的模块（uptime 包、Init/Handler、路由注册、单元测试）
