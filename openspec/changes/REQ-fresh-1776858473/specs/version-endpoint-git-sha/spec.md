# Specification: /version Endpoint with Git SHA

## Endpoint Definition

**HTTP Method:** GET  
**Path:** `/version`  
**Host:** Management Server (default port 8888)  

## Response Schema

### 2xx Success Response
**Status Code:** `200 OK`  
**Content-Type:** `application/json`

```json
{
  "version": "<string>",
  "module": "<string>",
  "go_os": "<string>",
  "go_arch": "<string>",
  "git_sha": "<string>"
}
```

**Fields:**
- `version` (string): Build version from go build info. Example: `"v1.0.0"` or `"(devel)"` if unset.
- `module` (string): Module path. Value: `"github.com/phona/ubox-crosser"`.
- `go_os` (string): Target OS from runtime.GOOS. Example: `"linux"`, `"darwin"`, `"windows"`.
- `go_arch` (string): Target architecture from runtime.GOARCH. Example: `"amd64"`, `"arm64"`.
- `git_sha` (string): Git commit SHA-1 hash (full or short). Format: 40-character hex for full SHA or 7+ for short.

### 4xx Error Responses
- **405 Method Not Allowed:** Returned for POST, PUT, DELETE, PATCH, etc.
- **404 Not Found:** Endpoint not implemented (pre-implementation)

## Scenarios

### Scenario: REQ-fresh-1776858473-S1
**Title:** GET /version returns 200 with all required fields

**Given:** Management server is running  
**When:** Client sends GET request to `/version`  
**Then:**
- Response status code is 200
- Response Content-Type header is `application/json`
- Response body contains JSON with all required fields: `version`, `module`, `go_os`, `go_arch`, `git_sha`

### Scenario: REQ-fresh-1776858473-S2
**Title:** GET /version includes git_sha field with valid format

**Given:** Management server is running  
**When:** Client sends GET request to `/version`  
**Then:**
- Response contains `git_sha` field
- `git_sha` is a string
- `git_sha` is not empty
- `git_sha` contains only hexadecimal characters (0-9, a-f)

### Scenario: REQ-fresh-1776858473-S3
**Title:** GET /version git_sha remains consistent across requests

**Given:** Management server is running  
**When:** Client sends multiple GET requests to `/version`  
**Then:**
- `git_sha` value is identical in all responses
- Same git_sha across requests confirms it's derived from static build info

### Scenario: REQ-fresh-1776858473-S4
**Title:** GET /version response includes module field

**Given:** Management server is running  
**When:** Client sends GET request to `/version`  
**Then:**
- Response contains `module` field
- `module` value is `"github.com/phona/ubox-crosser"`

### Scenario: REQ-fresh-1776858473-S5
**Title:** GET /version response includes go_os and go_arch

**Given:** Management server is running  
**When:** Client sends GET request to `/version`  
**Then:**
- Response contains `go_os` field (e.g., "linux", "darwin")
- Response contains `go_arch` field (e.g., "amd64", "arm64")
- Both fields are non-empty strings

### Scenario: REQ-fresh-1776858473-S6
**Title:** POST /version returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends POST request to `/version`  
**Then:**
- Response status code is 405
- Response indicates Method Not Allowed

### Scenario: REQ-fresh-1776858473-S7
**Title:** PUT /version returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends PUT request to `/version`  
**Then:**
- Response status code is 405

### Scenario: REQ-fresh-1776858473-S8
**Title:** DELETE /version returns 405 Method Not Allowed

**Given:** Management server is running  
**When:** Client sends DELETE request to `/version`  
**Then:**
- Response status code is 405

## Edge Cases

### Edge Case: EC1 - Git SHA Format
**Condition:** Binary built from git repository  
**Expected:** `git_sha` is populated with commit hash (40 char for full SHA, or short form)

### Edge Case: EC2 - Development Build
**Condition:** Binary built without git metadata (e.g., outside repo or unclean state)  
**Expected:** `git_sha` has a meaningful fallback (empty string or special value like "unknown")

### Edge Case: EC3 - All Fields Present
**Condition:** Normal operation  
**Expected:** All five fields (version, module, go_os, go_arch, git_sha) are present in every response

## Response Examples

### Example 1: Production build with git metadata
```json
{
  "version": "v1.2.3",
  "module": "github.com/phona/ubox-crosser",
  "go_os": "linux",
  "go_arch": "amd64",
  "git_sha": "abc123def456..."
}
```

### Example 2: Development build
```json
{
  "version": "(devel)",
  "module": "github.com/phona/ubox-crosser",
  "go_os": "darwin",
  "go_arch": "arm64",
  "git_sha": "1a2b3c4d5e6f..."
}
```

## OpenAPI Definition

```yaml
paths:
  /version:
    get:
      summary: Version and build information including git SHA
      operationId: getVersion
      responses:
        '200':
          description: Version information
          content:
            application/json:
              schema:
                type: object
                properties:
                  version:
                    type: string
                    description: Build version from go build info
                  module:
                    type: string
                    description: Module path
                  go_os:
                    type: string
                    description: Target OS
                  go_arch:
                    type: string
                    description: Target architecture
                  git_sha:
                    type: string
                    description: Git commit SHA-1 hash
                    minLength: 7
                required:
                  - version
                  - module
                  - go_os
                  - go_arch
                  - git_sha
        '405':
          description: Method not allowed
```
