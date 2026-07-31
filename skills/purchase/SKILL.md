---
name: purchase
description: Hypersku 采购订单管理能力。当用户需要查询采购订单详情、分页搜索订单、采购日志、国际物流轨迹时使用。
version: 1.0.0
author: owen
tags:
  - 订单商品
  - 分页查询
  - 国际物流
  - 采购日志
---

# 采购订单管理 Skill

通过 `hypersku-cli purchase` 子命令管理 Hypersku 采购订单，支持按单号查询订单详情、按条件分页搜索订单、查看采购日志及国际物流轨迹。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询订单详情 | 根据订单号查询商品明细、金额、订单类型/状态、仓库等 | `hypersku-cli purchase info <orderId>` | [info.md](references/info.md) |
| 分页搜索订单 | 按日期/交易号/物流单号分页查询订单列表（仅概要，不含商品明细） | `hypersku-cli purchase info page [flags]` | [page.md](references/page.md) |
| 查询采购日志 | 根据订单号查询订单状态流转日志 | `hypersku-cli purchase log <orderId>` | [log.md](references/log.md) |
| 查询国际物流 | 根据订单号查询包裹国际段物流轨迹（承运商/单号/轨迹） | `hypersku-cli purchase logistics <orderId>` | [logistics.md](references/logistics.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"订单详情/订单信息/商品/买了什么 + 具体订单号"时，执行 `info <orderId>` 展示订单概要及商品明细。
- 用户提到"搜索/查找/列表/分页/有哪些订单 + 日期/交易号/物流单号"时，执行 `info page` 并通过 flags 指定筛选条件分页查询。
- 用户提到"日志/状态/进度/到哪一步了/流程"时，执行 `log <orderId>` 展示状态流转时间线。
- 用户提到"物流/轨迹/快递/运输/到哪了/包裹"时，执行 `logistics <orderId>` 展示包裹国际物流轨迹。
- 若未提供必要参数，提示用户补全。

## 注意事项

1. **orderId 必填**：所有子命令都需要传入订单号参数，这个订单号既可以是采购订单号也可以是客户订单号，未提供时会显示帮助信息并退出。
3. **国际化**：地址和物流信息可能包含英文内容，具体取决于订单来源国家。
4. **多包裹**：国际物流查询可能返回多个包裹的轨迹信息，每个包裹独立展示。


