# DG-0003 Project License Choice

## Recommendation (给你先看这个)
推荐 `Option B (Apache-2.0)`。

理由：
- 明确专利授权条款，对开源 CLI 项目更稳妥
- 社区接受度高，商业友好
- 相比 MIT，法律边界更清晰

## Options (A/B/C)
- A: MIT
  - 最简洁，传播广
  - 不含明确专利授权条款
- B: Apache-2.0（推荐）
  - 包含专利授权与贡献者条款
  - 文本相对更长
- C: BSD-3-Clause
  - 宽松许可证，和 MIT 接近
  - 条款略多于 MIT，专利条款仍不如 Apache 明确

## What Changes (agent 写)
- 选择后执行：
  - 新增 `LICENSE`
  - 更新 `README.md` 的 License 段落
  - 在 `TODO.md` 标记 T-009 DONE 并记录 commit

## Decision (你只填这一行)
Decision: A|B|C
