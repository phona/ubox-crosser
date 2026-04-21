# REQ-983: Tasks

## Stage: contract-tests (owner: contract-spec-agent)
- [x] 定义 stamp-readme contract spec 场景（追加时间戳、幂等性、README 缺失时的行为）
  - specs/stamp-readme/spec.md: REQ-983-S1, REQ-983-S2, REQ-983-S3
- [x] 编写 contract 测试验证 spec 场景
  - tests/contract/stamp_readme_test.go

## Stage: acceptance-tests (owner: acceptance-spec-agent) [SKIP]
- [ ] SKIP: FEATURE-A1: make build 后 README.md 末尾出现 Built at 时间戳行 (`specs/stamp-readme/spec.md::FEATURE-A1`)
- [ ] SKIP: FEATURE-A2: 时间戳为合法的 ISO 8601 UTC 格式且与当前时间偏差不超过 60s (`specs/stamp-readme/spec.md::FEATURE-A2`)
- [ ] SKIP: FEATURE-A3: 多次构建不会重复追加时间戳行（幂等性）(`specs/stamp-readme/spec.md::FEATURE-A3`)
- [ ] SKIP: FEATURE-A4: README.md 不存在时构建仍然成功 (`specs/stamp-readme/spec.md::FEATURE-A4`)
- [ ] SKIP: FEATURE-A5: stamp-readme target 可独立调用 (`specs/stamp-readme/spec.md::FEATURE-A5`)

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 Makefile 新增 stamp-readme target（sed 删旧行 + echo 追加新时间戳，README 不存在则跳过）
- [ ] TODO: 在 build target 末尾追加 stamp-readme 调用
- [ ] TODO: 验证 contract 测试全部通过
