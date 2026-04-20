## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 列出要覆盖的 API 契约点（状态码、Content-Type、响应体结构、非 GET 方法拒绝）

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] [REQ-736-A1] Normal uptime response: GET /uptime returns 200 with JSON {"uptime_seconds": <int>}
- [x] [REQ-736-A2] Uptime value increases over time (monotonic)
- [x] [REQ-736-A3] Uptime resets to near-zero after service restart
- [x] [REQ-736-A4] Non-GET methods rejected with 405
- [x] [REQ-736-A5] Endpoint accessible without authentication

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 列出要实现的模块（uptime 包、Init/Handler、路由注册、单元测试）
