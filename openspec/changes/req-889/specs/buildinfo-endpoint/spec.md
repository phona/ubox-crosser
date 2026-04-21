---
change_id: req-889
title: "buildinfo-endpoint"
---

## ADDED Requirements

### Requirement: GET /buildinfo returns JSON build metadata

The admin HTTP server SHALL expose `GET /buildinfo` returning HTTP 200 with `Content-Type: application/json` and a JSON body containing `version`, `commit`, and `build_time` fields. The response is identical to `GET /version`.

#### Scenario: REQ-889-S1 Successful buildinfo response

- **GIVEN** the server is running with the admin HTTP listener
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
- **AND** the body is a JSON object with keys `version`, `commit`, `build_time`
- **AND** all values are non-empty strings

#### Scenario: REQ-889-S2 Response schema matches BuildInfo contract

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** the response body contains exactly 3 keys: `version`, `commit`, `build_time`
- **AND** no additional keys are present
- **AND** the `version` field matches semantic version format

---

### Requirement: Non-GET methods on /buildinfo are rejected

The `/buildinfo` endpoint is registered with `GET /buildinfo` pattern on `http.ServeMux`. Non-GET methods SHALL return HTTP 405 Method Not Allowed.

#### Scenario: REQ-889-S3 POST /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: REQ-889-S4 PUT /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: REQ-889-S5 DELETE /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: /buildinfo response is identical to /version

Since `/buildinfo` reuses `version.Handler`, the response body and headers SHALL be byte-for-byte identical to `GET /version`.

#### Scenario: REQ-889-S6 Buildinfo matches version endpoint output

- **GIVEN** admin HTTP listener is running with both `/version` and `/buildinfo` registered
- **WHEN** a client sends `GET /buildinfo` and `GET /version`
- **THEN** both responses have status `200 OK`
- **AND** both response bodies are identical JSON objects
