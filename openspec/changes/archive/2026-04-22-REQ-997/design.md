# REQ-997 设计 — crosser-api 控制面

## 1. 高层架构
```
+----------+        REST/JSON          +-------------+         SQLite
|  Admin   |  <----------------------> |  crosser-api |  <-->  data.db
+----------+                            +-------------+
                                              ^
                                              | register / heartbeat (machine token)
                                              |
                                       +-------------+
                                       | proxy-server |  (本变更不改造其行为)
                                       +-------------+
```

## 2. 选型与取舍

### 2.1 HTTP 框架：候选 chi（首选）
| 候选 | 优点 | 缺点 |
|---|---|---|
| **chi** ✅ | 标准 `net/http` 兼容、轻量、中间件链清晰 | 生态略小 |
| gin | 生态成熟、文档多 | 自己一套 ctx，与标准库脱节 |

**决策**：chi。NFR-03.01 强调轻量、单二进制，chi 更贴合。

### 2.2 SQLite 驱动：纯 Go（首选）
| 候选 | 优点 | 缺点 |
|---|---|---|
| **modernc.org/sqlite** ✅ | 纯 Go，免 cgo，可静态编译 | 体积稍大，性能略低 |
| mattn/go-sqlite3 | 性能高 | 需 cgo，破坏单二进制目标 |

**决策**：modernc.org/sqlite。

### 2.3 JWT 库
`github.com/golang-jwt/jwt/v5`（社区维护活跃，API 稳定）。

### 2.4 配置 & 启动
`koanf`（与 pma-go 习惯一致）+ flag/env 覆盖；首次启动迁移内嵌（`embed` SQL）。

## 3. API 契约要点（详见 contract.spec.yaml，由 contract-spec-agent 产出）
- 统一响应：`{"code": 0, "message": "success", "data": ...}`
- 错误：HTTP 状态码 + 同结构 `code != 0`
- 鉴权：`Authorization: Bearer <jwt>`（管理员）；`X-Proxy-Token: <token>`（proxy 机器凭证）

## 4. 数据模型（与原 REQ-03 一致，做 1 处修订）
按原 schema 落地。修订点：
- `services` 表新增 `expose_address TEXT NOT NULL DEFAULT ''` 字段，吸收 `auth_server.json` 的同名字段
- `connection_stats` 表保留原结构；增长治理留给后续 REQ

## 5. 配置导出兼容性
`GET /api/v1/services/:name/config` 输出**与现有 `server.json` 完全字节级兼容**的片段：
```json
{
  "common": { "key":"...", "method":"...", "address":"...",
               "login_password":"...", "auth_password":"...",
               "log_file":"", "log_level":"info" },
  "<service_name>": { "key":"..." }
}
```
acceptance test 必须用现有 proxy-server 直接吃这份输出做 smoke。

## 6. 安全
- 密码：`bcrypt`（cost=12）
- 密钥：DB 存明文（管理面需要导出原值给 proxy 用）；API 返回时按 `?reveal=true` 显式开关脱敏 / 明文，默认脱敏（前 4 + `****` + 后 4）
- 默认管理员首密随机；ENV `CROSSER_ADMIN_INITIAL_PASSWORD` 可指定
- proxy token：每个 proxy 实例注册时签发一次性长 token，存 `proxy_instances.token_hash`

## 7. 风险与缓解
| 风险 | 影响 | 缓解 |
|---|---|---|
| modernc.org/sqlite 性能 | 心跳写入压力 | 心跳写聚合表 + 异步 flush |
| 现有 proxy 不会主动调 register | FR-03.04 验证不闭环 | 本 REQ 只做 server 端接收；调用方改造单开 REQ；契约测试用 mock proxy |
| 默认管理员策略与 BKD 自动化冲突 | CI 登录失败 | 提供 ENV 兜底 + README 明确 |

## 8. 不做（Out of Scope）
- 管理 UI、SSO、审计日志、密钥轮换灰度、stats 聚合 / 清理、proxy 端配套改造
