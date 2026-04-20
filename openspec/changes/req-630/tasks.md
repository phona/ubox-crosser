## Stage 0 — Spec & Contract Lock (accept-spec)

- [x] 审查 proposal.md、design.md、specs/version-endpoint/spec.md、contract.spec.yaml 一致性
- [x] 修复：contract.spec.yaml 补充 405 响应（POST/PUT/DELETE/PATCH）+ Allow: GET header
- [x] 修复：contract.spec.yaml 补充 catch-all 404 路径
- [x] 修复：contract.spec.yaml 升级至 OpenAPI 3.1.0（与 REQ-601 一致）
- [x] 修复：spec.md FEATURE-S6 补充 Allow: GET header 要求
- [x] 修复：spec.md 新增 FEATURE-S8/S9/S10（unknown path 404 + trailing slash 404）
- [x] 修复：spec.md FEATURE-S7 重编号为 S11（避免与新增 S8-S10 冲突）
- [x] BKD issue 确认 `layer:backend` tag
- [x] BKD issue move → `review`

### 审查备注

> **与 REQ-601 HTTP listener 协调**：REQ-601 design 已规划独立 health HTTP listener（`--health-address`），REQ-630 规划独立 `--http-addr` listener。两者默认端口均为 `:8080`。实现阶段应将 `/version` 注册到 REQ-601 的 health mux 上共享同一 HTTP listener，而非创建第二个独立 listener。此决策留待 dev-spec 阶段确认。

## 1. Version Package

- [ ] 1.1 Create `internal/version/version.go` with `Version`, `Commit`, `BuildTime` vars (defaults: `"0.1.0"`, `"unknown"`, `"unknown"`)
- [ ] 1.2 Create `internal/version/handler.go` with `Handler()` returning an `http.HandlerFunc` that writes JSON `{"version","commit","build_time"}`

## 2. HTTP Server Integration

- [ ] 2.1 Add `--http-addr` flag (default `:8080`) to `cmd/server/server.go`（或复用 REQ-601 的 `--health-address`，见 Stage 0 备注）
- [ ] 2.2 Register `GET /version` on `http.ServeMux` and start HTTP listener in a goroutine before the proxy loop

## 3. Build Pipeline

- [ ] 3.1 Update `Makefile` build target to inject `Version`, `Commit`, `BuildTime` via `-ldflags -X`
- [ ] 3.2 Verify `make build && ./bin/server --help` shows no regression

## 4. Unit Tests

- [ ] 4.1 Create `internal/version/handler_test.go` — test HTTP 200, Content-Type, JSON body fields (FEATURE-S1, S2, S5, S11)
- [ ] 4.2 Add test for POST → 405 + Allow header (FEATURE-S6)
- [ ] 4.3 Add tests for unknown path 404 and trailing slash 404 (FEATURE-S8, S9, S10)

## 5. Contract Tests

- [ ] 5.1 Create contract test — validate response against `contract.spec.yaml` schema
- [ ] 5.2 `go test ./...` 全绿
- [ ] 5.3 `golangci-lint run ./...` 无新增告警

## Stage: Contract Test (owner: contract-test-agent)

- [x] [FEATURE-S1] TestVersion_S1_Returns200WithJSON — verify HTTP 200, Content-Type application/json, exactly 3 required string fields (version, commit, build_time), no extra fields
- [x] [FEATURE-S5] TestVersion_S5_NoAuthRequired — verify unauthenticated GET /version returns 200
- [x] [FEATURE-S6] TestVersion_S6_PostReturns405 — verify POST /version returns 405 + Allow: GET header
- [x] [FEATURE-S6b] TestVersion_S6b_PutReturns405 — verify PUT /version returns 405 + Allow: GET header
- [x] [FEATURE-S8] TestVersion_S8_RootPathReturns404 — verify GET / returns 404
- [x] [FEATURE-S9] TestVersion_S9_UnknownPathReturns404 — verify GET /metrics returns 404
- [x] [FEATURE-S10] TestVersion_S10_TrailingSlashReturns404 — verify GET /version/ returns 404
- [x] [FEATURE-S11] TestVersion_S11_DefaultFieldsPresent — verify all fields present and non-empty in default build
