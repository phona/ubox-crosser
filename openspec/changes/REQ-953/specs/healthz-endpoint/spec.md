---
capability: healthz-endpoint
change_id: REQ-953
status: LOCKED
---

## ADDED

### Scenario: FEATURE-A1 — GET /api/healthz returns 200 with JSON health status

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a GET request to /api/healthz
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "status" equals "ok"
```

### Scenario: FEATURE-A2 — Docker healthcheck uses HTTP /api/healthz to determine service readiness

```gherkin
Given the docker-compose test environment is started
  and proxy-server healthcheck is configured to use HTTP GET /api/healthz
When the proxy-server finishes initializing
Then docker reports proxy-server as "healthy"
  and dependent services (client) start after proxy-server becomes healthy
```

### Scenario: FEATURE-A3 — /api/healthz is reachable from within docker-compose network

```gherkin
Given the docker-compose test environment is running
  and proxy-server is in healthy state
When the test-runner sends a GET request to http://proxy-server:8080/api/healthz
Then the response status code is 200
  and the response JSON field "status" equals "ok"
```

### Scenario: FEATURE-A4 — Non-GET methods on /api/healthz return appropriate response

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a POST request to /api/healthz
Then the response status code is 405
  or the response status code is 200
```

### Scenario: FEATURE-A5 — /api/healthz responds quickly under normal conditions

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a GET request to /api/healthz
Then the response is received within 500 milliseconds
  and the response status code is 200
```
