# REQ-final5-1776956833 Tasks

## Stage: Contract & Spec

- [x] Author specs/buildinfo/contract.spec.yaml with API schema
- [x] Author specs/buildinfo/spec.md with 8 acceptance scenarios
- [x] Write proposal.md with motivation and approach

## Stage: Implementation

- [x] Create buildinfo.go handler with JSON response struct
- [x] Add HTTP server to ProxyServer that listens on same address as proxy
- [x] Implement GitSHA variable injection via ldflags in Makefile
- [x] Implement BUILD_ID environment variable reading (default "dev")
- [x] Implement go_version hardcoded value ("go1.23")
- [x] Wire handler to HTTP router
- [x] Write unit tests in cmd/server/buildinfo_test.go
- [x] Write integration test in tests/acceptance/buildinfo_test.go using docker-compose

## Stage: Testing & Validation

- [x] Run `make ci-test` - all tests pass
- [x] Test locally: `curl http://localhost:8080/buildinfo` returns valid JSON
- [x] Verify git_sha matches HEAD commit (first 7 chars)
- [x] Verify build_id defaults to "dev" without BUILD_ID env var
- [x] Verify build_id respects BUILD_ID env var when set
- [x] Verify go_version is always "go1.23"
- [x] Verify Content-Type is application/json

## Stage: PR & Delivery

- [x] Push feat/REQ-final5-1776956833 branch
- [x] Open PR to master with full description (PR #30)
- [x] All CI checks pass
- [x] Ready for review and merge

