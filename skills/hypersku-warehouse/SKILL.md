---
name: hypersku-warehouse
display_name: 仓库与物流管理查询
display_name_en: Warehouse & Logistics Management Query
description_zh: 通过 hypersku-cli warehouse 子命令查询仓库物流轨迹、按名称查询仓库列表、按名称查询物流列表。
description_en: Query warehouse logistics tracking, search warehouses by name, and search logistics by name via hypersku-cli warehouse subcommands.
description: HyperSKU 仓库与物流管理查询。当用户提到仓库物流/仓库签收/入库/仓库操作/到仓/仓库列表/查仓库/物流列表/查物流时，通过 hypersku-cli warehouse 子命令查询对应数据。
version: 3.0.0
author: owen
tags:
  - hypersku
  - cli
  - 仓库
---

# 仓库物流轨迹查询

通过 `hypersku-cli warehouse` 子命令查询包裹的仓库侧物流轨迹、按名称搜索仓库列表、按名称搜索物流列表。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询仓库物流轨迹 | 根据物流单号查询仓库物流轨迹（快递签收、仓库签收、入库、物流轨迹、仓库操作） | `hypersku-cli warehouse tracking <trackingNumber>` | [tracking.md](references/tracking.md) |
| 按名称查询仓库 | 按仓库名称分页查询仓库列表 | `hypersku-cli warehouse list --name <name>` | [list.md](references/list.md) |
| 按名称查询物流 | 按物流名称分页查询物流列表 | `hypersku-cli warehouse logistics-list --name <name>` | [logistics-list.md](references/logistics-list.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"仓库物流/仓库签收/入库/仓库操作/到仓"时，执行 `tracking <trackingNumber>` 展示仓库侧物流轨迹。
- 用户提到"仓库列表/查仓库/搜索仓库"并提供仓库名称时，执行 `list --name <name>` 查询仓库列表。
- 用户提到"物流列表/查物流/搜索物流/物流线路"并提供物流名称时，执行 `logistics-list --name <name>` 查询物流列表。
- 若未提供必要参数，提示用户提供物流单号或搜索名称。

## 注意事项

1. **运单号必填**：查询需要传入物流单号参数，未提供时会显示帮助信息并退出。
2. **与快递轨迹不同**：`tracking` 命令返回仓库侧的快递签收、仓库签收、入库及仓库操作记录，与物流模块的快递运输轨迹不同。
3. **空字段**：仓库签收时间、入库时间为空时不展示。
