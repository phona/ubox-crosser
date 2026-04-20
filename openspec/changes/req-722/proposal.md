---
change_id: req-722
title: "Add GET /whoami endpoint"
layers:
  - backend
status: draft
repos: [ubox-crosser]
---

## Why

运维和编排平台在多节点部署场景下，需要快速确认请求到达了哪台主机。`GET /whoami` 返回当前机器的主机名，便于调试路由、负载均衡以及日志关联。

## What Changes

- 新增 `whoami` 包，HTTP handler 调用 `os.Hostname()` 返回纯文本主机名
- 在 admin `http.ServeMux` 中注册 `GET /whoami`
- Go 1.22+ method-based routing 自动拒绝非 GET 请求（405）

## Capabilities

### New Capabilities
- `whoami-endpoint`: HTTP GET /whoami 返回当前主机名，纯文本，HTTP 200

### Modified Capabilities
（无）

## Impact

- `whoami/` — 新包（handler.go, handler_test.go）
- `cmd/server/server.go` — 注册 `GET /whoami` 到现有 admin mux
- 无新外部依赖（仅 stdlib `os.Hostname`）
