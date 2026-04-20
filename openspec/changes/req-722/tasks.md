## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 列出要覆盖的 API 契约点（GET /whoami 路径、状态码、Content-Type、响应体格式）

## Stage: acceptance-tests (owner: accept-test-agent)
- [ ] TODO: 列出要验的用户行为（正常返回主机名、非 GET 方法拒绝、现有端点不受影响）

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 实现 whoami 包（handler.go + handler_test.go）
- [ ] TODO: 注册路由到 admin mux（cmd/server/server.go）
