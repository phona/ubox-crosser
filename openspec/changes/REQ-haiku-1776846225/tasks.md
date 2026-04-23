# Tasks for REQ-haiku-1776846225 — /version Endpoint Acceptance Spec

## Overview
- **Feature**: GET /version endpoint for build information
- **Stage**: acceptance-spec
- **Status**: In Progress
- **Real Link Test**: E2E test with actual server startup and version endpoint validation

## Stage: acceptance-spec

### A.1. Write Acceptance Test Scenarios
- [x] Create spec.md with user-facing acceptance scenarios (FEATURE-A1..A6)
- [x] Define expected behavior: status codes, JSON format, default values
- [x] Cover happy path, error cases, custom configuration

### A.2. Implement Docker-based Integration Tests
- [ ] Create `tests/acceptance/version_endpoint_test.sh` script
- [ ] Use Docker Compose to spin up real ubox-crosser server
- [ ] Test real link: server startup → health check → /version request → response validation
- [ ] Validate JSON response structure and commit hash format
- [ ] Test custom admin address configuration
- [ ] Test non-GET method rejection (405)

### A.3. Lint and Validate OpenSpec
- [ ] Run `openspec validate openspec/changes/REQ-haiku-1776846225`
- [ ] Run `/opt/sisyphus/scripts/check-scenario-refs.sh` for reference consistency
- [ ] Ensure all scenario names follow FEATURE-A\<N\> pattern
- [ ] Verify spec.md ADDED block structure

### A.4. Complete and Push
- [ ] Commit changes to stage branch
- [ ] Push to origin/stage/REQ-haiku-1776846225-dev
- [ ] Move issue to review status with proper tags

## Definition of Done
- [x] Acceptance scenarios defined in spec.md (FEATURE-A1..A6)
- [ ] Integration tests pass in isolated container environment
- [ ] OpenSpec lint passes
- [ ] Changes pushed and ready for dev/contract stage feedback
