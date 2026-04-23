# Design: /healthz Endpoint Implementation

## Architecture

### HTTP Server Component

The ProxyServer initializes an additional HTTP listener alongside its existing TCP message listeners.

```
ProxyServer
├── TCP Listeners (existing message protocol)
│   ├── Port 7000: Encrypted Shadowsocks connections
│   └── Port 7001+: Additional configured ports
│
└── HTTP Listener (new)
    └── Port 8080: RESTful health check
        └── GET /healthz
```

### Data Flow

1. **Startup:** ProxyServer records `startTime := time.Now()` during initialization
2. **Health Check Request:** External client sends HTTP GET to `/healthz`
3. **Uptime Calculation:** Endpoint calculates `uptime := time.Now() - startTime`
4. **Response:** JSON response with status and uptime fields

## Implementation Details

### Configuration

**Environment Variables:**
- `PROXY_HTTP_PORT`: HTTP listener port (default: `8080`)
- `PROXY_HTTP_ADDR`: Full HTTP address (default: `0.0.0.0:8080`)

### Response Format

```json
{
  "status": "healthy",
  "uptime": {
    "seconds": 3661,
    "duration": "1h1m1s"
  }
}
```

### Error Handling

- If HTTP listener fails to bind, ProxyServer logs warning but continues (non-fatal)
- All network errors in handlers return 500 Internal Server Error
- Malformed requests return 400 Bad Request

### Concurrency

- HTTP listener runs in dedicated goroutine
- Uptime calculation is thread-safe (reads `startTime` only, no mutation)
- No contention with existing TCP message processing

## Contract Specification

See `specs/healthz/spec.md` for OpenAPI contract and test scenarios.

### Scenarios Tested

| Scenario ID | Test Name | Purpose |
|------------|-----------|---------|
| REQ-e2e-1776916220-S1 | TestHealthzS1 | Response format and status code validation |
| REQ-e2e-1776916220-S2 | TestHealthzS2 | Uptime progression over time |
| REQ-e2e-1776916220-S3 | TestHealthzS3 | Duration string format validity |

## Testing Strategy

### Contract Tests

Located in `tests/contract/healthz_test.go`, these tests:
- Run in Docker Compose environment with real ProxyServer
- Validate HTTP response code, headers, and JSON structure
- Verify uptime increases over time
- Check duration string format

### Integration Environment

Docker Compose provides:
- ProxyServer with HTTP listener enabled
- Environment variable `PROXY_HTTP_ADDR=proxy-server:8080`
- Test container with HTTP client

## Deployment Considerations

### Backwards Compatibility

- New HTTP listener is independent of existing TCP protocol
- Zero impact on client connection logic
- Non-breaking change

### Monitoring & Observability

The `/healthz` endpoint enables:
- Kubernetes liveness/readiness probes
- Load balancer health checks
- Container orchestration monitoring

## Security Considerations

- No authentication required (health checks are public by design)
- No sensitive data exposed (only uptime and status)
- Standard endpoint path `/healthz` follows convention
