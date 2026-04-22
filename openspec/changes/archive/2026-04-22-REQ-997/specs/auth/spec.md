---
capability: auth
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: REQ-997-S1 — POST /api/v1/auth/login with valid credentials returns JWT

```gherkin
Given the crosser-api server is running
  and a default admin user exists
When the client sends POST /api/v1/auth/login with valid username and password
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "code" is 0
  and the response JSON field "data.token" is a non-empty string
  and the response JSON field "data.expires_at" is a valid RFC3339 timestamp
```

### Scenario: REQ-997-S2 — POST /api/v1/auth/login with invalid credentials returns 401

```gherkin
Given the crosser-api server is running
When the client sends POST /api/v1/auth/login with wrong password
Then the response status code is 401
  and the response JSON field "code" is 1001
  and the response JSON field "message" is a non-empty string
```

### Scenario: REQ-997-S3 — Authenticated endpoints reject requests without valid JWT

```gherkin
Given the crosser-api server is running
When the client sends GET /api/v1/services without Authorization header
Then the response status code is 401
  and the response JSON field "code" is 1003
When the client sends GET /api/v1/services with an expired JWT
Then the response status code is 401
  and the response JSON field "code" is 1002
```
