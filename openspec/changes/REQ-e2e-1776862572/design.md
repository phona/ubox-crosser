# Contract Testing Design: /version Endpoint

## Architecture

### Docker Compose Stack

The contract tests run against a real multi-service docker-compose environment:

```
┌─────────────────────────────────────────────────────┐
│                  Docker Compose                      │
├─────────────────────────────────────────────────────┤
│                                                       │
│  ┌─────────────┐    ┌──────────────┐                │
│  │ test-runner │───→│ proxy-server │:8080 (admin)   │
│  │ (tests)     │    │ (ubox-crosser)                │
│  └─────────────┘    └──────────────┘                │
│                         ↓                            │
│                    ADMIN_SERVER_ADDR                │
│                    =proxy-server:8080               │
│                                                       │
└─────────────────────────────────────────────────────┘
```

### Environment Injection

```yaml
Services:
  test-runner:
    environment:
      - ADMIN_SERVER_ADDR=proxy-server:8080       # Primary admin server
      - CUSTOM_ADMIN_ADDR=<optional>              # For custom port testing
      - GOCOVERDIR=/coverage                      # Coverage collection
```

### Test Execution Flow

1. **Build Phase**: Compile ubox-crosser binary with git SHA injected
   - LDFLAGS: `-X main.Version=$(git rev-parse HEAD)`
   - Binary embedded in proxy-server image

2. **Startup Phase**: 
   - proxy-server starts HTTP admin server on :8080
   - test-runner waits for services ready
   - All containers have shared access to coverage directory

3. **Test Phase**:
   - test-runner runs: `go test -tags integration -timeout=120s ./tests/integration/...`
   - Tests connect to proxy-server:8080/version
   - Each test validates contract scenarios

4. **Verification**:
   - Exit code 0 = all tests pass (contracts satisfied)
   - Exit code non-0 = contract violation (tests fail)

## Contract Definitions

### Contract 1: HTTP Response Format

**Given**: proxy-server running with admin HTTP on :8080
**When**: GET /version request sent
**Then**:
  - Status: 200 OK
  - Content-Type: application/json
  - Body: JSON object with [version, module, go_os, go_arch, commit] fields

### Contract 2: Git SHA Validity

**Given**: Binary compiled with git commit SHA
**When**: GET /version response received
**Then**:
  - `commit` field contains 40-character hexadecimal string
  - Regex match: `^[0-9a-f]{40}$`

### Contract 3: Commit Immutability

**Given**: Server running without restart
**When**: Multiple GET /version requests in sequence
**Then**:
  - All responses contain identical `commit` value
  - SHA does not change during server lifetime

### Contract 4: Security (No Leakage)

**Given**: Server responding to HTTP requests
**When**: GET /version response examined
**Then**:
  - Response contains only expected fields
  - No sensitive patterns (PASSWORD, TOKEN, SECRET, LDFLAGS, buildinfo, paths)
  - No environment variables or build flags leaked

### Contract 5: Custom Admin Port

**Given**: Optional CUSTOM_ADMIN_ADDR environment variable set
**When**: GET request to custom port /version
**Then**:
  - Status: 200 OK
  - Response format matches Contract 1
  - Commit SHA valid per Contract 2

## Test File Organization

```
tests/
├── docker-compose.yml              # Stack definition (LOCKED)
├── Dockerfile.test                 # Build image (LOCKED)
├── integration/
│   ├── version_endpoint_test.go    # Version contract tests (LOCKED)
│   ├── tunnel_test.go              # Other integration tests
│   └── helpers.go                  # getEnv() utility
└── configs/                        # Service configurations
    ├── server.json
    ├── client.json
    └── auth_server.json
```

## Test Tags and Build Constraints

```go
//go:build integration    // Only run with: go test -tags integration

// Tests require:
// - ADMIN_SERVER_ADDR environment variable (default: localhost:8080)
// - Running docker-compose stack OR real server on that address
// - HTTP connectivity to the admin server
```

## CI/CD Integration

The contract tests are part of the **ci-int** (continuous integration) job:

```bash
# Local validation
docker compose -f tests/docker-compose.yml up --build --exit-code-from test-runner

# This ensures:
# - Docker images are built fresh
# - All services start correctly
# - Tests have real server to validate against
# - Coverage data collected
# - Exit code reflects test results
```

## Failure Modes

| Failure Mode | Root Cause | Resolution |
|---|---|---|
| Status code != 200 | HTTP handler missing or broken | Dev implements /version handler |
| Content-Type mismatch | Handler not returning JSON | Dev adds application/json header |
| Commit field empty | LDFLAGS not injected at build | Dev updates Makefile/Dockerfile |
| Commit not valid hex | Non-SHA value injected | Dev fixes build injection |
| SHA changes mid-test | Server modifying response | Dev makes field immutable |
| Sensitive data leaked | Handler returning extra fields | Dev sanitizes response |

## LOCKED Zones

The contract-spec stage defines these as LOCKED (dev-agent cannot modify):

- ✗ `tests/docker-compose.yml` - Stack definition is part of contract
- ✗ `tests/integration/version_endpoint_test.go` - Contract test scenarios
- ✗ `tests/Dockerfile.test` - Build process for contract environment

Dev-agent can modify:
- ✓ `cmd/` - Implementation code
- ✓ `Makefile` - Build flags (must inject COMMIT)
- ✓ `Dockerfile` - Build (must preserve LDFLAGS injection)

## Testing Best Practices

1. **Never Use httptest.NewServer**: Tests must validate real behavior
2. **Environment Variable Driven**: Use ADMIN_SERVER_ADDR, CUSTOM_ADMIN_ADDR
3. **Timeout Awareness**: Real servers may be slower than mocks
4. **Coverage Collection**: GOCOVERDIR enables coverage reporting
5. **Dependency Wait**: Use healthchecks and depends_on for service readiness
