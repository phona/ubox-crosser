---
capability: uptime-endpoint
change_id: REQ-736
---

## Scenario: REQ-736-A1 — Normal uptime response

**Given** the ubox-crosser service is running
**When** a user sends `GET /uptime`
**Then** the response status is `200 OK`
**And** the response body is JSON containing `{"uptime_seconds": <integer>}`
**And** `uptime_seconds` is a non-negative integer

## Scenario: REQ-736-A2 — Uptime value increases over time

**Given** the ubox-crosser service has been running for at least 2 seconds
**When** a user sends `GET /uptime` at time T1, waits 2 seconds, then sends `GET /uptime` at time T2
**Then** the `uptime_seconds` value at T2 is greater than or equal to the value at T1 + 1

## Scenario: REQ-736-A3 — Uptime resets after restart

**Given** the ubox-crosser service is restarted
**When** a user sends `GET /uptime` within 5 seconds of the restart
**Then** `uptime_seconds` is less than 5

## Scenario: REQ-736-A4 — Non-GET methods are rejected

**Given** the ubox-crosser service is running
**When** a user sends `POST /uptime`
**Then** the response status is `405 Method Not Allowed`

## Scenario: REQ-736-A5 — Endpoint does not require authentication

**Given** the ubox-crosser service is running
**When** an unauthenticated user sends `GET /uptime`
**Then** the response status is `200 OK` (not 401/403)
