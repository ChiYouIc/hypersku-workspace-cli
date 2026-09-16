---
name: hypersku-purchase-after-sales
display_name: 采购售后查询
display_name_en: Purchase After-Sales Query
description_zh: 通过 hypersku-cli purchase-after-sales 子命令按交易号或订单号查询采购售后工单列表，含工单状态、异常类型、商品信息等。
description_en: Query purchase after-sales work orders by trade ID or order ID via hypersku-cli purchase-after-sales, including work order status, abnormal type, and product info.
description: HyperSKU 采购售后查询。当用户提到采购售后/采购退款/采购异常工单/采购工单/交易号查售后/订单号查售后时，通过 hypersku-cli purchase-after-sales 子命令查询采购售后工单。
version: 1.0.0
author: owen
---

# 采购售后查询

通过 `hypersku-cli purchase-after-sales` 子命令按交易号或订单号查询采购售后工单列表，支持查看工单详情、异常类型、处理状态等信息。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 按交易号查询 | 根据交易号查询采购售后工单列表 | `hypersku-cli purchase-after-sales by-trade <tradeId>` | [by-trade.md](references/by-trade.md) |
| 按订单号查询 | 根据订单号查询采购售后工单列表 | `hypersku-cli purchase-after-sales by-order <orderId>` | [by-order.md](references/by-order.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"采购售后/采购退款/采购异常"并提供交易号（第三方订单号）时，执行 `by-trade <tradeId>` 查询采购售后工单。
- 用户提到"采购售后/采购工单"并提供订单号时，执行 `by-order <orderId>` 查询采购售后工单。
- 若未提供交易号或订单号，提示用户提供。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含写入/修改操作。
2. **必填参数**：`by-trade` 需要交易号（第三方订单号），`by-order` 需要订单号。
3. **与 1688 售后区别**：本命令查询 HyperSKU 内部采购售后工单；1688 售后查询请使用 `hypersku-cli after-sales` 命令。
4. **与包裹拦截/退件区别**：本命令查询采购侧售后工单；包裹拦截/退件请使用 `hypersku-cli customer order intercept` 或 `return-by-trade` 命令。
