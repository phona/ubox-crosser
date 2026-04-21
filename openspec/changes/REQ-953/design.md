# REQ-953: Design — /api/healthz 端点

## 高层方案

在已有的 `server.NewAdminMux()` HTTP mux 上增加 `GET /api/healthz` 路由，返回：

```json
{"status": "ok"}
```

HTTP 200 表示健康，非 200 表示异常。

### 选型决策

| 决策点 | 选择 | 理由 |
|---|---|---|
| 路由位置 | `server/admin.go` 的 `NewAdminMux()` | 复用已有 admin HTTP 基础设施，无需新建 listener |
| 响应格式 | JSON `{"status":"ok"}` | 标准做法，便于 docker healthcheck 和监控系统消费 |
| Admin server 启动 | `cmd/server/server.go` 中 goroutine 启动 `http.ListenAndServe` | 与 proxy TCP listener 并行，互不干扰 |
| Admin 监听地址 | 配置字段 `AdminAddress`（默认 `:8080`） | 避免与代理端口冲突，可配置 |
| Docker healthcheck | `curl -sf http://localhost:8080/api/healthz` | 替代 TCP-only check，语义更清晰 |

### 风险

- **端口冲突**：admin 默认 `:8080` 可能与其他服务冲突 → 通过配置可覆盖
- **依赖 curl**：Dockerfile.test 基于 `golang:1.23`，自带网络工具；或可继续用自编译 healthcheck 工具发 HTTP 请求
