# REQ-997 Tasks

> 三段式：contract-tests → acceptance-tests → implementation。
> 本文件由 analyze-agent 立骨架；细节 TODO 由对应 stage 的 owner agent 自行展开（章节归属见 owner 标注）。

## Stage: contract-tests (owner: contract-spec-agent)
- [x] OpenAPI spec `crosser-api/contract.spec.yaml` — 9 endpoints: auth/login, services CRUD (5), config export, proxy register, proxy heartbeat
- [x] Contract test suite `crosser-api/tests/contract/` — 15 scenarios (REQ-997-S1..S15) mapped 1:1 to spec
- [x] Unified response envelope `{code, message, data}` + error code enum (1001..9999) locked in spec

### Scenarios → Specs:
- REQ-997-S1: Login success → `specs/auth/spec.md`
- REQ-997-S2: Login invalid credentials → `specs/auth/spec.md`
- REQ-997-S3: JWT auth rejection → `specs/auth/spec.md`
- REQ-997-S4: Create service → `specs/services/spec.md`
- REQ-997-S5: List services → `specs/services/spec.md`
- REQ-997-S6: Get service detail → `specs/services/spec.md`
- REQ-997-S7: Update service → `specs/services/spec.md`
- REQ-997-S8: Delete service → `specs/services/spec.md`
- REQ-997-S9: Service not found → `specs/services/spec.md`
- REQ-997-S10: Duplicate service → `specs/services/spec.md`
- REQ-997-S11: Config export → `specs/services/spec.md`
- REQ-997-S12: Proxy register → `specs/proxy/spec.md`
- REQ-997-S13: Proxy heartbeat → `specs/proxy/spec.md`
- REQ-997-S14: Proxy token rejection → `specs/proxy/spec.md`
- REQ-997-S15: Envelope structure → `specs/proxy/spec.md`

## Stage: acceptance-tests (owner: acceptance-spec-agent)
- [x] 产出 FEATURE 列表（覆盖 AC-03.01 / 02 / 03 / 04 共 4 条原始验收）
- [x] 用 Gherkin 描述每条 FEATURE 的 Given/When/Then 场景（FEATURE-A1..A15）
- [x] 标注每条 FEATURE 与 contract.spec.yaml 端点的反向追溯

### Capability: auth (AC-03.01) → specs/auth/spec.md
- FEATURE-A1: Admin login with valid credentials returns JWT token → POST /api/v1/auth/login
- FEATURE-A2: Admin login with invalid credentials returns 401 → POST /api/v1/auth/login
- FEATURE-A3: Accessing protected endpoint without JWT returns 401 → GET /api/v1/services (auth middleware)
- FEATURE-A4: Accessing protected endpoint with valid JWT succeeds → GET /api/v1/services (auth middleware)

### Capability: services-crud (AC-03.02) → specs/services-crud/spec.md
- FEATURE-A5: Create a new service via API → POST /api/v1/services
- FEATURE-A6: List all services → GET /api/v1/services
- FEATURE-A7: Update an existing service → PUT /api/v1/services/:id
- FEATURE-A8: Delete a service → DELETE /api/v1/services/:id

### Capability: proxy-registration (AC-03.03) → specs/proxy-registration/spec.md
- FEATURE-A9: Proxy instance registers with valid service token → POST /api/v1/proxy/register
- FEATURE-A10: Proxy heartbeat updates instance status → POST /api/v1/proxy/heartbeat
- FEATURE-A11: Query proxy instances shows online status → GET /api/v1/proxy/status
- FEATURE-A12: Proxy registration without valid token returns 401 → POST /api/v1/proxy/register

### Capability: config-export (AC-03.04) → specs/config-export/spec.md
- FEATURE-A13: Config export produces server.json-compatible output → GET /api/v1/services/:name/config
- FEATURE-A14: Config export for non-existent service returns 404 → GET /api/v1/services/:name/config
- FEATURE-A15: Exported config can be consumed by existing proxy-server → GET /api/v1/services/:name/config

## Stage: implementation (owner: dev-agent)
- [ ] TODO: 落地 `crosser-api` module 骨架（go.mod / cmd / internal/{handler,service,repository,middleware,model,database}）
- [ ] TODO: 实现 SQLite 初始化 + embed 迁移 `001_init.sql`
- [ ] TODO: 实现 auth（登录 / 刷新 / JWT 中间件 / bcrypt / 默认管理员首密策略）
- [ ] TODO: 实现 services CRUD + 配置导出（与现有 server.json 字节级兼容）
- [ ] TODO: 实现 proxy register / heartbeat / status（含 X-Proxy-Token 鉴权）
- [ ] TODO: 通过 contract-tests 与 acceptance-tests 全绿
