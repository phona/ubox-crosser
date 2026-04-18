# CLAUDE.md

## 项目概述

UBox-Crosser 是 SOCKS5 反向代理系统，帮助 NAT/防火墙后的客户端暴露服务到公网。

| 组件 | 路径 | 说明 |
|------|------|------|
| ProxyServer | `server/` + `cmd/server/` | 中心代理中继 |
| Client | `client/` + `cmd/client/` | NAT 后客户端 |
| AuthServer | `server/auth_server.go` + `cmd/auth_server/` | 公网认证网关 |

## 常用命令

```bash
make build              # 构建 3 个二进制
make test               # 运行所有测试
make unit-test          # 仅单元测试
make unit-test-coverage # 单元测试 + 覆盖率
make lint               # golangci-lint
make fmt                # 格式化
make vet                # 静态检查
make test-integration   # Docker Compose 集成测试
make test-help          # 查看测试命令帮助
make sonar              # 本地 SonarQube 扫描
```

## Git 提交规范

```
<type>(<scope>): <subject>
```

类型：feat, fix, docs, style, refactor, perf, test, build, ci, chore

## 编码规范

详细规范按路径自动加载（`.claude/rules/`）：

| 规则文件 | 作用域 | 内容 |
|----------|--------|------|
| `go.md` | `**/*.go` | 架构、命名、关键约束 |
| `pr-precheck.md` | PR 创建前 | 自动检查流程 |
