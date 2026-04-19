# DEV SPEC — REQ-495: GET /version 接口

## 文件结构

### 新增文件

| 文件 | 用途 |
|------|------|
| `version/version.go` | 版本变量定义（`Version`, `Commit`, `BuildTime`） |
| `version/handler.go` | HTTP handler，返回版本 JSON |

### 修改文件

| 文件 | 变更内容 |
|------|----------|
| `models/config/config.go` | `ServerConfig` 新增 `HTTPAddress` 字段 |
| `cmd/server/server.go` | 新增 `--http-address` flag；当 HTTPAddress 非空时启动 HTTP server goroutine |
| `Makefile` | build target 的 ldflags 增加 `-X` 注入版本变量 |
| `Dockerfile` | 构建阶段增加版本 ldflags（通过 ARG 传入或 Makefile 复用） |

---

## 函数签名与职责

### `version/version.go`

```go
package version

var (
    Version   string = "dev"
    Commit    string = "unknown"
    BuildTime string = "unknown"
)
```

- 三个包级变量，由 `-ldflags -X` 在编译时注入
- 默认值保证未注入时返回合理值

### `version/handler.go`

```go
package version

import "net/http"

func Handler() http.HandlerFunc
```

- 返回一个 `http.HandlerFunc`
- 序列化 `VersionInfo{Version, Commit, BuildTime}` 为 JSON
- 设置 `Content-Type: application/json`，写入 HTTP 200 响应
- 内部定义响应结构体：

```go
type VersionInfo struct {
    Version   string `json:"version"`
    Commit    string `json:"commit"`
    BuildTime string `json:"buildTime"`
}
```

### `models/config/config.go` 变更

```go
type ServerConfig struct {
    LoginPass   string `json:"login_password"`
    AuthPass    string `json:"auth_password"`
    Address     string `json:"address"`
    HTTPAddress string `json:"http_address"`  // 新增
    Config
}
```

- `HTTPAddress` 为空字符串时不启动 HTTP server（向后兼容）
- JSON tag `http_address` 与 config file 字段名一致

### `cmd/server/server.go` 变更

新增 CLI flag：

```go
cmd.Flags().StringVar(&cmdConfig.HTTPAddress, "http-address", "", "HTTP listen address for version endpoint (e.g. :8080)")
```

新增 HTTP server 启动逻辑（在 `Run` 函数内、proxy 启动前后均可）：

```go
func startHTTPServer(addr string)
```

- 签名：接收监听地址字符串
- 职责：创建 `http.ServeMux`，注册 `GET /version` handler，在新 goroutine 中调用 `http.ListenAndServe`
- 仅当 `HTTPAddress` 非空时调用
- HTTP server 的 `ListenAndServe` 错误仅记录日志（`logrus.Errorf`），不 panic、不退出主进程

逻辑位置：在 `Run` 函数内，configs 解析完成后、proxy 启动前，遍历 configs 找到任一有效 `HTTPAddress` 即启动。如果多个 config 块有不同 HTTPAddress，取 common 或第一个非空值（单进程只需一个 HTTP listener）。

### `Makefile` 变更

```makefile
VERSION    ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
COMMIT     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w \
  -X github.com/phona/ubox-crosser/version.Version=$(VERSION) \
  -X github.com/phona/ubox-crosser/version.Commit=$(COMMIT) \
  -X github.com/phona/ubox-crosser/version.BuildTime=$(BUILD_TIME)
```

- 替换现有 build target 中的 `-ldflags="-s -w"` 为 `-ldflags="$(LDFLAGS)"`
- 三个二进制（client、server、auth_server）统一注入相同版本信息
- `ci-build` 中的 `docker build` 命令通过 `--build-arg` 传递 LDFLAGS

### `Dockerfile` 变更

```dockerfile
ARG LDFLAGS="-s -w"

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="${LDFLAGS}" -o /app/crosser ./cmd/${BINARY}
```

- 新增 `ARG LDFLAGS`，默认 `"-s -w"`（无 ARG 时行为不变）
- `ci-build` 通过 `--build-arg LDFLAGS="..."` 传入完整 ldflags

---

## 依赖项

- **无新外部依赖**
- 仅使用 Go 标准库：`net/http`、`encoding/json`
- 已有依赖 `logrus`（记录 HTTP server 错误日志）

---

## 错误处理策略

| 场景 | 处理方式 |
|------|----------|
| `http.ListenAndServe` 返回错误 | `logrus.Errorf` 记录日志，不退出进程 |
| `json.Marshal` 失败（理论上不会） | 返回 HTTP 500 + 错误消息 |
| HTTPAddress 为空 | 不启动 HTTP server，无错误 |
| HTTPAddress 格式错误 | `ListenAndServe` 会立即返回错误，由日志记录 |

---

## 边界条件

1. **HTTPAddress 为空字符串**：不启动 HTTP server，proxy 行为与当前完全一致
2. **HTTPAddress 端口被占用**：`ListenAndServe` 返回 `bind: address already in use`，记录日志
3. **未注入 ldflags**：返回默认值 `{"version":"dev","commit":"unknown","buildTime":"unknown"}`
4. **并发请求**：`version` 包变量在初始化后只读，handler 天然并发安全
5. **非 GET 方法访问 /version**：`http.ServeMux` 默认不区分方法，返回相同 JSON（符合契约，契约仅定义 GET 200，未禁止其他方法）
6. **多 config 块的 HTTPAddress 冲突**：取第一个非空值，忽略其余，记录日志提示

---

## 存储/DB 需求

无。版本信息为编译时常量，无运行时持久化需求。

---

## 实现顺序建议

1. `version/version.go` — 纯数据，无依赖
2. `version/handler.go` — 依赖 step 1
3. `models/config/config.go` — 新增字段
4. `cmd/server/server.go` — 集成 HTTP server，依赖 step 1-3
5. `Makefile` — ldflags 注入
6. `Dockerfile` — ARG LDFLAGS 支持
