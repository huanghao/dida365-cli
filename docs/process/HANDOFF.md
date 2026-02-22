# Handoff

## Current
- Active Task: `T-004 实现 token refresh（验证与设计）`
- Status: `PENDING`
- Branch: `main`

## Ready-To-Run Next Steps
1. 对照 Dida 官方文档复核 refresh token 是否官方支持。
2. 用真实响应验证当前 `auth refresh` 行为是否可用。
3. 若不可用，输出设计方案（是否降级为 re-auth 流程）。
4. 更新 `docs/guides/dida-auth-token-flow.md`（记录 access token 过期时长观测值）。

## Blockers
- 无

## If Blocked Later
- 创建 `docs/process/design-gates/DG-XXXX-<topic>.md`
- 在 `TASK-BOARD.md` 将对应任务改为 `BLOCKED_DESIGN`
- 在本文件写清恢复条件
