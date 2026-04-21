---
change_id: REQ-924
title: "Add GET /echo endpoint"
repos: [ubox-crosser]
layers:
  - backend
status: draft
---

## Why

提供一个轻量级的 `/echo` 端点，接收 `msg` 查询参数并原样返回。用于网络连通性验证、调试请求链路，以及确认 admin HTTP 服务正常响应。比 `/ping` 更灵活——调用方可以验证特定字符串是否被完整传递。

## What Changes

- 新建 `echo/` 包，实现 `Handler`：读取 `?msg=xxx` 查询参数，以 `text/plain` 返回原文
- 在 `cmd/server/server.go` admin mux 注册 `GET /echo`
- `msg` 缺失时返回空 body（HTTP 200），不报错

### Design Decision: 独立包 vs. 内联

与 `/buildinfo`（复用 `version.Handler`）不同，`/echo` 有自己的逻辑（读取 query param），应遵循项目惯例（`ping/`、`health/`、`version/` 各自独立包）建立 `echo/` 包。

## Capabilities

### New Capabilities
- `echo-endpoint`: HTTP GET /echo?msg=xxx 以 text/plain 返回 msg 参数值，HTTP 200

### Modified Capabilities
- None

## Impact

- 新增 `echo/handler.go` — Handler 实现
- `cmd/server/server.go` — 添加一行路由注册
- 无新依赖
