# 协作执行机制（TODO 驱动）

本机制用于你与任意 agent（包括跨会话、跨 agent）协作完成 `TODO.md`。

## 1. 单一任务源
- 任务源唯一：`TODO.md`
- 执行顺序：从上到下，按第一个未 `DONE` 的任务开始
- 任务完成后必须：
  - 在 `TODO.md` 标记 `DONE`
  - 写入完成时间
  - 写入 commit hash
  - 更新 `CHANGELOG.md`

## 2. 执行状态机
- `PENDING`：未开始
- `IN_PROGRESS`：正在实现
- `BLOCKED_DESIGN`：需要你确认设计
- `DONE`：已完成并提交

状态同步位置：`docs/process/TASK-BOARD.md`

## 3. 遇到设计分歧时的规则
默认策略：`先继续下一个不依赖该设计的任务`。

仅在以下情况必须停下并等待确认：
- 该设计会影响多个后续任务的基础模型（例如数据模型、接口契约、目录结构）
- 该设计存在不可逆迁移成本（数据格式、公开命令参数）
- 不确认就会造成高概率返工

如果满足“必须停下”条件：
- 创建设计闸门文档：`docs/process/design-gates/DG-XXXX-<topic>.md`
- 将任务状态改为 `BLOCKED_DESIGN`
- 在 `docs/process/HANDOFF.md` 写明阻塞点和恢复条件

如果不满足“必须停下”条件：
- 在设计闸门文档记录假设
- 标记为“临时决策，可回滚”
- 继续执行下一个任务

## 4. 设计闸门确认方式（极简）
你只需要做一件事：
- 在 DG 文档最后写一行：`Decision: A`（或 `B` / `C`）

或在对话中直接回复：
- `Decision: A`
- `CONFIRM DG-XXXX: A`

其余内容（背景、影响、回滚）由 agent 负责维护。

## 5. 跨 agent 交接规范
每次暂停或完成阶段后，更新：
- `docs/process/HANDOFF.md`

必须写清：
- 当前任务 ID 与状态
- 已完成内容
- 未完成内容
- 下一步可直接执行命令
- 若阻塞，写明需要你的确认点

## 6. 文档职责边界
- `TODO.md`：任务列表与完成打勾
- `docs/process/TASK-BOARD.md`：任务状态流转
- `docs/process/design-gates/*.md`：设计决策与确认记录
- `docs/process/HANDOFF.md`：跨会话交接
- `CHANGELOG.md`：完成结果与提交记录

## 7. 快速执行规则（给 agent）
1. 读取 `TODO.md` 和 `TASK-BOARD.md`
2. 执行第一个未完成任务
3. 只在“必须停下”条件下发起设计闸门
4. 完成后更新 `TODO.md`、`CHANGELOG.md`、`HANDOFF.md`
