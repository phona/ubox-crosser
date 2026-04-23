# 设计文档：/healthz 端点实现

## 高层设计

### 1. 数据模型变更
在 `ProxyServer` 结构体中添加启动时间记录：
```go
type ProxyServer struct {
    // ... existing fields ...
    startTime  time.Time  // 服务启动时间
    gitSHA     string     // 现有字段
}
```

### 2. 端点实现
遵循现有 `/version` 端点的实现模式：
- **路径**：`GET /healthz`
- **返回格式**：JSON `{"uptime": <seconds>}`
- **HTTP 状态**：200 OK
- **Content-Type**：application/json

### 3. 实现位置
在 `server/server.go` 的 `handleHTTPRequest` 方法中：
1. 检测 HTTP 请求是否为 `/healthz`（添加 `isHealthzRequest` 方法）
2. 计算运行时长 = 当前时间 - 启动时间，转换为秒数
3. 组织 JSON 响应并返回

### 4. 初始化时间
在 `NewProxyServer` 函数中记录 `startTime = time.Now()`

## 技术选型

### Time 库选择
- **方案**：使用标准库 `time.Time` 和 `time.Since()`
- **理由**：Go 标准库充分，无需外部依赖；精度满足需求（纳秒级）

### JSON 响应
- 继续使用 `fmt.Sprintf` 组织 HTTP 响应（保持与 `/version` 一致）
- `{"uptime": 3600}` 表示运行 3600 秒

## 风险与缓解

| 风险 | 影响 | 缓解 |
|------|------|------|
| 服务重启时 uptime 重置 | 监控告警可能误判 | 文档说明，使用时配合服务启动日志 |
| 时间精度溢出 | 理论上运行 >= 292 年时 int64 溢出 | 无需处理，实践中服务不会运行这么久 |
| 并发读写 startTime | 极低概率竞态 | startTime 初始化后只读，无需锁 |

## 与现有功能的兼容性
✅ 完全兼容
- 不修改现有代理逻辑
- `/version` 端点保持不变
- 新增端点不影响性能（仅简单计算）
