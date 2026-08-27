---
name: hypersku-cli
description: HyperSKU CLI 统一入口（所有子 skill 入口合并于此）。当用户需要查询采购订单、客户订单、客户档案（是否已支付/绑定店铺/联系方式/姓名/DS 经验/广告预算/细分市场/客户标签/渠道来源）、物流轨迹、仓库物流、1688售后工单、国内第三方交易异常订单，或需要处理物流异常（排查/处理）、申请售后（物流查询/赔偿），或需要查询客户画像（订单统计/日订单数量/交易统计/交易流水），或需要生成客户画像/线索评估/跟进优先级/转化策略/销售话术/促绑店/促首单时使用。覆盖采购管理、客户管理、物流追踪、仓库管理、售后管理、物流异常管理、异常处理、申请售后八大能力域及客户画像分析。
version: 2.0.0
author: owen
tags:
  - hypersku
  - cli
  - 入口
---

# HyperSKU CLI

`hypersku-cli` 已集成到当前系统环境中（可执行文件位于当前用户目录下的 `~/.hypersku-cli/hypersku-cli.exe`），可直接在终端调用。它是连接 Hypersku 平台的命令行工具，通过子命令按需查询采购订单、客户订单、物流轨迹等业务数据，无需打开网页后台。

本文件是 Hypersku 全部能力的**唯一入口**，包含「能力总览」与「检索引导」：前者索引八大能力域的 CLI 命令与详细文档；后者按用户问法关键词定位应使用的模块、命令与文档。各子目录下的 `README.md` 为能力域详细引导，`references/` 为命令输出与字段说明参考，均由本入口索引。

## 使用方式

在终端中直接输入命令即可，格式为：

```
hypersku-cli <模块> <操作> <参数>
```

常用示例：

```bash
# 查询采购订单详情
hypersku-cli purchase info PO202401010001

# 查询客户订单收货地址
hypersku-cli customer order address CUST202401010001

# 查询国内快递物流轨迹
hypersku-cli logistics tracking 1Z999AA10123456784
```

> **环境说明**：`hypersku-cli` 已安装到当前用户目录 `~/.hypersku-cli/` 并配置完毕，所有 API 鉴权和网络连接均已就绪，无需额外配置即可使用。

## 能力总览（统一入口）

八大能力域统一索引如下。**查询类**能力直接调用 CLI 命令；**处理类**能力（异常处理、申请售后）为流程引导，配合查询类命令使用。

| # | 模块 | 用途 | CLI 命令前缀 | 详细文档 |
|---|------|------|-------------|----------|
| 1 | 采购管理 | 查询采购订单详情、分页搜索、状态日志、国际物流轨迹 | `hypersku-cli purchase` | [purchase/README.md](purchase/README.md) |
| 2 | 客户管理 | 查询客户基础信息（等级/店铺/预算）与基础档案（邮箱/客户经理/注册/登录），客户订单详情（商品/金额/币种）、物流单号、收货地址，及客户画像（订单统计/日订单数量/交易统计/交易流水） | `hypersku-cli customer` | [customer/README.md](customer/README.md) |
| 3 | 物流追踪 | 根据运单号查询国内快递完整物流轨迹 | `hypersku-cli logistics` | [logistics/README.md](logistics/README.md) |
| 4 | 仓库管理 | 根据运单号查询仓库侧物流轨迹（快递签收/仓库签收/入库/仓库操作） | `hypersku-cli warehouse` | [warehouse/README.md](warehouse/README.md) |
| 5 | 售后管理 | 查询1688售后工单、售后商品、退款详情、留言记录 | `hypersku-cli after-sales` | [after-sales/README.md](after-sales/README.md) |
| 6 | 物流异常管理 | 查询国内第三方交易异常订单（丢包裹/丢件/未签收等）及留言记录 | `hypersku-cli domestic-third-trade-exception` | [domestic-third-trade-exception/README.md](domestic-third-trade-exception/README.md) |
| 7 | 异常处理 | 排查/处理异常订单（未发货/假发货/未到货/假签收/未签收/退件/丢件/未入库/丢包裹/无货） | 流程引导（配合 6 查询） | [domestic-exception-handling/README.md](domestic-exception-handling/README.md) |
| 8 | 申请售后 | 国际物流未签收发起售后（未签收 < 30 天走物流查询，> 30 天走物流赔偿） | 流程引导（配合 2/3 查询） | [after-sales-apply/README.md](after-sales-apply/README.md) |
| 9 | 客户画像分析 | 分析已注册未消费客户（评估线索质量/意向热度/跟进优先级/转化策略/销售话术/促绑店/促首单） | 分析引导（配合 2 查询） | [customer-profile-analysis/README.md](customer-profile-analysis/README.md) |

## 检索引导

按用户问法中的关键词或意图，定位到对应模块 → 执行命令 → 查阅文档。**先匹配关键词，再执行命令；参数不足时向用户询问补全。**

### 查询类

| 用户意图 / 关键词 | 模块 | 执行命令 | 文档 |
|------|------|------|------|
| 采购订单详情、订单信息、商品明细、"买了什么" + 订单号 | 采购管理 | `hypersku-cli purchase info <orderId>` | [purchase/README.md](purchase/README.md) |
| 采购订单搜索、列表、分页、有哪些订单 + 日期/交易号/物流单号 | 采购管理 | `hypersku-cli purchase info page [flags]` | [purchase/README.md](purchase/README.md) |
| 采购日志、订单状态、进度、到哪一步了 | 采购管理 | `hypersku-cli purchase log <orderId>` | [purchase/README.md](purchase/README.md) |
| 国际物流、国际段包裹轨迹、运输状态 | 采购管理 | `hypersku-cli purchase logistics <orderId>` | [purchase/README.md](purchase/README.md) |
| 客户档案、是否有已支付订单、绑定店铺、联系方式、姓名、DS 经验、广告预算、细分市场、客户标签、渠道来源 + 客户ID | 客户管理 | `hypersku-cli customer detail <customerId>` | [customer/README.md](customer/README.md) |
| 客户订单详情、客户买了什么、金额、币种 | 客户管理 | `hypersku-cli customer order info <orderId>` | [customer/README.md](customer/README.md) |
| 客户订单物流单号、承运商（不含轨迹） | 客户管理 | `hypersku-cli customer order logistics <orderId>` | [customer/README.md](customer/README.md) |
| 收货地址、收件人、税号、VAT、邮编 | 客户管理 | `hypersku-cli customer order address <orderId>` | [customer/README.md](customer/README.md) |
| 客户订单退件、退件工单、退货记录 + 客户订单号 | 客户管理 | `hypersku-cli customer order return <customerOrderId>` | [customer/README.md](customer/README.md) |
| 客户画像、订单统计、日均订单、总订单数 + 客户ID/日期范围 | 客户管理 | `hypersku-cli customer profile order count <customerId> --start <date> --end <date>` | [customer/README.md](customer/README.md) |
| 日订单数量、每天下单多少、付款/履约/超期/退款订单数 + 客户ID/日期范围 | 客户管理 | `hypersku-cli customer profile order daily <customerId> --start <date> --end <date>` | [customer/README.md](customer/README.md) |
| 交易统计、交易额、成交额、退款金额、客单价 + 客户ID/日期范围 | 客户管理 | `hypersku-cli customer profile transaction count <customerId> --start <date> --end <date>` | [customer/README.md](customer/README.md) |
| 交易流水、每日交易明细、手续费 + 客户ID/日期范围 | 客户管理 | `hypersku-cli customer profile transaction bills <customerId> --start <date> --end <date>` | [customer/README.md](customer/README.md) |
| 国内快递轨迹、运单号 tracking、包裹到哪了 | 物流追踪 | `hypersku-cli logistics tracking <trackingNumber>` | [logistics/README.md](logistics/README.md) |
| 仓库物流、仓库签收、入库、仓库操作、到仓 | 仓库管理 | `hypersku-cli warehouse tracking <trackingNumber>` | [warehouse/README.md](warehouse/README.md) |
| 1688售后工单、退款列表 + 1688订单号/交易号 | 售后管理 | `hypersku-cli after-sales 1688 <thirdOrderId>` | [after-sales/README.md](after-sales/README.md) |
| 售后商品、退款商品明细 + 订单号/退款ID | 售后管理 | `hypersku-cli after-sales 1688 goods <thirdOrderId> <refundId>` | [after-sales/README.md](after-sales/README.md) |
| 售后详情、退款详情、退款原因、卖家信息 + 退款ID | 售后管理 | `hypersku-cli after-sales 1688 detail <refundId>` | [after-sales/README.md](after-sales/README.md) |
| 售后留言、沟通记录 + 退款ID | 售后管理 | `hypersku-cli after-sales 1688 message <refundId>` | [after-sales/README.md](after-sales/README.md) |
| 异常订单、丢包裹、丢件、未签收、假发货等 + 异常状态 | 物流异常管理 | `hypersku-cli domestic-third-trade-exception page-list --hypersku-status <主状态> --hypersku-sub-status <子状态>` | [domestic-third-trade-exception/README.md](domestic-third-trade-exception/README.md) |
| 异常订单留言、跟进记录 + 监控订单ID/物流ID | 物流异常管理 | `hypersku-cli domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>` | [domestic-third-trade-exception/README.md](domestic-third-trade-exception/README.md) |

### 处理类（流程引导，无独立 CLI 命令）

| 用户意图 / 关键词 | 模块 | 配合命令 | 文档 |
|------|------|------|------|
| 异常是什么意思、为什么丢包裹/假发货、怎么排查、下一步怎么办、要不要退款/补发 | 异常处理 | `domestic-third-trade-exception page-list` / `message-list` | [domestic-exception-handling/README.md](domestic-exception-handling/README.md) |
| 物流不动了、轨迹不更新、申请物流查询（未签收 < 30 天） | 申请售后 | `customer order info` / `logistics`、`logistics tracking` | [after-sales-apply/README.md](after-sales-apply/README.md) |
| 快递丢了要赔偿、索赔、理赔（未签收 > 30 天） | 申请售后 | `customer order info`、`customer order return` | [after-sales-apply/README.md](after-sales-apply/README.md) |
| 生成画像、画像分析、线索评估、意向热度、跟进优先级、转化策略、销售话术、促绑店、促首单 + 客户ID | 客户画像分析 | `customer detail <customerId>` | [customer-profile-analysis/README.md](customer-profile-analysis/README.md) |

## 详细文档索引

- [purchase/README.md](purchase/README.md) — 采购订单管理（详情/分页/日志/国际物流）
- [customer/README.md](customer/README.md) — 客户信息、客户订单管理（详情/物流/地址/退件）与客户画像（订单统计/日订单/交易统计/交易流水）
- [logistics/README.md](logistics/README.md) — 国内快递物流轨迹
- [warehouse/README.md](warehouse/README.md) — 仓库物流轨迹（签收/入库/仓库操作）
- [after-sales/README.md](after-sales/README.md) — 1688 售后管理（工单/商品/详情/留言）
- [domestic-third-trade-exception/README.md](domestic-third-trade-exception/README.md) — 国内第三方交易异常订单查询
- [domestic-exception-handling/README.md](domestic-exception-handling/README.md) — 异常订单处理（排查/处理流程）
- [after-sales-apply/README.md](after-sales-apply/README.md) — 申请售后（物流查询/赔偿）
- [customer-profile-analysis/README.md](customer-profile-analysis/README.md) — 客户画像分析（新用户转化：线索质量/意向热度/跟进优先级/转化策略/销售话术）

## 约束

1. **查不到即止**：如果通过某个接口查询不到订单或数据，不要换用其他接口反复尝试，直接告知用户无结果。
2. **参数必填**：所有子命令的必填参数（如 orderId、trackingNumber）缺失时，命令会自动显示帮助信息，需提示用户补全。
3. **按需选择模块**：根据用户意图匹配对应模块——采购用 `purchase`，客户用 `customer order`，国内快递轨迹用 `logistics`，物流异常用 `domestic-third-trade-exception`，不要混用。
4. **一次一问**：每次只执行一条命令，等待结果返回后再决定下一步，不要连续并行调用多个查询。