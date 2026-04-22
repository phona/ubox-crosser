# REQ-997 Tasks

> 三段式：contract-tests → acceptance-tests → implementation。
> 本文件由 analyze-agent 立骨架；细节 TODO 由对应 stage 的 owner agent 自行展开（章节归属见 owner 标注）。

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: 产出 `crosser-api/contract.spec.yaml` —— 覆盖 auth / services CRUD / proxy register|heartbeat|status / services config 导出 共 9 个端点
- [ ] TODO: 产出对应契约测试套件骨架（路径、用例编号），与 contract.spec.yaml 一一映射
- [ ] TODO: 锁定统一响应壳 `{code,message,data}` 与错误码枚举

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] TODO: 产出 FEATURE 列表（覆盖 AC-03.01 / 02 / 03 / 04 共 4 条原始验收）
- [ ] TODO: 用 Gherkin 描述每条 FEATURE 的 Given/When/Then 场景
- [ ] TODO: 标注每条 FEATURE 与 contract.spec.yaml 端点的反向追溯

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 落地 `crosser-api` module 骨架（go.mod / cmd / internal/{handler,service,repository,middleware,model,database}）
- [ ] TODO: 实现 SQLite 初始化 + embed 迁移 `001_init.sql`
- [ ] TODO: 实现 auth（登录 / 刷新 / JWT 中间件 / bcrypt / 默认管理员首密策略）
- [ ] TODO: 实现 services CRUD + 配置导出（与现有 server.json 字节级兼容）
- [ ] TODO: 实现 proxy register / heartbeat / status（含 X-Proxy-Token 鉴权）
- [ ] TODO: 通过 contract-tests 与 acceptance-tests 全绿
