# REQ-975: 设计方案

## 方案概述
在 Makefile 中新增 `stamp-readme` target，构建完成后在 README.md 末尾追加一行
`Built at: <ISO-8601 UTC 时间戳>`。`build` target 依赖 `stamp-readme` 以自动触发。

## 选型
- **追加方式**: `echo` + shell `date` 命令，零依赖
- **时间戳格式**: ISO 8601 UTC（`date -u +%Y-%m-%dT%H:%M:%SZ`），通用且可排序
- **幂等性**: 每次构建先移除已有的 `Built at:` 行，再追加新行，避免重复累积

## 风险
- README.md 被 build 修改后不应被意外提交；可通过 `.gitignore` 或团队约定规避
- 如果 README.md 不存在，`stamp-readme` 应优雅失败或跳过
