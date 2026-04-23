# /buildinfo Endpoint Specification

## Overview

The `/buildinfo` endpoint provides build metadata as JSON. It's accessible via HTTP GET and requires no authentication.

## Acceptance Scenarios

### [ubox-crosser-S1] Successful buildinfo request
**Given** the ubox-crosser server is running on port 7000  
**When** I send a GET request to `http://localhost:7000/buildinfo`  
**Then** I receive HTTP 200  
**And** the response contains valid JSON with fields: git_sha, build_id, go_version  

### [ubox-crosser-S2] git_sha format is correct
**Given** the server is running  
**When** I request `/buildinfo`  
**Then** git_sha matches the pattern `^[0-9a-f]{7}$`  
**And** git_sha matches the first 7 characters of HEAD commit  

### [ubox-crosser-S3] build_id defaults to 'dev'
**Given** the server is running with no BUILD_ID environment variable set  
**When** I request `/buildinfo`  
**Then** build_id field equals "dev"  

### [ubox-crosser-S4] build_id uses environment variable
**Given** the server is running with BUILD_ID=myapp-20240423-v1.0 environment variable  
**When** I request `/buildinfo`  
**Then** build_id field equals "myapp-20240423-v1.0"  

### [ubox-crosser-S5] go_version is always go1.23
**Given** the server is running  
**When** I request `/buildinfo`  
**Then** go_version field equals "go1.23"  

### [ubox-crosser-S6] Response is valid JSON
**Given** the server is running  
**When** I request `/buildinfo`  
**Then** the response can be parsed as valid JSON  
**And** all required fields are present  

### [ubox-crosser-S7] Curl command works
**Given** the server is running on port 7000  
**When** I execute `curl -fsS http://localhost:7000/buildinfo`  
**Then** the exit code is 0  
**And** the output is the JSON response body  

### [ubox-crosser-S8] Content-Type is application/json
**Given** the server is running  
**When** I request `/buildinfo` and check response headers  
**Then** Content-Type header is "application/json"  

## API Contract Reference

See [contract.spec.yaml](contract.spec.yaml) for detailed API schema.

