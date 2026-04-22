# Tasks: REQ-clean-1776863811 /version 端点实现

## Stage: contract-tests (owner: contract-spec-agent)
- [x] 定义 /version 端点的 API 契约规范 (openspec/changes/REQ-clean-1776863811/specs/version-endpoint/contract.md)
- [x] 验证端点返回 HTTP 200 状态码 (C1 scenario)
- [x] 验证响应包含 `sha` 和 `commit` 字段 (C2-C3 scenarios)
- [x] 验证响应格式为有效 JSON (C1 scenario)
- [x] 验证端点支持 GET 方法 (C4 scenario)
- [x] 验证端点不接受其他 HTTP 方法 (POST/DELETE/PUT) (C4 scenario)
- [x] 验证 SHA 格式为 40 字符十六进制或 "unknown" (C3 scenario)
- [x] 验证查询参数被忽略 (C5 scenario)
- [x] 实现契约测试 (tests/contract/version_endpoint_test.go)

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [x] 定义 acceptance spec 中的验收场景（FEATURE-A1 到 A5）
- [x] 场景 A1: 部署后端点立即可访问
- [x] 场景 A2: 返回的 git SHA 与部署代码一致
- [x] 场景 A3: 验证多次调用端点结果一致（幂等性）
- [x] 场景 A4: 验证端点响应时间满足性能要求（<100ms）
- [x] 场景 A5: 验证错误处理和边界条件

## Stage: implementation (owner: dev-agent)
- [x] 在 HTTP router 中注册 /version 端点处理函数
- [x] 实现版本信息结构体和 JSON 序列化
- [x] 在构建时通过 `-ldflags` 注入 git SHA
- [x] 更新 Makefile 编译流程以支持版本注入
- [x] 更新 CI/CD 流程（GitHub Actions）以传递 git SHA
- [x] 为 /version 端点添加单元测试
- [x] 验证编译和部署流程正确注入版本信息
