---
change_id: REQ-924
title: "GET /echo endpoint — design"
---

## Context

Admin HTTP 服务已有 `/ping`（固定返回 "pong"）、`/healthz`（JSON 健康检查）、`/version` 和 `/buildinfo`（构建信息）。需求是增加一个可回显任意输入的端点。

## Goals

1. 通过 `GET /echo?msg=xxx` 回显调用方传入的字符串
2. 遵循现有项目结构惯例（每个端点一个包）

## Decision

**新建 `echo/` 包**，实现 `Handler` 函数：

```go
func Handler(w http.ResponseWriter, r *http.Request) {
    msg := r.URL.Query().Get("msg")
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(msg))
}
```

路由注册：

```go
mux.HandleFunc("GET /echo", echo.Handler)
```

## Behavior

| 场景 | 请求 | 响应 |
|------|------|------|
| 正常 | `GET /echo?msg=hello` | 200, body: `hello` |
| 空 msg | `GET /echo?msg=` | 200, body: (空) |
| 缺少 msg | `GET /echo` | 200, body: (空) |
| 错误方法 | `POST /echo` | 405 Method Not Allowed |

## Risks / Tradeoffs

| Risk | Mitigation |
|------|-----------|
| msg 可能包含超长字符串 | Admin 端口仅内网暴露，Go stdlib 已有默认请求大小限制 |
| 回显可能被用于反射攻击 | Admin 端口不对外，Content-Type 为 text/plain 不会被浏览器执行 |

## Dependencies

- 无新依赖，仅使用 `net/http` 标准库
