---
capability: proxy
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: REQ-997-S12 — POST /api/v1/proxy/register returns a proxy token

```gherkin
Given the crosser-api server is running
  and a service "test-svc" exists
  and a valid service token is available
When the client sends POST /api/v1/proxy/register with service_name="test-svc", instance_id="inst-01", address="10.0.0.1:8388"
  and the request includes header X-Proxy-Token with the service token
Then the response status code is 200
  and the response JSON field "code" is 0
  and the response JSON field "data.token" is a non-empty string
```

### Scenario: REQ-997-S13 — POST /api/v1/proxy/heartbeat updates proxy status

```gherkin
Given the crosser-api server is running
  and a proxy instance "inst-01" is registered for service "test-svc"
  and the proxy has a valid token
When the client sends POST /api/v1/proxy/heartbeat with instance_id="inst-01", connection_count=5
  and the request includes header X-Proxy-Token with the proxy token
Then the response status code is 200
  and the response JSON field "code" is 0
```

### Scenario: REQ-997-S14 — Proxy endpoints reject requests without valid X-Proxy-Token

```gherkin
Given the crosser-api server is running
When the client sends POST /api/v1/proxy/register without X-Proxy-Token header
Then the response status code is 401
  and the response JSON field "code" is 3001
When the client sends POST /api/v1/proxy/heartbeat with an invalid X-Proxy-Token
Then the response status code is 401
  and the response JSON field "code" is 3001
```

### Scenario: REQ-997-S15 — All responses follow unified envelope {code, message, data}

```gherkin
Given the crosser-api server is running
When any endpoint returns a success response
Then the response JSON has integer field "code" with value 0
  and the response JSON has string field "message"
  and the response JSON has field "data"
When any endpoint returns an error response
Then the response JSON has integer field "code" with non-zero value
  and the response JSON has string field "message"
  and the response JSON does not have field "data"
```
