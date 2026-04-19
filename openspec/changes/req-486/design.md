## Context

ubox-crosser 是一个基于 TCP 的 SOCKS5 反向代理系统，包含三个二进制（server、client、auth_server）。当前没有任何 HTTP 端点。运维需要一种标准方式查询服务版本。

现有构建系统使用 `-ldflags="-s -w"` 但未注入版本信息。项目已有 git tag（v1.0, v2.0, v2.1）可作为版本来源。

## Goals / Non-Goals

**Goals:**
- 在 proxy server 中提供 `GET /version` HTTP 端点返回版本号
- 构建时通过 `-ldflags` 注入版本信息（version、commit、build time）
- HTTP listener 地址可配置，有合理默认值

**Non-Goals:**
- 不为 client 和 auth_server 添加 HTTP 端点
- 不实现健康检查、metrics 等其他运维端点（后续可扩展）
- 不引入 HTTP 框架，使用标准库 `net/http`

## Decisions

### 1. 版本信息注入方式：ldflags

通过 `go build -ldflags "-X pkg.Var=value"` 在编译时注入版本号。

**备选方案**: 读取嵌入文件（`embed`）或运行时调用 `git describe`。ldflags 是 Go 生态最通用的做法，无运行时依赖，CI 友好。

### 2. 版本包位置：`internal/version`

新建 `internal/version/version.go`，包含 `Version`、`GitCommit`、`BuildTime` 变量及 `Info()` 方法。放在 `internal/` 下因为这是内部实现，不需要外部导入。

### 3. HTTP listener 嵌入 proxy server

在 `ProxyServer` 启动时额外启动一个 goroutine 运行 `net/http` server。

**备选方案**: 独立二进制或独立进程。嵌入方式更简单，减少部署复杂度，且 `/version` 足够轻量不会影响主服务。

### 4. 配置方式：ServerConfig 新增 HTTPAddress 字段

在 `ServerConfig` 的 common 段新增 `http_address` 字段（如 `":8080"`）。为空时不启动 HTTP listener。

### 5. 响应格式：JSON

```json
{"version": "2.1.0", "git_commit": "abc1234", "build_time": "2026-04-19T10:00:00Z"}
```

JSON 便于程序解析，也足够人类可读。

## Risks / Trade-offs

- **[风险] HTTP 端口暴露安全面** → 默认仅绑定 localhost；文档说明生产环境不应暴露到公网。HTTP address 为空时不启动，零开销。
- **[风险] 忘记在 CI 中传递 ldflags 导致版本为空** → 设置默认值 `"dev"` 作为未注入时的兜底，并在 Makefile 中使用 `git describe` 自动获取版本。
