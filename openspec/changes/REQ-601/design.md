# Design — REQ-601 GET /health 健康检查接口

## Context

`ubox-crosser` 代理服务（`server/server.go`）通过 `ProxyServer` 管理 TCP 长连接，使用自定义二进制协议（可选叠加 Shadowsocks 加密）。整个进程目前没有 HTTP 表面。`cmd/server/server.go` 是入口：构造 `ProxyServer`、`go proxy.Process()`、然后阻塞在错误上报循环。配置经 `utils/conf.ParseServerConfigFile` 加载到 `models/config.ServerConfig`，CLI flag 通过 cobra 绑定。

K8s liveness/readiness 与负载均衡器需要 HTTP 健康检查；裸 TCP 端口探测无法区分"端口在监听"与"业务真正可用"。

## Goals / Non-Goals

**Goals：**

- 提供 `GET /health` 端点，返回 `{"status":"ok"}` + HTTP 200。
- 监听地址通过配置文件 `health_address` 与 CLI `--health-address` 可配置；不配置则不启动。
- 复用 `ProxyServer` 的生命周期与错误通道，不引入新的 supervisor 抽象。
- 与现有 TCP 代理端口物理隔离（不同端口、不同 listener）。

**Non-Goals：**

- Readiness 检查（不探测 TCP listener 状态、controller 连接数等业务可用性）。
- Metrics / tracing / 其他 admin endpoint。
- `auth_server`、`client` 二进制的健康检查。
- TLS、鉴权、限流。
- Graceful shutdown（不持有 `*http.Server` 引用）。

## Decisions

### 1. 使用 stdlib `net/http`，不引入 router 框架

**Decision**：直接用 `net/http.ServeMux` + `http.HandleFunc`。

**Rationale**：endpoint 极简（一个路由 + 静态响应）。引入 Chi/Gin 只为单一路径属于过度设计。Go 1.22+ 的 `ServeMux` 原生支持 `"METHOD /path"` 模式，已能精确区分方法。

**Alternatives**：
- *Chi/Gin*：拒绝——单路由不值得新增依赖。
- *复用 TCP dispatcher*：拒绝——TCP 用自定义协议 + 加密，无法承载 HTTP，强行混合会污染连接抽象。

### 2. 独立 HTTP 监听器、独立端口

**Decision**：health server 监听独立 `host:port`，与 TCP 代理 listener 完全解耦。

**Rationale**：TCP listener 用自定义协议；HTTP 不能复用同一 socket。独立端口也允许通过防火墙/SecurityGroup 把健康检查只暴露给内网或 sidecar。

**默认端口建议**：`:8080`（仅作文档建议，实际默认空 = 不启动；用户必须显式设置才会监听）。

### 3. 路径用 `/health` 而非 `/healthz`

**Decision**：使用 `/health`。

**Rationale**：本次需求标题与产品方明确指定 `/health`。`/healthz` 是 Kubernetes 历史习惯，但 `/health` 在 Spring Boot Actuator、Consul、Nomad、Traefik、AWS ELB 等生态里同样是事实标准，且语义更直观。Kubernetes manifest 可显式指定 `path: /health`，无功能差异。

**Trade-off**：与社区 REQ-567 引入的 `/healthz` 不一致。如果未来希望两条路径并存，可在同一 mux 上注册别名，本次不做（YAGNI）。

### 4. 配置项 `health_address`，空值即禁用

**Decision**：`ServerConfig` 新增 `HealthAddress string \`json:"health_address"\``；CLI 新增 `--health-address`。空字符串 = 不启动 HTTP 监听器。

**Rationale**：与现有 `address` 字段同级、同风格。空值禁用保持向后兼容——已有部署不必任何变更即可继续运行。

### 5. 生命周期挂在 `ProxyServer` 内

**Decision**：在 `NewProxyServer` 中遍历所有 config，取首个非空 `HealthAddress`，启动 goroutine 跑 health server；启动错误（`ListenAndServe` 返回非 `ErrServerClosed`）写入 `ProxyServer.errs`。

**Rationale**：与现有错误通道复用，运维侧无需新增日志通道。Health server 是 per-process 的（不是 per-config-entry），即便配置文件里多个 section 都填了 `health_address`，也只启一个。

### 6. 严格的方法与路径处理

**Decision**：仅 `GET /health` 返回 200；`/health` 上其他方法返回 405 + `Allow: GET`；其他路径返回 404。

**Rationale**：遵守 HTTP 语义，避免误用与缓存代理的歧义。

### 7. 响应体使用字面量，不走 `encoding/json`

**Decision**：`w.Write([]byte(\`{"status":"ok"}\`))`。

**Rationale**：响应永远固定，避免 marshal 开销与 error 分支。简化测试——可以做精确字节比对。

## Risks / Trade-offs

- **[端口冲突]** → 默认空（不启动）规避；如显式配置端口与现有服务冲突，由 `net.Listen` 错误自然暴露并写入 errs channel。
- **[只 liveness 不 readiness]** → 探针只能反映"进程未死"，不能反映"controller 已连上"。在 Non-Goals 中显式声明；后续可加 `/ready` 端点。
- **[未加密 HTTP]** → 健康检查典型部署在内网，不引入 TLS 复杂度。如需要可前置 ingress/sidecar 终止 TLS。
- **[路径选型与 REQ-567 历史 `/healthz` 分叉]** → 本次按需求方明确 `/health` 实现；如未来需统一，可在同一 mux 上同时注册两条路径 alias，成本极低。
