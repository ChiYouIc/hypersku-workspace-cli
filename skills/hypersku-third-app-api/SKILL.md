---
name: hypersku-third-app-api
display_name: 第三方平台数据查询
display_name_en: Third-Party Platform Data Query
description_zh: 通过 hypersku-cli third-app-api 子命令查询 1688 和 Shopify 等第三方平台数据：订单物流、订单详情、退款单、产品信息、供应商信息、店铺产品与库存。
description_en: Query third-party platform data (1688, Shopify) via hypersku-cli third-app-api: order logistics, order details, refund records, product info, supplier info, shop products and inventory.
description: HyperSKU 第三方平台数据查询。当用户提到 1688物流/1688订单/1688退款/1688产品/1688供应商/Shopify店铺/店铺产品/店铺库存/店铺订单/混批/子账号时，通过 hypersku-cli third-app-api 子命令查询对应数据。
version: 1.0.0
author: owen
---

# 第三方平台数据查询

通过 `hypersku-cli third-app-api` 子命令查询 1688、Shopify 等第三方平台的聚合数据，仅提供只读查询能力。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 1688 物流信息 | 根据交易号查询 1688 订单物流信息 | `hypersku-cli third-app-api ali-logistics <orderId>` | [ali-logistics.md](references/ali-logistics.md) |
| 1688 物流轨迹 | 根据交易号查询 1688 完整物流轨迹 | `hypersku-cli third-app-api ali-logistics-trace <orderId>` | [ali-logistics-trace.md](references/ali-logistics-trace.md) |
| 1688 订单详情 | 根据交易号查询 1688 订单详情 | `hypersku-cli third-app-api ali-order-detail <orderId>` | [ali-order-detail.md](references/ali-order-detail.md) |
| 1688 退款详情 | 根据退款单号查询 1688 退款详情 | `hypersku-cli third-app-api ali-refund-detail <refundId>` | [ali-refund-detail.md](references/ali-refund-detail.md) |
| 1688 退款列表 | 根据交易号查询 1688 退款单列表 | `hypersku-cli third-app-api ali-refund-list <orderId>` | [ali-refund-list.md](references/ali-refund-list.md) |
| 1688 退款操作记录 | 根据退款单号查询操作记录 | `hypersku-cli third-app-api ali-refund-operations <refundId>` | [ali-refund-operations.md](references/ali-refund-operations.md) |
| 1688 产品信息 | 按 ID 或 URL 查询 1688 产品信息 | `hypersku-cli third-app-api ali-product --id <id>` | [ali-product.md](references/ali-product.md) |
| 1688 供应商信息 | 按登录 ID 或店铺地址查询供应商 | `hypersku-cli third-app-api ali-supplier --login-id <id>` | [ali-supplier.md](references/ali-supplier.md) |
| 1688 混批设置 | 查询卖家混批设置 | `hypersku-cli third-app-api ali-mix-config --member-id <id>` | [ali-mix-config.md](references/ali-mix-config.md) |
| 1688 子账号列表 | 查询当前账号下所有子账号 | `hypersku-cli third-app-api ali-sub-accounts` | [ali-sub-accounts.md](references/ali-sub-accounts.md) |
| 店铺产品查询 | 根据店铺 ID 查询产品 | `hypersku-cli third-app-api shop-product --store-id <id>` | [shop-product.md](references/shop-product.md) |
| 店铺订单查询 | 根据店铺 ID 查询订单列表 | `hypersku-cli third-app-api shop-order-list --store-id <id>` | [shop-order-list.md](references/shop-order-list.md) |
| 店铺库存查询 | 根据店铺 ID 查询库存 | `hypersku-cli third-app-api shop-inventory --store-id <id>` | [shop-inventory.md](references/shop-inventory.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"1688物流/1688快递/1688运单"并提供交易号时，执行 `ali-logistics <orderId>` 查询物流信息。
- 用户提到"1688物流轨迹/1688轨迹/1688到哪了"时，执行 `ali-logistics-trace <orderId>` 查询完整轨迹。
- 用户提到"1688订单详情/1688订单/1688买了什么"时，执行 `ali-order-detail <orderId>` 查询订单详情。
- 用户提到"1688退款详情/1688退款单"并提供退款单号时，执行 `ali-refund-detail <refundId>` 查询退款详情。
- 用户提到"1688退款列表/1688退款"并提供交易号时，执行 `ali-refund-list <orderId>` 查询退款列表。
- 用户提到"1688退款操作/1688退款记录"时，执行 `ali-refund-operations <refundId>` 查询操作记录。
- 用户提到"1688产品/1688商品/1688链接"时，执行 `ali-product` 查询产品信息（支持 --id 或 --url）。
- 用户提到"1688供应商/供应商信息"时，执行 `ali-supplier` 查询供应商信息（支持 --login-id 或 --domain）。
- 用户提到"混批/混批设置"时，执行 `ali-mix-config` 查询混批设置。
- 用户提到"子账号/子账号列表"时，执行 `ali-sub-accounts` 查询子账号。
- 用户提到"店铺产品/Shopify产品"并提供店铺 ID 时，执行 `shop-product` 查询店铺产品。
- 用户提到"店铺订单/Shopify订单"时，执行 `shop-order-list` 查询店铺订单。
- 用户提到"店铺库存/Shopify库存"时，执行 `shop-inventory` 查询店铺库存。
- 若未提供必要参数（如 orderId、storeId），提示用户提供。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含写入/修改操作。
2. **1688 订单详情耗时较长**：`ali-order-detail` 接口响应时间较长，请耐心等待。
3. **退款单号格式**：1688 退款单逻辑主键格式为 `TQ+数字`（如 TQ123456）。
4. **店铺查询必填**：`shop-product`、`shop-order-list`、`shop-inventory` 均需要 `--store-id` 参数。
