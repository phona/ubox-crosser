---
change_id: REQ-736
title: "GET /uptime design"
---

## High-Level Design

### 方案

在 `uptime` 包内维护一个 package-level 的启动时间变量，由 `Init()` 在服务启动时设置为 `time.Now()`。`Handler` 计算 `time.Since(startTime).Seconds()` 并以 JSON 返回整数秒数。

### 选型理由

- **JSON 响应**：与 `/version` 保持一致，方便机器解析
- **整数秒**：运维场景不需要亚秒精度，整数更直观
- **Init() 模式**：让调用方控制 startTime 的设置时机，避免 `init()` 隐式行为，也便于测试时注入不同的启动时间

### 依赖

- 仅 Go stdlib（`time`, `encoding/json`, `net/http`）

### 风险

- 无。这是一个只读、无副作用的纯信息端点。
