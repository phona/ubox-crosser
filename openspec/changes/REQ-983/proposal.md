# REQ-983: README 末尾追加构建时间戳行 (skip_accept 验整链路)

## 需求描述
在项目构建时，自动在 README.md 末尾追加一行构建时间戳，格式如 `Built at: 2026-04-21T12:00:00Z`。
本 REQ 与 REQ-975 功能相同，但本次跳过 acceptance-tests 阶段（skip_accept），用于验证 contract → implementation 的简化链路是否畅通。

## 影响范围
- **Makefile** — 新增 stamp-readme target，build 完成后自动调用追加时间戳
- **README.md** — 末尾动态追加一行（构建产物，幂等）

## 支持性判定
- 改动范围极小（Makefile 新增 ~5 行），无新依赖
- 不影响 Go 源码、Dockerfile、CI pipeline
- acceptance-tests 阶段标记 skip，仅走 contract-tests + implementation
- **支持实施**
