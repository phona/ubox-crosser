# REQ-975: Tasks

## Stage: contract-tests (owner: contract-spec-agent)
- [x] Contract spec scenarios REQ-975-S1..S3 in specs/stamp-readme/spec.md
- [x] Contract test — make stamp-readme appends valid ISO 8601 UTC timestamp to README.md last line (REQ-975-S1)
- [x] Contract test — running stamp-readme twice yields exactly one Built at line, idempotent (REQ-975-S2)
- [x] Contract test — stamp-readme exits 0 gracefully when README.md absent (REQ-975-S3)

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] FEATURE-A1: make build 后 README.md 末尾出现 Built at 时间戳行 (`specs/stamp-readme/spec.md::FEATURE-A1`)
- [ ] FEATURE-A2: 时间戳为合法的 ISO 8601 UTC 格式且与当前时间偏差不超过 60s (`specs/stamp-readme/spec.md::FEATURE-A2`)
- [ ] FEATURE-A3: 多次构建不会重复追加时间戳行（幂等性）(`specs/stamp-readme/spec.md::FEATURE-A3`)
- [ ] FEATURE-A4: README.md 不存在时构建仍然成功 (`specs/stamp-readme/spec.md::FEATURE-A4`)
- [ ] FEATURE-A5: stamp-readme target 可独立调用 (`specs/stamp-readme/spec.md::FEATURE-A5`)

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 在 Makefile 中新增 stamp-readme target
- [ ] TODO: 修改 build target 依赖 stamp-readme
- [ ] TODO: 验证幂等性（多次构建不重复追加）
