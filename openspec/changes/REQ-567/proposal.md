# REQ-567: GET /healthz JSON 健康检查接口

## 背景

ubox-crosser 当前仅使用 TCP 端口探测做健康检查（docker-compose 中 `healthcheck` 工具），无法区分"端口在监听"和"服务真正就绪"。需要一个标准 HTTP 健康检查端点，返回结构化 JSON，供 Docker / K8s / 监控系统消费。

## 目标

在 proxy-server 上新增一个 HTTP 服务（独立端口），提供 `GET /healthz` 接口，返回 JSON 格式的健康状态信息。

## 方案

- 在 `ProxyServer` 启动时，额外启动一个 `net/http` 服务监听健康检查端口（默认 `:8080`，可通过配置 `health_address` 覆盖）。
- `GET /healthz` 返回 HTTP 200 + JSON body，包含：
  - `status`: `"ok"` 表示健康
  - `uptime_seconds`: 服务运行秒数
  - `services`: 当前已注册的服务名列表
- 仅支持 GET 方法，其他方法返回 405。
- 非 `/healthz` 路径返回 404。

## 非目标

- 不做深度健康检查（如检查下游连通性）。
- 不做认证/鉴权。
- 不替换现有 TCP 健康检查，而是作为补充。
