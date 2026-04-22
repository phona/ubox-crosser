# Contract Testing Specification: /version Endpoint

## Overview

This specification formalizes the **contract tests** for the `/version` endpoint that returns git commit SHA. These are blackbox integration tests validating the endpoint behavior under real deployment conditions (docker-compose stack).

## Test Environment Contract

**Stack**: Docker Compose with real ubox-crosser binary
**Admin Server**: Accessible at ADMIN_SERVER_ADDR (default: proxy-server:8080)
**Build Requirements**: Git SHA injected via LDFLAGS
**Test Framework**: Go testing with `//go:build integration` tag

## Contract Scenarios

All scenarios are tested in `tests/integration/version_endpoint_test.go`

### REQ-e2e-1776862572-S1: HTTP Response Status and Content-Type

**Test Function**: `TestVersionEndpointReturnsGitSHA/returns 200 OK` + `returns application/json content type`

**Contract**:
- Endpoint: `GET /version` on admin server
- Response Status: **200 OK**
- Response Header: `Content-Type: application/json`
- No authentication required

**Validation**:
```go
resp.StatusCode == http.StatusOK
resp.Header.Get("Content-Type") == "application/json"
```

---

### REQ-e2e-1776862572-S2: Git SHA Field Format

**Test Function**: `TestVersionEndpointReturnsGitSHA/response contains valid commit SHA`

**Contract**:
- Response Body: Valid JSON object
- Required Field: `commit`
- Field Type: String
- Format: 40-character hexadecimal (git SHA)
- Pattern: `^[0-9a-f]{40}$`

**Validation**:
```go
var resp VersionResponse
json.Unmarshal(body, &resp)
regexp.MustCompile(`^[0-9a-f]{40}$`).MatchString(resp.Commit)
```

---

### REQ-e2e-1776862572-S3: Commit SHA Immutability

**Test Function**: `TestVersionEndpointReturnsGitSHA/commit SHA is immutable across requests`

**Contract**:
- During server lifetime (no restart)
- Multiple sequential requests to `/version`
- `commit` field value MUST remain identical
- Response must be deterministic

**Validation**:
```go
resp1.Commit == resp2.Commit  // Multiple requests return same SHA
```

**Why**: Ensures build-time injection is immutable at runtime

---

### REQ-e2e-1776862572-S4: JSON Structure Preservation

**Test Function**: `TestVersionEndpointReturnsGitSHA/other version fields are preserved`

**Contract**:
- Response must include all expected fields:
  - `version` (string, non-empty)
  - `module` (string, non-empty)
  - `go_os` (string, non-empty)
  - `go_arch` (string, non-empty)
  - `commit` (string, 40-char hex)

**Validation**:
```go
var resp VersionResponse
resp.Version != ""
resp.Module != ""
resp.GoOS != ""
resp.GoArch != ""
resp.Commit matches ^[0-9a-f]{40}$
```

---

### REQ-e2e-1776862572-S5: Custom Admin Port

**Test Function**: `TestVersionEndpointWithCustomAdminPort`

**Contract**:
- Optional CUSTOM_ADMIN_ADDR environment variable
- When set: `/version` available on alternate port
- All other contracts (S1-S4, S6-S7) apply to custom port
- Must not break when env var not set

**Validation**:
```go
CUSTOM_ADMIN_ADDR = "custom-host:9000"  // From environment
GET http://custom-host:9000/version
// Expect same response format as S1-S4
```

---

### REQ-e2e-1776862572-S6: Security - No Sensitive Data Leakage

**Test Function**: `TestVersionEndpointSecurityNoLeakage`

**Contract**:
- Response must only contain expected JSON fields
- No sensitive keywords in response body:
  - PASSWORD, TOKEN, SECRET (case-insensitive)
  - LDFLAGS, buildinfo
  - System paths: /root, /home
  - Environment variable names
  - Build timestamps

**Validation**:
```go
// Only expected fields present
for field := range resp {
  assert(field in ["version", "module", "go_os", "go_arch", "commit"])
}

// No sensitive patterns
sensitivePatterns := []string{"PASSWORD", "TOKEN", "SECRET", "LDFLAGS", ...}
for _, pattern := range sensitivePatterns {
  assert(!regexp.MatchString(`(?i)` + pattern, responseBody))
}
```

**Why**: Prevent accidental leakage of build environment or configuration secrets

---

### REQ-e2e-1776862572-S7: Response Time SLA

**Test Function**: Implicit in all tests (via 5-second timeout)

**Contract**:
- HTTP client timeout: 5 seconds
- Expected response time: < 100ms
- No blocking I/O in /version handler
- Handler should be fast (~1ms for JSON marshaling)

**Why**: Admin server endpoints must be responsive for health checks

---

## Test Execution Contract

**Test Command**:
```bash
docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner
```

**Exit Codes**:
- `0`: All contract scenarios pass
- Non-zero: One or more scenarios fail (contract violation)

**Artifacts**:
- **Logs**: Container logs visible in stdout
- **Coverage**: Collected in /tmp/crosser-coverage (if COVERAGE_HOST_DIR set)
- **Test Report**: Go test verbose output

---

## Build-Time Contract

The feature depends on proper build-time configuration:

**Makefile/Dockerfile Requirement**:
```makefile
COMMIT := $(shell git rev-parse HEAD)
GO_LDFLAGS := -ldflags="-X main.Version=$(COMMIT)"
```

**When this contract is violated**:
- Scenario S2 fails: `commit` field is empty or not 40-char hex
- Scenario S3 fails: `commit` field is "unknown" or ""
- Impact: Feature is non-functional

---

## Related Contracts

- **Parent**: REQ-rbac-1776859928/specs/version_endpoint_git_sha/ (acceptance specs)
- **Owner**: contract-spec-agent (this stage)
- **Implementation**: dev-agent (next stage)
- **Validation**: Automated CI (ci-int job)

---

## Contract Status

All scenarios defined above are **currently tested** by existing code in:
- `tests/integration/version_endpoint_test.go`
- `tests/docker-compose.yml`

**Status**: READY FOR IMPLEMENTATION

The tests are currently RED (before /version handler exists or is incomplete).
After dev-agent implements the handler, tests should be GREEN.
