---
spec_id: version-endpoint
change_id: req-653
title: "GET /version returns build metadata as JSON"
layer: backend
---

## ADDED Requirements

### Requirement: GET /version returns 200 with JSON body
The server SHALL respond to `GET /version` with HTTP 200, `Content-Type: application/json`, and a JSON body containing `version`, `commit`, and `build_time` string fields.

#### Scenario: Successful version query
- **WHEN** a client sends `GET /version` to the HTTP server
- **THEN** the response status code SHALL be 200
- **THEN** the `Content-Type` header SHALL be `application/json`
- **THEN** the response body SHALL be valid JSON with exactly three fields: `version` (string), `commit` (string), `build_time` (string)

#### Scenario: Version field matches compiled constant
- **WHEN** a client sends `GET /version`
- **THEN** the `version` field in the response SHALL match the `Version` constant defined in the `version` package

### Requirement: Non-GET methods return 405
The server SHALL reject non-GET requests to `/version` with HTTP 405 Method Not Allowed.

#### Scenario: POST to /version
- **WHEN** a client sends `POST /version`
- **THEN** the response status code SHALL be 405

#### Scenario: PUT to /version
- **WHEN** a client sends `PUT /version`
- **THEN** the response status code SHALL be 405

#### Scenario: DELETE to /version
- **WHEN** a client sends `DELETE /version`
- **THEN** the response status code SHALL be 405

### Requirement: Build metadata injected at compile time
The `Commit` and `BuildTime` variables SHALL be populated via `-ldflags -X` during `go build`. When not injected, they SHALL default to `"unknown"`.

#### Scenario: Default values without ldflags
- **WHEN** the binary is built without `-ldflags`
- **THEN** `GET /version` SHALL return `commit: "unknown"` and `build_time: "unknown"`

#### Scenario: Injected values with ldflags
- **WHEN** the binary is built with `-ldflags` setting `Commit` and `BuildTime`
- **THEN** `GET /version` SHALL return the injected commit hash and build timestamp
