# 客户管理

> 本文为 [HyperSKU CLI 统一入口](../SKILL.md) 下的「客户管理」能力域详细引导。

通过 `hypersku-cli customer` 子命令管理 Hypersku 客户订单与客户画像，支持查询客户基础信息、订单详情（商品、金额）、物流信息、收货地址及客户画像（订单统计、交易数据）。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询客户扩展信息 | 根据客户 ID 查询等级、近 30 天订单数、店铺数、Dropshipping 时长、周广告预算等 | `hypersku-cli customer extends <customerId>` | [customer-extends.md](references/customer-extends.md) |
| 查询客户档案 | 根据客户 ID 查询注册时间、注册区域、最近登录时间、总订单数与总购买金额等基础档案 | `hypersku-cli customer detail <customerId>` | [customer-detail.md](references/customer-detail.md) |
| 查询订单详情 | 根据客户订单号查询订单概要 + 商品明细（含金额、币种、仓库、采购状态等） | `hypersku-cli customer order info <orderId>` | [info.md](references/info.md) |
| 查询订单物流 | 根据客户订单号查询物流单号及承运商（不含轨迹详情） | `hypersku-cli customer order logistics <orderId>` | [logistics.md](references/logistics.md) |
| 查询收货地址 | 根据客户订单号查询收货人、地址、税号、VAT 等信息 | `hypersku-cli customer order address <orderId>` | [address.md](references/address.md) |
| 查询订单退件 | 根据客户订单号查询退件工单详情 | `hypersku-cli customer order return <customerOrderId>` | [return.md](references/return.md) |
| 查询订单统计 | 根据客户 ID 查询时间范围内的总订单数、日均/日最大/日最小订单数 | `hypersku-cli customer profile order count <customerId> --start --end` | [profile-order-count.md](references/profile-order-count.md) |
| 查询日订单数量 | 根据客户 ID 查询时间范围内每日订单数（付款/履约/超期/退款） | `hypersku-cli customer profile order daily <customerId> --start --end` | [profile-order-daily.md](references/profile-order-daily.md) |
| 查询交易统计 | 根据客户 ID 查询时间范围内的总交易额、实际成交额、退款金额、客单价 | `hypersku-cli customer profile transaction count <customerId> --start --end` | [profile-transaction-count.md](references/profile-transaction-count.md) |
| 查询交易流水 | 根据客户 ID 查询时间范围内每日交易流水（多币种金额） | `hypersku-cli customer profile transaction bills <customerId> --start --end` | [profile-transaction-bills.md](references/profile-transaction-bills.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"客户信息/客户等级/店铺数量/广告预算/客户背景"时，执行 `extends <customerId>` 展示客户扩展信息。
- 用户提到"客户档案/注册时间/注册区域/最近登录/总订单数/总购买金额"时，执行 `detail <customerId>` 展示客户基础档案。
- 用户提到"客户订单/订单详情/商品/买了什么/金额"时，执行 `order info <orderId>` 展示订单概要及商品明细。
- 用户提到"物流/单号/承运商/快递/tracking"时，执行 `order logistics <orderId>` 以表格展示物流单号与承运商。
- 用户提到"地址/收货/税号/VAT/邮编/收件人"时，执行 `order address <orderId>` 展示收货地址详情。
- 用户提到"退件/退货工单/退件记录"时，执行 `order return <customerOrderId>` 展示退件工单详情。
- 用户提到"订单统计/总订单数/日均订单/客户下了多少单"时，执行 `profile order count <customerId> --start --end` 展示订单统计。
- 用户提到"每天订单/日订单/每日下单量/付款/履约/超期订单数"时，执行 `profile order daily <customerId> --start --end` 展示每日订单数量。
- 用户提到"交易统计/交易额/成交额/退款金额/客单价"时，执行 `profile transaction count <customerId> --start --end` 展示交易统计。
- 用户提到"交易流水/每日交易/流水明细/手续费"时，执行 `profile transaction bills <customerId> --start --end` 展示每日交易流水。
- 若未提供 orderId / customerId，提示用户提供客户订单号或客户 ID。

## 注意事项

1. **orderId / customerId 必填**：`order` 下子命令需传入客户订单号，`info` 与 `profile` 下子命令需传入客户 ID，未提供时会显示帮助信息并退出。
2. **画像日期参数**：`profile` 下子命令的 `--start` / `--end` 为必填，格式 `yyyy-MM-dd`；开始日期自动补 `00:00:00`，结束日期自动补 `23:59:59`。
3. **最多 90 天**：`profile order daily` 与 `profile transaction bills` 每次最多返回 90 天的数据，超出时仅返回前 90 天。
4. **采购状态码**：`info` 输出中的采购状态使用与采购订单相同的状态码对照表（见 [purchase/README.md](../purchase/README.md)）。
5. **物流 vs 物流轨迹**：`customer` 的 logistics 子命令只返回物流单号和承运商，不包含具体轨迹；如需查询国内快递轨迹，请使用 `hypersku-cli logistics tracking`。
6. **国际化**：地址信息可能包含英文内容，具体取决于客户所在国家。
7. **币种**：订单金额字段附带币种符号和币种代码，支持多币种展示；`profile transaction bills` 的金额列依次为 CNY/USD/EUR。
