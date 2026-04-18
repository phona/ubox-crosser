## ADDED Requirements

### Requirement: Echo endpoint returns message
The system SHALL expose a `GET /api/echo` HTTP endpoint that accepts a `msg` query parameter and returns its value in a JSON response body.

#### Scenario: Successful echo with message
- **WHEN** a GET request is made to `/api/echo?msg=hello`
- **THEN** the response status code SHALL be 200
- **AND** the response `Content-Type` header SHALL be `application/json`
- **AND** the response body SHALL be `{"message": "hello"}`

#### Scenario: Echo with empty message parameter
- **WHEN** a GET request is made to `/api/echo?msg=`
- **THEN** the response status code SHALL be 200
- **AND** the response body SHALL be `{"message": ""}`

#### Scenario: Echo with missing message parameter
- **WHEN** a GET request is made to `/api/echo` without a `msg` query parameter
- **THEN** the response status code SHALL be 400
- **AND** the response body SHALL be `{"error": "msg parameter is required"}`

### Requirement: Echo endpoint rejects non-GET methods
The system SHALL reject requests to `/api/echo` that use HTTP methods other than GET.

#### Scenario: POST request to echo endpoint
- **WHEN** a POST request is made to `/api/echo?msg=hello`
- **THEN** the response status code SHALL be 405
- **AND** the response body SHALL indicate method not allowed
