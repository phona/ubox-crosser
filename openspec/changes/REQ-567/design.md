# REQ-567 Design: GET /healthz

## 关键决策

1. **独立端口**：健康检查 HTTP 服务运行在独立端口（默认 `:8080`），与主 TCP 代理端口分离，避免协议冲突。
2. **配置字段**：在 `ServerConfig` 新增 `health_address` 字段，为空时默认 `:8080`。
3. **最小依赖**：仅使用 `net/http` 标准库，不引入第三方 HTTP 框架。
4. **启动时间戳**：`ProxyServer` 记录启动时间，用于计算 `uptime_seconds`。

## 响应格式

```json
{
  "status": "ok",
  "uptime_seconds": 123,
  "services": ["test_service", "protocol_test"]
}
```

## 文件变更

| 文件 | 变更 |
|------|------|
| `models/config/config.go` | `ServerConfig` 新增 `HealthAddress string` |
| `server/server.go` | `ProxyServer` 新增 `startedAt` 字段，新增 `startHealthServer()` 方法 |
| `server/health.go` | `/healthz` handler 实现 |
| `tests/configs/server.json` | 新增 `health_address` 配置 |
| `tests/docker-compose.yml` | 暴露健康检查端口 |
