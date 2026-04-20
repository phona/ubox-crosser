## ADDED Requirements

### Requirement: Version endpoint returns build metadata
The proxy-server binary SHALL expose a `GET /version` HTTP endpoint that returns JSON containing the application version, git commit SHA, and build timestamp.

#### Scenario: FEATURE-S1 Successful version response
- **WHEN** a client sends `GET /version` to the HTTP listener
- **THEN** the server SHALL respond with HTTP 200
- **AND** the `Content-Type` header SHALL be `application/json`
- **AND** the response body SHALL be a JSON object containing exactly three fields: `version` (string), `commit` (string), `build_time` (string in ISO 8601 format)

#### Scenario: FEATURE-S2 Version field reflects compiled constant
- **WHEN** a client sends `GET /version`
- **THEN** the `version` field SHALL match the value of the `Version` variable in the `internal/version` package

#### Scenario: FEATURE-S3 Commit field reflects build-time injection
- **WHEN** the binary is built with `-ldflags -X .../internal/version.Commit=abc1234`
- **THEN** a `GET /version` response SHALL contain `"commit": "abc1234"`

#### Scenario: FEATURE-S4 Build time field reflects build-time injection
- **WHEN** the binary is built with `-ldflags -X .../internal/version.BuildTime=2024-01-15T10:30:00Z`
- **THEN** a `GET /version` response SHALL contain `"build_time": "2024-01-15T10:30:00Z"`

### Requirement: Version endpoint requires no authentication
The `GET /version` endpoint SHALL be accessible without any authentication credentials or tokens.

#### Scenario: FEATURE-S5 Unauthenticated access succeeds
- **WHEN** a client sends `GET /version` without any authentication headers
- **THEN** the server SHALL respond with HTTP 200 and the version JSON body

### Requirement: Version endpoint rejects non-GET methods
The `GET /version` endpoint SHALL only accept the GET HTTP method.

#### Scenario: FEATURE-S6 POST to version endpoint is rejected
- **WHEN** a client sends `POST /version`
- **THEN** the server SHALL respond with HTTP 405 Method Not Allowed

### Requirement: Default values when ldflags not provided
When the binary is built without ldflags injection, the version fields SHALL use sensible defaults.

#### Scenario: FEATURE-S7 Development build defaults
- **WHEN** the binary is built with `go build` (no ldflags)
- **AND** a client sends `GET /version`
- **THEN** `version` SHALL be `"0.1.0"`
- **AND** `commit` SHALL be `"unknown"`
- **AND** `build_time` SHALL be `"unknown"`
