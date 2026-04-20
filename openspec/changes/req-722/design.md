---
change_id: req-722
title: "GET /whoami endpoint — Design"
---

## Context

ubox-crosser 的 admin HTTP listener（`--admin-addr`，默认 `:8080`）已经提供了 `GET /version`（构建元数据）和 `GET /ping`（连通性检查）。在多实例部署中，运维需要知道请求到达了哪个节点，尤其是通过负载均衡器或 tunnel 访问时。

现有端点遵循 package-per-handler 模式：每个端点一个独立包，导出 `Handler` 函数，在 `cmd/server/server.go` 中注册到 admin mux。

## Goals / Non-Goals

**Goals:**
- 在现有 admin HTTP listener 上暴露 `GET /whoami`，返回纯文本主机名
- 遵循已有的 package-per-handler 模式（`ping/`, `version/`）
- 仅接受 GET，其他方法返回 405

**Non-Goals:**
- 返回 IP 地址或其他系统信息（超出范围）
- 认证/鉴权（与 `/ping`, `/version` 一致，依赖 `--admin-addr` 网络层控制）
- 主机名缓存（`os.Hostname()` 开销极小，无需缓存）

## Decisions

### Decision 1: 纯文本响应（非 JSON）

主机名是单个字符串值，JSON 包装没有实际价值。纯文本与 `/ping` 保持一致，便于 `curl` 和脚本直接使用。

**选择:** 返回 `text/plain; charset=utf-8`，使用 `io.WriteString`。

| Option | Pros | Cons |
|--------|------|------|
| 纯文本主机名 | 简洁，与 `/ping` 一致，脚本友好 | 与 `/version` 的 JSON 不一致 |
| JSON `{"hostname":"..."}` | 与 `/version` 一致 | 对单值场景过度设计 |

### Decision 2: 直接调用 `os.Hostname()`

每次请求调用 `os.Hostname()`。该函数在 Linux 上是 `uname` syscall，开销极小（微秒级）。如果返回错误，回退到 `"unknown"`。

### Decision 3: 独立 `whoami` 包

沿用 `ping/`, `version/` 的模式 — 独立包内 `handler.go` + `handler_test.go`，保持端点解耦和可独立测试。

## Risks / Trade-offs

- **[主机名泄露]** → 在受信任的 admin 网络中可接受，与 `/version` 暴露构建信息风险级别一致。通过 `--admin-addr` 绑定控制访问。
- **[os.Hostname 失败]** → 极端罕见（容器环境下主机名总是可用），回退到 `"unknown"` 确保不会 500。
