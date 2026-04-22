# REQ-final2-1776868985: /version 端点返回 git SHA

## 一句话需求

为 ubox-crosser 新增 HTTP `/version` 端点，返回当前构建的 git commit SHA。

## 影响范围

**代码层面：**
- API 路由层：新增 `/version` 端点处理器
- 编译配置：通过 `-ldflags` 在构建时注入 git SHA
- 可能涉及 Makefile 调整，确保 CI/CD pipeline 能正确获取和注入 SHA

**测试层面：**
- Contract spec：定义 `/version` 端点的 HTTP 接口（状态码、响应格式）
- Acceptance spec：验证端点返回正确的 git SHA，且在不同构建场景下一致

**无 breaking changes：**
- 仅新增端点，不修改现有 API

## 支持性判定

✅ **可顺利实施，风险低**

- **技术成熟度**：Go 的 `-ldflags` 注入版本信息是业界标准做法
- **项目现状**：ubox-crosser 已具备 contract-spec 和 acceptance-spec 框架（见最近提交），无需从头建设
- **依赖状况**：不依赖任何新的外部库，仅涉及原生 Go 功能
- **工作量估算**：～100 lines Go code + 测试，分三阶段约 0.5-1 day

## 核心任务

1. **Contract Spec**：定义 `/version` 的 HTTP 接口规范（JSON 格式、必返字段）
2. **Acceptance Spec**：编写 e2e 场景，验证端点在多种构建方式下的行为
3. **Implementation**：添加路由、编译时注入、CI 配置调整
