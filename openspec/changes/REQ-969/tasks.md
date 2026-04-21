---
change_id: REQ-969
title: "version endpoint — tasks"
---

## Stage: contract-tests (owner: contract-spec-agent)
- [x] Define OpenAPI spec for GET /api/version (200 VersionResponse schema) in contract.spec.yaml
- [x] Contract test — GET /api/version returns 200 with JSON {"commit":"<string>"} (REQ-969-S1)
- [x] Contract test — POST/PUT/DELETE /api/version returns 405 Method Not Allowed (REQ-969-S3)
- [x] Contract test — response Content-Type is application/json, schema validation (REQ-969-S2)

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] FEATURE-A1: GET /api/version 返回 200 + JSON `{"commit":"<hash>"}` (`specs/version-endpoint/spec.md::FEATURE-A1`)
- [ ] FEATURE-A2: commit 字段为合法 40 字符 hex git hash 或 "unknown" (`specs/version-endpoint/spec.md::FEATURE-A2`)
- [ ] FEATURE-A3: 非 GET 方法对 /api/version 返回 405 (`specs/version-endpoint/spec.md::FEATURE-A3`)
- [ ] FEATURE-A4: docker-compose 网络内 test-runner 可达 proxy-server:8080/api/version (`specs/version-endpoint/spec.md::FEATURE-A4`)
- [ ] FEATURE-A5: /api/version 在正常条件下 500ms 内响应 (`specs/version-endpoint/spec.md::FEATURE-A5`)

## Stage: implementation (owner: dev-agent)
- [x] version/handler.go — package-level `var Commit string` + `Handler` (GET-only → 200 JSON `{"commit":"..."}`, non-GET → 405, empty Commit defaults to `"unknown"`)
- [x] server/admin.go — import `version` package, register `mux.HandleFunc("/api/version", version.Handler)` in `NewAdminMux()`
- [x] Makefile — add `COMMIT` and `LDFLAGS` variables with `-X $(MODULE)/version.Commit=$(COMMIT)`, use in all `go build` commands
- [x] version/handler_test.go — unit tests: GET returns 200 + correct commit, empty commit → "unknown", schema validation (single `commit` field), POST/PUT/DELETE → 405
