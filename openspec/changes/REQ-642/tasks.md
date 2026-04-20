---
change_id: REQ-642
title: "Tasks: GET /version endpoint (v2)"
---

# Tasks

## Stage 0 — Spec & Contract Lock (accept-spec)

- [x] 审查 proposal.md、design.md、specs/FEATURE-S1.md、contract.spec.yaml 一致性
- [x] contract.spec.yaml 定义 GET /version 200 响应 + 405 响应
- [x] FEATURE-S1 覆盖：200 + JSON body 三字段 + Non-GET 405
- [x] BKD analyze issue (#642) 确认 `layer:backend` tag

### 审查备注

> **与 REQ-601 的关系**：REQ-601 规划了独立 health HTTP listener（`--health-address`），但 REQ-601 尚未合入 health endpoint 实现。REQ-642 作为独立特性先行实现 `--http-addr` listener，仅挂载 `/version`。若后续 REQ-601 合入，需协调将 `/version` 迁移到共享 mux，此为后续 REQ 范围，不影响 REQ-642 交付。

> **contract.spec.yaml 精简**：当前 contract 仅覆盖 200 和 405，未定义 404（unknown path）和 trailing-slash 行为。REQ-642 scope 限定为单一 `/version` 端点，路由兜底行为由 `net/http` 默认 mux 决定，不纳入 contract 约束。

## Stage 1 — Dev-Spec Decisions

- [x] **确认包路径**：`version/`（repo root level），非 `internal/version`。理由：design.md 明确 "repo root level"，且 version 信息可被其他 cmd 二进制复用。
- [x] **确认 HTTP listener 方式**：独立 `--http-addr` flag（默认 `:8080`），在 `cmd/server` 中创建新 `http.ServeMux`。理由：REQ-601 health mux 尚未合入，不阻塞本次交付。
- [x] **确认 JSON 序列化**：`encoding/json.NewEncoder(w).Encode()`，配合 `Info` struct tag 保证字段名与 contract 一致。
- [x] **确认 method 检查方式**：handler 内 `r.Method != http.MethodGet` 直接返回 405，不使用中间件。理由：单端点无需路由框架。
- [x] **确认 ldflags 注入目标**：`$(MODULE)/version.Commit` 和 `$(MODULE)/version.BuildTime`，Version 常量保持硬编码 `0.1.0`。
- [x] **确认无第三方依赖**：handler 测试用 `net/http/httptest`，不引入断言库。

## Stage 2 — Backend Dev

- [x] 新增 `version/version.go`：声明 `const Version = "0.1.0"` + `var (Commit = "unknown"; BuildTime = "unknown")`
- [x] 新增 `version/handler.go`：实现 `Handler(w, r)` 函数，GET 返回 200 + JSON，非 GET 返回 405
- [x] 修改 `cmd/server/server.go`：添加 `--http-addr` flag，启动 goroutine 运行 HTTP server，注册 `/version` handler
- [x] 更新 `Makefile`：添加 `GIT_COMMIT`、`BUILD_TIME`、`VERSION_PKG`、`LDFLAGS` 变量，build target 加 `-ldflags`
- [x] 更新 `Dockerfile`：build stage 传入 ldflags 参数
- [x] 本地编译通过：`make build` / `go build ./...`

## Stage 3 — Unit Test

- [x] 新增 `version/handler_test.go`：
  - `TestHandler_GET_ReturnsVersionJSON`：验证 GET 返回 200 + Content-Type application/json + 三字段 JSON body
  - `TestHandler_POST_Returns405`：验证 POST 返回 405
  - `TestHandler_VersionFieldMatchesConstant`：验证 response 中 version 字段 == `Version` 常量
- [x] `go test ./version/...` 全绿
- [x] `go vet ./...` 无告警

## Stage 4 — Integration Test

- [ ] 新增 `tests/integration/version_test.go`：Docker Compose 启动 server，`curl GET /version` 验证 200 + JSON body
- [ ] 覆盖 FEATURE-S1 完整场景：status code、Content-Type、三字段存在性、Non-GET 405
- [ ] `make test-integration` 在本地 docker compose 中通过

## Stage 5 — Verify

- [ ] CI 流水线（unit + lint）全部 pass
- [ ] 在 vm-node04 上部署，`curl -i http://<addr>:8080/version` 返回 200 + 正确 body
- [ ] `curl -X POST http://<addr>:8080/version` 返回 405
- [ ] 确认 body 中 `version`、`commit`、`build_time` 字段值与构建参数一致

## Stage 6 — Accept

- [ ] 验证日志/截图归档到 BKD issue
- [ ] BKD issue move → `review`（人工确认后 → `done`）
