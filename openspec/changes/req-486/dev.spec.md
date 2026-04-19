# DEV SPEC — REQ-486: GET /version 端点

## 文件结构

### 新增文件

| 文件 | 职责 |
|------|------|
| `internal/version/version.go` | 版本信息包：变量定义 + `Info()` 方法 |
| `internal/version/version_test.go` | version 包单元测试 |
| `server/http.go` | HTTP handler 和 listener 启动逻辑 |
| `server/http_test.go` | HTTP handler 单元测试 |

### 修改文件

| 文件 | 变更内容 |
|------|----------|
| `models/config/config.go` | `ServerConfig` 新增 `HTTPAddress` 字段 |
| `cmd/server/server.go` | 读取 common 配置中的 `http_address`，启动 HTTP listener |
| `Makefile` | `-ldflags` 注入 version、commit、build time |
| `Dockerfile` | 构建命令同步 ldflags 变更 |

---

## 函数签名与职责

### `internal/version/version.go`

```go
package version

// 由 -ldflags -X 在构建时注入，未注入时使用零值/默认值
var (
    Version   string = "dev"
    GitCommit string
    BuildTime string
)

// VersionInfo 是 JSON 响应的序列化结构体
type VersionInfo struct {
    Version   string `json:"version"`
    GitCommit string `json:"git_commit"`
    BuildTime string `json:"build_time"`
}

// Info 返回当前版本信息结构体
func Info() VersionInfo
```

**ldflags 注入路径**：`github.com/phona/ubox-crosser/internal/version.Version` 等。

### `server/http.go`

```go
package server

import "net/http"

// VersionHandler 返回处理 GET /version 的 http.HandlerFunc
// 响应 Content-Type: application/json，状态码 200
// 调用 version.Info() 获取数据，json.Marshal 序列化后写入 response body
func VersionHandler() http.HandlerFunc

// StartHTTPServer 在指定地址启动 HTTP listener
// addr 为空时直接返回（不启动）
// 注册路由 GET /version -> VersionHandler()
// 使用 goroutine 运行 http.ListenAndServe
// 监听失败时通过 errs channel 报告错误
func StartHTTPServer(addr string, errs chan<- error)
```

### `models/config/config.go` 变更

```go
type ServerConfig struct {
    LoginPass   string `json:"login_password"`
    AuthPass    string `json:"auth_password"`
    Address     string `json:"address"`
    HTTPAddress string `json:"http_address"` // 新增：HTTP 监听地址，为空则不启动
    Config
}
```

**注意**：`HTTPAddress` 从 `common` 段读取，通过现有 `Config.Update()` 方法继承。但 `Update()` 方法基于反射只处理 `Config` 结构体自身的字段，不会处理 `ServerConfig` 的直属字段。需要验证 `ParseServerConfigFile` 中 `json.Unmarshal` 是否能正确解析 `http_address` 到 `ServerConfig.HTTPAddress`——答案是**可以**，因为 `json.Unmarshal` 直接解析到 `ServerConfig`，`Update()` 只用于将 common 的 `Config` 嵌入字段覆盖到其他配置段。`HTTPAddress` 作为 `ServerConfig` 的直属字段会在 `json.Unmarshal` 阶段正确解析。但 common 段的 `HTTPAddress` 不会自动继承到其他配置段（因为 `Update()` 只处理嵌入的 `Config` 字段）。这是**符合预期的**——HTTP listener 只需要一个全局地址，应从 common 段读取。

### `cmd/server/server.go` 变更

在 `Run` 函数中，配置解析完成后、`proxy.Process()` 之前：

```go
// 从 common 配置读取 http_address
if commonConfig, ok := configs[conf.CommonConfigName]; ok && commonConfig.HTTPAddress != "" {
    server.StartHTTPServer(commonConfig.HTTPAddress, /* errs channel */)
}
```

需要将 `StartHTTPServer` 的 errs 与 proxy 的 errs 整合。方案：`StartHTTPServer` 内部 log.Fatal 或单独 log 错误即可，因为 HTTP listener 失败不应影响主代理服务。

---

## 依赖项

- **无新外部依赖**。仅使用 Go 标准库：
  - `net/http`：HTTP server 和路由
  - `encoding/json`：JSON 序列化
- 内部依赖：`internal/version` 被 `server/http.go` 和 `Makefile`（ldflags）引用

---

## 错误处理策略

| 场景 | 处理方式 |
|------|----------|
| `http_address` 为空 | 不启动 HTTP listener，静默跳过 |
| HTTP listener 绑定失败（端口占用等） | `log.Errorf` 记录错误，不中断主代理服务 |
| `json.Marshal` 失败（理论上不会发生） | 返回 HTTP 500 + 错误信息 |
| 非 GET 方法请求 /version | 标准库 `http.ServeMux` 默认允许所有方法，返回相同 JSON 响应（契约只定义了 GET，但不需要主动拒绝其他方法） |
| 请求未知路径 | 标准库默认返回 404 |

---

## 边界条件

1. **版本未注入**：`Version` 默认值为 `"dev"`，`GitCommit` 和 `BuildTime` 为空字符串 `""`
2. **`http_address` 仅在 common 段有效**：非 common 配置段的 `http_address` 被忽略，HTTP listener 只启动一个
3. **配置文件为空或不存在**：命令行模式下无 common 段，HTTP listener 不启动（除非通过命令行参数扩展，本期不做）
4. **并发安全**：版本信息变量在 init 后只读，无需加锁
5. **优雅关闭**：本期不实现 HTTP server 的优雅关闭（与现有 proxy server 行为一致，都是进程级别退出）

---

## 存储/DB 需求

无。版本信息全部在编译时注入，运行时只读。

---

## Makefile ldflags 变更

```makefile
# 新增变量
VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS    := -s -w \
    -X github.com/phona/ubox-crosser/internal/version.Version=$(VERSION) \
    -X github.com/phona/ubox-crosser/internal/version.GitCommit=$(GIT_COMMIT) \
    -X github.com/phona/ubox-crosser/internal/version.BuildTime=$(BUILD_TIME)

# build target 修改
build: $(SOURCES)
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/client ./cmd/client
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/server ./cmd/server
	CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/auth_server ./cmd/auth_server
```

### Dockerfile 变更

```dockerfile
# Builder stage 中增加相同的 ldflags
ARG VERSION=dev
RUN VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev") && \
    GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown") && \
    BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ) && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w \
        -X github.com/phona/ubox-crosser/internal/version.Version=${VERSION} \
        -X github.com/phona/ubox-crosser/internal/version.GitCommit=${GIT_COMMIT} \
        -X github.com/phona/ubox-crosser/internal/version.BuildTime=${BUILD_TIME}" \
    -o /app/crosser ./cmd/${BINARY}
```

---

## 测试计划

### 单元测试

1. **`internal/version/version_test.go`**
   - 验证 `Info()` 返回默认值（`Version="dev"`, `GitCommit=""`, `BuildTime=""`）
   - 验证 `VersionInfo` JSON 序列化字段名正确

2. **`server/http_test.go`**
   - 使用 `httptest.NewRecorder` 测试 `VersionHandler`
   - 验证响应状态码 200
   - 验证 Content-Type 为 `application/json`
   - 验证响应体 JSON 包含 `version`、`git_commit`、`build_time` 三个字段
   - 验证默认情况下 `version` 字段值为 `"dev"`

### 集成测试（后续，不在本 dev spec 范围）

- 构建二进制后 curl `GET /version` 验证完整流程
