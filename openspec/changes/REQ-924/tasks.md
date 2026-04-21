---
change_id: REQ-924
title: "GET /echo endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [x] [REQ-924-S1] GET /echo?msg=hello returns 200 with body "hello" and Content-Type text/plain
- [x] [REQ-924-S2] GET /echo?msg= returns 200 with empty body
- [x] [REQ-924-S3] GET /echo (no msg param) returns 200 with empty body
- [x] [REQ-924-S4] POST/PUT/DELETE /echo returns 405 Method Not Allowed
- [x] OpenAPI contract spec: contract.spec.yaml
- [x] Contract test suite: tests/contract/echo_test.go

## Stage: acceptance-tests (owner: accept-test-agent)
- [x] [FEATURE-A1] 正常回显：GET /echo?msg=hello 返回 200 + text/plain body "hello"
- [x] [FEATURE-A2] 特殊字符回显：GET /echo?msg=hello%20world%21 返回 URL 解码后的原文
- [x] [FEATURE-A3] 空 msg 参数：GET /echo?msg= 返回 200 + 空 body
- [x] [FEATURE-A4] 缺失 msg 参数：GET /echo 返回 200 + 空 body
- [x] [FEATURE-A5] 非法方法 POST：POST /echo 返回 405
- [x] [FEATURE-A6] 非法方法 PUT：PUT /echo 返回 405
- [x] [FEATURE-A7] 共存验证：GET /ping 仍正常
- [x] [FEATURE-A8] 共存验证：GET /healthz 仍正常
- [x] [FEATURE-A9] 共存验证：GET /version 仍正常

## Stage: implementation (owner: dev-agent)
- [x] echo/handler.go — Handler 实现（读 msg query param，返回 text/plain 200）
- [x] cmd/server/server.go — 注册 GET /echo 路由到 admin mux
- [x] echo/handler_test.go — 单元测试（WithMsg、SpecialChars、EmptyMsg、NoMsgParam）
- [x] go vet + make build 编译通过
- [x] make ci-unit-test 全部绿色（echo coverage 100%）
- [x] 契约测试 tests/contract/echo_test.go 全部通过
