---
name: hypersku-logistics
display_name: 物流轨迹查询
display_name_en: Logistics Tracking
description_zh: 通过 hypersku-cli logistics tracking 命令根据运单号查询包裹的国内物流轨迹，返回完整运输节点记录（时间、城市、事件）。
description_en: Query domestic parcel logistics tracking by tracking number via hypersku-cli logistics tracking, returning the full shipment timeline with time, city, and events.
description: HyperSKU 物流轨迹查询。当用户提到物流/轨迹/快递/运输/到哪了/包裹/tracking 并提供运单号时，执行 hypersku-cli logistics tracking 查询包裹的国内物流轨迹。
version: 2.0.0
author: owen
tags:
  - hypersku
  - cli
  - 物流
---

# 物流轨迹查询

通过 `hypersku-cli logistics` 子命令查询包裹的国内物流轨迹，支持根据运单号获取完整的运输节点记录。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询物流轨迹 | 根据运单号查询包裹的完整物流轨迹（时间、城市、事件） | `hypersku-cli logistics tracking <trackingNumber>` | [get-tracking.md](references/get-tracking.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"物流/轨迹/快递/运输/到哪了/包裹/tracking"时，执行 `tracking <trackingNumber>` 以表格展示物流时间线。
- 用户提到"单号/运单/tracking number"但未提供完整单号时，提示用户提供运单号。
- 若未提供 trackingNumber，提示用户提供物流单号。

## 注意事项

1. **运单号必填**：所有查询都需要传入物流单号参数，未提供时会显示帮助信息并退出。
2. **物流商**：当前支持的物流商有限，具体取决于 hypersku 平台对接的物流渠道。
3. **轨迹顺序**：轨迹按时间正序排列，从揽件到签收，可能包含多条记录。
