## Context

ubox-crosser 是一个基于 TCP 自定义协议的 SOCKS5 反向代理，目前没有任何 HTTP 接口。三个二进制（client、server、auth_server）均通过 cobra CLI 启动，构建时仅 strip 符号（`-ldflags="-s -w"`），未注入版本信息。

运维需要一种标准方式查询运行中实例的版本号，HTTP `GET /version` 是最通用、最易集成的方案。

## Goals / Non-Goals

**Goals:**
- 在 proxy server 中提供 `GET /version` HTTP 端点
- 通过 build-time ldflags 注入版本号、commit、构建时间
- 版本信息集中管理，可被多个二进制复用

**Non-Goals:**
- 不为 client 或 auth_server 添加 HTTP 端点（后续可扩展）
- 不实现健康检查、metrics 等其他 HTTP 端点（本次仅 version）
- 不引入外部 HTTP 框架（标准库 `net/http` 足够）

## Decisions

### D1: 新增 `version` 包存放版本变量

在项目根目录新增 `version/version.go`，包含 `Version`、`Commit`、`BuildTime` 三个包级变量，由 ldflags 注入。

**备选方案**: 放在 `main` 包中 → 无法跨二进制复用，排除。

### D2: HTTP server 内嵌在 proxy server 进程

在 `cmd/server/server.go` 中启动一个额外的 goroutine 运行 `net/http` server，与 TCP proxy 共存于同一进程。

**备选方案**: 独立的 HTTP 进程 → 增加部署复杂度，对于单一端点过度设计，排除。

### D3: HTTP 监听地址通过配置指定

`ServerConfig` 新增 `HTTPAddress` 字段（config file: `http_address`，CLI: `--http-address`）。未指定时不启动 HTTP server，保持向后兼容。

**备选方案**: 硬编码端口 → 不灵活，端口冲突风险，排除。

### D4: Makefile 自动注入版本信息

```makefile
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w \
  -X github.com/phona/ubox-crosser/version.Version=$(VERSION) \
  -X github.com/phona/ubox-crosser/version.Commit=$(COMMIT) \
  -X github.com/phona/ubox-crosser/version.BuildTime=$(BUILD_TIME)
```

## Risks / Trade-offs

- **[端口冲突]** HTTP 端口可能与其他服务冲突 → 通过配置化解决，不启动时无风险。
- **[安全暴露]** 版本信息可能暴露给攻击者 → 版本端点通常仅在内网监听，且信息量有限，风险可接受。
- **[HTTP server 生命周期]** HTTP server 的错误不应导致 proxy 主进程退出 → 使用 goroutine 运行，仅记录错误日志。
