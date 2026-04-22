# 设计方案：/version 端点

## 方案概述
在 ubox-crosser 中新增 HTTP GET 端点 `/version`，通过在构建时注入 git SHA 实现版本信息的返回。

## 架构选择

### 1. 版本信息注入方案
**选择**: 构建时通过 `-ldflags` 注入（Go 标准做法）
- ✅ 无运行时开销
- ✅ 不依赖环境变量或文件
- ✅ 部署时无需额外配置

```go
// 编译时: go build -ldflags "-X main.GitSHA=<commit-sha>"
var GitSHA string = "unknown"
```

### 2. 端点实现位置
**选择**: 在 main HTTP router 中注册，与现有端点同级
- 位置: `cmd/` 目录的 main 或 route 定义处
- 优势: 启动快、无依赖、易测试

### 3. 响应格式
**选择**: 返回 JSON 对象，包含 sha 和可选的 build_time/version
```json
{
  "sha": "abc123def456...",
  "build_time": "2026-04-22T12:00:00Z"  // 可选
}
```

## 关键实现细节

### 路由注册
- 路径: `/version`
- 方法: GET
- 响应码: 200 (正常)

### 依赖项
- 零额外依赖
- 仅使用 Go 标准库

### 版本更新流程
1. 代码提交到 git
2. CI 流程（Makefile/GitHub Actions）获取 `git rev-parse HEAD`
3. 编译时传入 `-ldflags "-X main.GitSHA=<SHA>"`
4. 部署时自动包含正确的版本信息

## 风险与缓解

| 风险 | 缓解措施 |
|------|--------|
| 版本信息过期 | CI 流程保证每次编译都注入当前 SHA |
| 端点泄露内部信息 | SHA 本身是公开的，无安全风险 |
| 未正确注入时降级 | 默认值 "unknown" 保证启动不失败 |

## 测试策略

### Contract 测试
- 验证 API 签名: GET /version 返回 200
- 验证响应结构: 包含 `sha` 字段
- 验证响应格式: 有效的 JSON

### Acceptance 测试 (E2E)
- 部署后调用 /version，验证返回有效的 git SHA
- 验证 SHA 与当前部署代码一致
- 验证多次调用结果一致（幂等性）

## 交付范围
- ✅ 实现 `/version` 端点和版本注入机制
- ✅ 更新 Makefile 和 CI 流程以支持版本注入
- ✅ Contract 和 Acceptance 测试应该已定义
- ❌ 文档更新（超出本 REQ 范围）
- ❌ Client SDK 更新（超出本 REQ 范围）
