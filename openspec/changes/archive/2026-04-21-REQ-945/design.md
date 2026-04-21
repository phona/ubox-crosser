---
change_id: REQ-945
title: "POST /webhook-debug endpoint — design"
---

## Context

Admin HTTP 服务已有 `/ping`（固定返回 "pong"）、`/echo`（回显 query param）、`/healthz`（JSON 健康检查）、`/version` 和 `/buildinfo`（构建信息）。这些端点都只接受 GET 请求，无法用于调试 webhook 回调（通常是 POST + JSON body）。

## Goals

1. 提供 `/webhook-debug` 端点，接收任意 HTTP 方法的请求
2. 以 JSON 格式返回完整请求详情，便于开发者检查 webhook 回调内容
3. 遵循现有项目结构惯例（每个端点一个包）

## Decision

**新建 `webhookdebug/` 包**，实现 `Handler` 函数：

```go
type RequestInfo struct {
    Method  string              `json:"method"`
    Path    string              `json:"path"`
    Query   map[string][]string `json:"query"`
    Headers map[string][]string `json:"headers"`
    Body    string              `json:"body"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
    body, _ := io.ReadAll(r.Body)
    info := RequestInfo{
        Method:  r.Method,
        Path:    r.URL.Path,
        Query:   r.URL.Query(),
        Headers: r.Header,
        Body:    string(body),
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(info)
}
```

路由注册（不限定 HTTP 方法）：

```go
mux.HandleFunc("/webhook-debug", webhookdebug.Handler)
```

## Behavior

| 场景 | 请求 | 响应 |
|------|------|------|
| POST + JSON body | `POST /webhook-debug` body: `{"event":"push"}` | 200, JSON 含 method/path/headers/body |
| GET + query params | `GET /webhook-debug?foo=bar` | 200, JSON 含 query params |
| PUT + form body | `PUT /webhook-debug` body: `key=val` | 200, JSON 含 body 原文 |
| 无 body 请求 | `DELETE /webhook-debug` | 200, JSON body 字段为空字符串 |

## Response Format

```json
{
  "method": "POST",
  "path": "/webhook-debug",
  "query": {},
  "headers": {
    "Content-Type": ["application/json"],
    "X-Webhook-Secret": ["abc123"]
  },
  "body": "{\"event\":\"push\"}"
}
```

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| body 可能很大，ReadAll 会消耗内存 | Admin 端口仅内网暴露，可考虑后续加 body 大小限制（如 1MB），MVP 阶段不需要 |
| 返回 headers 可能泄露敏感信息 | Admin 端口不对外暴露，调试用途可接受 |
| 不限 HTTP 方法，可能影响路由优先级 | Go 1.22+ ServeMux 中无方法前缀的路由优先级最低，不影响其他端点 |

## Dependencies

- 无新依赖，仅使用 `net/http`、`encoding/json`、`io` 标准库
