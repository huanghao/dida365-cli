# Dida Task 字段与 CLI 展示映射

更新时间：2026-02-23

目标：
- 对齐 Dida OpenAPI Task 字段与当前 `dida` 的展示行为
- 识别高价值但未展示字段
- 标记当前代码里“暂未用于表格输出但仍有价值”的字段

## 1. 字段映射总览

| API 字段 | list(table) | show(table) | `--json` | 说明 |
|---|---|---|---|---|
| `id` | ✅ | ✅ | ✅ | 任务主键 |
| `projectId` | - | ✅ | ✅ | 项目归属 |
| `title` | ✅ | ✅ | ✅ | 任务标题 |
| `status` | ✅（枚举） | ✅（枚举） | ✅（原值） | 0/2 等状态值 |
| `completedTime` | 间接用于 completed 列 | ✅ | ✅ | 完成时间 |
| `dueDate` | ✅ | ✅ | ✅ | 高频字段 |
| `priority` | ✅（枚举） | ✅（枚举） | ✅（原值） | 0/1/3/5 |
| `content` | ✅（截断预览） | ✅ | ✅ | 高频字段 |
| `desc` | - | ✅ | ✅ | 长描述 |
| `startDate` | - | ✅ | ✅ | 起始时间 |
| `isAllDay` | - | ✅ | ✅ | 全天任务 |
| `timeZone` | - | ✅ | ✅ | 时区 |
| `repeatFlag` | - | ✅ | ✅ | 重复规则 |
| `sortOrder` | - | - | ✅ | 排序相关 |
| `reminders` | - | - | ✅ | 提醒配置 |
| `items` | - | - | ✅ | 检查项/子项 |

说明：
- `list` 设计为“高密度总览”，不扩展过多列，避免 agent 解析复杂度上升。
- `show` 设计为“单任务详情”，补齐了 `isAllDay/timeZone/repeatFlag`。

## 2. 高频字段覆盖结论

根据当前使用场景（内容、完成状态、due、priority）：
- 已覆盖：`content`、`status`、`completedTime`、`dueDate`、`priority`
- 通过 `show --json` 可获得全量 Task 字段，适合 agent 精确处理

## 3. 未展示字段评估（不是无用字段）

以下字段当前未进表格，但不建议删除：
- `reminders`：后续“提醒管理”能力需要
- `items` 与 `TaskItem.*`：后续子任务/清单能力需要
- `sortOrder`：列表排序与同步对齐可用

结论：
- 当前代码中没有“明显无用且应删除”的 Task 字段定义。
- 未进表格输出不等于无用；这些字段仍通过 `--json` 暴露，且为后续功能预留。

## 4. 建议

1. 对 agent 场景继续坚持：写操作和精确读操作优先 `--json`。
2. 若后续新增筛选命令（如按 `dueDate` 或 `priority`），优先在命令层做过滤，不先扩展 `list` 列数。
3. 如要进一步压缩 agent 解析成本，可在后续新增 `dida show --brief`（固定小字段集）。
