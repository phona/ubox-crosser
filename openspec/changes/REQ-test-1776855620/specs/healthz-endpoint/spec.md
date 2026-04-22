# Specification: /healthz Endpoint with Uptime

## Endpoint Definition

**HTTP Method:** GET  
**Path:** `/healthz`  
**Host:** Management Server (default port 8888)  

## Response Schema

### 2xx Success Response
**Status Code:** `200 OK`  
**Content-Type:** `application/json`

```json
{
  "status": "ok",
  "uptime_seconds": <integer>
}
```

**Fields:**
- `status` (string): Current service status. Value: `"ok"` when healthy.
- `uptime_seconds` (integer): Seconds elapsed since service started.

### 4xx Error Responses
- **405 Method Not Allowed:** Returned for POST, PUT, DELETE, PATCH, etc.
- **404 Not Found:** Endpoint not implemented (pre-implementation)

## Scenarios

### Scenario: REQ-test-1776855620-S1
**Title:** GET /healthz returns 200 with status ok

**Given:** Management server is running  
**When:** Client sends GET request to `/healthz`  
**Then:**
- Response status code is 200
- Response Content-Type header is `application/json`
- Response body contains JSON with `status` field
- `status` field value is `"ok"`

### Scenario: REQ-test-1776855620-S2
**Title:** GET /healthz includes uptime_seconds field

**Given:** Management server has been running for at least 1 second  
**When:** Client sends GET request to `/healthz`  
**Then:**
- Response contains `uptime_seconds` field
- `uptime_seconds` is an integer
- `uptime_seconds` value is >= 1
- `uptime_seconds` value increases between consecutive calls

### Scenario: REQ-test-1776855620-S3
**Title:** POST /healthz returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends POST request to `/healthz` with empty body  
**Then:**
- Response status code is 405
- Response indicates Method Not Allowed

### Scenario: REQ-test-1776855620-S4
**Title:** PUT /healthz returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends PUT request to `/healthz`  
**Then:**
- Response status code is 405

### Scenario: REQ-test-1776855620-S5
**Title:** DELETE /healthz returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends DELETE request to `/healthz`  
**Then:**
- Response status code is 405

### Scenario: REQ-test-1776855620-S6
**Title:** /healthz response time is acceptable

**Given:** Management server is running  
**When:** Client sends GET request to `/healthz`  
**Then:**
- Response time is less than 100 milliseconds
- Uptime calculation does not cause noticeable latency

## Edge Cases

### Edge Case: EC1 - Initial Startup
**Condition:** Server just started (< 1 second)  
**Expected:** `uptime_seconds` is 0

### Edge Case: EC2 - Long Running Service
**Condition:** Server has been running for days  
**Expected:** `uptime_seconds` reflects accurate elapsed time

### Edge Case: EC3 - Rapid Requests
**Condition:** Multiple requests within milliseconds  
**Expected:** `uptime_seconds` may be identical or increment by 1, showing monotonic increase

## Response Examples

### Example 1: Server running for ~5 seconds
```json
{
  "status": "ok",
  "uptime_seconds": 5
}
```

### Example 2: Server running for ~1 hour
```json
{
  "status": "ok",
  "uptime_seconds": 3600
}
```

## OpenAPI Definition

```yaml
paths:
  /healthz:
    get:
      summary: Health check with uptime
      operationId: getHealthz
      responses:
        '200':
          description: Service is healthy
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    enum: [ok]
                  uptime_seconds:
                    type: integer
                    minimum: 0
                required:
                  - status
                  - uptime_seconds
        '405':
          description: Method not allowed
```
