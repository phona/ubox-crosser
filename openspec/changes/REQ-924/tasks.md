---
change_id: REQ-924
title: "GET /echo endpoint — tasks"
---

## Stage: contract-tests (owner: contract-test-agent)
- [ ] TODO: 列出要覆盖的 API 契约点（路径、方法、query 参数、响应格式、状态码、错误方法处理）

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
- [ ] TODO: 列出要实现的模块（echo 包 handler、路由注册、单元测试）
