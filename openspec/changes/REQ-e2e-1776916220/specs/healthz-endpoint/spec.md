# /healthz Endpoint - Health Check with Uptime

## Summary
Add a `/healthz` HTTP endpoint to the ubox-crosser service that returns the current health status and service uptime information. This endpoint allows external monitoring tools and load balancers to check if the service is running and how long it has been active.

## Scenarios

### FEATURE-A1: Health Check Endpoint Returns 200 OK
**Given** the ubox-crosser service is running with HTTP health check enabled  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the endpoint responds with status code 200 OK  
**And** the response body contains valid JSON with service health information  

### FEATURE-A2: Endpoint Reports Service Uptime in Seconds
**Given** the ubox-crosser service started at time T0  
**When** a client makes an HTTP GET request to `/healthz` at time T1  
**Then** the response includes an `uptime_seconds` field  
**And** the `uptime_seconds` value equals T1 - T0 (the service's running duration)  

### FEATURE-A3: Uptime Increases Over Time
**Given** the ubox-crosser service is running  
**When** a client makes two sequential requests to `/healthz` with a 5-second delay between them  
**Then** the second response's `uptime_seconds` is approximately 5 seconds greater than the first  
**And** the difference is within ±1 second (accounting for processing latency)  

### FEATURE-A4: Health Endpoint Available After Service Startup
**Given** the ubox-crosser service is starting up  
**When** the service has completed initialization and is listening on the configured health check port  
**Then** clients can successfully connect to the health check endpoint  
**And** the endpoint responds with uptime matching the time since service startup  

### FEATURE-A5: Response Format Includes Status Field
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/healthz`  
**Then** the response body is JSON with at least:
- `status` field set to "healthy"
- `uptime_seconds` field with an integer value ≥ 0
- `timestamp` field with the current Unix timestamp  

### FEATURE-A6: Endpoint Works Concurrently
**Given** the ubox-crosser service is running  
**When** multiple clients make concurrent HTTP GET requests to `/healthz`  
**Then** all requests complete successfully  
**And** each response contains a valid uptime value  
**And** the uptime values differ only by the time it takes to process the requests  

### FEATURE-A7: Service Uptime Resets on Restart
**Given** the ubox-crosser service is running with uptime U1  
**When** the service is restarted  
**And** the service becomes ready to handle requests  
**Then** the `uptime_seconds` value starts from 0 and increments upward  
**And** the new uptime is less than the previous uptime before restart  

## Additional Notes

- The `/healthz` endpoint should be accessible via HTTP GET requests
- The endpoint should not require authentication
- Response time should be minimal (< 100ms) to be suitable for load balancer health checks
- The uptime counter should be based on when the service fully initializes, not when the binary starts
