# REQ-clean-1776863811: /version 端点返回 git SHA

## 一句话需求
实现 `/version` HTTP 端点，返回当前部署代码的 git SHA 和版本信息。

## 影响范围
- **API**: 新增 `/version` 端点（HTTP GET）
- **Deployment**: 版本信息需在构建/部署时注入
- **Client**: 可用于健康检查、版本验证、灰度发布等场景

## 支持度判定
✅ **可支持** — 已有完整的 contract 测试、acceptance 测试规范，代码改动范围明确，不涉及数据库迁移或外部依赖

## 设计要点
1. **版本数据来源**: git commit SHA（构建时注入）
2. **端点签名**: `GET /version` → 返回 JSON 格式的版本信息
3. **测试覆盖**: Contract spec（API 规范）+ Acceptance spec（e2e 验证）
4. **可靠性**: 幂等、无外部依赖、启动时可用
