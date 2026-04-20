---
change_id: REQ-736
title: "Add GET /uptime endpoint"
layers:
  - backend
status: draft
repos: [ubox-crosser]
---

## Why

运维和监控平台需要快速确认服务已运行多久，用于判断是否刚刚重启、计算可用率以及辅助故障排查。`GET /uptime` 返回服务启动至今的秒数，轻量且易于集成到健康检查和仪表盘中。

## What Changes

- 新增 `uptime` 包，提供 `Init()` 记录启动时间和 `Handler` 返回 JSON `{"uptime_seconds": N}`
- 在 admin `http.ServeMux` 中注册 `GET /uptime`
- `cmd/server/server.go` 启动时调用 `uptime.Init()` 记录启动时刻
- Go 1.22+ method-based routing 自动拒绝非 GET 请求（405）

## Capabilities

### New Capabilities
- `uptime-endpoint`: HTTP GET /uptime 返回 JSON `{"uptime_seconds": <int>}`，HTTP 200

### Modified Capabilities
（无）

## Impact

- `uptime/` — 新包（handler.go, handler_test.go）
- `cmd/server/server.go` — 调用 `uptime.Init()` + 注册 `GET /uptime` 到现有 admin mux
- 无新外部依赖（仅 stdlib `time`）
