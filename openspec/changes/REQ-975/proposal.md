# REQ-975: 在 README 末尾追加构建时间戳行

## 需求描述
在项目构建时，自动在 README.md 末尾追加一行构建时间戳，格式如 `Built at: 2026-04-21T12:00:00Z`。

## 影响范围
- **Makefile** — 新增或修改 build target，在编译完成后追加时间戳到 README.md
- **README.md** — 末尾会动态追加一行（构建产物，不入版本控制）

## 支持性判定
- 改动范围小（仅 Makefile + README 约定），无需新增依赖
- 不影响现有 Go 源码、Dockerfile、CI pipeline
- **支持实施**
