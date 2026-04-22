# Tasks for REQ-test-1776855620 — /healthz Endpoint with Service Uptime

## Overview
- **Feature**: GET /healthz endpoint returning service uptime
- **Current Stage**: contract-tests
- **Status**: In Progress
- **Requirement**: Add /healthz endpoint that returns server uptime in seconds

## Stage: contract-tests

### C.1. Write Contract Test Suite
- [x] Create `tests/contract/healthz_test.go` with comprehensive test coverage
  - [x] Test GET /healthz returns 200 OK
  - [x] Test response is valid JSON with status and uptime_seconds fields
  - [x] Test uptime_seconds increases over time (monotonic)
  - [x] Test POST returns 405 Method Not Allowed
  - [x] Test PUT returns 405 Method Not Allowed
  - [x] Test DELETE returns 405 Method Not Allowed
  - [x] Test Content-Type is application/json
  - [x] Test response time is acceptable (< 100ms)
  - [x] Test status is always "ok"

### C.2. Set Up Docker Compose Environment
- [x] Create `tests/docker-compose.contract.yml` for isolated contract testing
  - [x] Configure proxy server with management endpoints exposed
  - [x] Configure test-runner container with proper environment variables
  - [x] Ensure management server healthcheck passes before tests run

### C.3. Verify Tests are RED (Endpoint Not Implemented)
- [ ] Run contract tests: `docker compose -f tests/docker-compose.contract.yml up --build --exit-code-from test-runner`
- [ ] Confirm tests fail (RED) since /healthz endpoint not yet implemented
- [ ] Document baseline failure state

### C.4. Specification and Design Documentation
- [x] Update proposal.md with clear acceptance criteria
- [x] Update design.md with implementation approach
- [x] Update spec.md with OpenAPI definition and detailed scenarios
- [x] Ensure scenario names follow REQ-test-1776855620-S<N> pattern

### C.5. Lint and Validate OpenSpec
- [ ] Run `openspec validate openspec/changes/REQ-test-1776855620`
- [ ] Verify all scenario references are correct
- [ ] Check spec.md structure and formatting

### C.6. Commit and Push Spec Changes
- [ ] Stage all spec and test files
- [ ] Commit with message: "chore(contract-spec): add /healthz endpoint contract tests for REQ-test-1776855620"
- [ ] Push to `stage/REQ-test-1776855620-dev` branch
- [ ] Update BKD issue tags to indicate contract-spec completion

## Definition of Done (Contract-Tests Stage)
- [x] Contract test suite created and RED (tests fail before implementation)
- [x] Docker Compose test environment configured
- [x] Complete specification with scenarios
- [x] Implementation design documented
- [ ] Changes committed and pushed to stage branch
- [ ] BKD issue updated with contract-spec tag
- [ ] Ready for dev-agent to implement the endpoint

## Next Stage: dev
- dev-agent will implement /healthz endpoint in server/management.go
- Tests should transition from RED to GREEN when implementation is complete
