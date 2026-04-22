---
capability: proxy-registration
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: FEATURE-A9 — Proxy instance registers successfully with valid service token

```gherkin
Given the crosser-api server is running
  and a service named "test-svc" exists
  and a valid service token has been generated for "test-svc"
When a proxy sends a POST request to /api/v1/proxy/register with the service token in X-Proxy-Token header
Then the response status code is 200
  and the response JSON field "code" equals 0
  and the response JSON field "data.instance_id" is a non-empty string
```

### Scenario: FEATURE-A10 — Proxy heartbeat updates instance status to online

```gherkin
Given the crosser-api server is running
  and a proxy instance has been registered for service "test-svc"
When the proxy sends a POST request to /api/v1/proxy/heartbeat with its instance_id and token
Then the response status code is 200
  and the response JSON field "code" equals 0
```

### Scenario: FEATURE-A11 — Query proxy instances shows online status after heartbeat

```gherkin
Given the crosser-api server is running
  and the admin is authenticated with a valid JWT
  and a proxy instance has registered and sent a heartbeat within the last 90 seconds
When the admin sends a GET request to /api/v1/proxy/status?service=test-svc
Then the response status code is 200
  and the response JSON field "data" is a non-empty array
  and the first item field "status" equals "online"
```

### Scenario: FEATURE-A12 — Proxy registration without valid token returns 401

```gherkin
Given the crosser-api server is running
When a proxy sends a POST request to /api/v1/proxy/register without X-Proxy-Token header
Then the response status code is 401
  and the response JSON field "code" is not 0
```
