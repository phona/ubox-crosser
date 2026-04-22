# Version Endpoint Git SHA Specification

## Feature: GET /version Endpoint Returns Git SHA

### ADDED

#### Scenario: FEATURE-A1 — GET /version returns commit field with git SHA

**Given**
- ubox-crosser server is built from a git repository
- Server is running with HTTP admin server enabled
- Binary is compiled with git commit SHA injected via build flags

**When**
- Client sends `GET /version` HTTP request to admin server (default `:8080`)

**Then**
- Response status code is `200 OK`
- Response content-type is `application/json`
- Response JSON object contains `commit` field
- `commit` field value is a valid 40-character hexadecimal SHA (git commit hash)

**Example Response**
```json
{
  "commit": "23e96fb5a7d0c1a2b3c4d5e6f7g8h9i0j1k2l3m4"
}
```

#### Scenario: FEATURE-A2 — /version with custom admin server port

**Given**
- ubox-crosser server started with `--admin-addr :9090`
- Server built with git commit SHA in binary
- Admin HTTP server listening on custom port `9090`

**When**
- Client sends `GET /version` to `localhost:9090/version`

**Then**
- Response status code is `200 OK`
- Response JSON `commit` field contains valid 40-character git SHA
- SHA matches the actual build commit

#### Scenario: FEATURE-A3 — /version responds with correct SHA after rebuild

**Given**
- ubox-crosser binary is rebuilt after a new git commit
- Previous binary returned git SHA `abc123def456...`
- New commit creates new binary with different SHA `fed654cba321...`

**When**
- Old binary is stopped and new binary is started
- Client sends `GET /version` to both instances

**Then**
- Old binary response contains original SHA
- New binary response contains new SHA
- Both SHAs are valid 40-character git hashes
- Each response is consistent across multiple requests

#### Scenario: FEATURE-A4 — /version with missing commit SHA during build

**Given**
- ubox-crosser server is compiled without git commit SHA injection
- Build process did not inject COMMIT environment variable via LDFLAGS

**When**
- Client sends `GET /version` HTTP request

**Then**
- Response status code is `200 OK`
- Response JSON contains `commit` field
- `commit` value is either:
  - Empty string `""`, OR
  - Default value `"unknown"`, OR
  - Null value `null`
- Behavior is consistent and documented

#### Scenario: FEATURE-A5 — /version SHA is immutable during server lifetime

**Given**
- ubox-crosser server is running
- Server has been running for some time with stable git SHA

**When**
- Client sends multiple `GET /version` requests at intervals
- Server continues running without restart between requests

**Then**
- All responses contain identical `commit` SHA value
- Response format is consistent
- No unexpected changes to commit field occur

#### Scenario: FEATURE-A6 — End-to-end integration test with git SHA verification

**Given**
- Real ubox-crosser server binary built with actual git repository commit
- Server is running with admin HTTP server enabled
- Test environment has access to git repository metadata

**When**
- E2E test starts the server
- Validates connectivity to admin server on expected port
- Sends `GET /version` request
- Compares returned SHA with actual `git rev-parse HEAD` from repo

**Then**
- Response status code is `200 OK`
- Response contains valid JSON with `commit` field
- Response `commit` value is 40-character hexadecimal string
- Response SHA exactly matches actual git HEAD commit hash
- Response time is within acceptable latency bounds (< 100ms)
- Test passes with reproducible results

#### Scenario: FEATURE-A7 — /version endpoint security - no sensitive data leakage

**Given**
- ubox-crosser server with /version endpoint enabled
- Server is exposed to network clients

**When**
- Clients with no authentication send requests to `/version`
- Various HTTP headers and parameters are sent

**Then**
- Endpoint returns only git SHA in `commit` field
- No build flags, internal paths, or other metadata exposed
- Response does not contain environment variables
- Response does not contain system information beyond git commit

### RELATED

- **Endpoint**: `GET /version` on HTTP admin server
- **Build Integration**: Makefile/Dockerfile must inject `COMMIT` via LDFLAGS
- **Related Feature**: REQ-haiku-1776846225 — Initial /version endpoint implementation
- **Integration Tests**: Docker Compose-based blackbox tests in `tests/acceptance/version/`
- **Git Integration**: Binary should be built with `-ldflags="-X main.Version=$(git rev-parse HEAD)"`
