# Tasks: REQ-clean-1776863811 /version 端点实现

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] 定义 /version 端点的 API 契约规范
- [ ] 验证端点返回 HTTP 200 状态码
- [ ] 验证响应包含 `sha` 字段
- [ ] 验证响应格式为有效 JSON
- [ ] 验证端点支持 GET 方法
- [ ] 验证端点不接受其他 HTTP 方法

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] 部署代码后验证 /version 端点可访问
- [ ] 验证返回的 git SHA 与当前部署代码一致
- [ ] 验证多次调用端点结果一致（幂等性）
- [ ] 验证端点响应时间满足性能要求
- [ ] 验证端点在服务启动后立即可用

## Stage: implementation (owner: dev-agent)
- [ ] 在 HTTP router 中注册 /version 端点处理函数
- [ ] 实现版本信息结构体和 JSON 序列化
- [ ] 在构建时通过 `-ldflags` 注入 git SHA
- [ ] 更新 Makefile 编译流程以支持版本注入
- [ ] 更新 CI/CD 流程（GitHub Actions）以传递 git SHA
- [ ] 为 /version 端点添加单元测试
- [ ] 验证编译和部署流程正确注入版本信息
