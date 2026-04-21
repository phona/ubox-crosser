---
change_id: REQ-945
title: "webhook-debug endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [x] REQ-945-S1: POST /webhook-debug with JSON body → 200 + JSON with method/path/headers/body (`tests/contract/webhook_debug_test.go::TestPostWithJSONBody`)
- [x] REQ-945-S2: GET /webhook-debug?foo=bar&baz=qux → 200 + JSON with query params (`tests/contract/webhook_debug_test.go::TestGetWithQueryParams`)
- [x] REQ-945-S3: PUT /webhook-debug with form body → 200 + JSON with body content (`tests/contract/webhook_debug_test.go::TestPutWithBody`)
- [x] REQ-945-S4: DELETE /webhook-debug → 200 + JSON with empty body field (`tests/contract/webhook_debug_test.go::TestDeleteEmptyBody`)
- [x] REQ-945-S5: Response JSON schema validation — all 5 required fields with correct types (`tests/contract/webhook_debug_test.go::TestResponseSchemaValidation`)
- [x] OpenAPI contract spec: `contract.spec.yaml`
- [x] Contract test suite: `tests/contract/webhook_debug_test.go`

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
