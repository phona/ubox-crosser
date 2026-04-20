---
layers:
  - backend
---

## Why

容器编排平台（Kubernetes、Docker Swarm）和负载均衡器需要 HTTP 健康检查端点来判断 ubox-crosser 代理服务的存活状态。当前服务只暴露自定义二进制协议的 TCP 监听端口，运维方只能用 raw TCP 探测，无法区分"端口监听中"与"进程实际可用"。需要新增一个标准 `GET /health` HTTP 端点，与 TCP 代理监听器解耦，便于独立暴露给探活系统。

## What Changes

- 在代理服务进程内新增独立的 HTTP 监听器，专用于健康检查。
- 新增 `GET /health` 端点：返回 HTTP 200，body `{"status":"ok"}`，header `Content-Type: application/json`。
- `/health` 路径上的非 GET 请求返回 HTTP 405，并设置 `Allow: GET` header。
- 任何非 `/health` 路径返回 HTTP 404。
- `models/config.ServerConfig` 新增 `health_address` 字段（JSON key `"health_address"`），同步新增 `--health-address` CLI flag；空值表示不启动 health 监听器。
- health 监听器的启动失败（端口占用等）通过现有 `ProxyServer.errs` channel 上报。

## Capabilities

### New Capabilities

- `health-endpoint`：在代理服务进程上提供标准 `GET /health` HTTP 健康检查端点，定义请求/响应契约、状态码、JSON schema 与配置项。

### Modified Capabilities

（无已有 spec 需要修改）

## Impact

- **代码**：新增 `server/health.go`；修改 `models/config/config.go`、`server/server.go`、`cmd/server/server.go`。
- **API**：在独立端口上暴露新的 HTTP API（`GET /health`），与现有 TCP 代理端口隔离。
- **配置**：`ServerConfig` 新增可选字段 `health_address`；CLI 新增可选 flag `--health-address`。空值保持向后兼容，不启动 HTTP 监听。
- **依赖**：无新增依赖，仅使用 Go 标准库 `net/http`。
- **部署**：Docker / Kubernetes 健康探针可改用 HTTP probe 指向 `/health`。
