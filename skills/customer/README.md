# 客户订单管理

> 本文为 [HyperSKU CLI 统一入口](../SKILL.md) 下的「客户管理」能力域详细引导。

通过 `hypersku-cli customer order` 子命令管理 Hypersku 客户订单，支持查询订单详情（商品、金额）、物流信息及收货地址。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询订单详情 | 根据客户订单号查询订单概要 + 商品明细（含金额、币种、仓库、采购状态等） | `hypersku-cli customer order info <orderId>` | [info.md](references/info.md) |
| 查询订单物流 | 根据客户订单号查询物流单号及承运商（不含轨迹详情） | `hypersku-cli customer order logistics <orderId>` | [logistics.md](references/logistics.md) |
| 查询收货地址 | 根据客户订单号查询收货人、地址、税号、VAT 等信息 | `hypersku-cli customer order address <orderId>` | [address.md](references/address.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"客户订单/订单详情/商品/买了什么/金额"时，执行 `info <orderId>` 展示订单概要及商品明细。
- 用户提到"物流/单号/承运商/快递/tracking"时，执行 `logistics <orderId>` 以表格展示物流单号与承运商。
- 用户提到"地址/收货/税号/VAT/邮编/收件人"时，执行 `address <orderId>` 展示收货地址详情。
- 若未提供 orderId，提示用户提供客户订单号。

## 注意事项

1. **orderId 必填**：所有子命令都需要传入客户订单号参数，未提供时会显示帮助信息并退出。
2. **采购状态码**：`info` 输出中的采购状态使用与采购订单相同的状态码对照表（见 [purchase/README.md](../purchase/README.md)）。
3. **物流 vs 物流轨迹**：`customer` 的 logistics 子命令只返回物流单号和承运商，不包含具体轨迹；如需查询国内快递轨迹，请使用 `hypersku-cli logistics tracking`。
4. **国际化**：地址信息可能包含英文内容，具体取决于客户所在国家。
5. **币种**：订单金额字段附带币种符号和币种代码，支持多币种展示。
