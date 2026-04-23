# REQ-final14-1776975275 Tasks

## Stage: contract / spec

- [x] author specs/buildinfo-endpoint/contract.spec.yaml with OpenAPI 3.0 schema
- [x] author specs/buildinfo-endpoint/spec.md with 7 UBOX-S scenarios in openspec delta format (SHALL/MUST in body paragraph)

## Stage: implementation

- [x] add server/management.go — ManagementServer with /healthz and /buildinfo handlers
- [x] update cmd/server/server.go — var GitSHA, --management-addr flag, start management server goroutine
- [x] update Makefile — add ci-test and dev-cross-check targets
- [x] fix tests/acceptance/healthz_test.go — add //go:build acceptance tag to prevent go test ./... failures
- [x] write unit tests in server/management_test.go (pure, no build tag, uses httptest)

## Stage: PR

- [x] git push feat/REQ-final14-1776975275
- [x] gh pr create
