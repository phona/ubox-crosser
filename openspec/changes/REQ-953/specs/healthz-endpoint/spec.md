---
capability: healthz-endpoint
change_id: REQ-953
status: LOCKED
---

## ADDED

### Scenario: REQ-953-S1 — GET /api/healthz returns 200 with JSON health status

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /api/healthz
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "status" equals "ok"
```

### Scenario: REQ-953-S2 — GET /api/healthz response schema validation

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /api/healthz
Then the response status code is 200
  and the response JSON contains required field: status (string, value "ok")
  and the response JSON contains no extra unexpected fields
```

### Scenario: REQ-953-S3 — Non-GET methods return 405 Method Not Allowed

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /api/healthz
Then the response status code is 405
When the user sends a PUT request to /api/healthz
Then the response status code is 405
When the user sends a DELETE request to /api/healthz
Then the response status code is 405
```
