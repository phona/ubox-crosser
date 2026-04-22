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
- [x] Create `tests/integration/version_endpoint_test.go` with Go-based integration tests
- [x] Use Docker Compose environment for testing against real running server
- [x] Test HTTP GET /version endpoint with response validation
- [x] Validate JSON response contains valid 40-character git SHA field
- [x] Test immutability: multiple requests return identical commit SHA
- [x] Test custom admin address configuration
- [x] Test security: verify no sensitive data leakage
- [x] Update docker-compose.yml with ADMIN_SERVER_ADDR environment variable
- [x] Update Dockerfile.test to inject git commit SHA via LDFLAGS during build

### A.3. Lint and Validate OpenSpec
- [x] Verify spec.md ADDED block structure (scenarios FEATURE-A1..A7 defined)
- [x] Verify spec.md RELATED section with endpoint and integration details
- [x] Ensure all scenario names follow FEATURE-A<N> pattern
- [x] Integration tests properly formatted and follow Docker-based testing patterns

### A.4. Complete and Push
- [x] Commit changes to contract-spec branch
- [x] Push to origin with proper branch naming
- [ ] Update BKD issue tags to reflect contract-spec completion

## Definition of Done
- [x] Acceptance scenarios defined in spec.md (FEATURE-A1..A7)
- [x] Integration tests implemented: tests/integration/version_endpoint_test.go
- [x] Docker configuration updated for admin server testing
- [x] Dockerfile updated to inject git SHA via LDFLAGS
- [x] Git SHA format validation (40 hex characters) verified in test assertions
- [x] Changes committed and ready for validation
- [ ] Integration tests run successfully (requires runner pod environment)
