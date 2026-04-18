---
description: 生成提交消息
argument-hint: [type|message]
allowed-tools: Bash(git status:*), Bash(git diff:*), Bash(git commit:*), Bash(git branch:*), Bash(git log:*)
---

# 上下文

- 暂存状态: !`git diff --staged --name-only`
- 当前分支: !`git branch --show-current`
- Git 状态: !`git status`
- 暂存更改: !`git diff --staged`
- 最近提交: !`git log --oneline -5`

# 任务

基于上述暂存的更改生成规范的 Git 提交消息。

## 前置检查

**如果暂存状态为空（无暂存更改），直接退出并告知用户先使用 `git add` 暂存文件，不要尝试自动 git add。**

## 用户提示

可选参数 `$1` 作为提交类型或描述提示：
- 如提供，作为生成提交消息的参考
- 如未提供，完全自动分析

## 提交类型

必须是以下之一：
- `feat` - 新功能
- `fix` - Bug 修复
- `docs` - 文档更新
- `style` - 代码格式
- `refactor` - 重构
- `test` - 测试
- `perf` - 性能优化
- `build` - 构建系统/依赖更新
- `ci` - CI 配置
- `chore` - 其他杂项

## 提交消息格式

格式: `{type}({scope}): {description}` 或 `{type}: {description}`

- `{type}`: 必须是上述类型之一
- `{scope}`: 可选，受影响的包（如 server、client、connector）
- `{description}`: 1-72 字符的简短描述

## 生成策略

1. 分析暂存文件的路径，识别受影响的包作为 scope
2. 根据代码变更类型确定 type
3. 用简洁的中文描述变更内容（1-72 字符）
4. 参考最近提交的风格保持一致性
5. 执行 `git commit` 提交
