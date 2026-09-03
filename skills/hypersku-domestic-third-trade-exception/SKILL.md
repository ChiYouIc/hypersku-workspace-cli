---
name: hypersku-domestic-third-trade-exception
display_name: 国内第三方交易异常订单管理
display_name_en: Domestic Third-Party Trade Exception
description_zh: 通过 hypersku-cli domestic-third-trade-exception 子命令分页查询国内第三方交易异常订单（含监控单号、交易号、物流单号、物流明细）及异常订单留言记录。
description_en: Page through domestic third-party trade exception orders with logistics details, and query message records of exception orders via hypersku-cli domestic-third-trade-exception commands.
description: HyperSKU 国内第三方交易异常订单管理。当用户提到异常订单/丢包裹/丢件/未签收/假发货/假签收/退件/未入库/无货，或需要查询异常订单留言/跟进记录时，通过 hypersku-cli domestic-third-trade-exception 子命令分页查询对应数据。
version: 2.0.0
author: owen
tags:
  - hypersku
  - cli
  - 异常订单
---

# 国内第三方交易异常订单管理

通过 `hypersku-cli domestic-third-trade-exception` 子命令管理 Hypersku 国内第三方交易异常订单（HyperSKU 在第三方平台如 1688/淘宝 采购产生的采购单，其国内段物流异常在此监控），支持按异常主/子状态分页查询异常订单及物流明细、查看异常订单的留言记录。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 分页查询异常订单 | 按异常主/子状态分页查询国内第三方交易异常订单（含监控单号、交易号、物流单号、物流公司、状态等） | `hypersku-cli domestic-third-trade-exception page-list [flags]` | [page-list.md](references/page-list.md) |
| 查询异常订单留言 | 根据监控订单ID与监控物流ID查询异常订单的留言记录 | `hypersku-cli domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>` | [message-list.md](references/message-list.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"异常订单/丢包裹/丢件/未签收/假发货/假签收/退件/未入库/无货 + 异常状态"时，执行 `page-list` 分页查询异常订单及物流明细。
- 用户提到"留言/备注/跟进/沟通记录 + 监控订单ID 和 监控物流ID"时，执行 `message-list <monitorOrderId> <monitorLogisticsId>` 展示留言记录。
- 若未提供 `--hypersku-status` / `--hypersku-sub-status` 等必填参数，提示用户提供异常状态后重试。

## 注意事项

1. **业务语义**：这里的"异常订单"是 HyperSKU 在第三方平台（1688/淘宝等）的采购单，其国内段物流发生异常（丢包裹/丢件/未签收等）在此监控，**不是** HyperSKU 客户订单；输出中的"交易号/1688订单号"即第三方平台的交易号。
2. **page-list 必填参数**：`--hypersku-status`（异常主状态）与 `--hypersku-sub-status`（异常子状态列表）必填，缺失时会打印可选值并显示帮助。
3. **物流维度输出**：`page-list` 输出按物流记录展开，每个订单的每条物流占一行；订单无物流记录时整单不展示。
4. **状态枚举**：主状态 `--hypersku-status` 可选 1-未发货/2-假发货/3-未到货/4-假签收/5-未签收/6-退件/7-丢件/8-未入库/9-丢包裹/10-无货；子状态 `--hypersku-sub-status` 可选 1-待处理/2-处理中/3-已处理/4-已关闭/5-已拒绝（完整对照见 [page-list.md](references/page-list.md)）。
5. **message-list 参数**：两个位置参数顺序为监控订单ID、监控物流ID，可从 `page-list` 输出中获取。
