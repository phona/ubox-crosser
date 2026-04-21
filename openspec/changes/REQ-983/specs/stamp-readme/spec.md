---
capability: stamp-readme
change_id: REQ-983
status: LOCKED
---

## ADDED

### Scenario: REQ-983-S1 — make stamp-readme appends Built-at timestamp to README.md

```gherkin
Given a README.md file exists in the project root
When the user runs "make stamp-readme"
Then README.md contains exactly one line matching "Built at: <ISO-8601-UTC>"
  and the timestamp format matches "^Built at: [0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$"
```

### Scenario: REQ-983-S2 — stamp-readme is idempotent (no duplicate lines)

```gherkin
Given a README.md file exists in the project root
  and README.md already contains a "Built at:" line
When the user runs "make stamp-readme" again
Then README.md contains exactly one line matching "Built at:"
  and the old timestamp line is replaced by the new one
```

### Scenario: REQ-983-S3 — stamp-readme exits 0 when README.md does not exist

```gherkin
Given no README.md file exists in the project root
When the user runs "make stamp-readme"
Then the exit code is 0
  and no README.md file is created
```
