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
- [ ] TODO: verify GET /buildinfo returns 200 with correct JSON schema
- [ ] TODO: verify non-GET methods return 405

## Stage: implementation (owner: dev-agent)
- [ ] TODO: register GET /buildinfo route in cmd/server/server.go
- [ ] TODO: add mux-level routing tests for /buildinfo
