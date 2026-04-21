# REQ-983: Tasks

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: 定义 stamp-readme contract spec 场景（追加时间戳、幂等性、README 缺失时的行为）
- [ ] TODO: 编写 contract 测试验证 spec 场景

## Stage: acceptance-tests (owner: acceptance-spec-agent) [SKIP]
- [ ] SKIP: 本 REQ 跳过 acceptance-tests 阶段，用于验证 skip_accept 简化链路

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 Makefile 新增 stamp-readme target（sed 删旧行 + echo 追加新时间戳，README 不存在则跳过）
- [ ] TODO: 在 build target 末尾追加 stamp-readme 调用
- [ ] TODO: 验证 contract 测试全部通过
