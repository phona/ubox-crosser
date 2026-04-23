# Tasks for REQ-final4-1776953521: /buildinfo HTTP Endpoint

## Stage: Specification
- [x] Write contract.spec.yaml with endpoint schema
- [x] Write spec.md with acceptance scenarios

## Stage: Implementation
- [x] Implement /buildinfo handler in server/healthcheck.go
- [x] Add /healthz handler alongside /buildinfo
- [x] Implement HTTP health check server with proper startup
- [x] Add git_sha injection via ldflags in Makefile
- [x] Add unit tests for buildinfo handler in server/healthcheck_test.go
- [x] Add integration tests for /buildinfo endpoint in tests/acceptance/buildinfo_test.go
- [x] Update docker-compose.yml to set BUILD_ID environment variable
- [x] Update Dockerfile.test to inject git SHA and install curl for healthcheck
- [x] Integrate HTTP server startup in cmd/server/server.go

## Stage: PR & Review
- [ ] Push to feat/REQ-final4-1776953521
- [ ] Create pull request with clear description
- [ ] All tests passing in CI
