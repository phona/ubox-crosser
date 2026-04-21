---
capability: version-endpoint
change_id: REQ-969
status: LOCKED
---

## ADDED

### Scenario: REQ-969-S1 — GET /api/version returns 200 with JSON commit field

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /api/version
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "commit" is a non-empty string
```

### Scenario: REQ-969-S2 — GET /api/version response schema validation

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /api/version
Then the response status code is 200
  and the response JSON contains required field: commit (string)
  and the response JSON contains no extra unexpected fields
```

### Scenario: REQ-969-S3 — Non-GET methods return 405 Method Not Allowed

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /api/version
Then the response status code is 405
When the user sends a PUT request to /api/version
Then the response status code is 405
When the user sends a DELETE request to /api/version
Then the response status code is 405
```

### Scenario: FEATURE-A1 — GET /api/version returns 200 with JSON containing commit hash

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a GET request to /api/version
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "commit" is a non-empty string
```

### Scenario: FEATURE-A2 — commit field contains a valid 40-character hex git hash or "unknown"

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a GET request to /api/version
Then the response status code is 200
  and the response JSON field "commit" matches pattern "^([0-9a-f]{40}|unknown)$"
```

### Scenario: FEATURE-A3 — Non-GET methods on /api/version return 405 Method Not Allowed

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a POST request to /api/version
Then the response status code is 405
```

### Scenario: FEATURE-A4 — /api/version is reachable from within docker-compose network

```gherkin
Given the docker-compose test environment is running
  and proxy-server is in healthy state
When the test-runner sends a GET request to http://proxy-server:8080/api/version
Then the response status code is 200
  and the response JSON field "commit" is a non-empty string
```

### Scenario: FEATURE-A5 — /api/version responds quickly under normal conditions

```gherkin
Given the proxy-server is running with admin HTTP server on port 8080
When the user sends a GET request to /api/version
Then the response is received within 500 milliseconds
  and the response status code is 200
```
