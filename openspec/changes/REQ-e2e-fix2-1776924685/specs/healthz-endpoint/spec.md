# /healthz Endpoint - Health Check with Uptime

## Summary
Add a `/healthz` HTTP endpoint to the ubox-crosser service that returns the current health status and service uptime information. This endpoint allows external monitoring tools and load balancers to check if the service is running and how long it has been active.

## Scenarios

### REQ-e2e-fix2-1776924685-S1: Health Check Endpoint Returns 200 OK
**Given** the ubox-crosser service is running with HTTP health check enabled  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the endpoint responds with status code 200 OK  
**And** the response body contains valid JSON with service health information  

### REQ-e2e-fix2-1776924685-S2: Endpoint Reports Service Uptime in Seconds
**Given** the ubox-crosser service started at time T0  
**When** a client makes an HTTP GET request to `/healthz` at time T1  
**Then** the response includes an `uptime_seconds` field  
**And** the `uptime_seconds` value equals T1 - T0 (the service's running duration)  

### REQ-e2e-fix2-1776924685-S3: Uptime Increases Over Time
**Given** the ubox-crosser service is running  
**When** a client makes two sequential requests to `/healthz` with a 5-second delay between them  
**Then** the second response's `uptime_seconds` is approximately 5 seconds greater than the first  
**And** the difference is within ±1 second (accounting for processing latency)  

### REQ-e2e-fix2-1776924685-S4: Response Format Includes Required Fields
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the response body is JSON with at least:
- `status` field set to "healthy"
- `uptime_seconds` field with an integer value ≥ 0
- `timestamp` field with the current Unix timestamp  

### REQ-e2e-fix2-1776924685-S5: Endpoint Works Concurrently
**Given** the ubox-crosser service is running  
**When** multiple clients make concurrent HTTP GET requests to `/healthz`  
**Then** all requests complete successfully  
**And** each response contains a valid uptime value  
**And** the uptime values differ only by the time it takes to process the requests  

## Additional Notes

- The `/healthz` endpoint should be accessible via HTTP GET requests
- The endpoint should not require authentication
- Response time should be minimal (< 100ms) to be suitable for load balancer health checks
- The uptime counter should be based on when the service fully initializes, not when the binary starts
