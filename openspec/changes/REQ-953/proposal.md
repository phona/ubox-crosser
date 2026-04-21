# REQ-953: 为 e2e 测试环境添加 HTTP health 检查端点 /api/healthz

## 一句话需求

在 proxy-server 的 admin HTTP mux 上添加 `/api/healthz` 端点，返回 JSON 健康状态，替代现有的 TCP-only healthcheck，让 docker-compose 和集成测试可以通过 HTTP 判断服务就绪。

## 影响范围

| 文件/目录 | 变更类型 | 说明 |
|---|---|---|
| `server/admin.go` | 修改 | 在 `NewAdminMux()` 中注册 `/api/healthz` 路由 |
| `cmd/server/server.go` | 修改 | 启动 admin HTTP server（如尚未启动） |
| `models/config/config.go` | 修改（可能） | 添加 admin listen address 配置字段 |
| `tests/docker-compose.yml` | 修改 | proxy-server healthcheck 改用 HTTP `/api/healthz` |
| `tests/Dockerfile.test` | 修改（可能） | 加入 curl 或保留 TCP healthcheck 作为 fallback |

## 支持性判定

**支持实现**。变更集中在 admin HTTP 层，不影响核心代理隧道逻辑。`server/admin.go` 已存在 HTTP mux 基础设施（`/webhook-debug`），只需追加路由。风险低。
