## ADDED Requirements

### Requirement: Version endpoint returns current version
The system SHALL expose an HTTP `GET /version` endpoint that returns the application version as a JSON object.

#### Scenario: Successful version request
- **WHEN** a client sends an HTTP `GET` request to `/version`
- **THEN** the server SHALL respond with HTTP status `200 OK`, `Content-Type: application/json`, and a body of `{"version":"v3"}`

#### Scenario: Wrong HTTP method
- **WHEN** a client sends an HTTP request to `/version` using a method other than `GET` (e.g., `POST`, `PUT`, `DELETE`)
- **THEN** the server SHALL respond with HTTP status `405 Method Not Allowed`

### Requirement: HTTP server listens on configurable address
The system SHALL start an HTTP server on a configurable listen address (default `:8080`) when the proxy server starts.

#### Scenario: Default HTTP address
- **WHEN** the server starts without an explicit `--http-addr` flag
- **THEN** the HTTP server SHALL listen on `:8080`

#### Scenario: Custom HTTP address
- **WHEN** the server starts with `--http-addr :9090`
- **THEN** the HTTP server SHALL listen on `:9090`
