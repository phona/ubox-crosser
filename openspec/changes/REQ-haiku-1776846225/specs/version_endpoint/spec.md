# Version Endpoint Specification

## Feature: GET /version Endpoint for Build Information

### ADDED

#### Scenario: FEATURE-A1 — GET /version returns 200 with version info

**Given**
- ubox-crosser server is running
- HTTP admin server listening on default address `:8080`
- Server compiled with commit hash in build metadata

**When**
- Client sends `GET /version` HTTP request to admin server

**Then**
- Response status code is `200 OK`
- Response content-type is `application/json`
- Response body contains JSON object with `commit` field
- `commit` field contains the git commit hash (40 hex chars)

**Example Response**
```json
{
  "commit": "6d6fa485f0d0c1a2b3c4d5e6f7g8h9i0j1k2l3m4"
}
```

#### Scenario: FEATURE-A2 — GET /version with empty/unknown commit

**Given**
- ubox-crosser server is running
- Server compiled without commit hash (COMMIT var not injected)

**When**
- Client sends `GET /version` HTTP request

**Then**
- Response status code is `200 OK`
- Response JSON contains `commit: "unknown"` as default value

#### Scenario: FEATURE-A3 — Non-GET method returns 405

**Given**
- ubox-crosser server is running
- HTTP admin server listening on port `:8080`

**When**
- Client sends `POST /version` (or PUT/DELETE) HTTP request

**Then**
- Response status code is `405 Method Not Allowed`
- Response indicates only GET is allowed

#### Scenario: FEATURE-A4 — /version requires admin server to be running

**Given**
- ubox-crosser server is started without `--admin-addr` flag
- Default admin HTTP server should be on `:8080`

**When**
- Client connects to `localhost:8080/version`

**Then**
- Connection succeeds
- GET /version returns version information

#### Scenario: FEATURE-A5 — /version with custom admin address

**Given**
- ubox-crosser server is started with `--admin-addr :9090`
- Admin HTTP server listening on custom port `9090`

**When**
- Client sends `GET localhost:9090/version`

**Then**
- Response status code is `200 OK`
- Response contains valid version JSON

#### Scenario: FEATURE-A6 — End-to-end real link test with actual server

**Given**
- Real ubox-crosser server instance built with git commit hash
- Server running with admin HTTP server enabled
- Network connectivity to admin server

**When**
- E2E test script starts the server
- Waits for admin server ready (health check)
- Sends GET /version request
- Validates response JSON structure

**Then**
- Version endpoint responds within 100ms
- Response contains valid commit hash format
- Response is reproducible across multiple requests
- Commit hash matches the actual build commit

### RELATED

- **Feature Implemented By**: `server/admin.go` — HTTP mux with /version route
- **Handler**: `version/handler.go` — GET-only handler returning JSON with commit
- **Build Integration**: Makefile/Dockerfile — LDFLAGS injection of git commit
- **Unit Tests**: `version/handler_test.go` — Unit test coverage for handler
