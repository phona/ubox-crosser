# Tasks: REQ-e2e-1776916220 /healthz Endpoint

## Stage: contract-spec

**Objective:** Define the contract and write comprehensive contract tests for the `/healthz` endpoint.

### ✅ Task 1: Write OpenAPI Specification

**Status:** COMPLETED

**Deliverables:**
- `openspec/changes/REQ-e2e-1776916220/specs/healthz/spec.md`
  - Endpoint definition: `GET /healthz`
  - Request/response contract
  - HTTP status codes and JSON schema
  - Implementation requirements
  - Test scenarios (S1, S2, S3)

### ✅ Task 2: Write Contract Tests

**Status:** COMPLETED

**Deliverables:**
- `tests/contract/healthz_test.go`
  - `TestHealthzS1`: Validates response format and status code
  - `TestHealthzS2`: Validates uptime progression
  - `TestHealthzS3`: Validates duration string format
  - Uses environment variable `PROXY_HTTP_ADDR` for endpoint address
  - Follows integration test skill conventions (build tag, timeouts, error handling)

**Test Execution:**
```bash
docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner
```

Expected result: **RED** (implementation not yet complete)

### ✅ Task 3: Define Test Scenarios

**Status:** COMPLETED

**Scenarios:**
- **REQ-e2e-1776916220-S1:** Healthz returns valid 200 response with healthy status and numeric uptime
- **REQ-e2e-1776916220-S2:** Uptime increases when checking multiple times
- **REQ-e2e-1776916220-S3:** Duration format matches Go time.Duration string representation

### ✅ Task 4: Document Implementation Requirements

**Status:** COMPLETED

**Documents:**
- `proposal.md`: Problem statement, success criteria, scope
- `design.md`: Architecture, data flow, configuration, testing strategy
- `specs/healthz/spec.md`: Technical contract and scenarios

---

## Stage: dev (NOT YET STARTED)

**Objective:** Implement HTTP listener and /healthz endpoint in ProxyServer.

### Task 1: Add HTTP Server to ProxyServer

- [ ] Add `startTime` field to ProxyServer struct
- [ ] Initialize startTime during `NewProxyServer()`
- [ ] Create HTTP listener on configurable port (env: `PROXY_HTTP_PORT`, default 8080)
- [ ] Run HTTP listener in separate goroutine (non-blocking)
- [ ] Handle listener startup errors gracefully

### Task 2: Implement /healthz Handler

- [ ] Create HTTP handler for `GET /healthz`
- [ ] Calculate uptime: `time.Now() - startTime`
- [ ] Format duration string using `time.Duration.String()`
- [ ] Return JSON with status="healthy", uptime.seconds, uptime.duration
- [ ] Set Content-Type header to application/json

### Task 3: Environment Configuration

- [ ] Accept `PROXY_HTTP_PORT` environment variable
- [ ] Accept `PROXY_HTTP_ADDR` environment variable
- [ ] Set sensible defaults (port 8080, bind to 0.0.0.0)

### Task 4: Run Contract Tests

- [ ] Execute contract tests: `docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner`
- [ ] Verify all tests pass (S1, S2, S3)
- [ ] Check response format matches specification exactly

---

## Stage: accept (NOT YET STARTED)

**Objective:** Verify implementation against contract.

### Task 1: Functional Verification

- [ ] Manual test: `curl http://localhost:8080/healthz` returns 200 OK
- [ ] Verify JSON contains all required fields
- [ ] Check uptime is non-zero after service has run

### Task 2: Contract Compliance

- [ ] Run full contract test suite: `docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner`
- [ ] All 3 scenarios (S1, S2, S3) must PASS
- [ ] No regressions in existing functionality

### Task 3: Documentation

- [ ] Update README with /healthz endpoint documentation
- [ ] Document environment variables
- [ ] Add curl examples

---

## Success Criteria (All Stages Combined)

1. ✅ Contract specification locked in `specs/healthz/spec.md`
2. ✅ Contract tests written in `tests/contract/healthz_test.go` (currently RED, will be GREEN after dev)
3. ⏳ HTTP endpoint implemented in ProxyServer
4. ⏳ All contract tests pass (GREEN)
5. ⏳ No regressions in existing proxy functionality
6. ⏳ Documentation updated

---

## Files Modified/Created

### contract-spec stage
- ✅ `openspec/changes/REQ-e2e-1776916220/proposal.md` (NEW)
- ✅ `openspec/changes/REQ-e2e-1776916220/design.md` (NEW)
- ✅ `openspec/changes/REQ-e2e-1776916220/specs/healthz/spec.md` (NEW)
- ✅ `openspec/changes/REQ-e2e-1776916220/tasks.md` (NEW)
- ✅ `tests/contract/healthz_test.go` (NEW)

### dev stage (pending)
- `server/server.go` (MODIFY: add HTTP listener and startTime)
- `Makefile` (possibly: update test targets)

### accept stage (pending)
- `README.md` (MODIFY: document /healthz endpoint)

---

## Testing Matrix

| Test | Scenario | Status | Expected Result |
|------|----------|--------|-----------------|
| TestHealthzS1 | Response validation | ⚠️ RED | Will PASS when dev completes |
| TestHealthzS2 | Uptime progression | ⚠️ RED | Will PASS when dev completes |
| TestHealthzS3 | Duration format | ⚠️ RED | Will PASS when dev completes |

**Note:** Tests are intentionally RED at contract-spec stage (behavior-driven development). They turn GREEN once the dev agent implements the feature.
