---
layers:
  - api
  - build
---

## Why

当前 ubox-crosser 没有暴露任何 HTTP 接口，运维和监控无法通过标准方式查询服务运行版本。需要一个轻量的 `GET /version` HTTP 端点，使运维工具和健康检查系统能够快速获取当前部署的版本号。

## What Changes

- 新增 `version` 包，通过 `-ldflags` 在构建时注入版本号（git tag / commit）
- 在 proxy server 二进制中新增一个内嵌的 HTTP listener，提供 `GET /version` 端点
- 更新 Makefile 的 `-ldflags` 以注入版本信息
- server 配置中新增可选的 `http_address` 字段用于指定 HTTP 监听地址

## Capabilities

### New Capabilities

- `version-endpoint`: 提供 GET /version HTTP 端点，返回构建时注入的版本号

### Modified Capabilities

（无）

## Impact

- **代码**: 新增 `internal/version` 包；`server` 包增加 HTTP listener 启动逻辑；`models/config` 增加 HTTP 地址配置字段
- **构建**: Makefile ldflags 变更，CI build 命令同步更新
- **API**: 新增 HTTP 端点，不影响现有 TCP 协议
- **依赖**: 仅使用 Go 标准库 `net/http`，无新外部依赖
