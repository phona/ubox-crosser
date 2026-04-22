# Healthz Endpoint Specification

## Feature: GET /healthz Endpoint with Service Uptime

### ADDED

#### Scenario: FEATURE-A1 — GET /healthz returns 200 with uptime information

**Given**
- ubox-crosser server is running
- HTTP admin server listening on default address `:8080`
- Server has been running for some time

**When**
- Client sends `GET /healthz` HTTP request to admin server

**Then**
- Response status code is `200 OK`
- Response content-type is `application/json`
- Response body contains JSON object with `uptime` field
- `uptime` field contains a positive integer representing seconds since server start

**Example Response**
```json
{
  "uptime": 3600
}
```

#### Scenario: FEATURE-A2 — GET /healthz returns uptime in seconds

**Given**
- ubox-crosser server is running
- Server started at known time T0

**When**
- Client sends `GET /healthz` HTTP request
- Current time is T0 + 120 seconds

**Then**
- Response status code is `200 OK`
- Response JSON contains `uptime` field with value approximately 120 (±2 seconds tolerance)
- `uptime` is a valid integer (seconds since server startup)

#### Scenario: FEATURE-A3 — GET /healthz uptime increases over time

**Given**
- ubox-crosser server is running and has been for N seconds
- Previous uptime reading was U1

**When**
- Client makes first request and gets uptime U1
- Waits 5 seconds
- Client makes second request and gets uptime U2

**Then**
- Response status code is `200 OK` for both requests
- U2 > U1 (uptime monotonically increases)
- U2 - U1 is approximately 5 seconds (±1 second tolerance)

#### Scenario: FEATURE-A4 — Non-GET method returns 405

**Given**
- ubox-crosser server is running
- HTTP admin server listening on port `:8080`

**When**
- Client sends `POST /healthz` (or PUT/DELETE/PATCH) HTTP request

**Then**
- Response status code is `405 Method Not Allowed`
- Response indicates only GET is allowed

#### Scenario: FEATURE-A5 — /healthz with custom admin address

**Given**
- ubox-crosser server is started with `--admin-addr :9090`
- Admin HTTP server listening on custom port `9090`
- Server has been running

**When**
- Client sends `GET localhost:9090/healthz`

**Then**
- Response status code is `200 OK`
- Response contains valid uptime JSON with positive integer
- Uptime matches server runtime

#### Scenario: FEATURE-A6 — /healthz immediately after server start

**Given**
- ubox-crosser server has just started (< 1 second)
- Admin HTTP server is ready to accept requests

**When**
- Client sends `GET /healthz` within first second of server startup

**Then**
- Response status code is `200 OK`
- Response JSON contains `uptime` field with value 0 or 1
- No errors in response

#### Scenario: FEATURE-A7 — /healthz response format validation

**Given**
- ubox-crosser server is running

**When**
- Client sends `GET /healthz` HTTP request
- Parses response as JSON

**Then**
- Response status code is `200 OK`
- Response content-type is `application/json` or `application/json; charset=utf-8`
- Response JSON is valid (parseable)
- Response contains exactly `uptime` field with integer value
- `uptime` is >= 0 (non-negative integer)

#### Scenario: FEATURE-A8 — End-to-end real server uptime tracking test

**Given**
- Real ubox-crosser server instance freshly started
- Server running with admin HTTP server enabled on default `:8080`
- Network connectivity to admin server

**When**
- E2E test script starts the server
- Records current timestamp as T_start
- Waits for admin server ready (health check)
- Makes GET /healthz request at time T_check
- Waits 3 seconds
- Makes second GET /healthz request at time T_check2

**Then**
- First /healthz response returns status `200 OK`
- First response `uptime` field is close to elapsed time from start
- Second /healthz response returns `200 OK`
- Second response `uptime` is greater than first
- Uptime difference matches approximate wait time (3 seconds ±1)
- Response latency is < 50ms for each request

### RELATED

- **Feature Implemented By**: `server/management.go` — HTTP mux with /healthz route
- **Handler**: `management.go:handleHealthz` — GET-only handler returning uptime
- **Server Startup**: `server/management.go:NewManagementServer()` — Records server start time
- **Unit Tests**: `server/management_test.go` — Unit test coverage for healthz handler
- **Integration Tests**: `tests/contract/healthz_*.go` or `tests/acceptance/healthz_*.go` — E2E test scenarios
