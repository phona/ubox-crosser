---
capability: webhook-debug-endpoint
change_id: REQ-945
---

## ADDED

### Scenario: FEATURE-A1 — POST webhook with JSON body returns complete request info

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /webhook-debug
  with header "Content-Type: application/json"
  and body '{"event":"push","repo":"ubox-crosser"}'
Then the response status code is 200
  and the response Content-Type is "application/json; charset=utf-8"
  and the response JSON field "method" equals "POST"
  and the response JSON field "path" equals "/webhook-debug"
  and the response JSON field "body" equals '{"event":"push","repo":"ubox-crosser"}'
  and the response JSON field "headers" contains key "Content-Type"
```

### Scenario: FEATURE-A2 — GET with query params populates query field

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /webhook-debug?foo=bar&baz=123
Then the response status code is 200
  and the response JSON field "method" equals "GET"
  and the response JSON field "query" contains key "foo" with value ["bar"]
  and the response JSON field "query" contains key "baz" with value ["123"]
  and the response JSON field "body" equals ""
```

### Scenario: FEATURE-A3 — Custom headers are captured in response

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /webhook-debug
  with header "X-Webhook-Secret: abc123"
  and header "X-Custom-Event: deploy"
Then the response status code is 200
  and the response JSON field "headers" contains key "X-Webhook-Secret" with value ["abc123"]
  and the response JSON field "headers" contains key "X-Custom-Event" with value ["deploy"]
```

### Scenario: FEATURE-A4 — Empty body request returns empty string body

```gherkin
Given the admin HTTP server is running
When the user sends a DELETE request to /webhook-debug
Then the response status code is 200
  and the response JSON field "method" equals "DELETE"
  and the response JSON field "body" equals ""
```

### Scenario: FEATURE-A5 — Large body is handled correctly

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /webhook-debug
  with a body of 64KB of repeated "x" characters
Then the response status code is 200
  and the response JSON field "body" has length 65536
  and the response JSON field "method" equals "POST"
```

### Scenario: FEATURE-A6 — Existing endpoints still work after adding webhook-debug

```gherkin
Given the admin HTTP server is running with /webhook-debug registered
When the user sends a GET request to /ping
Then the response status code is 200
  and the response body is "pong"

When the user sends a GET request to /healthz
Then the response status code is 200
  and the response Content-Type is "application/json"

When the user sends a GET request to /echo?msg=hello
Then the response status code is 200

When the user sends a GET request to /version
Then the response status code is 200
```
