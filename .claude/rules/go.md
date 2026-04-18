# Go 编码规范

## 项目架构

三组件 SOCKS5 反向代理系统：

| 组件 | 路径 | 说明 |
|------|------|------|
| Client | `client/` + `cmd/client/` | NAT 后客户端，主动连接 ProxyServer |
| Server | `server/` + `cmd/server/` | 中心代理，管理 controller 和 worker |
| AuthServer | `server/auth_server.go` + `cmd/auth_server/` | 公网入口网关 |

## 通信协议

- 消息格式：JSON + `\n` 分隔符
- 加密层：Shadowsocks（可选）
- 消息类型定义在 `models/message/messages.go`

## 命名规范

- 包名/文件名：snake_case
- 使用 `any` 替代 `interface{}`
- 使用 `return error`，不使用 `panic`

## 关键约束

- `fmt.Errorf` / `log.Errorf` 中 uint8 类型用 `%d` 不用 `%s`
- 结构体字面量必须使用 keyed fields
- 不使用 `io/ioutil`（已废弃），用 `io` 或 `os` 替代
- error 返回值必须检查（errcheck lint）

## 常用命令

```bash
make build              # 构建 3 个二进制
make test               # 运行测试
make unit-test          # 仅单元测试
make lint               # golangci-lint
make fmt                # 格式化
make vet                # 静态检查
make test-integration   # Docker Compose 集成测试
```
