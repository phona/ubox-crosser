---
change_id: REQ-991
title: "README 构建时间戳 — design"
---

## High-Level Approach

在 README.md 末尾追加一行 `Built: <timestamp>`，作为 e2e-lowdisk 流水线验证的最小变更载体。

## Why Minimal

此 REQ 的目的是验证流水线 skip 机制（SKIP_CI_INT + SKIP_ACCEPT），不是交付功能。变更必须足够小以避免磁盘/时间开销，同时足够"真实"以触发完整的 CI pipeline。

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| 时间戳导致每次构建 diff 不同 | 预期行为，用于验证流水线确实执行了变更 |
| README 格式被破坏 | 仅追加一行，不修改已有内容 |
