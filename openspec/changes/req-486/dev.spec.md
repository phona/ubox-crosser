# DEV SPEC — REQ-486: GET /version 端点

## 文件结构

### 新增文件

| 文件 | 职责 |
|------|------|
| `internal/version/version.go` | 版本信息包：编译时注入变量 + 结构体 + JSON 序列化 |
| `internal/version/version_test.go` | version 包单元测试 |
| `server/http.go` | HTTP handler 和 HTTP server 启动逻辑 |
| `server/http_test.go` | HTTP handler 单元测试 |
| `tests/contract/version_contract_test.go` | API 契约测试（验证响应格式符合 contract.spec.yaml） |
| `tests/integration/version_endpoint_test.go` | 集成测试（Docker Compose 环境下端到端验证） |

### 修改文件

| 文件 | 变更内容 |
|------|----------|
| `models/config/config.go` | `ServerConfig` 新增 `HTTPAddress string` 字段（json tag: `http_address`） |
| `server/server.go` | `NewProxyServer` 中根据 `HTTPAddress` 配置启动 HTTP listener goroutine |
| `Makefile` | 新增 `VERSION`、`GIT_COMMIT`、`BUILD_TIME` 变量；`LDFLAGS` 增加 `-X` 注入 |
| `Dockerfile` | builder 阶段增加 `VERSION`、`GIT_COMMIT`、`BUILD_TIME` ARG 并传入 ldflags |
| `tests/configs/server.json` | 新增 `http_address` 配置字段 |
| `tests/docker-compose.yml` | 暴露 HTTP 端口供集成测试访问 |

## 函数签名与职责

### `internal/version` 包

```go
// 编译时通过 -ldflags -X 注入的包级变量
var (
    Version   string  // 默认 "dev"
    GitCommit string  // 默认 ""
    BuildTime string  // 默认 ""
)

// Info 版本信息结构体，对应 API 响应 JSON 格式
type Info struct {
    Version   string `json:"version"`
    GitCommit string `json:"git_commit"`
    BuildTime string `json:"build_time"`
}

// GetInfo 返回当前注入的版本信息
func GetInfo() Info

// JSON 将 Info 序列化为 JSON 字节切片
func (i Info) JSON() []byte
```

### `server` 包 — HTTP 层

```go
// versionHandler 处理 GET /version 请求，返回 JSON 版本信息
// 设置 Content-Type: application/json，写入 version.GetInfo().JSON()
func versionHandler(w http.ResponseWriter, r *http.Request)

// startHTTPServer 启动 HTTP listener
// 注册 "GET /version" 路由到 versionHandler，阻塞调用 http.ListenAndServe
// 错误时记录日志（不 panic）
func startHTTPServer(addr string)
```

### `server/server.go` — 启动集成

`NewProxyServer` 遍历 configs 时，对每个非空 `HTTPAddress`：
- 用 `httpStarted map[string]bool` 去重（避免多个 config 共享同一 HTTP 地址时重复监听）
- 以 goroutine 调用 `startHTTPServer(config_.HTTPAddress)`

### `models/config/config.go` — 配置

```go
type ServerConfig struct {
    // ... 现有字段 ...
    HTTPAddress string `json:"http_address"` // 新增，可选，为空时不启动 HTTP listener
}
```

## 依赖项

- **无新外部依赖**。全部使用 Go 标准库：
  - `net/http` — HTTP server 和 handler
  - `encoding/json` — JSON 序列化
  - `net/http/httptest` — 测试用 HTTP recorder
- 现有依赖 `github.com/sirupsen/logrus` 用于 HTTP server 日志

## 错误处理策略

| 场景 | 处理方式 |
|------|----------|
| `http.ListenAndServe` 失败（端口被占、权限不足等） | `log.Errorf` 记录错误，不 panic/crash 主进程 |
| `json.Marshal` 失败 | `Info.JSON()` 忽略错误（结构体为纯 string 字段，Marshal 不会失败） |
| `HTTPAddress` 为空 | 不启动 HTTP listener，零开销 |
| 构建时未传 `-ldflags` | Version 默认 "dev"，GitCommit/BuildTime 默认空字符串 |

## 边界条件

1. **多 ServerConfig 共享同一 `http_address`**：通过 `httpStarted` map 去重，只启动一次
2. **`http_address` 全部为空**：不启动任何 HTTP listener，行为与变更前完全一致
3. **非 GET 请求访问 /version**：Go 1.22+ `ServeMux` 的 `"GET /version"` 模式自动返回 405 Method Not Allowed
4. **访问未注册路径**：`ServeMux` 默认返回 404
5. **Version 字段未注入**：返回 `{"version":"dev","git_commit":"","build_time":""}`
6. **HTTP server goroutine 退出**：仅记录日志，不影响主 TCP proxy 功能

## 构建注入

Makefile 变量：

```makefile
VERSION     ?= $(shell git describe --tags --always 2>/dev/null || echo "dev")
GIT_COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME  ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION_PKG := github.com/phona/ubox-crosser/internal/version

LDFLAGS := -s -w \
    -X $(VERSION_PKG).Version=$(VERSION) \
    -X $(VERSION_PKG).GitCommit=$(GIT_COMMIT) \
    -X $(VERSION_PKG).BuildTime=$(BUILD_TIME)
```

Dockerfile 通过 `ARG` 接收并传入相同的 ldflags。

## 存储/DB 需求

无。版本信息为编译时常量，不涉及任何持久化存储。

## 测试覆盖矩阵

| 测试层 | 文件 | 验证项 |
|--------|------|--------|
| 单元测试 | `internal/version/version_test.go` | 默认值正确；JSON 序列化包含全部 3 个字段 |
| 单元测试 | `server/http_test.go` | handler 返回 200 + `application/json`；响应体包含所有必需字段 |
| 契约测试 | `tests/contract/version_contract_test.go` | 响应格式符合 contract.spec.yaml 定义 |
| 集成测试 | `tests/integration/version_endpoint_test.go` | Docker 环境中端到端 HTTP 调用验证 |
