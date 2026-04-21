---
change_id: REQ-945
title: "Add POST /webhook-debug endpoint"
repos: [ubox-crosser]
layers:
  - backend
status: draft
---

## Why

在对接第三方 webhook 回调时，经常需要一个"请求检查器"来查看实际收到的请求内容（method、headers、query params、body）。现有端点（/ping、/echo、/healthz）都只处理 GET 请求且返回固定格式，无法满足 webhook 调试需求。

新增 `/webhook-debug` 端点，接收任意 HTTP 方法的请求，将完整请求信息以 JSON 格式返回，方便开发者快速验证 webhook 回调的内容是否正确。

## What Changes

- 新建 `webhookdebug/` 包，实现 `Handler`：
  - 读取并返回请求的 method、URL path、query params、headers、body
  - 以 `application/json` 格式返回
  - 支持 GET/POST/PUT/PATCH/DELETE 等任意 HTTP 方法
- 在 `cmd/server/server.go` admin mux 注册 `/webhook-debug` 路由（不限定方法）
- body 读取后原样放入 JSON 的 `body` 字段（string 类型）

### Design Decision: 独立包 vs. 内联

遵循项目惯例（`ping/`、`health/`、`version/`、`echo/` 各自独立包），新建 `webhookdebug/` 包。

## Capabilities

### New Capabilities
- `webhook-debug-endpoint`: HTTP /webhook-debug 接受任意方法，返回 JSON 格式的完整请求信息（method、path、query、headers、body），HTTP 200

### Modified Capabilities
- None

## Impact

- 新增 `webhookdebug/handler.go` — Handler 实现
- `cmd/server/server.go` — 添加一行路由注册
- 无新依赖，仅使用 `net/http`、`encoding/json`、`io` 标准库
