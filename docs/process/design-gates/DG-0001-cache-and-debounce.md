# DG-0001 Cache and Debounce Strategy

## Recommendation (给你先看这个)
推荐 `Option B`：
- 只对只读查询命令做短 TTL 本地缓存
- 对写操作做“幂等键 + 短窗口防抖”
- 提供全局开关禁用缓存

理由：
- 降低 API 抖动，避免 agent 高频重复读取
- 减少误触导致的短时间重复创建
- 风险和复杂度可控，不改变主要命令语义

## Options (A/B/C)
- A: 不做本地缓存，只做写操作防抖
  - 优点：实现最简单
  - 风险：读请求抖动问题仍在
- B: 读缓存 + 写防抖（推荐）
  - 读命令（`projects list`, `list`, `show`）缓存 10s
  - 写命令（`add`, `update`, `done`, `delete`, `projects create`）做 3s 防抖
  - 支持 `--no-cache` 或 `DIDA_NO_CACHE=1`
- C: 全量缓存（读写都缓存）
  - 风险：一致性复杂，错误概率高，不推荐

## What Changes (agent 写)
- 新增本地缓存文件：`~/.config/dida365-cli/cache.json`
- 缓存键：`method + path + normalized_query + body_hash`
- 缓存范围：仅 GET
- TTL：默认 10 秒，可配置
- 写操作防抖：
  - key = `command + normalized args`
  - 窗口默认 3 秒
  - 命中时返回结构化提示（JSON 模式下可机器读）
- 失效策略：
  - 任何写操作成功后，清理相关读缓存（简单策略：清空全部）

## Decision (你只填这一行)
Decision: A
