# REQ-final15-1776989948 Tasks

## Stage: contract / spec

- [x] author specs/buildinfo-endpoint/contract.spec.yaml with OpenAPI 3.0 schema
- [x] author specs/buildinfo-endpoint/spec.md with 7 UBOX-S scenarios in openspec delta format (SHALL/MUST in body paragraph)

## Stage: implementation

- [x] add server/management.go — ManagementServer with /healthz and /buildinfo handlers
- [x] update cmd/server/server.go — var GitSHA, --management-addr flag, start management server goroutine
- [x] update Makefile — add GIT_SHA ldflags in build target, ci-test and dev-cross-check targets
- [x] write unit tests in server/management_test.go (pure, no build tag, uses httptest)
- [x] add acceptance tests in tests/acceptance/buildinfo_test.go (//go:build acceptance)

## Stage: PR

- [x] git push feat/REQ-final15-1776989948
- [x] gh pr create
