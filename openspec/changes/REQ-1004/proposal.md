---
change_id: REQ-1004
title: "Add crosser-api control plane service"
repos: [ubox-crosser]
layers:
  - backend
status: draft
---

## Why

当前 ubox-crosser 的代理服务配置完全依赖手动编辑 JSON 文件，缺少集中管理界面。
运维人员无法实时查看 proxy 实例在线状态和 client 连接拓扑。
需要一个轻量控制面 API 服务，提供 REST 接口管理代理服务、用户认证、密钥管理，
并接收 proxy 的注册和状态上报。

## What Changes

- 新建 `crosser-api/` Go module（独立 go.mod），作为管理控制面
- 使用 SQLite 作为嵌入式存储（零外部依赖）
- 分层架构：handler → service → repository
- 主要功能模块：
  1. **用户认证**：JWT 登录 + Token 鉴权中间件，默认管理员 admin/admin
  2. **服务管理 CRUD**：代理服务增删改查（name, address, cipher method, passwords）
  3. **密钥管理**：每个服务的加密密钥自动生成、轮换、查看（脱敏）
  4. **Proxy 注册/心跳**：proxy-server 启动时注册，定期心跳上报连接统计
  5. **在线状态查询**：返回所有 proxy 实例 + 关联 client 连接数
  6. **配置导出**：导出与现有 JSON 配置格式兼容的配置文件

### Design Decision: 独立 module vs. 主模块内包

选择独立 `crosser-api/` Go module：
- crosser-api 引入 SQLite (CGO)、JWT 等依赖，与现有纯网络代理模块职责完全不同
- 独立 module 可独立构建、独立发版，不影响现有 proxy 二进制体积
- 与现有 cmd/{server,client,auth_server} 三组件解耦

### Design Decision: SQLite vs. 外部数据库

选择 SQLite 嵌入式：
- 符合 NFR-03.01（单二进制、零外部依赖）
- 管理控制面数据量小（服务/用户/proxy 实例均为个位到百位级）
- 部署简单，无需运维额外数据库

## Capabilities

### New Capabilities
- `auth-login`: POST /api/v1/auth/login — JWT 登录认证
- `auth-refresh`: POST /api/v1/auth/refresh — Token 刷新
- `service-crud`: 代理服务增删改查 REST API
- `service-config-export`: GET /api/v1/services/:name/config — 兼容格式导出
- `proxy-register`: POST /api/v1/proxy/register — proxy 实例注册
- `proxy-heartbeat`: POST /api/v1/proxy/heartbeat — proxy 心跳上报
- `proxy-status`: GET /api/v1/proxy/status — 在线状态查询

### Modified Capabilities
- None（独立 module，不修改现有代码）

## Impact

- 新增 `crosser-api/` 目录（独立 Go module）
- 新增依赖：mattn/go-sqlite3 (CGO)、golang-jwt/jwt/v5、golang.org/x/crypto (bcrypt)
- 新增构建目标：`bin/crosser-api`
- Makefile 添加 crosser-api 构建/测试/lint 目标
- Dockerfile 新增 crosser-api 构建 stage
- 无对现有 proxy/client/auth_server 代码的修改
