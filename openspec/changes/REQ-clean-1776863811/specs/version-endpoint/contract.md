# Contract Spec: /version 端点 HTTP API

## 概述

本文档定义了 `/version` 端点的 HTTP API 契约。契约测试验证 API 的基础签名和响应格式，与具体的部署环境无关。

---

## API 端点定义

### GET /version

**描述**: 返回当前部署代码的版本信息，包括 git commit SHA。

**URL**: `GET /version`

**响应状态码**:
- `200 OK` — 请求成功，返回版本信息
- `405 Method Not Allowed` — 使用了不支持的 HTTP 方法（如 POST、DELETE）

**响应头**:
- `Content-Type: application/json` — 响应内容为 JSON 格式

**响应体** (JSON):
```json
{
  "sha": "abc123def456...",
  "version": "1.0.0",
  "module": "github.com/example/ubox-crosser",
  "go_os": "linux",
  "go_arch": "amd64",
  "commit": "abc123def456..."
}
```

**响应字段说明**:

| 字段 | 类型 | 必需 | 描述 |
|------|------|------|------|
| `sha` | string | ✅ | Git commit SHA（完整 40 字符 hex）或 "unknown" |
| `version` | string | ✅ | 应用版本号 |
| `module` | string | ✅ | Go module 名称 |
| `go_os` | string | ✅ | 编译目标操作系统 |
| `go_arch` | string | ✅ | 编译目标架构 |
| `commit` | string | ✅ | Git commit SHA（同 `sha` 字段） |

---

## 契约场景 (Contract Scenarios)

### C1: 基础 API 签名 (REQ-clean-1776863811-S1)

**验证要点**:
- GET 请求返回 HTTP 200
- 响应 Content-Type 为 application/json
- 响应体是有效的 JSON

**执行步骤**:
```
GET /version
Accept: application/json
```

**预期结果**:
```
HTTP/1.1 200 OK
Content-Type: application/json

{ ... valid JSON ... }
```

---

### C2: 响应结构合规性 (REQ-clean-1776863811-S2)

**验证要点**:
- 响应包含所有必需字段（sha, version, module, go_os, go_arch, commit）
- 字段值类型正确（都是字符串）
- 字段值非空

**执行步骤**:
1. 发送 GET /version
2. 解析响应 JSON
3. 验证每个必需字段存在且非空

**预期结果**:
```json
{
  "sha": "non-empty string",
  "version": "non-empty string",
  "module": "non-empty string",
  "go_os": "non-empty string",
  "go_arch": "non-empty string",
  "commit": "non-empty string"
}
```

---

### C3: Git SHA 格式验证 (REQ-clean-1776863811-S3)

**验证要点**:
- `sha` 字段为 40 字符的十六进制字符串，或 "unknown"
- `commit` 字段与 `sha` 字段相同

**有效格式**:
- `abc123def456abc123def456abc123def456abc1` (40 chars, hex)
- `"unknown"` (默认值，当无法获取 SHA 时)

**执行步骤**:
1. 发送 GET /version
2. 检查 `sha` 和 `commit` 字段
3. 验证格式: `^([0-9a-f]{40}|unknown)$`

---

### C4: HTTP 方法支持 (REQ-clean-1776863811-S4)

**验证要点**:
- GET 返回 200 OK
- POST 被拒绝（405 or 404）
- DELETE 被拒绝（405 or 404）
- PUT 被拒绝（405 or 404）
- HEAD 支持（可选）

**执行步骤**:
1. `GET /version` → 应返回 200
2. `POST /version` → 应返回 405 或 404
3. `DELETE /version` → 应返回 405 或 404
4. `PUT /version` → 应返回 405 或 404

---

### C5: 查询参数处理 (REQ-clean-1776863811-S5)

**验证要点**:
- GET /version?foo=bar 返回 200（查询参数被忽略）
- 响应与无参数请求相同

**执行步骤**:
1. `GET /version` → 获取响应 A
2. `GET /version?foo=bar&baz=qux` → 获取响应 B
3. 验证响应 A 和 B 完全相同

---

### C6: 请求头容错 (REQ-clean-1776863811-S6)

**验证要点**:
- 自定义请求头不影响响应
- 缺少常见请求头（Accept, User-Agent）不影响响应

**执行步骤**:
1. `GET /version` (无额外请求头) → 200
2. `GET /version` 带 `Accept: text/plain` → 仍返回 application/json
3. `GET /version` 带自定义请求头 `X-Custom: value` → 200

---

## 契约违反检测

### 不符合契约的响应示例

❌ **缺少必需字段**:
```json
{
  "sha": "abc123..."
  // 缺少 version, module, go_os, go_arch, commit
}
```

❌ **字段类型错误**:
```json
{
  "sha": "abc123...",
  "version": 1.0,  // 应为字符串，不是数字
  ...
}
```

❌ **SHA 格式错误**:
```json
{
  "sha": "abc123"  // 只有 6 字符，应为 40 或 "unknown"
}
```

❌ **错误的 Content-Type**:
```
HTTP/1.1 200 OK
Content-Type: text/plain

{"sha": "..."}  // 应为 application/json
```

---

## 技术要求

### 实现约束
- 端点必须无状态（幂等）
- 不依赖外部服务或数据库
- 响应时间应 < 100ms

### 扩展性
- 未来可能添加新字段（`build_time` 等）
- 客户端应忽略未知字段

---

## 变更历史

- **2026-04-22**: 初始契约定义
