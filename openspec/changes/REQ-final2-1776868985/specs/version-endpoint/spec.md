# Acceptance Spec: /version Endpoint Git SHA

## Overview

The `/version` endpoint should return the git commit SHA of the current build. This spec verifies the endpoint works correctly across different build scenarios and returns properly formatted responses.

## Scenarios

### FEATURE-A1: GET /version returns git SHA in normal build

```gherkin
Given the ubox-crosser binary is built with git SHA injected
When a client sends GET request to /version
Then the response status code is 200
And the response body is JSON with format {"sha": "<sha>"}
And the "sha" field contains a 40-character hexadecimal string
```

**Test Steps:**
1. Build binary with `make build` (which includes -ldflags injection)
2. Start the server
3. Send GET /version request
4. Verify HTTP 200 response
5. Parse JSON response
6. Validate sha field is exactly 40 hex chars

**Expected Result:** Response is `{"sha": "<full-git-sha>"}` where sha matches the HEAD commit SHA at build time.

---

### FEATURE-A2: GET /version returns consistent SHA across container builds

```gherkin
Given a Docker image is built from the repository
And the Dockerfile includes the build step with proper git SHA injection
When a container is started and requests /version
Then the response sha field matches the git commit SHA from the build context
And the sha value is consistent across multiple container startups
```

**Test Steps:**
1. Note the current HEAD SHA: `git rev-parse HEAD`
2. Build Docker image: `docker build -t ubox-crosser:test .`
3. Start container and call /version
4. Verify returned SHA matches the noted HEAD SHA
5. Stop and restart container
6. Verify SHA is still the same (not dynamic)

**Expected Result:** SHA is the same across container restarts, proving it was injected at build time.

---

### FEATURE-A3: GET /version handles development environment (unknown SHA)

```gherkin
Given the binary is built without proper git SHA injection
When a client sends GET request to /version
Then the response status code is 200
And the "sha" field value is "unknown" (fallback default)
```

**Test Steps:**
1. Build binary without ldflags: `go build -o ./crosser ./cmd/main.go`
2. Start server
3. Call /version
4. Verify response contains `{"sha": "unknown"}`

**Expected Result:** Returns "unknown" gracefully instead of error, allowing development builds to work.

---

### FEATURE-A4: /version endpoint is idempotent and fast

```gherkin
Given the server is running with injected git SHA
When a client sends multiple rapid GET requests to /version
Then each response returns 200
And all response bodies are identical
And response time is < 10ms (no I/O operations)
```

**Test Steps:**
1. Build and start server
2. Send 100 sequential requests to /version
3. Measure response time for each
4. Verify all responses are identical
5. Verify no response takes > 10ms

**Expected Result:** All requests return the exact same response body, confirming no runtime git operations.

---

### FEATURE-A5: /version endpoint is not broken by empty or invalid git state

```gherkin
Given a clone that might not have a proper .git directory
When the binary is built with failed git revision resolution
Then the /version endpoint still returns 200 with a fallback value
And the endpoint does not crash or hang
```

**Test Steps:**
1. Clone repo in a state without `.git` (e.g., `git archive` source)
2. Attempt build with git revision that fails
3. Build should either:
   - Inject "unknown" as fallback
   - Use a pre-set default value
4. Verify binary still runs and /version returns 200

**Expected Result:** Graceful handling of missing git info, returning either "unknown" or a build-time constant.

---

## Integration Test Coverage

The acceptance spec is verified by integration tests in `tests/integration/version_test.go`, which:

1. **Test HTTP Response Format**
   - Verify status code 200
   - Verify JSON parsing
   - Verify "sha" field exists and is a string

2. **Test SHA Validity**
   - Verify SHA is 40 hex characters (or "unknown")
   - Verify no other unexpected fields

3. **Test Consistency**
   - Verify multiple calls return identical responses
   - Verify performance (should be < 10ms)

4. **Test Fallback Behavior**
   - Verify "unknown" is returned when SHA is not injected

## Docker Compose Integration

The `docker-compose.yml` includes the `/version` endpoint in the proxy-server healthcheck (optional, for reference):
- A separate request can be added to validate version endpoint is healthy
- This ensures the endpoint is available before integration tests run

## Acceptance Criteria

✅ All scenarios pass with the implemented /version endpoint
✅ Integration tests execute successfully
✅ SHA value is correctly injected at build time (not runtime)
✅ Fallback behavior works when SHA is unavailable
✅ Response format is consistent across builds and environments
