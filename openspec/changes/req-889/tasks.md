---
change_id: req-889
title: "GET /buildinfo endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: define OpenAPI contract for GET /buildinfo (same schema as /version)

## Stage: acceptance-tests (owner: accept-test-agent)
- [ ] TODO: verify GET /buildinfo returns 200 with correct JSON schema
- [ ] TODO: verify non-GET methods return 405

## Stage: implementation (owner: dev-agent)
- [ ] TODO: register GET /buildinfo route in cmd/server/server.go
- [ ] TODO: add mux-level routing tests for /buildinfo
