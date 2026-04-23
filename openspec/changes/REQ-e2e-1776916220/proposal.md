# 需求：添加 /healthz 端点返回服务 uptime

## 一句话需求
为 ubox-crosser 代理服务添加 `/healthz` HTTP 端点，返回服务启动以来的运行时长（uptime）。

## 影响范围
- **单个 source repo**：phona/ubox-crosser
- **无 integration repo**：纯代码改动，不涉及外部依赖或集成环境

## 主要改动点
1. **server/server.go**
   - 在 `ProxyServer` 结构体中添加 `startTime` 字段记录启动时间
   - 在 `handleHTTPRequest` 中添加 `/healthz` 请求处理，返回 JSON `{"uptime": <seconds>}`

2. **cmd/main.go 或初始化代码**
   - 在创建 `ProxyServer` 时记录启动时间

3. **测试**
   - 集成测试验证 `/healthz` 端点返回正确的 uptime（秒级精度）

## 支持性判定
✅ **可行** — 与现有 `/version` 端点实现模式一致，改动范围明确，无技术阻碍
