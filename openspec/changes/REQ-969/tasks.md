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
- [ ] TODO: Acceptance test — GET /api/version returns 200 with valid commit hash format
- [ ] TODO: Acceptance test — non-GET method returns 405
- [ ] TODO: Acceptance test — response JSON schema validation (commit field present, string type)

## Stage: implementation (owner: dev-agent)
- [ ] TODO: version/handler.go — var Commit + Handler function (GET-only, return JSON with commit)
- [ ] TODO: server/admin.go — register /api/version route in NewAdminMux()
- [ ] TODO: Makefile — add LDFLAGS with -X version.Commit=$(git rev-parse HEAD)
- [ ] TODO: version/handler_test.go — unit tests (GET 200, non-GET 405, empty commit defaults to "unknown")
