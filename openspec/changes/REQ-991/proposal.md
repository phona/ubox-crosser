---
change_id: REQ-991
title: "README 加构建时间戳行 (e2e-lowdisk)"
repos: [ubox-crosser]
layers:
  - docs
status: draft
skip_ci_integration: true
skip_acceptance: true
---

## Why

验证 CI 流水线在 low-disk 场景下的端到端运行能力。通过一个极简变更（README 加一行构建时间戳）触发完整流水线，同时使用 SKIP_CI_INT + SKIP_ACCEPT 跳过集成测试和验收测试，确保 skip 机制正常工作。

## What Changes

- `README.md` — 追加一行构建时间戳，格式：`Built: <ISO-8601 timestamp>`
- 不涉及任何代码、测试、配置变更

## Capabilities

### New Capabilities
- None（纯文档变更）

### Modified Capabilities
- None

## Impact

- `README.md` — 追加 1 行
- 无新依赖、无代码变更、无测试变更
- contract-tests：跳过（SKIP_CI_INT）
- acceptance-tests：跳过（SKIP_ACCEPT）
