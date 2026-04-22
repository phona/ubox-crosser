---
capability: services
change_id: REQ-997
status: LOCKED
---

## ADDED

### Scenario: REQ-997-S4 — POST /api/v1/services creates a service and returns 201

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
When the client sends POST /api/v1/services with name="test-svc", key="k1", method="aes-256-cfb", address=":8388"
Then the response status code is 201
  and the response JSON field "code" is 0
  and the response JSON field "data.name" is "test-svc"
  and the response JSON field "data.created_at" is a valid RFC3339 timestamp
```

### Scenario: REQ-997-S5 — GET /api/v1/services returns service list

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and at least one service exists
When the client sends GET /api/v1/services
Then the response status code is 200
  and the response JSON field "code" is 0
  and the response JSON field "data.services" is a non-empty array
  and the response JSON field "data.total" is a positive integer
  and each item in "data.services" has required fields: name, key, method, address
```

### Scenario: REQ-997-S6 — GET /api/v1/services/{name} returns service details with proxy instances

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and a service "test-svc" exists
When the client sends GET /api/v1/services/test-svc
Then the response status code is 200
  and the response JSON field "code" is 0
  and the response JSON field "data.service.name" is "test-svc"
  and the response JSON field "data.proxy_instances" is an array
```

### Scenario: REQ-997-S7 — PUT /api/v1/services/{name} updates service fields

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and a service "test-svc" exists
When the client sends PUT /api/v1/services/test-svc with key="new-key"
Then the response status code is 200
  and the response JSON field "code" is 0
  and the response JSON field "data.key" is "new-key"
  and the response JSON field "data.updated_at" differs from created_at
```

### Scenario: REQ-997-S8 — DELETE /api/v1/services/{name} deletes a service

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and a service "test-svc" exists
When the client sends DELETE /api/v1/services/test-svc
Then the response status code is 200
  and the response JSON field "code" is 0
When the client sends GET /api/v1/services/test-svc
Then the response status code is 404
```

### Scenario: REQ-997-S9 — GET /api/v1/services/{name} for non-existent service returns 404

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
When the client sends GET /api/v1/services/no-such-svc
Then the response status code is 404
  and the response JSON field "code" is 2001
```

### Scenario: REQ-997-S10 — POST /api/v1/services with duplicate name returns 409

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and a service "test-svc" exists
When the client sends POST /api/v1/services with name="test-svc"
Then the response status code is 409
  and the response JSON field "code" is 2002
```

### Scenario: REQ-997-S11 — GET /api/v1/services/{name}/config returns server.json compatible format

```gherkin
Given the crosser-api server is running
  and the client has a valid admin JWT
  and a service "test-svc" exists with key="k1", method="aes-256-cfb", address=":8388"
When the client sends GET /api/v1/services/test-svc/config
Then the response status code is 200
  and the response JSON field "code" is 0
  and the response JSON field "data.common.key" is "k1"
  and the response JSON field "data.common.method" is "aes-256-cfb"
  and the response JSON field "data.common.address" is ":8388"
  and the response JSON field "data.common" has fields: login_password, auth_password, log_file, log_level
  and the response JSON has a key "data.test-svc" with field "key"
```
