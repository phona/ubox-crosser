# Tasks — REQ-630 GET /version 版本信息端点

> 多 Stage 骨架。每个 Stage 由 n8n 按 layer 分发到对应执行 agent。
> 复选框留待执行 agent 勾选。

## Stage 1 — Spec & Contract Lock

- [ ] 确认 `proposal.md`、`design.md`、`specs/version-endpoint/spec.md`、`contract.spec.yaml` 内容一致（路径、方法、状态码、响应字段命名）
- [ ] BKD issue 上挂 `layer:backend` tag
- [ ] BKD issue move -> `review`，等待 spec 审核

## Stage 2 — Backend Dev

> **前置依赖**：REQ-601 health endpoint 已合入（`server/health.go` + `newHealthMux()` 存在）。

- [ ] 新增 `internal/version/version.go`：声明 `var (Version = "0.1.0"; Commit = "unknown"; BuildTime = "unknown")`
- [ ] 新增 `internal/version/handler.go`：实现 `Handler() http.HandlerFunc`，用 `json.Marshal` 序列化 `VersionResponse{Version, Commit, BuildTime}` 并写入响应，设置 `Content-Type: application/json`
- [ ] 修改 `server/health.go`：在 `newHealthMux()` 中注册 `GET /version` 到 `version.Handler()`，非 GET 返回 405 + `Allow: GET`
- [ ] 更新 `Makefile` 的 `build` target：ldflags 添加 `-X $(MODULE)/internal/version.Version=$(VERSION) -X $(MODULE)/internal/version.Commit=$(shell git rev-parse --short HEAD) -X $(MODULE)/internal/version.BuildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ)`
- [ ] 本地编译通过：`make build` / `go build ./...`

## Stage 3 — Unit Test

- [ ] 新增 `internal/version/handler_test.go`：覆盖 FEATURE-S1（GET 200 + JSON body 三字段）、FEATURE-S2 ~ S4（ldflags 注入值反映在响应中）、FEATURE-S7（默认值）
- [ ] 新增方法测试：覆盖 FEATURE-S6（POST → 405 + `Allow: GET`）
- [ ] 新增路径测试：`GET /version/` trailing slash → 404
- [ ] `go test ./internal/version/... ./server/...` 全绿
- [ ] `golangci-lint run ./...` 无新增告警

## Stage 4 — Integration Test (Contract Lock)

- [ ] 新增 `tests/integration/version_contract_test.go`，按 `contract.spec.yaml` 锁定路径、方法、状态码、header、body 字段
- [ ] 覆盖 FEATURE-S5（无认证访问成功）
- [ ] 覆盖 health_address 为空时 `/version` 不可访问
- [ ] `make test-integration` 在本地 docker compose 中通过

## Stage 5 — Verify

- [ ] CI 流水线（unit + integration + lint）全部 pass
- [ ] 在 vm-node04 上拉镜像/二进制，启动 `--health-address :8080`，`curl -i http://vm-node04:8080/version` 返回 200 + 正确 body
- [ ] `curl -X POST http://vm-node04:8080/version` 返回 405 + `Allow: GET`
- [ ] 确认 body 中 `version`、`commit`、`build_time` 字段值与构建参数一致

## Stage 6 — Accept

- [ ] 截图/日志归档到 BKD issue
- [ ] BKD issue move -> `review`（人工确认后人工 -> `done`）
- [ ] 触发 `openspec archive REQ-630`，把 spec-delta merge 进主 spec

## Stage: Contract Test (owner: contract-test-agent)

- [x] [FEATURE-S1] TestVersion_S1_Returns200WithJSON — verify HTTP 200, Content-Type application/json, exactly 3 required string fields (version, commit, build_time), no extra fields
- [x] [FEATURE-S5] TestVersion_S5_NoAuthRequired — verify unauthenticated GET /version returns 200
- [x] [FEATURE-S6] TestVersion_S6_PostReturns405 — verify POST /version returns 405 + Allow: GET header
- [x] [FEATURE-S6b] TestVersion_S6b_PutReturns405 — verify PUT /version returns 405 + Allow: GET header
- [x] [FEATURE-S8] TestVersion_S8_RootPathReturns404 — verify GET / returns 404
- [x] [FEATURE-S9] TestVersion_S9_UnknownPathReturns404 — verify GET /metrics returns 404
- [x] [FEATURE-S10] TestVersion_S10_TrailingSlashReturns404 — verify GET /version/ returns 404
- [x] [FEATURE-S11] TestVersion_S11_DefaultFieldsPresent — verify all fields present and non-empty in default build
