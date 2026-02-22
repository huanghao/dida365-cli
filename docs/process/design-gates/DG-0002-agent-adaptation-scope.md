# DG-0002 Agent Adaptation Scope

## Recommendation (给你先看这个)
推荐 `Option B`：
- 增加一个轻量 `stats` 子命令给 agent 做读侧统计
- 统一关键错误提示格式（保持可机读）
- 明确 agent 操作限制清单（写操作防抖、输入长度等）

理由：
- 你当前核心目标是“agent 稳定执行高频操作”，B 的投入和收益最均衡
- 不引入重型遥测/审计系统，避免复杂度快速膨胀

## Options (A/B/C)
- A: 只补文档，不加新命令
  - 优点：最快
  - 风险：agent 仍需自行拼装统计，稳定性提升有限
- B: 增加轻量 stats + 统一限制/错误契约（推荐）
  - 新增 `dida stats --project <id> --json`
  - 输出任务总数、完成数、未完成数、优先级分布、近期待办计数（可从 project data 本地计算）
  - 文档固化 agent 限制与返回约定
- C: 完整遥测与执行日志体系
  - 优点：观测最完整
  - 风险：范围过大，偏离当前 CLI MVP 目标

## What Changes (agent 写)
- 命令层：
  - 新增 `dida stats`（先支持 `--project`）
- 输出层：
  - 默认 table + `--json` 结构化输出
- 文档层：
  - 更新 agent 手册与 man 文档，补“限制与统计”章节

## Decision (你只填这一行)
Decision: A|B|C
