# Capability: version-endpoint

## ADDED Requirements

### Requirement: GET /version returns build metadata as JSON

Admin HTTP server SHALL expose `GET /version` returning HTTP 200 with `Content-Type: application/json` and a JSON body containing `version`, `commit`, and `build_time` string fields.

#### Scenario: ACCEPT-S1 Successful version query

- **GIVEN** the server is running with `--admin-addr :8080`
- **WHEN** a client sends `GET /version` to the admin HTTP listener
- **THEN** the response status is `200 OK`
- **AND** the response header `Content-Type` equals `application/json`
- **AND** the body is valid JSON containing keys `version`, `commit`, `build_time`
- **AND** `version` equals `0.1.0`

#### Scenario: ACCEPT-S2 Commit field is populated at build time

- **GIVEN** the binary is built with `-ldflags -X github.com/phona/ubox-crosser/version.Commit=abc1234`
- **WHEN** a client sends `GET /version`
- **THEN** the `commit` field in the response body equals `abc1234`

#### Scenario: ACCEPT-S3 BuildTime field is populated at build time

- **GIVEN** the binary is built with `-ldflags -X github.com/phona/ubox-crosser/version.BuildTime=2026-04-20T12:00:00Z`
- **WHEN** a client sends `GET /version`
- **THEN** the `build_time` field in the response body equals `2026-04-20T12:00:00Z`

#### Scenario: ACCEPT-S4 Default values when not injected

- **GIVEN** the binary is built without ldflags injection
- **WHEN** a client sends `GET /version`
- **THEN** `commit` equals `unknown`
- **AND** `build_time` equals `unknown`

---

### Requirement: Non-GET methods on /version are rejected

The `/version` endpoint is registered with `GET /version` pattern on `http.ServeMux`. Non-GET methods SHALL return HTTP 405 Method Not Allowed.

#### Scenario: ACCEPT-S5 POST /version returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `POST /version`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S6 PUT /version returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `PUT /version`
- **THEN** the response status is `405 Method Not Allowed`

#### Scenario: ACCEPT-S7 DELETE /version returns 405

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `DELETE /version`
- **THEN** the response status is `405 Method Not Allowed`

---

### Requirement: Unknown paths return 404

The admin HTTP server SHALL return HTTP 404 for any path other than `/version`.

#### Scenario: ACCEPT-S8 Request to root path returns 404

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /`
- **THEN** the response status is `404 Not Found`

#### Scenario: ACCEPT-S9 Request to unknown path returns 404

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /metrics`
- **THEN** the response status is `404 Not Found`

---

### Requirement: Admin address is configurable via CLI flag

The server SHALL accept `--admin-addr` flag to configure the admin HTTP listener address, defaulting to `:8080`.

#### Scenario: ACCEPT-S10 Custom admin address via CLI flag

- **GIVEN** the server is started with `--admin-addr :9090`
- **WHEN** a client sends `GET /version` to port 9090
- **THEN** the response status is `200 OK`
- **AND** port 8080 is not listening

#### Scenario: ACCEPT-S11 Default admin address

- **GIVEN** the server is started without `--admin-addr` flag
- **WHEN** a client sends `GET /version` to port 8080
- **THEN** the response status is `200 OK`

---

### Requirement: Build metadata injected via Makefile and Dockerfile

The `Makefile` SHALL define ldflags that inject `version.Commit` (from `git rev-parse --short HEAD`) and `version.BuildTime` (from `date -u`) at build time. The `Dockerfile` SHALL pass these values as build args.

#### Scenario: ACCEPT-S12 make build injects commit hash

- **GIVEN** the repository has at least one commit
- **WHEN** `make build` is executed
- **THEN** `bin/server` binary reports a `commit` value matching `git rev-parse --short HEAD`

#### Scenario: ACCEPT-S13 make build injects build time

- **GIVEN** the build environment has a valid system clock
- **WHEN** `make build` is executed
- **THEN** `bin/server` binary reports a `build_time` value in ISO 8601 UTC format

---

### Requirement: Response body schema matches OpenAPI contract

The JSON response body SHALL conform to the `VersionInfo` schema defined in `contract.spec.yaml`: an object with required string fields `version`, `commit`, and `build_time`.

#### Scenario: ACCEPT-S14 Response matches VersionInfo schema

- **GIVEN** admin HTTP listener is running
- **WHEN** a client sends `GET /version`
- **THEN** the response body is a JSON object
- **AND** the object has exactly three keys: `version`, `commit`, `build_time`
- **AND** all three values are strings
