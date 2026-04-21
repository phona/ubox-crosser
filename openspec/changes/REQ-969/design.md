---
change_id: REQ-969
title: "version endpoint — design"
---

## High-Level Approach

采用 Go 标准的 build-time injection 模式，通过 `go build -ldflags "-X pkg.Var=value"` 在编译时将 git commit hash 写入包级变量。

## Architecture

```
version/
  handler.go      # Handler + var Commit string
```

Handler 逻辑：
1. 检查 HTTP 方法，非 GET 返回 405
2. 构造 `{"commit": "<hash>"}` JSON 响应
3. 返回 200 + Content-Type: application/json

## Build Integration

Makefile 添加：
```makefile
COMMIT := $(shell git rev-parse HEAD)
LDFLAGS := -X github.com/phona/ubox-crosser/version.Commit=$(COMMIT)
```

若构建时未传 ldflags，Commit 默认为空字符串，Handler 返回 `"unknown"`。

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| 构建时忘记传 ldflags | Commit 默认 "unknown"，不会 panic |
| 端点暴露敏感信息 | commit hash 是公开仓库信息，无安全风险 |

## Alternatives Considered

- **运行时 `exec.Command("git")`**：需要容器内有 git，增加镜像体积和运行时开销，已否决
- **嵌入 version 文件**：需要额外构建步骤生成文件，ldflags 更简洁
