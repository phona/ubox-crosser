---
layers:
  - api
  - build
---

## Why

ubox-crosser 目前没有暴露任何 HTTP 接口，运维无法通过标准 HTTP 请求查询运行中实例的版本号。新增 `GET /version` 接口后，运维和监控系统可以方便地获取版本信息，用于部署验证和版本一致性检查。

## What Changes

- 在 proxy server 二进制中新增一个轻量级 HTTP 端点 `GET /version`，返回 JSON 格式的版本号。
- 在构建阶段通过 `-ldflags` 注入版本号（git tag 或手动指定）。
- 新增一个共享的 `version` 包，存放版本变量，供 HTTP handler 和将来的 CLI `--version` 复用。

## Capabilities

### New Capabilities
- `version-endpoint`: 提供 HTTP `GET /version` 接口，返回 JSON 格式的版本信息（版本号、commit、构建时间）。

### Modified Capabilities

（无已有 capability 需要修改）

## Impact

- **代码**: 新增 `version/` 包；`cmd/server/server.go` 增加 HTTP listener 启动逻辑。
- **构建**: `Makefile` 的 ldflags 增加 `-X` 变量注入版本信息；`Dockerfile` 同步更新。
- **配置**: `ServerConfig` 可能需要新增 `HTTPAddress` 字段指定 HTTP 监听地址（可选，默认端口如 `:8080`）。
- **依赖**: 仅使用 Go 标准库 `net/http`，无新外部依赖。
