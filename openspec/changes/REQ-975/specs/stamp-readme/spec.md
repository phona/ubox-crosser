---
capability: stamp-readme
change_id: REQ-975
status: LOCKED
---

## ADDED

### Scenario: FEATURE-A1 — make build 后 README.md 末尾出现 Built at 时间戳行

```gherkin
Given the project source tree is checked out
  and README.md exists in the project root
When the user runs make build
Then README.md ends with a line matching pattern "^Built at: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$"
```

### Scenario: FEATURE-A2 — 时间戳为合法的 ISO 8601 UTC 格式

```gherkin
Given the project source tree is checked out
When the user runs make build
Then the "Built at:" line in README.md contains a valid ISO 8601 UTC timestamp
  and the timestamp is within 60 seconds of the current UTC time
```

### Scenario: FEATURE-A3 — 多次构建不会重复追加时间戳行（幂等性）

```gherkin
Given the project source tree is checked out
When the user runs make build twice consecutively
Then README.md contains exactly one line starting with "Built at:"
  and that line reflects the timestamp of the second build
```

### Scenario: FEATURE-A4 — README.md 不存在时构建仍然成功

```gherkin
Given the project source tree is checked out
  and README.md does not exist in the project root
When the user runs make build
Then the build completes successfully with exit code 0
  and no "Built at:" line is created (README.md remains absent)
```

### Scenario: FEATURE-A5 — stamp-readme target 可独立调用

```gherkin
Given the project source tree is checked out
  and README.md exists in the project root
When the user runs make stamp-readme
Then README.md ends with a line matching pattern "^Built at: \d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$"
  and no binary artifacts are produced
```
