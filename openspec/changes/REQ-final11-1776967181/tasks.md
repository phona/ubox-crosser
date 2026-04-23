# REQ-final11-1776967181: /buildinfo Endpoint Tasks

## Stage: contract / spec
- [x] author `specs/buildinfo/contract.spec.yaml` — OpenAPI 3.0 schema for /buildinfo response
- [x] author `specs/buildinfo/spec.md` — 7 scenarios in openspec delta format (UBOX-S1..S7)
- [x] author `proposal.md` — design rationale and file inventory

## Stage: implementation
- [x] `models/config/config.go` — add `HTTPAddr` field to `ServerConfig`
- [x] `server/http.go` — HTTP management server with `/buildinfo` + `/healthz` handlers
- [x] `cmd/server/server.go` — add `GitSHA` ldflags var; start HTTP server at boot
- [x] `Makefile` — inject `GIT_SHA` ldflags in `build` target
- [x] `Dockerfile` — pass `GIT_SHA` build-arg through ldflags

## Stage: unit tests
- [x] `server/http_test.go` — 6 unit tests covering: default BUILD_ID, custom BUILD_ID, Content-Type, empty GitSHA, healthz 200, healthz Content-Type

## Stage: PR
- [x] git push `feat/REQ-final11-1776967181`
- [x] gh pr create
