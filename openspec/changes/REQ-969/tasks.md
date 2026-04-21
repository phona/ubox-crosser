---
change_id: REQ-969
title: "version endpoint — tasks"
---

## Stage: contract-tests (owner: contract-spec-agent)
- [ ] TODO: Define OpenAPI spec for GET /api/version (200 response schema, 405 for non-GET)
- [ ] TODO: Contract test — GET /api/version returns 200 with JSON {"commit":"<string>"}
- [ ] TODO: Contract test — POST /api/version returns 405 Method Not Allowed
- [ ] TODO: Contract test — response Content-Type is application/json

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [ ] FEATURE-A1: GET /api/version 返回 200 + JSON `{"commit":"<hash>"}` (`specs/version-endpoint/spec.md::FEATURE-A1`)
- [ ] FEATURE-A2: commit 字段为合法 40 字符 hex git hash 或 "unknown" (`specs/version-endpoint/spec.md::FEATURE-A2`)
- [ ] FEATURE-A3: 非 GET 方法对 /api/version 返回 405 (`specs/version-endpoint/spec.md::FEATURE-A3`)
- [ ] FEATURE-A4: docker-compose 网络内 test-runner 可达 proxy-server:8080/api/version (`specs/version-endpoint/spec.md::FEATURE-A4`)
- [ ] FEATURE-A5: /api/version 在正常条件下 500ms 内响应 (`specs/version-endpoint/spec.md::FEATURE-A5`)

## Stage: implementation (owner: dev-agent)
- [ ] TODO: version/handler.go — var Commit + Handler function (GET-only, return JSON with commit)
- [ ] TODO: server/admin.go — register /api/version route in NewAdminMux()
- [ ] TODO: Makefile — add LDFLAGS with -X version.Commit=$(git rev-parse HEAD)
- [ ] TODO: version/handler_test.go — unit tests (GET 200, non-GET 405, empty commit defaults to "unknown")
