# REQ-e2e-fix2-1776924685: /healthz Endpoint with Service Uptime

## Overview
Implement a `/healthz` HTTP endpoint that returns current health status and service uptime information to support external health checks and load balancer monitoring.

## Implementation Tasks

### Phase 1: Core Implementation
- [ ] Add HTTP server to ubox-crosser (if not already present)
- [ ] Implement `/healthz` endpoint handler
- [ ] Track service startup/initialization time
- [ ] Calculate and return uptime_seconds
- [ ] Return JSON response with status, uptime_seconds, timestamp
- [ ] Ensure endpoint accessible without authentication
- [ ] Configure listening port (default: 8080)

### Phase 2: Integration
- [ ] Add health check port configuration
- [ ] Ensure endpoint available immediately after service init
- [ ] Handle concurrent requests efficiently
- [ ] Verify response time < 100ms

### Phase 3: Testing
- [ ] Unit tests for uptime calculation
- [ ] Integration tests for HTTP endpoint behavior
- [ ] Concurrent request handling tests
- [ ] Service restart/uptime reset verification

## Stage: acceptance-tests

### Acceptance Criteria

#### FEATURE-A1: Health Check Endpoint Returns 200 OK
- [x] Spec written in `/specs/healthz-endpoint/spec.md`
- [ ] Integration test validates endpoint responds with 200
- [ ] Load balancer can successfully call the endpoint

#### FEATURE-A2: Endpoint Reports Service Uptime in Seconds
- [x] Spec written
- [ ] Integration test validates uptime_seconds field exists
- [ ] Uptime value matches actual service running duration

#### FEATURE-A3: Uptime Increases Over Time
- [x] Spec written
- [ ] Integration test validates uptime increments correctly
- [ ] Uptime difference matches elapsed time (±1 second tolerance)

#### FEATURE-A4: Health Endpoint Available After Service Startup
- [x] Spec written
- [ ] Integration test validates endpoint availability after startup
- [ ] Service startup time is captured correctly

#### FEATURE-A5: Response Format Includes Required Fields
- [x] Spec written
- [ ] Integration test validates JSON structure
- [ ] All required fields present: status, uptime_seconds, timestamp

#### FEATURE-A6: Endpoint Works Concurrently
- [x] Spec written
- [ ] Integration test validates 10+ concurrent requests succeed
- [ ] No race conditions or dropped requests

#### FEATURE-A7: Service Uptime Resets on Restart
- [x] Spec written
- [ ] Service restart triggers uptime reset
- [ ] New uptime starts from 0 after restart

### Deliverables
- [x] Acceptance specification: `openspec/changes/REQ-e2e-fix2-1776924685/specs/healthz-endpoint/spec.md`
- [ ] Acceptance tests: `tests/acceptance/healthz_e2e_test.go` (or similar)
- [ ] Docker Compose integration test setup: `tests/acceptance/docker-compose.yml`
- [x] This tasks file: `openspec/changes/REQ-e2e-fix2-1776924685/tasks.md`

## Validation Checklist

- [x] All 7 acceptance scenarios have clear Given/When/Then format
- [x] Scenarios are written from user/operator perspective
- [x] Test cases cover happy path and edge cases
- [x] Concurrent request handling is tested
- [x] Uptime calculations and reset behavior are verified
- [ ] Spec file validates with `openspec validate`
- [ ] Integration tests pass with working implementation
- [ ] Docker Compose test setup is functional
- [ ] Endpoint response time < 100ms verified

## Notes

- The acceptance spec describes user-facing behavior from operator perspective
- Implementation details are handled in dev-impl stage
- Integration tests require actual `/healthz` endpoint implementation
- Uptime is measured from service readiness, not binary startup
- Endpoint is Kubernetes-probe compatible
