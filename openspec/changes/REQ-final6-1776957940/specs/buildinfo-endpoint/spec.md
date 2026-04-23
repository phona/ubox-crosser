# /buildinfo Endpoint

## Summary

Expose build metadata via an unauthenticated HTTP GET endpoint so operators can confirm which revision and CI build is running.

## Scenarios

### [UBOX-S1]: Returns 200 OK
**Given** the ubox-crosser server is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response status code is 200 OK

### [UBOX-S2]: Response is valid JSON
**Given** the ubox-crosser server is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response body is valid JSON

### [UBOX-S3]: Response includes git_sha field
**Given** the binary was built with `-X main.GitSHA=<sha>`  
**When** a client requests `/buildinfo`  
**Then** the response JSON contains a non-empty `git_sha` field

### [UBOX-S4]: Response includes build_id field
**Given** the server is running  
**When** a client requests `/buildinfo`  
**Then** the response JSON contains a non-empty `build_id` field

### [UBOX-S5]: go_version is go1.23
**Given** the server is running  
**When** a client requests `/buildinfo`  
**Then** `go_version` equals `"go1.23"`

### [UBOX-S6]: build_id defaults to "dev" when BUILD_ID env is absent
**Given** the server started without the `BUILD_ID` environment variable  
**When** a client requests `/buildinfo`  
**Then** `build_id` equals `"dev"`

### [UBOX-S7]: Endpoint requires no authentication
**Given** the server is running  
**When** a client makes a GET request to `/buildinfo` without any credentials  
**Then** the response is not 401 or 403
