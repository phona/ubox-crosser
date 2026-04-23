# Tasks for /healthz Endpoint Acceptance Tests

## Stage: acceptance-tests

### Task 1: Validate Endpoint Returns 200 OK
- Verify `/healthz` returns HTTP 200 status code
- Confirm response body is valid JSON
- Validate response contains expected health information fields

### Task 2: Validate Uptime Seconds Reporting
- Verify `uptime_seconds` field exists in response
- Confirm value matches elapsed time since service startup
- Test accuracy within expected tolerance (±1 second)

### Task 3: Validate Uptime Increases Over Time
- Make sequential requests with time delays
- Verify uptime values increase appropriately
- Validate time differences match observed delays (within ±1 second)

### Task 4: Validate Endpoint Available After Startup
- Test endpoint accessibility immediately after service initialization
- Verify uptime values are reported correctly from startup
- Confirm service readiness via health check

### Task 5: Validate Response Format
- Confirm `status` field equals "healthy"
- Verify `uptime_seconds` is a non-negative integer
- Validate `timestamp` field contains current Unix timestamp
- Check JSON structure matches specification

### Task 6: Validate Concurrent Request Handling
- Test multiple concurrent requests to `/healthz`
- Verify all requests succeed
- Confirm uptime values are consistent and accurate
- Validate response times are acceptable (< 100ms)

### Task 7: Validate Uptime Reset on Service Restart
- Record initial uptime value U1
- Restart the service
- Verify uptime counter resets and starts from 0
- Confirm new uptime is significantly less than U1
- Validate uptime increments correctly after restart
