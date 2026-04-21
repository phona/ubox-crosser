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

## ADDED Acceptance Scenarios

### Requirement: User can echo arbitrary text via GET /echo

From a user/operator perspective, the `/echo` endpoint reflects back any string passed via the `msg` query parameter, enabling network connectivity checks and request-chain debugging.

#### Scenario: FEATURE-A1 Successful echo with msg parameter

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `GET /echo?msg=hello`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`
- **AND** the response body is exactly `hello`

#### Scenario: FEATURE-A2 Echo with special characters

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `GET /echo?msg=hello%20world%21`
- **THEN** the response status is `200 OK`
- **AND** the response body is exactly `hello world!`

---

### Requirement: Missing or empty msg returns empty body

When the `msg` parameter is absent or empty, the endpoint returns HTTP 200 with an empty body rather than an error.

#### Scenario: FEATURE-A3 Empty msg parameter

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `GET /echo?msg=`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`
- **AND** the response body is empty

#### Scenario: FEATURE-A4 Missing msg parameter

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `GET /echo`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`
- **AND** the response body is empty

---

### Requirement: Non-GET methods are rejected by /echo

The `/echo` endpoint only accepts GET requests. Other HTTP methods return 405.

#### Scenario: FEATURE-A5 POST /echo returns 405

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `POST /echo`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: FEATURE-A6 PUT /echo returns 405

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `PUT /echo`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Existing admin endpoints unaffected by /echo addition

Adding `/echo` must not break existing admin endpoints.

#### Scenario: FEATURE-A7 GET /ping still works

- **GIVEN** the admin HTTP listener is running with `/echo` registered
- **WHEN** a client sends `GET /ping`
- **THEN** the response status is `200 OK`

#### Scenario: FEATURE-A8 GET /healthz still works

- **GIVEN** the admin HTTP listener is running with `/echo` registered
- **WHEN** a client sends `GET /healthz`
- **THEN** the response status is `200 OK`

#### Scenario: FEATURE-A9 GET /version still works

- **GIVEN** the admin HTTP listener is running with `/echo` registered
- **WHEN** a client sends `GET /version`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
