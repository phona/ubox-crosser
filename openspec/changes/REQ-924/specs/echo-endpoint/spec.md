---
capability: echo-endpoint
change_id: REQ-924
status: LOCKED
---

## Scenario: REQ-924-S1

**Title:** Echo endpoint returns msg parameter as plain text

```gherkin
Given the admin HTTP server is running
When a GET request is sent to /echo?msg=hello
Then the response status is 200
And the Content-Type header is "text/plain; charset=utf-8"
And the response body is "hello"
```

## Scenario: REQ-924-S2

**Title:** Echo endpoint returns empty body when msg is empty string

```gherkin
Given the admin HTTP server is running
When a GET request is sent to /echo?msg=
Then the response status is 200
And the Content-Type header is "text/plain; charset=utf-8"
And the response body is ""
```

## Scenario: REQ-924-S3

**Title:** Echo endpoint returns empty body when msg parameter is absent

```gherkin
Given the admin HTTP server is running
When a GET request is sent to /echo (no query parameters)
Then the response status is 200
And the Content-Type header is "text/plain; charset=utf-8"
And the response body is ""
```

## Scenario: REQ-924-S4

**Title:** Echo endpoint rejects non-GET methods with 405

```gherkin
Given the admin HTTP server is running
And the route is registered as "GET /echo"
When a POST request is sent to /echo
Then the response status is 405 Method Not Allowed
```
