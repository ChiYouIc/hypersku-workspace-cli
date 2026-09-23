---
name: hypersku-ai-analysis
display_name: AI 智能分析查询
display_name_en: AI Analysis Query
description_zh: 通过 hypersku-cli ai-analysis 子命令查询 HyperSKU 平台 AI 分析数据：国际物流异常与采购售后两个维度的分页分析结果及风险等级统计。
description_en: Query HyperSKU AI analysis data via hypersku-cli ai-analysis: paged analysis results and risk level statistics for international logistics anomalies and purchase after-sales.
description: HyperSKU AI 分析数据查询。当用户提到 AI分析/风险等级/国际物流异常/采购售后分析/风险统计时，通过 hypersku-cli ai-analysis 子命令查询 AI 分析结果（当前开放国际物流异常与采购售后两个维度）。
version: 1.1.0
author: owen
---

# AI 智能分析查询

通过 `hypersku-cli ai-analysis` 子命令查询 HyperSKU 平台 AI 智能分析结果，支持分页查询分析数据和风险等级统计。当前开放**国际物流异常**与**采购售后**两个维度。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 国际物流异常分析 | 分页查询国际物流异常 AI 分析结果（订单、物流、AI 摘要、风险等级） | `hypersku-cli ai-analysis intl-logistics [flags]` | [intl-logistics.md](references/intl-logistics.md) |
| 国际物流风险统计 | 统计国际物流异常风险等级分布 | `hypersku-cli ai-analysis intl-logistics-count [flags]` | [intl-logistics-count.md](references/intl-logistics-count.md) |
| 采购售后分析 | 分页查询采购售后 AI 分析结果（物流异常/1688售后来源） | `hypersku-cli ai-analysis purchase [flags]` | [purchase.md](references/purchase.md) |
| 采购风险统计 | 统计采购售后风险等级分布 | `hypersku-cli ai-analysis purchase-count [flags]` | [purchase-count.md](references/purchase-count.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"国际物流异常/AI分析/国际物流风险"时，执行 `intl-logistics` 分页查询国际物流异常 AI 分析结果。
- 用户提到"国际物流风险统计/国际物流风险等级"时，执行 `intl-logistics-count` 统计风险等级分布。
- 用户提到"采购售后分析/AI采购分析/采购风险"时，执行 `purchase` 分页查询采购售后分析。
- 用户提到"采购风险统计/采购风险等级"时，执行 `purchase-count` 统计采购风险等级分布。
- 若未提供必要参数，使用默认分页参数（page=1, limit=20）。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含任何写入/修改操作。
2. **分页参数**：所有分页命令支持 `--page`/`-p`（页码）和 `--limit`/`-l`（每页条数）参数。
3. **风险等级**：`intl-logistics` 的 `--risk-level` 支持 10=中风险、20=高风险（默认 20）。
4. **维度开放状态**：库存动销（inventory）、供应商（supplier）、客户订单（customer-order）、任务汇总（task-progress）四个维度的子命令当前版本**未启用**，请勿调用；用户问及这些维度时告知暂未开放。
