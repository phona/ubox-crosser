# /healthz Endpoint - Health Check with Service Uptime

## Summary
Add a `/healthz` HTTP endpoint to the ubox-crosser service that provides real-time health status and service uptime tracking. This endpoint enables external monitoring systems, load balancers, and Kubernetes probes to verify service health and availability.

## Scenarios

### FEATURE-A1: Health Check Endpoint Returns 200 OK
**Given** the ubox-crosser service is running with health check enabled  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the endpoint responds with HTTP status code 200 OK  
**And** the response body contains valid JSON with health information  

### FEATURE-A2: Endpoint Reports Service Uptime in Seconds
**Given** the ubox-crosser service started at time T0  
**When** a client makes an HTTP GET request to `/healthz` at time T1  
**Then** the response includes an `uptime_seconds` field  
**And** the value represents T1 - T0 (total service running duration in seconds)  

### FEATURE-A3: Uptime Increases Over Time
**Given** the ubox-crosser service is running  
**When** a client makes two sequential requests to `/healthz` with a 5-second delay  
**Then** the second response's `uptime_seconds` is approximately 5 seconds higher than the first  
**And** the difference is within ±1 second to account for processing time  

### FEATURE-A4: Health Endpoint Responds Immediately After Startup
**Given** the ubox-crosser service completes initialization  
**When** the service begins listening on the health check port  
**Then** clients can immediately connect and query the `/healthz` endpoint  
**And** the endpoint returns uptime_seconds matching the time since full service readiness  

### FEATURE-A5: Response Contains Required JSON Fields
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the response body includes these JSON fields:
- `status`: string value "healthy"
- `uptime_seconds`: non-negative integer
- `timestamp`: Unix timestamp of the response  

### FEATURE-A6: Handles Multiple Concurrent Requests
**Given** the ubox-crosser service is running  
**When** 10+ clients simultaneously make HTTP GET requests to `/healthz`  
**Then** all requests complete successfully with status 200  
**And** each response contains a valid uptime_seconds value  
**And** response times are consistent (no requests timeout or fail)  

### FEATURE-A7: Uptime Counter Resets on Service Restart
**Given** the ubox-crosser service is running with uptime U1 seconds  
**When** the service is stopped and restarted  
**And** the service completes initialization again  
**Then** the `uptime_seconds` value begins at 0 and increments upward  
**And** new uptime is less than the previous uptime before restart  

## Acceptance Criteria
- All 7 user-facing scenarios pass with a working implementation
- Endpoint response time is < 100ms for load balancer health checks
- No authentication required to access `/healthz`
- HTTP method accepted: GET (OPTIONS may also be supported)
- Response is in valid JSON format
- Uptime calculation is based on service readiness, not binary startup time

## Notes
- Scenarios are written from the user/operator perspective
- Endpoint must be available immediately after service initialization
- Uptime precision should be at the second level
- Endpoint is suitable for Kubernetes liveness and readiness probes
