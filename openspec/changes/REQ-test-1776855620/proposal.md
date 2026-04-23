# REQ-test-1776855620: Add /healthz Endpoint with Service Uptime

## Summary
Add a new `/healthz` endpoint to the management HTTP server that returns the service uptime in addition to health status. This endpoint enables health check monitoring systems to understand both the current health state and how long the service has been running.

## Problem Statement
The current management server provides a `/health` endpoint that only returns basic status (`ok`), and a `/version` endpoint for build info. There is no way to determine service uptime from HTTP endpoints, which is important for:
- Monitoring and alerting systems
- SLA compliance verification
- Deployment verification (ensuring new instances are running)
- Load balancer health checks that need uptime information

## Solution Overview
Add a new `/healthz` endpoint that:
1. Returns HTTP 200 OK
2. Includes service uptime in seconds since server start
3. Includes health status
4. Uses JSON response format
5. Supports only GET requests (POST/PUT/DELETE return 405)

## Acceptance Criteria
- [ ] `/healthz` GET request returns 200 OK
- [ ] Response includes `uptime_seconds` field (integer)
- [ ] Response includes `status` field with value "ok"
- [ ] Response is valid JSON
- [ ] Response Content-Type is `application/json`
- [ ] POST/PUT/DELETE requests return 405 Method Not Allowed
- [ ] Response time is acceptable (< 100ms)
- [ ] Uptime value increases with time

## Target Branch
`stage/REQ-test-1776855620-dev`
