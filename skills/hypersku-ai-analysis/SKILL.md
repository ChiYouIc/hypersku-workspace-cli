---
name: hypersku-ai-analysis
display_name: AI 智能分析查询
display_name_en: AI Analysis Query
description_zh: 通过 hypersku-cli ai-analysis 子命令查询 HyperSKU 平台各维度 AI 分析数据：国际物流异常、库存动销、采购售后、供应商、客户订单的分析结果与风险等级统计。
description_en: Query AI analysis results across HyperSKU dimensions via hypersku-cli ai-analysis: international logistics anomalies, inventory sell-through, purchase after-sales, supplier, and customer order analysis with risk level statistics.
description: HyperSKU AI 分析数据查询。当用户提到 AI分析/风险等级/国际物流异常/库存动销/采购售后分析/供应商分析/客户订单分析/风险统计时，通过 hypersku-cli ai-analysis 子命令查询 AI 分析结果。
version: 1.0.0
author: owen
---

# AI 智能分析查询

通过 `hypersku-cli ai-analysis` 子命令查询 HyperSKU 平台各维度 AI 智能分析结果，支持分页查询分析数据和风险等级统计。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 国际物流异常分析 | 分页查询国际物流异常 AI 分析结果 | `hypersku-cli ai-analysis intl-logistics [flags]` | [intl-logistics.md](references/intl-logistics.md) |
| 国际物流风险统计 | 统计国际物流异常风险等级分布 | `hypersku-cli ai-analysis intl-logistics-count [flags]` | [intl-logistics-count.md](references/intl-logistics-count.md) |
| 库存动销分析 | 分页查询库存动销 AI 分析结果 | `hypersku-cli ai-analysis inventory [flags]` | [inventory.md](references/inventory.md) |
| 库存风险统计 | 统计库存动销风险等级分布 | `hypersku-cli ai-analysis inventory-count [flags]` | [inventory-count.md](references/inventory-count.md) |
| 采购售后分析 | 分页查询采购售后 AI 分析结果 | `hypersku-cli ai-analysis purchase [flags]` | [purchase.md](references/purchase.md) |
| 采购风险统计 | 统计采购售后风险等级分布 | `hypersku-cli ai-analysis purchase-count [flags]` | [purchase-count.md](references/purchase-count.md) |
| 供应商 AI 分析 | 分页查询供应商 AI 分析结果 | `hypersku-cli ai-analysis supplier [flags]` | [supplier.md](references/supplier.md) |
| 供应商等级统计 | 统计供应商等级分布 | `hypersku-cli ai-analysis supplier-count [flags]` | [supplier-count.md](references/supplier-count.md) |
| 客户订单分析 | 分页查询客户订单 AI 分析结果 | `hypersku-cli ai-analysis customer-order [flags]` | [customer-order.md](references/customer-order.md) |
| 客户订单风险统计 | 统计客户订单风险等级分布 | `hypersku-cli ai-analysis customer-order-count [flags]` | [customer-order-count.md](references/customer-order-count.md) |
| AI 任务风险汇总 | 按分析类型统计 AI 任务风险等级分布 | `hypersku-cli ai-analysis task-progress <type>` | [task-progress.md](references/task-progress.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"国际物流异常/AI分析/国际物流风险"时，执行 `intl-logistics` 分页查询国际物流异常 AI 分析结果。
- 用户提到"国际物流风险统计/国际物流风险等级"时，执行 `intl-logistics-count` 统计风险等级分布。
- 用户提到"库存动销/库存AI分析/滞销/库存风险"时，执行 `inventory` 分页查询库存动销分析。
- 用户提到"库存风险统计/库存风险等级"时，执行 `inventory-count` 统计库存风险等级分布。
- 用户提到"采购售后分析/AI采购分析/采购风险"时，执行 `purchase` 分页查询采购售后分析。
- 用户提到"采购风险统计/采购风险等级"时，执行 `purchase-count` 统计采购风险等级分布。
- 用户提到"供应商AI分析/供应商评估/供应商风险"时，执行 `supplier` 分页查询供应商分析。
- 用户提到"供应商等级统计"时，执行 `supplier-count` 统计供应商等级分布。
- 用户提到"客户订单AI分析/客户风险/客户订单风险"时，执行 `customer-order` 分页查询客户订单分析。
- 用户提到"客户订单风险统计"时，执行 `customer-order-count` 统计客户订单风险等级分布。
- 用户提到"AI任务汇总/风险汇总/任务进度"时，执行 `task-progress <type>` 按类型统计风险等级。
- 若未提供必要参数，使用默认分页参数（page=1, limit=20）。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含任何写入/修改操作。
2. **分页参数**：所有分页命令支持 `--page`/`-p`（页码）和 `--limit`/`-l`（每页条数）参数。
3. **风险等级**：风险等级分为 0=无风险、1=低风险、2=中风险、3=高风险。
4. **task-progress 类型**：支持 `international_logistics`、`inventory`、`purchase`、`supplier`、`customer_order` 五种类型。
