---
name: purchase
description: Hypersku 采购订单管理能力。当用户需要查询采购订单商品、物流轨迹、采购日志、订单地址时使用。
version: 1.0.0
author: owen
tags:
  - 订单商品
  - 订单地址
  - 国际物流
  - 采购日志
---

# 采购订单管理 Skill

通过 `hypersku-cli purchase` 子命令管理 Hypersku 采购订单，支持查询订单详情、收货地址、采购日志及国际物流轨迹。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询订单详情 | 根据订单号查询商品明细、金额、订单类型/状态、仓库等 | `hypersku-cli purchase info <orderId>` | [info.md](references/info.md) |
| 查询订单地址 | 根据订单号查询收货地址（国家/省份/市/区/邮编） | `hypersku-cli purchase address <orderId>` | [address.md](references/address.md) |
| 查询采购日志 | 根据订单号查询订单状态流转日志 | `hypersku-cli purchase log <orderId>` | [log.md](references/log.md) |
| 查询国际物流 | 根据订单号查询包裹国际段物流轨迹（承运商/单号/轨迹） | `hypersku-cli purchase logistics <orderId>` | [logistics.md](references/logistics.md) |

## 意图判断决策树

当用户输入包含以下关键词或意图时，使用对应的子命令：

```
用户输入
├── 含 "订单详情" / "订单信息" / "商品" / "买了什么"
│   └── 执行 info <orderId> → 展示订单概要 + 商品明细列表
│
├── 含 "地址" / "收货地址" / "寄到哪里" / "收货人"
│   └── 执行 address <orderId> → 展示完整收货地址
│
├── 含 "日志" / "状态" / "进度" / "到哪一步了" / "流程"
│   └── 执行 log <orderId> → 展示时间线状态表
│
├── 含 "物流" / "轨迹" / "快递" / "运输" / "到哪了" / "包裹"
│   └── 执行 logistics <orderId> → 展示包裹物流轨迹
│
└── 未提供 orderId
    └── 提示用户提供订单号
```

## 注意事项

1. **orderId 必填**：所有子命令都需要传入订单号参数，未提供时会显示帮助信息并退出。
3. **国际化**：地址和物流信息可能包含英文内容，具体取决于订单来源国家。
4. **多包裹**：国际物流查询可能返回多个包裹的轨迹信息，每个包裹独立展示。


