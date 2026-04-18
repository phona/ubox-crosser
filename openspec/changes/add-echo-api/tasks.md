## 1. HTTP Handler

- [ ] 1.1 Create `server/api/echo.go` with the echo handler that reads `msg` query param and returns JSON
- [ ] 1.2 Return 400 with error JSON when `msg` param is missing
- [ ] 1.3 Return 405 for non-GET methods

## 2. HTTP Server Setup

- [ ] 2.1 Create `server/api/server.go` to register routes and start the HTTP server
- [ ] 2.2 Add HTTP listen address to server config in `models/config/`
- [ ] 2.3 Wire HTTP server startup into the existing server command in `cmd/`

## 3. Testing

- [ ] 3.1 Write unit tests for the echo handler in `server/api/echo_test.go`
- [ ] 3.2 Verify all spec scenarios pass (200 with msg, 200 with empty msg, 400 without msg, 405 for POST)
