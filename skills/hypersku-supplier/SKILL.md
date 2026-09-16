---
name: hypersku-supplier
display_name: 供应商管理查询
display_name_en: Supplier Management Query
description_zh: 通过 hypersku-cli supplier 子命令查询供应商详情、列表、采购次数、排行数据、SPU 信息、售后和履约数据。
description_en: Query supplier details, lists, purchase counts, rankings, SPU data, after-sales and fulfillment data via hypersku-cli supplier commands.
description: HyperSKU 供应商管理查询。当用户提到供应商详情/供应商列表/供应商采购次数/供应商排行/供应商售后/供应商履约/供应商SPU时，通过 hypersku-cli supplier 子命令查询对应数据。
version: 1.0.0
author: owen
---

# 供应商管理查询

通过 `hypersku-cli supplier` 子命令查询供应商管理数据，支持详情、列表、采购次数、排行、售后、履约等多维度只读查询。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 供应商详情 | 根据 ID 查询供应商完整信息 | `hypersku-cli supplier detail <supplierId>` | [detail.md](references/detail.md) |
| 供应商列表 | 分页查询供应商列表 | `hypersku-cli supplier list [flags]` | [list.md](references/list.md) |
| 采购次数 | 根据登录 ID 查询采购次数 | `hypersku-cli supplier pur-count <loginId1,loginId2>` | [pur-count.md](references/pur-count.md) |
| 供应商排行 | 查询供应商排行数据 | `hypersku-cli supplier ranking [flags]` | [ranking.md](references/ranking.md) |
| SPU 信息 | 查询供应商 SPU 维度数据 | `hypersku-cli supplier spu-list [flags]` | [spu-list.md](references/spu-list.md) |
| 售后信息 | 查询供应商售后数据 | `hypersku-cli supplier after-sales [flags]` | [after-sales.md](references/after-sales.md) |
| 履约信息 | 查询供应商履约数据 | `hypersku-cli supplier fulfillment [flags]` | [fulfillment.md](references/fulfillment.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"供应商详情/供应商信息/供应商联系人"并提供 ID 时，执行 `detail <supplierId>` 查询供应商完整信息。
- 用户提到"供应商列表/供应商查询"时，执行 `list` 分页查询供应商列表（支持 --name 筛选）。
- 用户提到"采购次数/采购了几单"并提供登录 ID 时，执行 `pur-count <loginIds>` 查询采购次数。
- 用户提到"供应商排行/供应商排名/供应商数据"时，执行 `ranking` 查询供应商排行数据。
- 用户提到"供应商SPU/供应商商品"时，执行 `spu-list` 查询供应商 SPU 信息。
- 用户提到"供应商售后/供应商退款"时，执行 `after-sales` 查询供应商售后数据。
- 用户提到"供应商履约/供应商发货/供应商时效"时，执行 `fulfillment` 查询供应商履约数据。
- 若未提供必要参数，提示用户提供供应商 ID 或登录 ID。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含写入/修改操作。
2. **分页参数**：`list`、`ranking`、`spu-list`、`after-sales`、`fulfillment` 支持 `--page`/`-p` 和 `--limit`/`-l` 参数。
3. **排行排序**：`ranking` 默认按采购量降序排列，可通过 `--sort` 自定义。
4. **采购次数**：`pur-count` 支持批量查询，多个登录 ID 用逗号分隔。
