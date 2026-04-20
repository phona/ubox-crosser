---
capability: uptime-endpoint
change_id: REQ-736
status: ADDED
---

## Scenario: REQ-736-S1 — GET /uptime returns 200 with JSON body

**Given** the server is running and uptime.Init() has been called
**When** a client sends `GET /uptime`
**Then** the response status code is `200`
**And** the `Content-Type` header is `application/json`
**And** the body is a JSON object with exactly one key `uptime_seconds` whose value is a non-negative integer

## Scenario: REQ-736-S2 — uptime_seconds reflects elapsed time

**Given** the server has been running for at least 1 second after Init()
**When** a client sends `GET /uptime`
**Then** `uptime_seconds` >= 1

## Scenario: REQ-736-S3 — Non-GET methods return 405

**Given** the server is running
**When** a client sends `POST /uptime`
**Then** the response status code is `405`

**When** a client sends `PUT /uptime`
**Then** the response status code is `405`

**When** a client sends `DELETE /uptime`
**Then** the response status code is `405`

## Scenario: REQ-736-S4 — Response body has no extra fields

**Given** the server is running
**When** a client sends `GET /uptime`
**Then** the JSON response contains exactly the key `uptime_seconds` and no additional keys
