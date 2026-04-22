---
change_id: REQ-1004
title: "crosser-api control plane — tasks"
---

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: Define OpenAPI spec for auth endpoints (POST /api/v1/auth/login, /api/v1/auth/refresh)
- [ ] TODO: Define OpenAPI spec for service CRUD endpoints (GET/POST /api/v1/services, GET/PUT/DELETE /api/v1/services/:name)
- [ ] TODO: Define OpenAPI spec for service config export (GET /api/v1/services/:name/config)
- [ ] TODO: Define OpenAPI spec for proxy endpoints (POST register/heartbeat, GET status)
- [ ] TODO: Define unified response schema {"code", "message", "data"}
- [ ] TODO: Contract tests for auth login (200 with token, 401 for wrong password)
- [ ] TODO: Contract tests for service CRUD (create/read/update/delete/list)
- [ ] TODO: Contract tests for proxy register/heartbeat/status
- [ ] TODO: Contract tests for JWT middleware (401 without token, 401 with expired token)

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] TODO: AC-03.01 — 默认管理员 admin/admin 登录返回 JWT token；错误密码返回 401
- [ ] TODO: AC-03.02 — 已登录管理员创建服务返回含自动生成 key 的服务详情；列表包含新服务
- [ ] TODO: AC-03.03 — proxy 注册后状态查询返回该实例为 online
- [ ] TODO: AC-03.04 — 通过 API 创建服务后 config 导出可作为 proxy-server 配置使用
- [ ] TODO: 密码 bcrypt 存储验证（不可逆）
- [ ] TODO: 密钥 API 返回脱敏验证（非全量显示）

## Stage: implementation (owner: dev-agent)
- [ ] TODO: crosser-api/go.mod — 初始化 module，声明依赖 (chi, sqlite3, jwt, bcrypt)
- [ ] TODO: crosser-api/internal/database/ — SQLite 初始化 + migration runner + 001_init.sql
- [ ] TODO: crosser-api/internal/model/ — 数据模型定义 (User, Service, ProxyInstance, ConnectionStats)
- [ ] TODO: crosser-api/internal/repository/ — user/service/proxy 数据访问层
- [ ] TODO: crosser-api/internal/service/ — auth/service/proxy 业务逻辑层
- [ ] TODO: crosser-api/internal/middleware/ — JWT 鉴权中间件
- [ ] TODO: crosser-api/internal/handler/ — auth/service/proxy HTTP handler
- [ ] TODO: crosser-api/cmd/api/main.go — 入口（配置解析、DB 初始化、路由注册、HTTP 启动）
- [ ] TODO: Makefile — 添加 crosser-api 构建/测试/lint 目标
- [ ] TODO: Dockerfile — 添加 crosser-api 构建 stage（CGO_ENABLED=1）
