# Handoff

## Current
- Active Task: `T-005 实现缓存，减少抖动`
- Status: `BLOCKED_DESIGN`
- Branch: `main`

## Ready-To-Run Next Steps
1. 设计短期缓存策略（查询命令缓存、写命令防抖）。
2. 明确缓存介质、TTL、失效策略与并发安全策略。
3. 输出设计闸门文档（若涉及跨命令一致性与不可逆改动）。 
4. 确认后实现并补测试。

## Blockers
- DG-0001 需要你确认缓存/防抖方案

## If Blocked Later
- 已阻塞：等待你在 `docs/process/design-gates/DG-0001-cache-and-debounce.md` 填写
  `Decision: A|B|C`
