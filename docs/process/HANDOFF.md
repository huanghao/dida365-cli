# Handoff

## Current
- Active Task: `T-007 限制创建 task 输入长度（<500字）`
- Status: `PENDING`
- Branch: `main`

## Ready-To-Run Next Steps
1. 在 `add/update` 输入校验阶段增加长度限制（<=500）。
2. 为校验逻辑补单元测试（边界值：499/500/501）。
3. 更新 README/man（明确限制仅作用于创建/更新输入）。
4. 更新 TODO 与 CHANGELOG。

## Blockers
- T-005 仍阻塞：DG-0001 需要你确认缓存/防抖方案

## If Blocked Later
- 若当前任务阻塞，优先继续下一个不依赖设计闸门的任务。
