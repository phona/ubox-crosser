# Design: /healthz Endpoint Implementation

## Architecture

### Server Startup Time Tracking
The `ManagementServer` struct will track the service start time:
```go
type ManagementServer struct {
    address   string
    mux       *http.ServeMux
    errs      chan error
    startTime time.Time  // NEW
}
```

### Endpoint Handler
New handler: `handleHealthz(w http.ResponseWriter, r *http.Request)`
- Validates HTTP method (only GET allowed)
- Calculates uptime: `time.Since(m.startTime).Seconds()`
- Returns JSON with uptime and status

### Response Format
```json
{
  "status": "ok",
  "uptime_seconds": 1234
}
```

### Integration Points
- Handler registered in `registerHandlers()` method
- Server start time initialized in `NewManagementServer()`
- No dependency on external state

## Implementation Details

### Changes to `server/management.go`
1. Add `startTime time.Time` field to `ManagementServer`
2. Initialize `startTime` in `NewManagementServer()` to `time.Now()`
3. Register `/healthz` handler in `registerHandlers()`
4. Implement `handleHealthz()` method
5. Add response struct type for type safety

### Type Definition
```go
type HealthzResponse struct {
    Status        string `json:"status"`
    UptimeSeconds int64  `json:"uptime_seconds"`
}
```

## Testing Strategy
- Unit tests with `httptest` (existing pattern in `management_test.go`)
- Contract tests with Docker Compose (integration testing with real stack)
- Acceptance scenarios validating uptime progression

## Backward Compatibility
- No changes to existing `/version` or `/health` endpoints
- New endpoint does not affect existing functionality
- Fully additive change
