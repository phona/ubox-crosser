# REQ-975: Tasks

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: 定义 stamp-readme 的契约测试（验证 README 末尾追加了 Built at 时间戳行）

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] TODO: 定义验收场景（构建后 README 末尾包含正确格式的时间戳行）

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 Makefile 中新增 stamp-readme target
- [ ] TODO: 修改 build target 依赖 stamp-readme
- [ ] TODO: 验证幂等性（多次构建不重复追加）
