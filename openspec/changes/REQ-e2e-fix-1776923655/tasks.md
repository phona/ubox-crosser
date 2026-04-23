# REQ-e2e-fix-1776923655: /healthz Endpoint Contract Tests

## Overview
Establish and validate the REST API contract for the `/healthz` endpoint to enable service health monitoring with uptime reporting.

## Stage: contract-tests

### Contract Test Specifications

#### REQ-e2e-fix-1776923655-S1: Health Check Endpoint Returns 200 OK
- [x] Spec written in `specs/healthz-endpoint/spec.md`
- [ ] Contract test validates endpoint responds with HTTP 200
- [ ] Response body is valid JSON
- **Test**: `tests/contract/healthz_test.go::TestREQe2efixS1HealthCheckEndpointReturns200OK`

#### REQ-e2e-fix-1776923655-S2: Endpoint Reports Service Uptime in Seconds
- [x] Spec written
- [ ] Contract test validates `uptime_seconds` field exists
- [ ] Field contains non-negative integer value
- **Test**: `tests/contract/healthz_test.go::TestREQe2efixS2EndpointReportsServiceUptimeInSeconds`

#### REQ-e2e-fix-1776923655-S3: Uptime Increases Over Time
- [x] Spec written
- [ ] Contract test validates uptime increments correctly
- [ ] Two requests 5 seconds apart show appropriate uptime difference (±1 second tolerance)
- **Test**: `tests/contract/healthz_test.go::TestREQe2efixS3UptimeIncreasesOverTime`

#### REQ-e2e-fix-1776923655-S4: Health Endpoint Available After Service Startup
- [x] Spec written
- [ ] Docker Compose startup sequence validates endpoint availability
- [ ] Service health check uses `/healthz` endpoint
- **Config**: `tests/contract/docker-compose.yml`

#### REQ-e2e-fix-1776923655-S5: Response Format Includes Status and Timestamp
- [x] Spec written
- [ ] Contract test validates JSON structure
- [ ] Required fields present: `status`, `uptime_seconds`, `timestamp`
- [ ] `status` field equals "healthy"
- **Test**: `tests/contract/healthz_test.go::TestREQe2efixS5ResponseFormatIncludesStatusField`

#### REQ-e2e-fix-1776923655-S6: Endpoint Handles Concurrent Requests
- [x] Spec written
- [ ] Contract test validates 10+ concurrent GET requests succeed
- [ ] No race conditions or dropped requests
- [ ] All responses contain valid uptime values
- **Test**: `tests/contract/healthz_test.go::TestREQe2efixS6EndpointHandlesConcurrentRequests`

#### REQ-e2e-fix-1776923655-S7: Service Uptime Resets on Restart
- [x] Spec written
- [ ] Acceptance test covers (in acceptance stage)
- [ ] Contract validates uptime field behavior
- **Note**: Manual or acceptance-level testing required

### Deliverables

- [x] API Contract Specification: `openspec/changes/REQ-e2e-fix-1776923655/specs/healthz-endpoint/spec.md`
  - Detailed endpoint contract with request/response format
  - Scenario-based requirements (S1-S7)
  - Performance and constraint specifications

- [x] Contract Test Suite: `tests/contract/healthz_test.go`
  - 6 contract tests covering S1, S2, S3, S5, S6
  - Integration build tag: `//go:build integration`
  - Environment variable configuration from Docker Compose
  - Proper timeout and error handling

- [x] Docker Compose Setup: `tests/contract/docker-compose.yml`
  - Service under test: ubox-crosser with /healthz on port 8080
  - Health check using `/healthz` endpoint
  - Test runner service with proper dependencies
  - Coverage data volume for CI integration

- [x] This tasks file: `openspec/changes/REQ-e2e-fix-1776923655/tasks.md`

### Testing Validation Checklist

- [x] Spec follows OpenAPI contract format
- [x] All scenario names follow convention: `REQ-e2e-fix-1776923655-S<N>`
- [x] Contract tests use proper build tag: `//go:build integration`
- [x] Docker Compose configured with appropriate health checks
- [x] Tests read configuration from environment variables
- [x] No hardcoded addresses (except localhost fallback)
- [x] Proper timeout handling on all network operations
- [x] JSON response structure matches spec
- [ ] Contract tests should run RED initially (endpoint not implemented)
- [ ] All tests pass once implementation is complete

### Running Tests

#### Run contract tests locally (after implementation):
```bash
docker compose -f tests/contract/docker-compose.yml up --build --exit-code-from test-runner
```

#### Run specific test:
```bash
docker compose -f tests/contract/docker-compose.yml up --build --exit-code-from test-runner && \
  go test -v -tags integration -run TestREQe2efixS1 ./tests/contract/...
```

### Contract Boundaries (LOCKED)

These are the LOCKED API agreements between client and service:

1. **Endpoint Path**: `/healthz` (exact match)
2. **HTTP Method**: GET only
3. **Port**: 8080 (configurable, environment-driven)
4. **Response Status**: HTTP 200 OK (always, when healthy)
5. **Content-Type**: `application/json`
6. **Response Fields**:
   - `status` (string): Always "healthy" when responding 200
   - `uptime_seconds` (integer): ≥ 0, increments with time
   - `timestamp` (integer): Unix epoch seconds
7. **Performance**: Response time < 100ms
8. **Authentication**: None required
9. **Concurrency**: Must handle 10+ parallel requests safely

### Notes

- Contract tests focus on the **API surface**, not implementation details
- Tests use Docker Compose to validate service startup and health check integration
- Uptime calculation accuracy is validated through time-based assertions
- These tests should fail (RED) until the dev-impl stage implements `/healthz`
- Once implementation is complete, all tests should pass (GREEN)
- Acceptance tests in later stage will cover integration with load balancers and restart scenarios
