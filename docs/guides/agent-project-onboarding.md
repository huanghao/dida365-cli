# Agent 项目导引（最小必读）

更新时间：2026-02-23

面向“新 agent 进入项目后 5 分钟内可开始执行任务”的导引。

## 1. 先读这 5 个文件

1. `TODO.md`
- 唯一任务源，按顺序执行

2. `docs/process/TASK-BOARD.md`
- 查看每个任务状态（`PENDING/IN_PROGRESS/BLOCKED_DESIGN/DONE`）

3. `docs/process/HANDOFF.md`
- 当前活跃任务和下一步操作

4. `docs/process/design-gates/README.md`
- 设计闸门确认方式（你只需写 `Decision: A|B|C`）

5. `docs/guides/agent-cli-quick-manual.md`
- 高频命令与强约束（默认 `--json`、写操作先 `--dry-run`）

## 2. 执行规则（必须）

- 不改动未确认的设计闸门结论
- 完成任务后必须同步：
  - `TODO.md`
  - `docs/process/TASK-BOARD.md`
  - `docs/process/HANDOFF.md`
- 若发现跨命令高影响设计项：新建 DG 并阻塞任务

## 3. 常用本地命令

```bash
go run ./cmd/dida --help
go test ./...
```

若 shell 未初始化 goenv，先执行：

```bash
export GOENV_ROOT="$HOME/.goenv"
export PATH="$GOENV_ROOT/bin:$PATH"
eval "$(goenv init -)"
```

## 4. 推荐新增文件（后续可选）

- `docs/process/release-checklist.md`：版本发布 checklist
- `docs/guides/troubleshooting.md`：高频错误排查
- `docs/guides/agent-stats-schema.md`：若启用 `stats`，定义固定 JSON 契约
