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

---

## ADDED Acceptance Scenarios

### Requirement: User can retrieve build information via /buildinfo

From a user/operator perspective, the `/buildinfo` endpoint provides build metadata for deployment verification and diagnostics.

#### Scenario: FEATURE-A1 Successful buildinfo retrieval

- **GIVEN** the server is running with `--admin-addr :8080`
- **WHEN** a client sends `GET /buildinfo` to the admin HTTP listener
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
- **AND** the body is a JSON object with keys `version`, `commit`, and `build_time`
- **AND** all three values are non-empty strings

#### Scenario: FEATURE-A2 Build info fields are meaningful

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo`
- **THEN** `version` is a non-empty string (e.g. `0.1.0`)
- **AND** `commit` is a hex string representing a git SHA
- **AND** `build_time` is a non-empty timestamp string

---

### Requirement: /buildinfo and /version return identical data

Since both endpoints serve the same handler, an operator can use either path interchangeably.

#### Scenario: FEATURE-A3 Buildinfo matches version response

- **GIVEN** the server is running
- **WHEN** a client sends `GET /buildinfo` and `GET /version`
- **THEN** both responses have status `200 OK`
- **AND** both response bodies are byte-identical JSON

---

### Requirement: Non-GET methods are rejected by /buildinfo

Operators should receive a clear error when using wrong HTTP methods.

#### Scenario: FEATURE-A4 POST /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: FEATURE-A5 PUT /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: FEATURE-A6 DELETE /buildinfo returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /buildinfo`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Existing admin endpoints unaffected

Adding `/buildinfo` must not break existing admin endpoints.

#### Scenario: FEATURE-A7 GET /version still works

- **GIVEN** admin HTTP listener is running with `/buildinfo` registered
- **WHEN** a client sends `GET /version`
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`

#### Scenario: FEATURE-A8 GET /healthz still works

- **GIVEN** admin HTTP listener is running with `/buildinfo` registered
- **WHEN** a client sends `GET /healthz`
- **THEN** the response status is `200 OK`
