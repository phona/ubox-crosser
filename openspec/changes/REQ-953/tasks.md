# REQ-953: Tasks

## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 定义 /api/healthz 端点的 OpenAPI spec（GET 200 响应 schema）
- [ ] TODO: 编写 /api/healthz 契约测试，验证响应状态码和 JSON body

## Stage: acceptance-tests (owner: accept-test-agent)
- [ ] TODO: 在 docker-compose 集成测试中验证 proxy-server 的 /api/healthz 端点可达且返回 200
- [ ] TODO: 验证 docker-compose healthcheck 使用 HTTP /api/healthz 后服务正确进入 healthy 状态

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 server/admin.go 的 NewAdminMux() 中注册 GET /api/healthz handler
- [ ] TODO: 在 models/config/config.go 的 ServerConfig 中添加 AdminAddress 字段
- [ ] TODO: 在 cmd/server/server.go 中启动 admin HTTP server goroutine
- [ ] TODO: 更新 tests/docker-compose.yml 的 proxy-server healthcheck 为 HTTP 方式
- [ ] TODO: 更新 tests/Dockerfile.test 支持 HTTP healthcheck（curl 或改造 healthcheck 工具）
- [ ] TODO: 更新 tests/configs/server.json 添加 admin_address 配置
