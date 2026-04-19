## ADDED Requirements

### Requirement: Version endpoint returns application version

The server SHALL expose an HTTP `GET /version` endpoint that returns the current application version in JSON format. The response body SHALL be a JSON object with a `version` field containing the version string. The response SHALL use HTTP status 200 and `Content-Type: application/json`.

#### Scenario: Successful version query

- **WHEN** a client sends `GET /version` to the management HTTP address
- **THEN** the server responds with HTTP 200 and body `{"version":"v2"}` with `Content-Type: application/json`

#### Scenario: Version reflects build-time value

- **WHEN** the binary is built with `-ldflags "-X <package>.Version=v2"`
- **THEN** the `GET /version` endpoint returns `{"version":"v2"}`

### Requirement: Management HTTP listener is configurable

The server SHALL accept a management listen address via configuration (CLI flag or config file). If no management address is configured, the management HTTP server SHALL NOT start.

#### Scenario: Management address provided via CLI flag

- **WHEN** the server is started with `--management-address 127.0.0.1:8080`
- **THEN** the HTTP management API listens on `127.0.0.1:8080`

#### Scenario: Management address omitted

- **WHEN** the server is started without a management address
- **THEN** no HTTP management listener is started and the server operates normally via TCP only

### Requirement: Management HTTP server does not block main proxy

The management HTTP server SHALL run in a separate goroutine and SHALL NOT interfere with the primary TCP proxy listener or its connection handling.

#### Scenario: Concurrent operation

- **WHEN** the management HTTP server is running and the TCP proxy is handling connections
- **THEN** both operate independently without blocking each other
