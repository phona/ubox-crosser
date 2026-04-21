# REQ-983: 高层方案

## 方案
与 REQ-975 方案一致：在 Makefile 中新增 `stamp-readme` target，逻辑为：
1. 用 sed 删除 README.md 中已有的 `Built at:` 行（保证幂等）
2. 用 echo 追加新的 `Built at: <ISO-8601-UTC>` 行
3. README.md 不存在时静默跳过（exit 0）
4. build target 末尾调用 `@$(MAKE) stamp-readme`

## 选型
- 纯 Makefile + shell，零依赖
- 时间戳用 `date -u +%Y-%m-%dT%H:%M:%SZ` 生成

## 风险
- 无实质风险，改动局限于 Makefile

## skip_accept 说明
本 REQ 故意跳过 acceptance-tests 阶段，用于端到端验证 contract → implementation 简化链路。
acceptance-tests section 保留但标记为 skip。
