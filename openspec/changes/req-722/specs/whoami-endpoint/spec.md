---
change_id: req-722
title: "whoami-endpoint"
---

## ADDED Requirements

### Requirement: GET /whoami returns plain-text hostname

The admin HTTP server SHALL expose `GET /whoami` returning HTTP 200 with `Content-Type: text/plain; charset=utf-8` and a body containing the current machine hostname.

#### Scenario: ACCEPT-A1 Successful whoami

- **GIVEN** the server is running with `--admin-addr :8080`
- **WHEN** a client sends `GET /whoami` to the admin HTTP listener
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`
- **AND** the body is a non-empty string representing a hostname

#### Scenario: ACCEPT-A2 Response body is a valid hostname

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /whoami`
- **THEN** the response body, after trimming trailing whitespace, is a non-empty string
- **AND** the value matches a valid hostname pattern (alphanumeric, hyphens, dots)

---

### Requirement: Non-GET methods on /whoami are rejected

The `/whoami` endpoint is registered with `GET /whoami` pattern on `http.ServeMux`. Non-GET methods SHALL return HTTP 405 Method Not Allowed.

#### Scenario: ACCEPT-A3 POST /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /whoami`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-A4 PUT /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /whoami`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-A5 DELETE /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /whoami`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: os.Hostname failure returns fallback "unknown"

If `os.Hostname()` returns an error, the handler SHALL return HTTP 200 with body `unknown` instead of failing with a 500.

#### Scenario: ACCEPT-A6 Fallback to "unknown" on hostname error

- **GIVEN** admin HTTP listener is running
- **AND** `os.Hostname()` returns an error
- **WHEN** a client sends `GET /whoami`
- **THEN** the response status is `200 OK`
- **AND** the body is the string `unknown`

---

## ADDED Contract Tests

### Requirement: GET /whoami response conforms to OpenAPI contract

The response of `GET /whoami` SHALL conform to the schema defined in `contract.spec.yaml`: HTTP 200, `Content-Type: text/plain; charset=utf-8`, body is a non-empty string.

#### Scenario: REQ-722-S1 GET /whoami returns 200 with plain-text content type

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /whoami`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`

#### Scenario: REQ-722-S2 GET /whoami body is a non-empty hostname string

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /whoami`
- **THEN** the response body is a non-empty string
- **AND** the body length is at least 1 character

#### Scenario: REQ-722-S3 os.Hostname failure returns "unknown"

- **GIVEN** admin HTTP listener is running
- **AND** `os.Hostname()` returns an error
- **WHEN** a client sends `GET /whoami`
- **THEN** the response status is `200 OK`
- **AND** the body is exactly `unknown`

---

### Requirement: Non-GET methods return 405

The `/whoami` endpoint SHALL reject non-GET methods with HTTP 405.

#### Scenario: REQ-722-S4 POST /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /whoami`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: REQ-722-S5 PUT /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /whoami`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: REQ-722-S6 DELETE /whoami returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /whoami`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Existing admin endpoints remain unaffected

Adding `GET /whoami` SHALL NOT change the behavior of existing admin endpoints (`GET /version`, `GET /healthz`, `GET /ping`).

#### Scenario: ACCEPT-A7 GET /ping still returns 200

- **GIVEN** admin HTTP listener is running with `/whoami` registered
- **WHEN** a client sends `GET /ping`
- **THEN** the response status is `200 OK`
- **AND** the body is `pong`

#### Scenario: ACCEPT-A8 GET /healthz still returns 200

- **GIVEN** admin HTTP listener is running with `/whoami` registered
- **WHEN** a client sends `GET /healthz`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
