# REQ-final5-1776956833 Tasks

## Stage: Contract & Spec

- [x] Author specs/buildinfo/contract.spec.yaml with API schema
- [x] Author specs/buildinfo/spec.md with 8 acceptance scenarios
- [x] Write proposal.md with motivation and approach

## Stage: Implementation

- [ ] Create buildinfo.go handler with JSON response struct
- [ ] Add HTTP server to ProxyServer that listens on same address as proxy
- [ ] Implement GitSHA variable injection via ldflags in Makefile
- [ ] Implement BUILD_ID environment variable reading (default "dev")
- [ ] Implement go_version hardcoded value ("go1.23")
- [ ] Wire handler to HTTP router
- [ ] Write unit tests in cmd/server/buildinfo_test.go
- [ ] Write integration test in tests/acceptance/buildinfo_test.go using docker-compose

## Stage: Testing & Validation

- [ ] Run `make ci-test` - all tests pass
- [ ] Test locally: `curl http://localhost:7000/buildinfo` returns valid JSON
- [ ] Verify git_sha matches HEAD commit (first 7 chars)
- [ ] Verify build_id defaults to "dev" without BUILD_ID env var
- [ ] Verify build_id respects BUILD_ID env var when set
- [ ] Verify go_version is always "go1.23"
- [ ] Verify Content-Type is application/json

## Stage: PR & Delivery

- [ ] Push feat/REQ-final5-1776956833 branch
- [ ] Open PR to master with full description
- [ ] All CI checks pass
- [ ] Ready for review and merge

