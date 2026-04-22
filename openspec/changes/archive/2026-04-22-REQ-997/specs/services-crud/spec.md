---
capability: services-crud
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: FEATURE-A5 — Create a new service via API

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
When the admin sends a POST request to /api/v1/services with a valid service payload
Then the response status code is 201
  and the response JSON field "code" equals 0
  and the response JSON field "data.name" matches the submitted service name
  and the response JSON field "data.id" is a positive integer
```

### Scenario: FEATURE-A6 — List all services

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
  and at least one service has been created
When the admin sends a GET request to /api/v1/services
Then the response status code is 200
  and the response JSON field "data" is a non-empty array
  and each item contains fields: id, name, key, method, address
```

### Scenario: FEATURE-A7 — Update an existing service

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
  and a service named "test-svc" exists
When the admin sends a PUT request to /api/v1/services/:id with updated address
Then the response status code is 200
  and the response JSON field "code" equals 0
  and the response JSON field "data.address" reflects the updated value
```

### Scenario: FEATURE-A8 — Delete a service

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
  and a service named "test-svc" exists
When the admin sends a DELETE request to /api/v1/services/:id
Then the response status code is 200
  and the response JSON field "code" equals 0
When the admin sends a GET request to /api/v1/services/:id
Then the response status code is 404
```
