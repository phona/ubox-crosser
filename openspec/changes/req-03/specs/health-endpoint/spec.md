## ADDED Requirements

### Requirement: Health endpoint returns status OK
The proxy server SHALL expose an HTTP `GET /health` endpoint that returns HTTP 200 with a JSON body `{"status":"ok"}` and `Content-Type: application/json` when the server process is running.

#### Scenario: Healthy server responds to GET /health
- **WHEN** the proxy server is running and a client sends `GET /health` to the health listen address
- **THEN** the server SHALL respond with HTTP status code `200` and body `{"status":"ok"}` with header `Content-Type: application/json`

#### Scenario: Non-GET method is rejected
- **WHEN** a client sends a request with a method other than GET (e.g., POST, PUT, DELETE) to `/health`
- **THEN** the server SHALL respond with HTTP status code `405 Method Not Allowed`

#### Scenario: Unknown path returns 404
- **WHEN** a client sends `GET /unknown` or any path other than `/health`
- **THEN** the server SHALL respond with HTTP status code `404 Not Found`

### Requirement: Health listen address is configurable
The proxy server SHALL accept a `--health-addr` CLI flag (and corresponding config field) that specifies the `host:port` for the health HTTP listener. The default value SHALL be `:8080`.

#### Scenario: Default health address
- **WHEN** the proxy server is started without the `--health-addr` flag
- **THEN** the health HTTP server SHALL listen on `:8080`

#### Scenario: Custom health address
- **WHEN** the proxy server is started with `--health-addr :9090`
- **THEN** the health HTTP server SHALL listen on `:9090`

### Requirement: Health server lifecycle
The health HTTP server SHALL start automatically when the proxy server starts. If the health listener fails to bind, the proxy server SHALL log the error and continue running without the health endpoint.

#### Scenario: Health server starts with proxy
- **WHEN** the proxy server starts successfully
- **THEN** the health HTTP endpoint SHALL be reachable at the configured address

#### Scenario: Health port bind failure is non-fatal
- **WHEN** the health listen port is already in use by another process
- **THEN** the proxy server SHALL log an error message and continue operating normally without the health endpoint
