# REQ-e2e-1776916220: /healthz Endpoint with Service Uptime

## Overview
Implement a `/healthz` HTTP endpoint that returns the current health status and service uptime information to support external health checks and monitoring.

## Implementation Tasks

### Phase 1: Core Implementation (TBD in dev-impl)
- [ ] Add HTTP server support to ubox-crosser
- [ ] Implement `/healthz` endpoint handler
- [ ] Track service startup time and calculate uptime
- [ ] Return JSON response with status, uptime_seconds, and timestamp fields
- [ ] Ensure endpoint doesn't require authentication
- [ ] Configure health check listening port (default: 8080)

### Phase 2: Integration
- [ ] Add configuration for health check port in server config
- [ ] Ensure `/healthz` endpoint is available immediately after service initialization
- [ ] Handle concurrent requests to `/healthz` efficiently
- [ ] Set response time SLA (< 100ms)

### Phase 3: Testing (TBD in test)
- [ ] Unit tests for uptime calculation logic
- [ ] Integration tests for HTTP endpoint behavior
- [ ] Load testing for concurrent request handling
- [ ] Restart/uptime reset verification tests

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

#### FEATURE-A5: Response Format Includes Status Field
- [x] Spec written
- [ ] Integration test validates JSON structure
- [ ] All required fields present: status, uptime_seconds, timestamp

#### FEATURE-A6: Endpoint Works Concurrently
- [x] Spec written
- [ ] Integration test validates 10+ concurrent requests succeed
- [ ] No race conditions or dropped requests

#### FEATURE-A7: Service Uptime Resets on Restart
- [x] Spec written (manual verification)
- [ ] Service restart triggers uptime reset
- [ ] New uptime starts from 0 after restart

### Deliverables
- [x] Acceptance specification: `openspec/changes/REQ-e2e-1776916220/specs/healthz-endpoint/spec.md`
- [x] Acceptance tests: `tests/acceptance/healthz_test.go`
- [x] Docker Compose integration test: `tests/acceptance/docker-compose.yml`
- [x] This tasks file: `openspec/changes/REQ-e2e-1776916220/tasks.md`

## Validation Checklist

- [x] All 7 acceptance scenarios have clear Given/When/Then format
- [x] Scenarios are written from user perspective (not implementation details)
- [x] Test cases cover happy path and edge cases
- [x] Concurrent request handling is tested
- [x] Uptime calculations are verified
- [x] Spec file validates with `openspec validate`
- [ ] Integration tests pass with working implementation
- [ ] Docker Compose test setup is functional

## Notes

- The spec describes user-facing behavior; implementation details are in dev-impl stage
- Integration tests require the actual `/healthz` endpoint implementation
- Docker Compose setup allows testing without production infrastructure
- Acceptance tests should run as part of CI/CD pipeline
