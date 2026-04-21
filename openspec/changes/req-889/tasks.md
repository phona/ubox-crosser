---
change_id: req-889
title: "GET /buildinfo endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [x] [REQ-889-S1] define OpenAPI contract for GET /buildinfo (contract.spec.yaml)
- [x] [REQ-889-S2] define BuildInfo schema with required fields: version, commit, build_time
- [x] [REQ-889-S3] specify 405 Method Not Allowed for POST /buildinfo
- [x] [REQ-889-S4] specify 405 Method Not Allowed for PUT /buildinfo
- [x] [REQ-889-S5] specify 405 Method Not Allowed for DELETE /buildinfo
- [x] [REQ-889-S6] specify /buildinfo response identical to /version

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] [FEATURE-A1] GET /buildinfo returns 200 with JSON containing version, commit, build_time
- [x] [FEATURE-A2] Build info fields are meaningful (non-empty, correct format)
- [x] [FEATURE-A3] GET /buildinfo and GET /version return identical responses
- [x] [FEATURE-A4] POST /buildinfo returns 405
- [x] [FEATURE-A5] PUT /buildinfo returns 405
- [x] [FEATURE-A6] DELETE /buildinfo returns 405
- [x] [FEATURE-A7] GET /version still returns 200 after /buildinfo added
- [x] [FEATURE-A8] GET /healthz still returns 200 after /buildinfo added

## Stage: implementation (owner: dev-agent)
- [x] register `GET /buildinfo` route on admin mux in `cmd/server/server.go` reusing `version.Handler`
- [x] add mux-level unit test: GET /buildinfo returns 200 with correct JSON body (S1, S2)
- [x] add mux-level unit test: POST/PUT/DELETE /buildinfo return 405 (S3, S4, S5)
- [x] add mux-level unit test: GET /buildinfo and GET /version return identical response (S6)
- [x] verify `go vet`, `go build`, and unit tests pass on CI
