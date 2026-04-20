# Dev Spec — REQ-567 GET /healthz 健康检查接口

## 文件结构

### 新增文件

| 文件 | 说明 |
|------|------|
| `server/health.go` | healthz HTTP handler 与 health server 启动逻辑 |

### 修改文件

| 文件 | 变更类型 | 说明 |
|------|----------|------|
| `models/config/config.go` | 修改 | `ServerConfig` 新增 `HealthAddress` 字段 |
| `server/server.go` | 修改 | `NewProxyServer` 中启动 health HTTP server |
| `cmd/server/server.go` | 修改 | 添加 `--health-address` CLI flag |

### 无需修改

| 文件 | 原因 |
|------|------|
| `utils/conf/utils.go` | `ParseServerConfigFile()` 解析 JSON 到 `ServerConfig`，新增字段自动反序列化。但注意 `Config.Update()` 只处理嵌入的 `Config` 字段，`HealthAddress` 在 `ServerConfig` 层，不会被 common→service 自动继承——在 `NewProxyServer` 中额外处理 |

---

## 函数签名与职责

### `server/health.go`（新增）

```go
package server

import "net/http"

func newHealthMux() *http.ServeMux
```
- 创建 `http.ServeMux`
- 注册 `/healthz` 路由，handler 逻辑：
  - `r.Method == "GET"` → 设置 `Content-Type: application/json`，写入 `{"status":"ok"}`，状态码 200
  - 其他方法 → 设置 `Allow: GET` header，返回 405
- 默认 handler（catch-all）：所有非 `/healthz` 路径返回 404
- Go 1.23 ServeMux 支持方法+路径模式，可用 `mux.HandleFunc("GET /healthz", ...)` 精确匹配 GET，再注册 `mux.HandleFunc("/healthz", ...)` 作为非 GET 的 fallback 返回 405

```go
func (p *ProxyServer) startHealthServer(addr string)
```
- 创建 `http.Server{Addr: addr, Handler: newHealthMux()}`
- 调用 `ListenAndServe()`
- 如果返回非 `http.ErrServerClosed` 错误，发送到 `p.errs` channel
- 由 `NewProxyServer` 在 goroutine 中调用：`go p.startHealthServer(healthAddr)`

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
- 字段位于 `ServerConfig` 上而非嵌入的 `Config` 上，与 `Address` 同级

### `server/server.go`（修改）

`NewProxyServer` 函数末尾追加逻辑：

```go
// 在 initWorker goroutine 启动之后、return server 之前
var healthAddr string
for _, cfg := range configs {
    if cfg.HealthAddress != "" {
        healthAddr = cfg.HealthAddress
        break
    }
}
if healthAddr != "" {
    go server.startHealthServer(healthAddr)
}
```

- 遍历所有 config，取第一个非空的 `HealthAddress`
- health server 是 per-process 的（只启动一个），不是 per-proxy-entry 的
- 优先取 common config 的值（遍历 map 顺序不确定，但 common 的值已通过 JSON 解析写入各 config entry）

### `cmd/server/server.go`（修改）

在现有 flag 定义之后添加：

```go
cmd.Flags().StringVar(&cmdConfig.HealthAddress, "health-address", "", "health check HTTP listen address (e.g. :8080)")
```

- 使用 `StringVar`（无短标志），与现有 `--log-file`、`--log-level`、`--config-file` 风格一致
- 无默认值——不设置则不启动 health server

---

## 依赖项

- **无新外部依赖**
- 仅使用 Go stdlib：`net/http`
- 响应体使用字面量 `[]byte(`{"status":"ok"}`)` 而非 `encoding/json`，避免不必要的序列化
- 现有依赖（cobra, logrus）照常使用，health server 中用 logrus 记录启动日志

---

## 错误处理策略

| 场景 | 处理方式 |
|------|----------|
| `health_address` 为空或未设置 | 跳过 health server 启动，不记录日志 |
| 端口被占用（`ListenAndServe` 失败） | 错误发送到 `p.errs` channel，通过主 goroutine 的 `proxy.Err()` 循环输出到 logrus |
| `http.ErrServerClosed` | 忽略，不发送到 errs channel（正常关闭场景） |
| 非 GET 请求到 `/healthz` | 返回 HTTP 405，设置 `Allow: GET` header，无 body |
| 请求到非 `/healthz` 路径 | 返回 HTTP 404（ServeMux 默认行为） |
| JSON 序列化失败 | 不会发生——响应体是固定字面量 `{"status":"ok"}`，直接 `w.Write` 写入 |

---

## 边界条件

1. **多个 config section 都设置了 `health_address`**：只启动一个 health server。遍历 configs map 取第一个非空值。Health server 是 per-process 的。

2. **`health_address` 格式无效**（如 `"abc"`）：`ListenAndServe` 返回错误，通过 errs channel 报告。不做预校验——让 `net.Listen` 自然失败。

3. **并发请求**：`http.Server` 天然支持并发，无需额外同步。

4. **`health_address` 与 proxy `address` 相同端口**：导致端口冲突。不做主动检测——属于配置错误，由 bind 错误自然暴露。

5. **响应 Content-Type**：`GET /healthz` 必须设置 `Content-Type: application/json`。405 和 404 响应无需 body。405 必须设置 `Allow: GET` header。

6. **路径匹配**：使用 Go 1.23 `http.ServeMux` 精确匹配。`/healthz/` 与 `/healthz` 视为不同路径。`/healthz/` 应返回 404。

7. **Trailing slash**：`/healthz/` 不匹配 `/healthz` 注册的路由，由 catch-all 返回 404。

8. **Config 继承限制**：`Config.Update()` 方法（`utils/conf/utils.go:36-72`）使用反射只处理嵌入的 `Config` struct 字段。`HealthAddress` 在 `ServerConfig` 层，不会被 common→service 自动合并。但 `ParseServerConfigFile` 先 `common.Update(v)` 合并 `Config` 嵌入字段，`HealthAddress` 的值来自 JSON 反序列化（`json.Unmarshal`），如果 common section 设置了 `health_address`，`newConfig := common` 的值拷贝已经包含了 `HealthAddress`。所以 common 的 `health_address` **会** 被继承到其他 config section。

---

## 存储/DB 需求

无。Health endpoint 返回静态响应，不查询任何状态。

---

## 实现注意事项

1. **响应体使用字面量**：`w.Write([]byte(`{"status":"ok"}`))` 而非 `json.Marshal`，避免不必要的序列化开销和 error handling。

2. **ServeMux 路由注册**（Go 1.23 语法）：
   ```go
   mux := http.NewServeMux()
   mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
       w.Header().Set("Content-Type", "application/json")
       w.Write([]byte(`{"status":"ok"}`))
   })
   mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
       w.Header().Set("Allow", "GET")
       w.WriteHeader(http.StatusMethodNotAllowed)
   })
   ```
   `"GET /healthz"` 精确匹配 GET 方法。其他方法落入 `"/healthz"` 的 fallback handler。
   未注册的路径由 ServeMux 默认返回 404。

3. **启动日志**：在 `startHealthServer` 中调用 `ListenAndServe` 之前用 `log.Infof("Health server listening on %s", addr)` 记录。

4. **不保存 `*http.Server` 引用**：当前 scope 不要求 graceful shutdown。如果未来需要，可添加 `healthServer *http.Server` 字段到 `ProxyServer`。
