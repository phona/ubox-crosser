# REQ-e2e-fix2-1776924685: /healthz Endpoint with Service Uptime

## Overview
Implement contract tests for `/healthz` HTTP endpoint that returns the current health status and service uptime information. These tests verify the contract between the endpoint and its consumers.

## Stage: contract-tests

### Test Scenarios
- [x] Spec written: `openspec/changes/REQ-e2e-fix2-1776924685/specs/healthz-endpoint/spec.md`
- [ ] Contract tests validate endpoint responds with 200 OK (REQ-e2e-fix2-1776924685-S1)
- [ ] Contract tests validate uptime_seconds field exists (REQ-e2e-fix2-1776924685-S2)
- [ ] Contract tests validate uptime increments correctly (REQ-e2e-fix2-1776924685-S3)
- [ ] Contract tests validate JSON structure (REQ-e2e-fix2-1776924685-S4)
- [ ] Contract tests validate concurrent request handling (REQ-e2e-fix2-1776924685-S5)

### Deliverables
- [x] Contract specification: `openspec/changes/REQ-e2e-fix2-1776924685/specs/healthz-endpoint/spec.md`
- [ ] Contract tests: `tests/contract/healthz_contract_test.go`
- [ ] Test helper functions in contract tests file
- [ ] This tasks file: `openspec/changes/REQ-e2e-fix2-1776924685/tasks.md`

### Test Environment
- Docker Compose setup for contract testing
- BASE_URL environment variable injection for service under test
- Integration test runner container

### Acceptance Criteria
- All 5 scenarios have clear Given/When/Then format
- Scenarios are written from contract perspective
- Test cases can be run independently
- Tests are RED before implementation (will fail without /healthz endpoint)
- Spec validates with `openspec validate`

## Notes

- Contract tests verify the API contract (request/response format)
- Tests should run against the actual service via HTTP
- No internal implementation details should leak into test assertions
- Tests are written as blackbox acceptance tests following the integration-test skill
