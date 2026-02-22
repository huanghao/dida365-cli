# Design Gates

命名规范：`DG-XXXX-<topic>.md`

你只需要做一件事：
- 在 DG 文档最后写：`Decision: A`（或 `B` / `C`）

其他内容（背景、影响、回滚）都由 agent 负责填写。

极简模板：

```md
# DG-XXXX <Title>

## Recommendation (给你先看这个)
...

## Options (A/B/C)
- A: ...
- B: ...
- C: ...

## What Changes (agent 写)
...

## Decision (你只填这一行)
Decision: A
```

对话里也可以直接回：
- `Decision: A`
- 或 `CONFIRM DG-XXXX: A`
