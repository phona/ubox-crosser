---
capability: webhook-debug-endpoint
change_id: REQ-945
status: LOCKED
---

## ADDED

### Scenario: REQ-945-S1 — POST webhook with JSON body returns complete request info

```gherkin
Given the admin HTTP server is running
When the user sends a POST request to /webhook-debug
  with header "Content-Type: application/json"
  and body '{"event":"push"}'
Then the response status code is 200
  and the response Content-Type starts with "application/json"
  and the response JSON field "method" equals "POST"
  and the response JSON field "path" equals "/webhook-debug"
  and the response JSON field "body" equals '{"event":"push"}'
  and the response JSON field "headers" contains key "Content-Type"
```

### Scenario: REQ-945-S2 — GET with query params populates query field

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /webhook-debug?foo=bar&baz=qux
Then the response status code is 200
  and the response JSON field "method" equals "GET"
  and the response JSON field "query" contains key "foo" with value ["bar"]
  and the response JSON field "query" contains key "baz" with value ["qux"]
  and the response JSON field "body" equals ""
```

### Scenario: REQ-945-S3 — PUT with form body returns body content

```gherkin
Given the admin HTTP server is running
When the user sends a PUT request to /webhook-debug
  with header "Content-Type: application/x-www-form-urlencoded"
  and body "key=val"
Then the response status code is 200
  and the response JSON field "method" equals "PUT"
  and the response JSON field "body" equals "key=val"
```

### Scenario: REQ-945-S4 — DELETE with no body returns empty body string

```gherkin
Given the admin HTTP server is running
When the user sends a DELETE request to /webhook-debug
Then the response status code is 200
  and the response JSON field "method" equals "DELETE"
  and the response JSON field "body" equals ""
```

### Scenario: REQ-945-S5 — Response JSON schema validation

```gherkin
Given the admin HTTP server is running
When the user sends a GET request to /webhook-debug
Then the response status code is 200
  and the response JSON contains required fields: method (string), path (string), query (map[string][]string), headers (map[string][]string), body (string)
```
