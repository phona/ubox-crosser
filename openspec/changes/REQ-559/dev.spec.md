# REQ-559 Dev Spec: 在 internal/config 加版本号常量

## 概述

在 `internal/config` 包中新增版本号常量，供所有三个二进制 (client, server, auth_server) 在启动时引用。同时通过 `-ldflags` 支持构建时注入 git commit / 构建时间等元数据。

## 文件结构

### 新增文件

| 文件 | 职责 |
|------|------|
| `internal/config/version.go` | 定义版本号常量与构建元数据变量 |

### 修改文件

| 文件 | 变更 |
|------|------|
| `cmd/client/client.go` | cobra.Command 添加 `Version` 字段，引用 `internal/config` |
| `cmd/server/server.go` | 同上 |
| `cmd/auth_server/server.go` | 同上 |
| `Makefile` | `go build -ldflags` 注入 `GitCommit` 和 `BuildTime` |
| `Dockerfile` | 构建命令同步更新 ldflags |

## 函数签名与职责

### `internal/config/version.go`

```go
package config

// 硬编码版本号，发版时手动更新
const Version = "0.1.0"

// 构建时通过 -ldflags 注入的元数据
var (
    GitCommit string // -X internal/config.GitCommit=$(git rev-parse --short HEAD)
    BuildTime string // -X internal/config.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)
)

// FullVersion 返回包含构建元数据的完整版本字符串
func FullVersion() string
```

- `FullVersion()` 返回格式：`0.1.0 (commit: abc1234, built: 2026-04-20T12:00:00Z)`
- 当 `GitCommit` 或 `BuildTime` 为空（本地 `go build` 未传 ldflags）时，返回 `0.1.0 (commit: unknown, built: unknown)`

## 依赖项

- **无新外部依赖**。仅使用 `fmt` 标准库。
- `internal/config` 是新包路径，与现有 `models/config` 不冲突。现有 `models/config` 保持不变。

## ldflags 注入

### Makefile 变更

```makefile
LDFLAGS := -s -w \
    -X github.com/phona/ubox-crosser/internal/config.GitCommit=$(shell git rev-parse --short HEAD) \
    -X github.com/phona/ubox-crosser/internal/config.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)

build:
    CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" -o bin/client ./cmd/client
    ...
```

### Dockerfile 变更

在 builder 阶段添加 git commit 参数：

```dockerfile
ARG GIT_COMMIT=unknown
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w \
      -X github.com/phona/ubox-crosser/internal/config.GitCommit=${GIT_COMMIT} \
      -X github.com/phona/ubox-crosser/internal/config.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /app/crosser ./cmd/${BINARY}
```

## cmd 变更

每个 `cmd/*/main()` 中的 `cobra.Command` 添加：

```go
import internalconfig "github.com/phona/ubox-crosser/internal/config"

cmd := &cobra.Command{
    Use:     "...",
    Version: internalconfig.FullVersion(),
    Run:     func(...) { ... },
}
```

这使得 `--version` 标志自动可用。

## 错误处理策略

- 本变更不引入任何可能失败的操作，无需错误处理。
- `GitCommit` / `BuildTime` 未注入时使用 `"unknown"` 兜底，不 panic。

## 边界条件

| 场景 | 行为 |
|------|------|
| `go build` 不带 `-ldflags` | `GitCommit=""`, `BuildTime=""` → FullVersion 显示 `unknown` |
| `go build` 带完整 `-ldflags` | 正常显示注入值 |
| `go run ./cmd/client --version` | 输出版本字符串并退出（cobra 内置行为） |
| 常量 `Version` 为空字符串 | 不允许，编译期不会检查，需代码审查保证 |

## 存储 / DB 需求

无。

## 测试要点（供测试 agent 参考）

1. `FullVersion()` 在 ldflags 未注入时返回含 `unknown` 的字符串
2. `FullVersion()` 在 ldflags 注入时返回正确格式
3. `make build` 后 `./bin/client --version` 输出包含版本号
4. `Version` 常量值为有效 semver 格式
