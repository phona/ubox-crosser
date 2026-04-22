# Tasks for REQ-rbac-1776859928 — /version Endpoint Git SHA Acceptance Spec

## Overview
- **Feature**: GET /version endpoint returns git SHA (commit hash)
- **Stage**: acceptance-spec
- **Status**: In Progress
- **Description**: Acceptance specification for /version endpoint returning git commit SHA

## Stage: acceptance-spec

### A.1. Write Acceptance Test Scenarios
- [x] Create spec.md with user-facing acceptance scenarios (FEATURE-A1..A7)
- [x] Define expected behavior: git SHA format, immutability, error handling
- [x] Cover happy path: valid SHA with correct format (40 hex chars)
- [x] Cover edge cases: missing SHA, custom admin port, multiple requests
- [x] Security: verify no sensitive data leakage beyond commit SHA

### A.2. Implement Docker-based Integration Tests
- [ ] Create `tests/acceptance/version_endpoint_git_sha_test.sh` script
- [ ] Use Docker Compose to spin up real ubox-crosser server
- [ ] Test real link: server startup → admin server health check → /version request
- [ ] Validate JSON response contains valid 40-character git SHA
- [ ] Verify SHA matches actual `git rev-parse HEAD` from repository
- [ ] Test immutability: multiple requests return identical SHA
- [ ] Test custom admin address with proper SHA response
- [ ] Test missing/unknown commit SHA fallback behavior

### A.3. Lint and Validate OpenSpec
- [ ] Run `openspec validate openspec/changes/REQ-rbac-1776859928`
- [ ] Run scenario reference check for consistency
- [ ] Ensure all scenario names follow FEATURE-A<N> pattern
- [ ] Verify spec.md ADDED block structure and RELATED section

### A.4. Complete and Push
- [x] Commit changes to acceptance-spec branch
- [x] Push to origin with proper branch naming
- [ ] Update BKD issue tags to reflect acceptance-spec completion

## Definition of Done
- [x] Acceptance scenarios defined in spec.md (FEATURE-A1..A7)
- [ ] Integration tests implemented and pass in isolated container
- [ ] OpenSpec lint validation passes
- [ ] Git SHA format validation (40 hex characters)
- [ ] Changes pushed and ready for dev stage implementation
