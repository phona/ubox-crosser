# /version Endpoint Spec

## Overview
The `/version` endpoint provides the git commit SHA of the current build.

## HTTP Interface Contract

### Request
- **Method:** GET
- **Path:** `/version`
- **Query Parameters:** None
- **Body:** Empty

### Response

#### Success Response (HTTP 200)
- **Content-Type:** `application/json`
- **Body:**
```json
{
  "sha": "abc1234567890def1234567890def1234567890"
}
```

#### Response Fields
| Field | Type | Description |
|-------|------|-------------|
| `sha` | string | Full 40-character git commit SHA (hex format) |

## Scenarios

### Scenario: REQ-final2-1776868985-S1
**Given:** Server is running with git SHA injected via -ldflags
**When:** Client sends GET request to `/version`
**Then:**
- Status code is 200
- Response Content-Type is `application/json`
- Response body contains JSON object with "sha" field
- SHA value is 40-character hex string matching `[a-f0-9]{40}`

### Scenario: REQ-final2-1776868985-S2
**Given:** Server is running
**When:** Client makes repeated requests to `/version`
**Then:**
- Each request returns HTTP 200
- SHA value is consistent across requests
- Response time is minimal (< 100ms)

### Scenario: REQ-final2-1776868985-S3
**Given:** Server is built with SHA injection
**When:** Server starts and receives GET `/version` request
**Then:**
- Response contains the correct git SHA from the build
- SHA is not "unknown" in production builds
- SHA matches git rev-parse HEAD from build time

## Design Decisions

1. **Full SHA (40 chars) vs Short SHA:** Using full SHA for uniqueness and production traceability
2. **JSON Response:** Standard REST API format, easy to extend with additional fields in future
3. **Read-only GET:** No side effects, idempotent operation
4. **No Authentication:** Version info is metadata that should be publicly available
