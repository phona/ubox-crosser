# Version Endpoint Git SHA E2E Acceptance Specification

## Feature: GET /version Endpoint Returns Git SHA (E2E Validation)

This specification documents end-to-end acceptance scenarios for the /version endpoint that returns git commit SHA. These scenarios validate the feature from a user/operator perspective during full system integration testing.

### ADDED

#### Scenario: FEATURE-A1 — E2E Server startup and endpoint accessibility

**Given**
- ubox-crosser source code is available in the test environment
- Real git repository with valid commit history exists
- Docker Compose environment is configured for e2e testing

**When**
- E2E test orchestrates server startup via Docker Compose
- Server initialization completes and admin HTTP server binds to configured port
- Test waits for server health check to pass

**Then**
- Server starts successfully within acceptable time (10 seconds)
- Admin HTTP server is accessible on localhost:8080 by default
- GET request to `/version` endpoint returns HTTP 200 OK
- Response headers include `Content-Type: application/json`
- Server healthcheck via `/version` succeeds

#### Scenario: FEATURE-A2 — E2E /version response contains valid git SHA matching build

**Given**
- ubox-crosser server is running via Docker Compose in e2e environment
- Binary was built with git commit SHA injected via LDFLAGS during docker build
- Current git repository HEAD commit is known from build environment

**When**
- E2E test retrieves the expected git HEAD: `git rev-parse HEAD` (from source at build time)
- E2E test sends GET request to `/version` endpoint
- E2E test parses JSON response body

**Then**
- Response JSON contains `commit` field
- `commit` value is exactly 40 hexadecimal characters (valid git SHA format)
- `commit` value exactly matches the git repository HEAD SHA at build time
- Response validates against expected JSON schema

#### Scenario: FEATURE-A3 — E2E /version response is consistent across multiple requests

**Given**
- ubox-crosser server is running in Docker Compose and stable
- Test environment can issue rapid HTTP requests to running server

**When**
- E2E test sends 10 consecutive GET requests to `/version` endpoint
- All requests complete within acceptable time

**Then**
- All 10 responses return HTTP 200 OK
- All responses contain identical `commit` SHA value
- Response format is identical and stable across all requests
- No intermittent failures or timeouts occur

#### Scenario: FEATURE-A4 — E2E /version response latency meets SLA

**Given**
- ubox-crosser server is running in Docker Compose
- Test measures HTTP request/response latency from client perspective

**When**
- E2E test sends GET request to `/version`
- Response is received and end-to-end timing is recorded
- Multiple requests are tested (sample size ≥ 5)

**Then**
- All individual response times are < 100 milliseconds
- Median response time is < 50 milliseconds
- No request timeout occurs (default 5-second timeout is sufficient)
- Response latency remains consistent across multiple requests

#### Scenario: FEATURE-A5 — E2E Admin server with custom configured port

**Given**
- ubox-crosser can be configured with custom admin port via environment variable
- Docker Compose environment can inject ADMIN_SERVER_ADDR or similar configuration
- Test environment has access to multiple ports

**When**
- E2E test configures server with custom admin port: 9090
- Server initialization completes with new port configuration
- Test sends GET request to `localhost:9090/version`

**Then**
- Server successfully binds to custom port 9090
- GET request to `localhost:9090/version` returns HTTP 200 OK
- Response contains valid git SHA in `commit` field
- Default port 8080 is not used or returns failure as expected

#### Scenario: FEATURE-A6 — E2E /version response contains only necessary fields

**Given**
- ubox-crosser server is running in e2e environment
- Response should be minimal and secure (no extra metadata leakage)

**When**
- E2E test sends GET request to `/version`
- Response JSON is parsed and fields are examined

**Then**
- Response JSON object contains expected fields (at minimum: `commit`)
- Additional expected fields like `version`, `module`, `go_os`, `go_arch` may be present
- No sensitive fields like build flags, file paths, or environment variables are exposed
- No build system metadata, compiler versions, or internal configuration leaks
- Response is structurally minimal and follows security best practices

#### Scenario: FEATURE-A7 — E2E Docker multi-stage build preserves correct SHA

**Given**
- ubox-crosser uses Docker multi-stage build (one stage for compilation, one for runtime)
- Build stage compiles with git SHA injected via LDFLAGS
- Runtime image executes the compiled binary

**When**
- E2E test builds docker image from Dockerfile using docker compose
- Build completes with git SHA properly captured at compile stage
- E2E test starts container from built image
- E2E test queries `/version` endpoint

**Then**
- Binary in runtime image contains the correct git SHA
- Response `/version` returns SHA matching build-time `git rev-parse HEAD`
- Multi-stage build optimization does not lose or alter the SHA
- Feature works with production-ready containerized deployment

#### Scenario: FEATURE-A8 — E2E /version endpoint security with external requests

**Given**
- ubox-crosser server is running in Docker Compose
- Admin HTTP server is accessible from test-runner service (separate container)
- Network communication happens within docker compose network

**When**
- E2E test (from separate container) sends GET request to `/version`
- Various HTTP headers and parameters are included in requests
- Client acts as external service accessing the endpoint

**Then**
- Endpoint returns only git SHA in `commit` field (or expected minimal fields)
- No container internal paths or docker metadata in response
- No pod environment variables or kubernetes metadata exposed
- Network service isolation is maintained
- Response does not contain sensitive infrastructure information

### RELATED

- **Feature**: GET /version endpoint returns git commit SHA (implementation from REQ-rbac-1776859928)
- **Scope**: End-to-end integration validation via Docker Compose (real running server)
- **Build Integration**: Server binary injected with git SHA via LDFLAGS in Dockerfile
- **Related Spec**: REQ-rbac-1776859928/specs/version_endpoint_git_sha/ (acceptance spec from contract phase)
- **Contract Spec**: REQ-e2e-1776862572/specs/version_endpoint_contract/ (implementation contracts)
- **Test Infrastructure**: Docker Compose blackbox testing with actual running ubox-crosser binary
- **Integration Tests**: `tests/acceptance/version/e2e_test.go` for e2e acceptance test implementation
- **Git Integration**: Tests verify that returned SHA matches `git rev-parse HEAD` from build environment
