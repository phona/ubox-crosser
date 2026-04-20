## Context

README.md 当前只有项目简介和架构图，没有版本记录。用户需要通过 git log 才能了解项目变更历史。

## Goals / Non-Goals

**Goals:**
- 在 README.md 中添加 `## Version History` 章节
- 版本条目包含版本号、日期、变更摘要
- 按时间倒序排列

**Non-Goals:**
- 不自动从 git log 生成版本记录
- 不引入 CHANGELOG.md 或其他独立文件
- 不改变现有 README 结构

## Decisions

1. **使用 Markdown 表格格式** — 表格比列表更结构化，便于快速扫描版本号和日期。列定义：Version | Date | Changes。
2. **章节位于 README 末尾** — 版本历史是参考信息，不应抢占项目简介和架构图的位置。
3. **日期格式 YYYY-MM-DD** — ISO 8601 格式，无歧义。

## Risks / Trade-offs

- [手动维护] → 版本记录需要手动更新，可能遗漏。可在 CI 中添加检查作为后续改进。
