---
name: hypersku-cli
description: HyperSKU CLI 统一入口。当用户需要查询采购订单、客户订单、物流轨迹、1688售后工单等 Hypersku 业务数据时使用。覆盖采购管理、客户管理、物流追踪、售后管理四大能力域。
version: 1.0.0
author: owen
tags:
  - hypersku
  - cli
  - 入口
---

# HyperSKU CLI

`hypersku-cli` 已集成到当前系统环境中（可执行文件位于当前用户目录下的 `~/.hypersku-cli/hypersku-cli.exe`），可直接在终端调用。它是连接 Hypersku 平台的命令行工具，通过子命令按需查询采购订单、客户订单、物流轨迹等业务数据，无需打开网页后台。

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

## 能力域

| 模块 | 说明 | 命令前缀 | 详细文档 |
|------|------|----------|----------|
| 采购管理 | 查询采购订单详情、按条件分页搜索、状态日志、国际物流轨迹 | `hypersku-cli purchase` | [purchase/SKILL.md](purchase/SKILL.md) |
| 客户管理 | 查询客户订单详情（商品/金额/币种）、物流单号、收货地址 | `hypersku-cli customer` | [customer/SKILL.md](customer/SKILL.md) |
| 物流追踪 | 根据运单号查询国内快递的完整物流轨迹 | `hypersku-cli logistics` | [logistics/SKILL.md](logistics/SKILL.md) |
| 售后管理 | 查询1688售后工单、售后商品、退款详情、留言记录 | `hypersku-cli after-sales` | [after-sales/SKILL.md](after-sales/SKILL.md) |

## 快速导航

- 用户问"采购订单"相关 → 参考 [purchase](purchase/SKILL.md)
- 用户问"客户订单"相关 → 参考 [customer](customer/SKILL.md)
- 用户问"快递/物流轨迹"相关 → 参考 [logistics](logistics/SKILL.md)
- 用户问"售后/退款/1688售后"相关 → 参考 [after-sales](after-sales/SKILL.md)

## 约束

1. **查不到即止**：如果通过某个接口查询不到订单或数据，不要换用其他接口反复尝试，直接告知用户无结果。
2. **参数必填**：所有子命令的必填参数（如 orderId、trackingNumber）缺失时，命令会自动显示帮助信息，需提示用户补全。
3. **按需选择模块**：根据用户意图匹配对应模块——采购用 `purchase`，客户用 `customer order`，国内快递轨迹用 `logistics`，不要混用。
4. **一次一问**：每次只执行一条命令，等待结果返回后再决定下一步，不要连续并行调用多个查询。