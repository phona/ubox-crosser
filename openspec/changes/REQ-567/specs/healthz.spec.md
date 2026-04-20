# Acceptance Spec: GET /healthz

## Scenario: Happy Path - 健康检查返回正常状态

**Given** proxy-server 已启动且健康检查 HTTP 服务正在监听
**When** 发送 `GET /healthz` 请求
**Then** 返回 HTTP 200
**And** Content-Type 为 `application/json`
**And** Body 包含 `{"status":"ok","uptime_seconds":<N>,"services":[...]}`
**And** `uptime_seconds` >= 0
**And** `services` 数组包含配置中的服务名

## Scenario: 非法方法 - POST 返回 405

**Given** 健康检查服务正在运行
**When** 发送 `POST /healthz` 请求
**Then** 返回 HTTP 405 Method Not Allowed

## Scenario: 未知路径 - 返回 404

**Given** 健康检查服务正在运行
**When** 发送 `GET /nonexistent` 请求
**Then** 返回 HTTP 404 Not Found

## Scenario: Uptime 随时间递增

**Given** 健康检查服务正在运行
**When** 连续发送两次 `GET /healthz`，间隔 1 秒
**Then** 第二次 `uptime_seconds` >= 第一次

## Scenario: services 字段非 null

**Given** proxy-server 启动时配置了至少一个服务
**When** 发送 `GET /healthz`
**Then** `services` 是非 null 的 JSON 数组
