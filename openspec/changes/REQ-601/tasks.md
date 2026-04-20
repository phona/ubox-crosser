# Tasks — REQ-601 GET /health 健康检查接口

> 多 Stage 骨架。每个 Stage 由 n8n 按 layer 分发到对应执行 agent。
> 复选框留待执行 agent 勾选。

## Stage 1 — Spec & Contract Lock

- [ ] 确认 `proposal.md`、`design.md`、`specs/health-endpoint/spec.md`、`contract.spec.yaml` 内容一致（路径、方法、状态码、响应字段命名）
- [ ] BKD issue 上挂 `layer:backend` tag
- [ ] BKD issue move → `review`，等待 spec 审核

## Stage 2 — Backend Dev

- [ ] 新增 `server/health.go`：实现 `newHealthMux()`，注册 `GET /health`、其他方法 fallback 405、catch-all 404
- [ ] 新增 `(p *ProxyServer) startHealthServer(addr string)`：`http.Server.ListenAndServe()`，过滤 `http.ErrServerClosed`，其他错误写入 `p.errs`
- [ ] `models/config/config.go`：`ServerConfig` 新增 `HealthAddress string \`json:"health_address"\``
- [ ] `server/server.go`：`NewProxyServer` 末尾遍历 configs 取首个非空 `HealthAddress`，非空则 `go server.startHealthServer(addr)`
- [ ] `cmd/server/server.go`：新增 `cmd.Flags().StringVar(&cmdConfig.HealthAddress, "health-address", "", "...")`
- [ ] 本地编译通过：`make build` / `go build ./...`

## Stage 3 — Unit Test

- [ ] 新增 `server/health_test.go`：覆盖 FEATURE-S1 ~ FEATURE-S7（GET 200 / 非 GET 405 + Allow / 未知路径 404 / trailing slash 404），使用 `net/http/httptest`
- [ ] 新增 config / cli flag 解析单测，覆盖 FEATURE-S8、FEATURE-S9
- [ ] `go test ./server/... ./models/config/...` 全绿
- [ ] `golangci-lint run ./...` 无新增告警

## Stage 4 — Integration Test (Contract Lock)

- [ ] 新增 `tests/integration/health_contract_test.go`，按 `contract.spec.yaml` 锁定路径、方法、状态码、header、body 字面量
- [ ] 覆盖 FEATURE-S10（空地址不监听）、FEATURE-S11（端口占用错误上报）
- [ ] `make test-integration` 在本地 docker compose 中通过

## Stage 5 — Verify

- [ ] CI 流水线（unit + integration + lint + sonar）全部 pass
- [ ] 在 vm-node04 上拉镜像/二进制，启动 `--health-address :8080`，`curl -i http://vm-node04:8080/health` 返回 200 + 正确 body
- [ ] `curl -X POST http://vm-node04:8080/health` 返回 405 + `Allow: GET`
- [ ] `curl http://vm-node04:8080/metrics` 返回 404

## Stage 6 — Accept

- [ ] 截图/日志归档到 BKD issue
- [ ] BKD issue move → `review`（人工确认后人工 → `done`）
- [ ] 触发 `openspec archive REQ-601`，把 spec-delta merge 进主 spec
