---
change_id: req-686
title: "ping-endpoint"
---

## ADDED Requirements

### Requirement: GET /ping returns plain-text "pong"

The admin HTTP server SHALL expose `GET /ping` returning HTTP 200 with `Content-Type: text/plain; charset=utf-8` and a body containing the string `pong`.

#### Scenario: ACCEPT-S1 Successful ping

- **GIVEN** the server is running with `--admin-addr :8080`
- **WHEN** a client sends `GET /ping` to the admin HTTP listener
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `text/plain; charset=utf-8`
- **AND** the body is the string `pong`

#### Scenario: ACCEPT-S2 Response body is exactly "pong"

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /ping`
- **THEN** the response body, after trimming trailing whitespace, equals `pong`

---

### Requirement: Non-GET methods on /ping are rejected

The `/ping` endpoint is registered with `GET /ping` pattern on `http.ServeMux`. Non-GET methods SHALL return HTTP 405 Method Not Allowed.

#### Scenario: ACCEPT-S3 POST /ping returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /ping`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S4 PUT /ping returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /ping`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S5 DELETE /ping returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /ping`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Response body conforms to OpenAPI contract

The plain-text response body SHALL conform to the `text/plain` schema defined in `contract.spec.yaml`: the constant string `pong`.

#### Scenario: ACCEPT-S6 Response matches contract schema

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /ping`
- **THEN** the response body is a plain-text string
- **AND** the value is exactly `pong`
- **AND** there is no JSON wrapping or additional formatting

---

### Requirement: Existing admin endpoints remain unaffected

Adding `GET /ping` SHALL NOT change the behavior of existing admin endpoints (`GET /version`, `GET /healthz`).

#### Scenario: ACCEPT-S7 GET /version still returns 200

- **GIVEN** admin HTTP listener is running with `/ping` registered
- **WHEN** a client sends `GET /version`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`

#### Scenario: ACCEPT-S8 Custom admin address via CLI flag

- **GIVEN** the server is started with `--admin-addr :9090`
- **WHEN** a client sends `GET /ping` to port 9090
- **THEN** the response status is `200 OK`
- **AND** the body is `pong`
