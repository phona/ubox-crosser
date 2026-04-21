---
change_id: REQ-945
title: "webhook-debug endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: POST /webhook-debug with JSON body returns 200 + application/json response containing method/path/headers/body
- [ ] TODO: GET /webhook-debug?foo=bar returns 200 + JSON with query params
- [ ] TODO: PUT /webhook-debug returns 200 + JSON with body content
- [ ] TODO: DELETE /webhook-debug returns 200 + JSON with empty body field
- [ ] TODO: response JSON schema validation (method, path, query, headers, body fields present)
- [ ] TODO: OpenAPI contract spec: contract.spec.yaml
- [ ] TODO: Contract test suite: tests/contract/webhook_debug_test.go

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] REQ-945-S1: POST webhook with JSON body, verify response contains correct method + body
- [x] REQ-945-S2: GET with query params, verify response query field populated
- [x] REQ-945-S3: PUT with form body, verify response body content
- [x] REQ-945-S4: DELETE with no body, verify body field is empty string
- [x] REQ-945-S5: Response JSON schema validation (all required fields present)

## Stage: implementation (owner: dev-agent)
- [ ] TODO: webhookdebug package (handler + RequestInfo struct)
- [ ] TODO: route registration in cmd/server/server.go
- [ ] TODO: unit tests for handler
