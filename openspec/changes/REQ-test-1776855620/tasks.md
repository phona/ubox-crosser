# Tasks for REQ-test-1776855620 — /healthz Endpoint with Service Uptime

## Overview
- **Feature**: GET /healthz endpoint returning service uptime
- **Stage**: acceptance-spec
- **Status**: In Progress
- **Requirement**: Add /healthz endpoint that returns server uptime in seconds

## Stage: acceptance-spec

### A.1. Write Acceptance Test Scenarios
- [x] Create spec.md with user-facing acceptance scenarios (FEATURE-A1..A8)
- [x] Define expected behavior: status codes, JSON format, uptime field
- [x] Cover happy path, error cases, uptime tracking, custom configuration
- [x] Include edge case: immediately after server start (near-zero uptime)

### A.2. Implement Docker-based Integration Tests
- [ ] Create `tests/acceptance/healthz_endpoint_test.sh` script
- [ ] Use Docker Compose to spin up real ubox-crosser server
- [ ] Test real link: server startup → health check → /healthz request → response validation
- [ ] Validate JSON response structure and uptime as integer
- [ ] Test uptime monotonically increases over time
- [ ] Test custom admin address configuration
- [ ] Test non-GET method rejection (405)
- [ ] Test uptime tracking across multiple requests with timing validation

### A.3. Lint and Validate OpenSpec
- [ ] Run `openspec validate openspec/changes/REQ-test-1776855620`
- [ ] Run `/opt/sisyphus/scripts/check-scenario-refs.sh` for reference consistency
- [ ] Ensure all scenario names follow FEATURE-A\<N\> pattern
- [ ] Verify spec.md ADDED block structure

### A.4. Complete and Push
- [ ] Commit changes to stage branch
- [ ] Push to origin/stage/REQ-test-1776855620-dev
- [ ] Move issue to review status with proper tags

## Definition of Done
- [x] Acceptance scenarios defined in spec.md (FEATURE-A1..A8)
- [ ] Integration tests pass in isolated container environment
- [ ] OpenSpec lint passes
- [ ] Changes pushed and ready for dev/contract stage feedback
