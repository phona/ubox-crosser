---
change_id: req-685
title: "healthz-endpoint"
---

## ADDED Requirements

### Requirement: GET /healthz returns liveness status as JSON

The admin HTTP server SHALL expose `GET /healthz` returning HTTP 200 with `Content-Type: application/json` and a JSON body containing a single `status` field with value `"ok"`.

#### Scenario: ACCEPT-S1 Successful health check

- **GIVEN** the server is running with `--admin-addr :8080`
- **WHEN** a client sends `GET /healthz` to the admin HTTP listener
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
- **AND** the body is valid JSON `{"status":"ok"}`

#### Scenario: ACCEPT-S2 Response body has exactly one key

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /healthz`
- **THEN** the response body is a JSON object with exactly one key: `status`
- **AND** the value of `status` is the string `"ok"`

---

### Requirement: Non-GET methods on /healthz are rejected

The `/healthz` endpoint is registered with `GET /healthz` pattern on `http.ServeMux`. Non-GET methods SHALL return HTTP 405 Method Not Allowed.

#### Scenario: ACCEPT-S3 POST /healthz returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /healthz`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S4 PUT /healthz returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /healthz`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S5 DELETE /healthz returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /healthz`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Response body schema matches OpenAPI contract

The JSON response body SHALL conform to the `HealthStatus` schema defined in `contract.spec.yaml`: an object with a single required string field `status`.

#### Scenario: ACCEPT-S6 Response matches HealthStatus schema

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /healthz`
- **THEN** the response body is a JSON object
- **AND** the object has exactly one key: `status`
- **AND** the value is a string
