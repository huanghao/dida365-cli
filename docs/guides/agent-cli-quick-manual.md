# dida cli 高频使用手册

适用场景：
- 人已完成认证（`auth init/login/token/status`）
- 后续由 agent 执行任务创建/查询/修改/完成/删除

强约束：
- Agent 默认使用 JSON 输出
- 所有写操作先用 `--dry-run` 预检，再执行正式命令

## 0. 前置检查

```bash
dida auth status --json
```

确认：
- `access_token_set: true`
- `client_id_set: true`
- `redirect_uri_set: true`

## 1. 高频命令

列项目（JSON）：
```bash
dida projects list --json
```

创建项目：
```bash
dida projects create --name "Roadmap" --view-mode list --kind TASK --json
```

按项目列任务（JSON）：
```bash
dida list --project <project_id> --json
```

查任务详情（JSON）：
```bash
dida show --project <project_id> --id <task_id> --json
```

关注字段：
- `title`
- `content`
- `status`
- `completedTime`
- `dueDate`
- `priority`

创建任务：
```bash
dida add --project <project_id> --title "Task title" --content "Task content" --due "2026-02-25T18:00:00+0800" --priority 3 --json
```

创建输入限制：
- `--title`、`--content`、`--desc` 需 `< 500` 字符。

更新任务：
```bash
dida update --project <project_id> --id <task_id> --title "New title" --content "New content" --due "2026-02-26T18:00:00+0800" --priority 5 --json
```

完成任务：
```bash
dida done --project <project_id> --id <task_id> --json
```

删除任务：
```bash
dida delete --project <project_id> --id <task_id> --json
```
