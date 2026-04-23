# /healthz Endpoint - Health Check with Uptime

## Summary

Add a `/healthz` HTTP endpoint to the ubox-crosser service that returns health status and service uptime. This endpoint enables external monitoring tools and load balancers to verify service health with minimal latency (<100ms).

## Scenarios

### REQ-e2e-fix-1776923655-S1: Health Check Endpoint Returns 200 OK
**Given** the ubox-crosser service is running with HTTP health check enabled  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the endpoint responds with HTTP status code 200 OK  
**And** the response body is valid JSON  

### REQ-e2e-fix-1776923655-S2: Endpoint Reports Service Uptime in Seconds
**Given** the ubox-crosser service started at time T0  
**When** a client makes an HTTP GET request to `/healthz` at time T1  
**Then** the response includes an `uptime_seconds` field  
**And** the `uptime_seconds` value is approximately (T1 - T0)  

### REQ-e2e-fix-1776923655-S3: Uptime Increases Over Time
**Given** the ubox-crosser service is running  
**When** a client makes two sequential requests to `/healthz` with a 5-second delay  
**Then** the second response's `uptime_seconds` is 5 seconds greater than the first  
**And** the difference is within ±1 second tolerance  

### REQ-e2e-fix-1776923655-S4: Health Endpoint Available After Service Startup
**Given** the ubox-crosser service is starting up  
**When** the service completes initialization and is listening  
**Then** clients can connect to `/healthz`  
**And** the endpoint responds with uptime matching time since startup  

### REQ-e2e-fix-1776923655-S5: Response Format Includes Status and Timestamp
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the response body is JSON with fields:
- `status` field set to "healthy"
- `uptime_seconds` field with integer value ≥ 0
- `timestamp` field with current Unix timestamp (seconds)

### REQ-e2e-fix-1776923655-S6: Endpoint Handles Concurrent Requests
**Given** the ubox-crosser service is running  
**When** multiple clients make concurrent GET requests to `/healthz` (10+ parallel)  
**Then** all requests complete successfully  
**And** each response contains valid uptime value  
**And** no race conditions occur  

### REQ-e2e-fix-1776923655-S7: Service Uptime Resets on Restart
**Given** the ubox-crosser service is running with uptime U1  
**When** the service restarts  
**And** the service becomes ready to handle requests  
**Then** the new `uptime_seconds` value starts from ~0 seconds  
**And** uptime increments normally after restart  

## API Contract

### Endpoint: GET /healthz

#### Request
```
GET /healthz HTTP/1.1
Host: <service>:8080
```

#### Response (200 OK)
```json
{
  "status": "healthy",
  "uptime_seconds": 42,
  "timestamp": 1713897600
}
```

#### Response Headers
- `Content-Type: application/json`
- `Content-Length: <appropriate length>`

#### Constraints
- Response time: < 100ms
- No authentication required
- HTTP GET method only
- Listen port: 8080 (configurable)

## Testing Strategy

### Unit Tests (Development Phase)
- Uptime calculation logic
- Timestamp generation
- JSON serialization

### Contract Tests (THIS PHASE)
- REST API contract validation
- Response format compliance
- Uptime field accuracy
- Concurrent request handling
- Docker Compose integration test

### Acceptance Tests (Validation Phase)
- End-to-end service behavior
- Service restart scenarios
- Load balancer integration
