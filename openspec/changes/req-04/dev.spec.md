# Dev Spec — REQ-04 Health Endpoint

## 文件结构

### 修改文件

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `models/config/config.go` | 修改 | ServerConfig 新增 `HealthAddress` 字段 |
| `server/server.go` | 修改 | ProxyServer 启动 health HTTP server |
| `cmd/server/server.go` | 修改 | 添加 `--health-address` CLI flag |
| `utils/conf/utils.go` | 无需修改 | 现有 `ParseServerConfigFile()` 自动解析新 JSON 字段 |

### 新增文件

| 文件 | 说明 |
|------|------|
| `server/health.go` | health HTTP handler 与 server 启动逻辑 |

---

## 函数签名与职责

### `server/health.go`（新增）

```go
func newHealthMux() *http.ServeMux
```
- 创建 `http.ServeMux`，注册 `/health` 路由
- `/health` handler：GET 返回 200 + `{"status":"ok"}` + `Content-Type: application/json`；非 GET 返回 405 + `Allow: GET` header
- 默认 handler（catch-all）：所有非 `/health` 路径返回 404

```go
func (p *ProxyServer) startHealthServer(addr string)
```
- 创建 `http.Server{Addr: addr, Handler: newHealthMux()}`
- 调用 `ListenAndServe()`
- 如果返回非 `http.ErrServerClosed` 错误，发送到 `p.errs` channel
- 由 `NewProxyServer` 在 goroutine 中调用

### `models/config/config.go`（修改）

```go
type ServerConfig struct {
    LoginPass     string `json:"login_password"`
    AuthPass      string `json:"auth_password"`
    Address       string `json:"address"`
    HealthAddress string `json:"health_address"`  // 新增
    Config
}
```
- 新增 `HealthAddress` 字段，JSON tag `"health_address"`
- 空字符串表示不启动 health server

### `server/server.go`（修改）

`NewProxyServer` 函数变更：
- 遍历 `configs` 中所有 config，提取第一个非空的 `HealthAddress`（或从 "common" 继承）
- 如果 `HealthAddress` 非空，调用 `go p.startHealthServer(healthAddr)`
- 位于现有 `initWorker` goroutine 启动之后

### `cmd/server/server.go`（修改）

```go
serverCmd.PersistentFlags().StringVarP(&cmdConfig.HealthAddress, "health-address", "", "", "health check HTTP listen address (e.g. :8080)")
```
- 新增 CLI flag `--health-address`
- 无短标志（`-`）以避免与现有单字母标志冲突
- Flag 绑定到 `cmdConfig.HealthAddress`

---

## 依赖项

- **无新外部依赖**。仅使用 Go stdlib：`net/http`、`encoding/json`（json.Marshal 或直接写 literal）
- 现有依赖（cobra, logrus）照常使用

---

## 错误处理策略

| 场景 | 处理方式 |
|------|----------|
| `health_address` 为空或未设置 | 跳过 health server 启动，无日志 |
| 端口被占用（`ListenAndServe` 失败） | 错误发送到 `p.errs` channel，通过主 goroutine 的 `proxy.Err()` 循环输出到 logrus |
| `http.ErrServerClosed` | 忽略，不发送到 errs channel（正常关闭场景） |
| 非 GET 请求到 `/health` | 返回 HTTP 405，设置 `Allow: GET` header |
| 请求到非 `/health` 路径 | 返回 HTTP 404 |
| JSON 序列化失败 | 不会发生——响应体是固定字面量 `{"status":"ok"}`，直接 `io.WriteString` 写入 |

---

## 边界条件

1. **多个 config section 都设置了 `health_address`**：只启动一个 health server。取 "common" section 的值；如果 "common" 没有设置，取第一个非空值。health server 是 per-process 的，不是 per-proxy-entry 的。

2. **`health_address` 格式无效**（如 `"abc"`）：`ListenAndServe` 返回错误，通过 errs channel 报告。不做预校验——让 `net.Listen` 自然失败。

3. **并发请求**：`http.Server` 天然支持并发，无需额外同步。

4. **`health_address` 与 proxy `address` 相同端口**：会导致端口冲突。不做主动检测——属于配置错误，由 bind 错误自然暴露。

5. **响应 Content-Type**：`GET /health` 必须设置 `Content-Type: application/json`。405 和 404 响应无需 body，但 405 必须设置 `Allow: GET` header。

6. **路径匹配**：使用 `http.ServeMux` 的精确匹配。Go 1.22+ ServeMux 支持 `GET /health` 方法+路径模式，可直接用来区分方法，但为兼容 Go 1.23 module 要求，确认 `go.mod` 中 Go 版本 ≥ 1.22 后可使用新路由语法。否则在 handler 内手动检查 `r.Method`。

7. **Trailing slash**：`/health/` 与 `/health` 视为不同路径。`/health/` 应返回 404。

---

## 存储/DB 需求

无。Health endpoint 返回静态响应，不查询任何状态。

---

## 实现注意事项

1. **响应体使用字面量**：直接 `w.Write([]byte(`{"status":"ok"}`))` 而非 `json.Marshal`，避免不必要的序列化开销和 error handling。

2. **ServeMux 路由注册**（Go 1.22+ 新语法）：
   - `mux.HandleFunc("GET /health", healthHandler)` — 精确匹配 GET 方法 + `/health` 路径
   - 非 GET 方法需要额外注册 `mux.HandleFunc("/health", methodNotAllowedHandler)` 作为 fallback
   - 或者只注册 `/health` 一个路由，在 handler 内部检查 `r.Method`

3. **http.Server 生命周期**：保存 `*http.Server` 到 ProxyServer 字段中可为未来 graceful shutdown 做准备，但当前 scope 不要求 shutdown 功能，可先不保存（design.md 提到 "can be cleanly shut down in the future if needed"）。如果实现时选择保存，添加一个 `healthServer *http.Server` 字段即可。

4. **Config 继承**：现有 `Config.Update()` 使用反射合并字段（`utils/conf/utils.go:48-56`），但 `HealthAddress` 在 `ServerConfig` 上而非嵌入的 `Config` 上。需确认 "common" section 的 `health_address` 能正确继承到其他 config section。查看 `ParseServerConfigFile` 逻辑——它先解析为 `map[string]ServerConfig`，然后 `common.Update()` 合并到各 config 中。`Update()` 只处理 `Config` 嵌入字段。**因此 `HealthAddress` 需要在 cmd/server 中从 "common" config 手动提取**，或者将 `HealthAddress` 的合并逻辑加到 `ServerConfig` 级别。最简单方案：在 `NewProxyServer` 中遍历 configs map 查找非空 `HealthAddress`。

5. **CLI flag 优先级**：现有模式是 CLI flag 覆盖 config file。`cmdConfig.HealthAddress` 通过 cobra flag 设置，在 config file 为空时生效（参考 `cmd/server/server.go` line 24-40 的 fallback 逻辑）。
