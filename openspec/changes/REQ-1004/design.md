---
change_id: REQ-1004
title: "crosser-api control plane — design"
---

## High-Level Approach

新建 `crosser-api/` 独立 Go module，采用经典分层架构（handler → service → repository），
使用 Chi router 做 HTTP 路由，SQLite 做嵌入式存储，JWT 做认证。

## Architecture

```
crosser-api/
├── go.mod
├── cmd/api/main.go                    # 入口：解析配置、初始化 DB、注册路由、启动 HTTP
├── internal/
│   ├── model/models.go                # 数据模型：User, Service, ProxyInstance, ConnectionStats
│   ├── handler/
│   │   ├── auth.go                    # POST /auth/login, /auth/refresh
│   │   ├── service.go                 # CRUD /services, /services/:name, /services/:name/config
│   │   └── proxy.go                   # POST /proxy/register, /proxy/heartbeat, GET /proxy/status
│   ├── service/
│   │   ├── auth.go                    # 登录验证、JWT 生成/刷新
│   │   ├── service.go                 # 服务 CRUD 业务逻辑、密钥生成
│   │   └── proxy.go                   # 注册、心跳、状态聚合
│   ├── repository/
│   │   ├── user.go                    # users 表 CRUD
│   │   ├── service.go                 # services 表 CRUD
│   │   └── proxy.go                   # proxy_instances + connection_stats 表操作
│   ├── middleware/
│   │   └── jwt.go                     # JWT 鉴权中间件
│   └── database/
│       ├── sqlite.go                  # DB 初始化 + migration runner
│       └── migrations/001_init.sql    # 建表 DDL
└── tests/
    └── api_test.go                    # 集成测试
```

## API Design

所有 API 以 `/api/v1` 为前缀，统一 JSON 响应格式：

```json
{"code": 0, "message": "success", "data": {...}}
```

错误时 code 非零，HTTP 状态码语义正确（400/401/404/500）。

### 认证流程
1. POST /api/v1/auth/login 带 username + password → 返回 JWT access_token（15min）+ refresh_token
2. 后续请求 Authorization: Bearer <token>
3. POST /api/v1/auth/refresh 用 refresh_token 换新 access_token

### 服务管理
- 创建服务时自动生成 cipher_key（32 字节 hex）
- 查询/列表时 cipher_key 脱敏（只返回前4+后4字符）
- 配置导出返回与现有 ServerConfig JSON 格式兼容的结构

### Proxy 注册/心跳
- 注册时记录 proxy 地址和服务列表，状态置为 online
- 心跳上报连接统计，更新 last_heartbeat
- 超过 3 倍心跳间隔未收到心跳的 proxy 标记为 offline

## Database Schema

4 张表：users, services, proxy_instances, connection_stats。
时间字段用 Unix timestamp（INTEGER），与现有项目风格一致。
详见 REQ-03 需求文档。

## Technology Choices

| 选型 | 选择 | 理由 |
|------|------|------|
| HTTP Router | go-chi/chi/v5 | 轻量、stdlib 兼容、middleware 生态好 |
| Database | mattn/go-sqlite3 | 嵌入式、零外部依赖、CGO |
| JWT | golang-jwt/jwt/v5 | Go 社区标准 JWT 库 |
| Password Hash | golang.org/x/crypto/bcrypt | 安全标准 |
| Config | 环境变量 + flags | 简单，符合 12-factor |

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| CGO 依赖（sqlite3）增加构建复杂度 | Dockerfile 使用 golang:alpine + gcc；CI 需 CGO_ENABLED=1 |
| JWT secret 硬编码 | 通过环境变量 JWT_SECRET 注入，启动时强制校验 |
| 默认 admin/admin 密码不安全 | 首次启动自动创建，日志提示修改；后续版本支持强制修改 |
| SQLite 并发写入限制 | WAL 模式 + busy_timeout；控制面并发极低，足够 |

## Alternatives Considered

- **Gin 框架**：功能全但偏重，Chi 更轻量且 stdlib 兼容
- **PostgreSQL/MySQL**：引入外部依赖，违反 NFR-03.01 单二进制要求
- **BoltDB/BadgerDB**：KV 存储，不适合关系查询（如 JOIN proxy_instances + connection_stats）
- **主模块内新包**：SQLite CGO 会污染现有 CGO_ENABLED=0 构建链
