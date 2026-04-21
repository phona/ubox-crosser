---
change_id: REQ-969
title: "Add GET /api/version endpoint"
repos: [ubox-crosser]
layers:
  - backend
status: draft
---

## Why

在运维和调试时，需要快速确认当前部署的服务版本。现有端点（/api/healthz、/webhook-debug）无法提供构建版本信息。

新增 `/api/version` 端点，返回当前部署的 git commit hash，方便运维人员和 CI/CD 流水线验证部署是否正确。

## What Changes

- 新建 `version/` 包，实现 `Handler`：
  - 返回 JSON `{"commit":"<git-hash>"}` 格式
  - commit hash 通过构建时 `-ldflags -X` 注入包级变量
  - 仅支持 GET 方法，其他方法返回 405
- 在 `server/admin.go` 的 `NewAdminMux()` 注册 `/api/version` 路由
- 在 `Makefile` / CI 构建命令中添加 `-ldflags` 注入 commit hash

### Design Decision: ldflags 注入 vs. 运行时 exec git

使用 `-ldflags -X` 在编译时注入 commit hash：
- 零运行时开销，无需容器内安装 git
- 符合 Go 社区标准做法
- 如果未注入则返回 `"unknown"`

## Capabilities

### New Capabilities
- `version-endpoint`: HTTP GET /api/version 返回 JSON `{"commit":"<40-char-hex>"}`, HTTP 200

### Modified Capabilities
- None

## Impact

- 新增 `version/handler.go` — Handler 实现 + 包级 Commit 变量
- `server/admin.go` — 添加一行路由注册
- `Makefile` — 构建命令添加 ldflags
- 无新依赖，仅使用 `net/http`、`encoding/json` 标准库
