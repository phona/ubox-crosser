---
spec_id: version-endpoint
change_id: req-669
title: "GET /version returns build metadata as JSON"
layer: backend
---

## ADDED Requirements

### Requirement: GET /version returns 200 with JSON body
The server SHALL respond to `GET /version` with HTTP 200, `Content-Type: application/json`, and a JSON body containing `version`, `commit`, and `build_time` string fields.

#### Scenario: Successful version query
- **WHEN** a client sends `GET /version` to the admin HTTP listener
- **THEN** the response status code SHALL be 200
- **THEN** the `Content-Type` header SHALL be `application/json`
- **THEN** the response body SHALL be valid JSON with exactly three fields: `version` (string), `commit` (string), `build_time` (string)

#### Scenario: Version field matches compiled constant
- **WHEN** a client sends `GET /version`
- **THEN** the `version` field in the response SHALL equal the `Version` constant defined in the `version` package (currently `0.1.0`)

### Requirement: Non-GET methods are rejected
The admin HTTP listener SHALL reject non-GET requests to `/version` with HTTP 405 Method Not Allowed, relying on Go 1.22+ method-based routing.

#### Scenario: POST to /version returns 405
- **WHEN** a client sends `POST /version`
- **THEN** the response status code SHALL be 405

#### Scenario: PUT to /version returns 405
- **WHEN** a client sends `PUT /version`
- **THEN** the response status code SHALL be 405

#### Scenario: DELETE to /version returns 405
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

### Requirement: Admin HTTP listener binds to configurable address
The server SHALL expose the admin HTTP listener on an address configurable via `--http-addr` flag, defaulting to `:8080`.

#### Scenario: Default admin address
- **WHEN** the server starts without `--http-addr` flag
- **THEN** the admin HTTP listener SHALL bind to `:8080`

#### Scenario: Custom admin address
- **WHEN** the server starts with `--http-addr 127.0.0.1:9090`
- **THEN** the admin HTTP listener SHALL bind to `127.0.0.1:9090`
