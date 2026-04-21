# REQ-953: Tasks

## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 定义 /api/healthz 端点的 OpenAPI spec（GET 200 响应 schema）
- [ ] TODO: 编写 /api/healthz 契约测试，验证响应状态码和 JSON body

## Stage: acceptance-tests (owner: accept-test-agent)
- [ ] FEATURE-A1: GET /api/healthz 返回 200 + JSON `{"status":"ok"}`（`specs/healthz-endpoint/spec.md::FEATURE-A1`）
- [ ] FEATURE-A2: docker-compose healthcheck 使用 HTTP /api/healthz 后 proxy-server 进入 healthy 状态（`specs/healthz-endpoint/spec.md::FEATURE-A2`）
- [ ] FEATURE-A3: docker-compose 网络内 test-runner 可达 proxy-server:8080/api/healthz（`specs/healthz-endpoint/spec.md::FEATURE-A3`）
- [ ] FEATURE-A4: 非 GET 方法对 /api/healthz 返回合理响应（`specs/healthz-endpoint/spec.md::FEATURE-A4`）
- [ ] FEATURE-A5: /api/healthz 在正常条件下 500ms 内响应（`specs/healthz-endpoint/spec.md::FEATURE-A5`）

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 server/admin.go 的 NewAdminMux() 中注册 GET /api/healthz handler
- [ ] TODO: 在 models/config/config.go 的 ServerConfig 中添加 AdminAddress 字段
- [ ] TODO: 在 cmd/server/server.go 中启动 admin HTTP server goroutine
- [ ] TODO: 更新 tests/docker-compose.yml 的 proxy-server healthcheck 为 HTTP 方式
- [ ] TODO: 更新 tests/Dockerfile.test 支持 HTTP healthcheck（curl 或改造 healthcheck 工具）
- [ ] TODO: 更新 tests/configs/server.json 添加 admin_address 配置
