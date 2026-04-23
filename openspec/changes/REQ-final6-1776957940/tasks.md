# REQ-final6-1776957940: /buildinfo Endpoint

## Stage: contract / spec

- [x] author specs/buildinfo-endpoint/contract.spec.yaml (OpenAPI schema)
- [x] author specs/buildinfo-endpoint/spec.md (7 acceptance scenarios [UBOX-S1..S7])

## Stage: implementation

- [x] server/http_server.go — HTTP server with /healthz + /buildinfo handlers
- [x] cmd/server/server.go — add GitSHA ldflag var; launch HTTP server goroutine
- [x] Dockerfile — accept GIT_SHA build arg; pass to ldflag
- [x] Makefile — GIT_SHA var; updated build + ci-build targets; ci-test target
- [x] server/http_server_test.go — unit tests for buildinfo + healthz handlers
- [x] tests/acceptance/buildinfo_test.go — acceptance tests via docker-compose

## Stage: PR

- [x] git push feat/REQ-final6-1776957940
- [x] gh pr create
