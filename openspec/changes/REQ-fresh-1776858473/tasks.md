# Tasks: REQ-fresh-1776858473 - /version Endpoint Git SHA

## Stage: contract-tests

### Contract Specification
- [x] Define `/version` endpoint contract with git_sha field
- [x] Document response schema with all 5 fields (version, module, go_os, go_arch, git_sha)
- [x] Specify git SHA format requirements (hexadecimal, 7-40 characters)
- [x] Create detailed scenarios (S1-S8) with acceptance criteria
- [x] Document edge cases and response examples
- [x] Define OpenAPI specification for the endpoint

### Contract Test Implementation
- [x] Write integration tests in `tests/contract/version_test.go`
- [x] Test git_sha presence and validity
- [x] Test git_sha format validation (hexadecimal regex)
- [x] Test git_sha consistency across requests
- [x] Test all required fields present
- [x] Test HTTP method restrictions (405 for POST/PUT/DELETE)
- [x] Test response time requirements (< 100ms)
- [x] Test Content-Type validation

### Spec Files Created
- [x] `openspec/changes/REQ-fresh-1776858473/proposal.md` - Problem statement and solution overview
- [x] `openspec/changes/REQ-fresh-1776858473/design.md` - Implementation architecture and design
- [x] `openspec/changes/REQ-fresh-1776858473/specs/version-endpoint-git-sha/spec.md` - Full contract specification
- [x] `tests/contract/version_test.go` - Contract test suite

### Validation
- [ ] Run `openspec validate openspec/changes/REQ-fresh-1776858473` (when tool available)
- [ ] Verify all scenario names follow pattern `REQ-fresh-1776858473-S*`
- [ ] Verify no breaking changes to existing endpoints

## Stage: dev

Implementation tasks:
- [x] Update `VersionInfo` struct to include `GitSha` field
- [x] Enhance `handleVersion()` to extract git SHA from build metadata
- [x] Implement git SHA extraction logic in `server/management.go`
- [ ] Run contract tests to verify RED→GREEN transition
- [ ] Ensure all tests pass before marking stage complete

## Notes

**Current Status:** Contract spec and tests complete. Ready for dev implementation.

**Blocked By:** None - this is the contract-spec stage and is self-contained.

**Related:** 
- REQ-haiku-1776846225: Original /version endpoint implementation
- REQ-test-1776855620: Similar contract-spec pattern for /healthz endpoint
