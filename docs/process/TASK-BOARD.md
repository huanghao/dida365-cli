# Task Board

状态枚举：`PENDING` | `IN_PROGRESS` | `BLOCKED_DESIGN` | `DONE`

| Task | Status | Owner | Last Update | Notes |
|---|---|---|---|---|
| T-001 实现文档目录结构，完成 dida.1.md | DONE | Codex | 2026-02-23 | 已完成，见 `TODO.md` |
| T-002 所有命令实现 `--json` 输出 | DONE | Codex | 2026-02-23 | 已完成，见 `TODO.md` |
| T-003 展示字段的枚举值，而不是裸数字 | DONE | Codex | 2026-02-23 | 枚举标签已落地（status/priority） |
| T-004 实现 token refresh（验证与设计） | DONE | Codex | 2026-02-23 | 已验证当前 app refresh 不可用，完成降级策略 |
| T-005 实现缓存，减少抖动 | DONE | Codex | 2026-02-23 | 按 DG-0001 A 实现写操作防抖（3 秒） |
| T-006 整理展示字段和 API 对应关系 | DONE | Codex | 2026-02-23 | 已补字段映射文档并增强 show 字段 |
| T-007 限制创建 task 输入长度（<500） | DONE | Codex | 2026-02-23 | add 命令已加字符长度校验与测试 |
| T-008 适配 agent 的功能/限制/统计 | PENDING | - | - | 可能涉及设计闸门 |
| T-009 项目增加 license 信息 | PENDING | - | - | |
| T-010 打包和发布 | PENDING | - | - | |
| T-011 面向 agent 的导引文件完善 | PENDING | - | - | |
