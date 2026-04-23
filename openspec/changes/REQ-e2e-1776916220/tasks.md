# 任务列表：/healthz 端点实现

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: 定义 /healthz 端点的 HTTP 契约规范（请求/响应格式）
- [ ] TODO: 编写契约验证测试，覆盖正常响应、边界场景（uptime = 0 时刻、大时间值）

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] TODO: 编写集成测试验证 /healthz 端点返回准确的 uptime（相对误差 < 1%）
- [ ] TODO: 验证 /healthz 与现有 /version 端点共存不冲突
- [ ] TODO: 验证 /healthz 在服务启动直后的快速响应

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 ProxyServer 结构体添加 startTime 字段
- [ ] TODO: 在 NewProxyServer 初始化时记录启动时间
- [ ] TODO: 实现 isHealthzRequest 方法判断 HTTP 请求
- [ ] TODO: 在 handleHTTPRequest 中添加 /healthz 处理逻辑
- [ ] TODO: 运行单元测试和集成测试确保功能正确
- [ ] TODO: 验证 HTTP Content-Length 计算准确
