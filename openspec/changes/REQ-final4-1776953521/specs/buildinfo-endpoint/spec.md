# /buildinfo Endpoint - Build Information

## Summary
Add a `/buildinfo` HTTP endpoint to the ubox-crosser service that returns build and version information. This endpoint allows clients to retrieve the Git commit SHA, build ID, and Go version used to build the service.

## Scenarios

### [ubox-crosser-S1]: Buildinfo Endpoint Returns 200 OK
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the endpoint responds with status code 200 OK  
**And** the response body contains valid JSON  

### [ubox-crosser-S2]: Buildinfo Response Contains Git SHA
**Given** the ubox-crosser service is running with git SHA injected via ldflags  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response includes a `git_sha` field  
**And** the `git_sha` is exactly 7 hexadecimal characters  
**And** the `git_sha` matches the short commit hash at build time  

### [ubox-crosser-S3]: Buildinfo Response Contains Build ID
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response includes a `build_id` field  
**And** when BUILD_ID environment variable is set, `build_id` contains that value  
**And** when BUILD_ID environment variable is not set, `build_id` equals "dev"  

### [ubox-crosser-S4]: Buildinfo Response Contains Go Version
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response includes a `go_version` field  
**And** the `go_version` value is "go1.23"  

### [ubox-crosser-S5]: Buildinfo Response Format Is Valid JSON
**Given** the ubox-crosser service is running  
**When** a client makes an HTTP GET request to `/buildinfo`  
**Then** the response body is valid JSON  
**And** the JSON object contains exactly three fields: `git_sha`, `build_id`, and `go_version`  
**And** all three fields are present and non-empty strings  

### [ubox-crosser-S6]: Buildinfo Endpoint Does Not Require Authentication
**Given** the ubox-crosser service is running  
**When** an unauthenticated client makes an HTTP GET request to `/buildinfo`  
**Then** the endpoint responds with status code 200 OK  
**And** the response contains the complete build information  

### [ubox-crosser-S7]: Buildinfo Endpoint Handles Concurrent Requests
**Given** the ubox-crosser service is running  
**When** multiple clients make concurrent HTTP GET requests to `/buildinfo`  
**Then** all requests complete successfully  
**And** each response contains identical build information  

## Additional Notes

- The `/buildinfo` endpoint should be accessible via HTTP GET requests
- The endpoint should not require authentication
- Response time should be minimal (< 50ms) for a simple build info retrieval
- The git_sha is injected at build time via ldflags and does not change during runtime
- The build_id can be overridden via the BUILD_ID environment variable
- The go_version is hardcoded to "go1.23"
