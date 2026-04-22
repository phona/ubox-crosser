# Contract Testing Tasks: /version Endpoint Git SHA

## Stage: contract-tests

This stage documents contract tests that validate the /version endpoint feature behavior under real deployment conditions.

### Task 1: Formalize Contract Test Scenarios

**Status**: ✓ COMPLETED

**Description**: Document all contract testing scenarios that must pass.

**Artifacts**:
- `specs/version_endpoint_contract/spec.md` - Contract scenario definitions mapped to tests

**Completed**: Contract test specifications documented with 7 scenarios (REQ-e2e-1776862572-S1 through S7)

---

### Task 2: Validate Test Coverage

**Status**: ✓ COMPLETED

All contract scenarios are covered by existing integration tests in `tests/integration/version_endpoint_test.go`:

- ✓ S1: HTTP 200 response + application/json content-type
- ✓ S2: Git SHA field format (40-char hex)
- ✓ S3: Commit SHA immutability across requests
- ✓ S4: JSON structure preservation (all fields)
- ✓ S5: Custom admin port support
- ✓ S6: Security - no sensitive data leakage
- ✓ S7: Response time SLA (5-second timeout)

---

### Task 3: Document Requirements

**Status**: ✓ COMPLETED

**Docker Compose Stack**:
- proxy-server (ubox-crosser binary) with admin HTTP on :8080
- test-runner service executing integration tests
- ADMIN_SERVER_ADDR environment variable injection

**Build Requirements**:
- COMMIT injected via LDFLAGS: `-X main.Version=$(git rev-parse HEAD)`
- 40-character git SHA available at compile time
- /version endpoint handler returns JSON

**Test Tags**:
- `//go:build integration` build tag
- Run with: `go test -tags integration`
- Environment-driven configuration

---

### Task 4: Define Locked Zones

**Status**: ✓ COMPLETED

**LOCKED Files** (contract-spec-agent owns, dev-agent cannot modify):
```
tests/
├── docker-compose.yml                       # Stack contract definition
├── integration/version_endpoint_test.go     # Contract test scenarios
└── Dockerfile.test                          # Test runner image
```

**Rationale**: Contract tests define acceptance criteria *before* implementation

---

## Summary

All contract-spec tasks are complete. The feature is ready for dev-agent implementation.

**Deliverables**:
- [x] Contract test specification (specs/version_endpoint_contract/spec.md)
- [x] Test coverage validation (all 7 scenarios covered)
- [x] Requirements documentation (design.md)
- [x] Locked zones definition (this file)

**Next Step**: dev-agent implements `/version` handler in cmd/ to make tests pass

**Validation**:
```bash
# Run contract tests
docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner

# Should be RED before implementation, GREEN after completion
```
