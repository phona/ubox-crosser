# REQ-997 — 控制面 API 服务（crosser-api）

## 一句话需求
新增独立 Go module `crosser-api`，提供 REST 控制面，用于管理代理服务（CRUD）、用户/密钥、proxy 实例注册与状态上报，替代当前手编 `server.json` / `auth_server.json` 的运维方式。

## 背景
当前 ubox-crosser 通过本地 JSON 配置文件启动 proxy-server / auth-server，运维需手动改文件并重启。无统一的服务管理、无在线状态可见性、无密钥轮换通道。

## 影响范围
- 新增目录 `crosser-api/`（独立 module，monorepo 内并列于现有代码）：
  - `cmd/api/main.go`、`internal/{handler,service,repository,middleware,model,database}/`
  - 嵌入 SQLite，初始迁移 `internal/database/migrations/001_init.sql`
- 现有 `server/`、`client/`、`cmd/server/`、`cmd/auth_server/`：本变更**不修改**它们的运行行为，仅通过 `GET /api/v1/services/:name/config` 输出与现有 JSON 完全兼容的配置（向后兼容）
- proxy 端调用 `/proxy/register` / `/proxy/heartbeat` 的改造拆为后续 REQ，不在本变更范围
- 新依赖：HTTP 框架（待选 chi / gin）、SQLite 驱动（modernc.org/sqlite，纯 Go 免 cgo）、JWT 库、bcrypt（golang.org/x/crypto/bcrypt）

## 支持性判定（是否做）
- ✅ 做：MVP（FR-03.01/02/07/08）—— 用户登录、服务 CRUD、配置导出、SQLite schema
- ✅ 做：FR-03.04/05/06 —— proxy 注册 / 心跳 / 在线状态查询的服务端接收方
- ⏸ 暂不做：proxy 端反向调用（独立 REQ）、密钥轮换灰度窗口、connection_stats 聚合 / 清理
- ⏸ 暂不做：管理 UI（仅 API）

## 关键决策（与原 REQ-03 分歧/补强）
1. **proxy ↔ api 鉴权**：单独的机器凭证（service token），不复用管理员 JWT
2. **service ↔ expose_address**：先做 1:1（与现有 `auth_server.json` 对齐），后续按需扩 1:N
3. **默认管理员**：首次启动随机生成首密并打到日志，避免硬编码 admin/admin
4. **online 判定**：心跳超时 90s 自动转 offline（NFR 待评审中确认）
5. **Repo 形态**：当前 repo 内 monorepo（`crosser-api/` 同级目录），共享类型与 proto

## 验收对齐
覆盖原 REQ-03 的 AC-03.01 (登录) / AC-03.02 (服务 CRUD) / AC-03.03 (proxy 注册) / AC-03.04 (配置兼容)。

## 风险
- SQLite 选 cgo 还是纯 Go：影响交叉编译与单二进制目标 → design.md 详述
- HTTP 框架选型：chi 体量轻 vs. gin 生态熟 → design.md 详述
- 默认管理员策略变更可能与 BKD 自动化测试期望（admin/admin 直登）冲突 → 提供 ENV `CROSSER_ADMIN_INITIAL_PASSWORD` 兜底
