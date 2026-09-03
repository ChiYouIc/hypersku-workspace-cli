---
name: hypersku-cli
display_name: HyperSKU CLI 能力总览与使用规则
display_name_en: HyperSKU CLI Overview & Usage Rules
description_zh: 介绍 hypersku-cli 的整体能力版图（登录、采购、客户、物流、仓库、售后、异常处理）与统一使用规则（登录前置、参数必填、查不到即止、一次一问），不承载具体命令细节。
description_en: Introduces the overall capability map of hypersku-cli (auth, purchase, customer, logistics, warehouse, after-sales, exception handling) and unified usage rules (login first, required params, stop on no result, one command at a time). No command details here.
description: HyperSKU CLI 能力总览与使用规则。当用户泛称 hypersku-cli/这个 CLI 都能干什么/有哪些能力/怎么用/使用规范，或不确定该用哪个能力域时，先读本技能了解能力版图与统一规则，再路由到对应能力域技能包（hypersku-auth/purchase/customer/logistics/warehouse/after-sales 等）。
version: 1.0.0
author: owen
tags:
  - hypersku
  - cli
  - 总览
---

# HyperSKU CLI 能力总览与使用规则

`hypersku-cli` 是 HyperSKU 的命令行工具，提供登录认证与采购、客户、物流、售后等业务查询能力。本技能只做**能力版图介绍 + 统一使用规则**，具体命令细节由各能力域技能包承载。

## 能力版图

| 能力域 | 用途 | 命令入口 | 承载技能包 |
|--------|------|----------|------------|
| 登录认证 | 登录状态远程校验（login/logout 暂未实现，凭证手动配置于 ~/.hypersku-cli/config.json） | `hypersku-cli auth login/status/logout` | `hypersku-auth` |
| 采购订单 | 订单详情、商品明细、订单搜索、采购日志、国际物流 | `hypersku-cli purchase ...` | `hypersku-purchase` |
| 客户管理 | 客户档案、客户订单、收货地址、退件工单、订单/交易统计 | `hypersku-cli customer ...` | `hypersku-customer` |
| 客户画像分析 | 新用户线索评估、跟进优先级、转化策略、销售话术 | 基于 `customer detail` 数据分析 | `hypersku-customer-profile-analysis` |
| 物流轨迹 | 国内快递轨迹查询（到哪了） | `hypersku-cli logistics tracking` | `hypersku-logistics` |
| 仓库物流 | 到仓签收、入库、仓库操作轨迹 | `hypersku-cli warehouse tracking` | `hypersku-warehouse` |
| 1688 售后 | 售后工单、退款列表/详情、售后商品、沟通留言 | `hypersku-cli after-sales 1688 ...` | `hypersku-after-sales` |
| 申请售后 | 物流查询/物流赔偿工单的判断与填写引导 | 流程引导（配合查询命令） | `hypersku-after-sales-apply` |
| 异常订单 | 异常订单（丢包裹/假发货等）分页查询与留言 | `hypersku-cli domestic-third-trade-exception ...` | `hypersku-domestic-third-trade-exception` |
| 异常处理 | 异常含义解释、根因排查、下一步动作建议 | 流程引导（配合查询命令） | `hypersku-domestic-exception-handling` |

## 统一使用规则

1. **登录前置**：所有业务查询都要求 token 有效（凭证保存于 `~/.hypersku-cli/config.json` 的 `api_token`，当前需手动配置）。查询报凭证错误时，先用 `auth status` 远程校验；未登录则引导用户更新配置文件中的 `api_token`（参见 `hypersku-auth` 技能包）。
2. **参数必填**：各命令的必填参数（如 orderId、trackingNumber、customerId）缺失时会显示帮助信息，需提示用户补全后再执行。
3. **查不到即止**：某个接口查不到数据时，不要换接口反复尝试，直接告知用户无结果。
4. **一次一问**：每次只执行一条命令，等结果返回后再决定下一步，不要连续并行调用多个查询。
5. **能力域路由**：按用户意图选择能力域——采购用 `purchase`、客户用 `customer`、国内快递轨迹用 `logistics`、仓库侧用 `warehouse`、售后用 `after-sales`，不要混用；具体子命令与参数请读取对应能力域技能包的 SKILL.md。
