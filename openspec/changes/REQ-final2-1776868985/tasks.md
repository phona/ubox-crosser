# REQ-final2-1776868985 Tasks

## Stage: contract-tests (owner: contract-spec-agent)

- [ ] 定义 `/version` 端点的 HTTP 接口规范：路由、方法、状态码、响应体格式
- [ ] 编写 contract test（Gherkin scenario 或类似格式）验证：
  - [ ] GET /version 返回 200
  - [ ] 响应体包含 "sha" 字段
  - [ ] "sha" 字段值为有效的 hex 字符串
- [ ] 确保 contract test 在 CI pipeline 中可执行
- [ ] 更新 TESTS.md 或相关文档说明 contract 的运行方式

## Stage: acceptance-tests (owner: acceptance-spec-agent)

- [ ] 编写 e2e acceptance scenario，验证：
  - [ ] /version 端点在正常构建的二进制中返回正确的 git SHA
  - [ ] SHA 值与构建时的 git HEAD 匹配
  - [ ] 不同构建（如 Docker 镜像构建）的 SHA 一致
- [ ] 集成测试：构建二进制、启动服务、调用端点、验证响应
- [ ] 测试在 `make ci-integration-test` 中涵盖或可独立运行
- [ ] 补充错误场景（如 SHA 为空或 unknown 时的行为）

## Stage: implementation (owner: dev-agent)

- [ ] 在 main.go 中声明 `var GitSHA = "unknown"` 
- [ ] 实现 /version 路由处理函数，返回 JSON：`{"sha": "<GitSHA>"}`
- [ ] 注册路由到 HTTP server（确定路由前缀和位置）
- [ ] 更新 Makefile，在 `go build` target 中注入 GIT_SHA：
  - [ ] 获取 `git rev-parse HEAD`
  - [ ] 传递 `-ldflags "-X main.GitSHA=..."`
- [ ] 测试本地构建和 Docker 镜像构建，验证注入成功
- [ ] 确保所有 test targets（ci-unit-test, ci-integration-test）通过
- [ ] 提交 commit 到分支 `stage/REQ-final2-1776868985-dev`
