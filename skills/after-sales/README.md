# 1688售后管理

> 本文为 [HyperSKU CLI 统一入口](../SKILL.md) 下的「售后管理」能力域详细引导。

通过 `hypersku-cli after-sales` 子命令管理 Hypersku 平台上的 1688 售后工单，支持查询工单列表、售后商品、退款详情及留言记录。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 查询售后工单 | 根据1688订单号/交易号查询售后工单列表 | `hypersku-cli after-sales 1688 <thirdOrderId>` | [1688.md](references/1688.md) |
| 查询售后商品 | 根据订单号和退款ID查询售后工单中的商品项 | `hypersku-cli after-sales 1688 goods <thirdOrderId> <refundId>` | [goods.md](references/goods.md) |
| 查询售后详情 | 根据退款ID查询售后详细信息（退款金额/原因/卖家） | `hypersku-cli after-sales 1688 detail <refundId>` | [detail.md](references/detail.md) |
| 查询售后留言 | 根据退款ID查询售后留言记录 | `hypersku-cli after-sales 1688 message <refundId>` | [message.md](references/message.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"售后工单/1688售后/退款列表 + 订单号/交易号"时，执行 `1688 <thirdOrderId>` 以表格展示售后工单列表。
- 用户提到"售后商品/退款商品 + 订单号/退款ID"时，执行 `1688 goods <thirdOrderId> <refundId>` 展示售后商品明细。
- 用户提到"售后详情/退款详情/退款原因/卖家信息 + 退款ID"时，执行 `1688 detail <refundId>` 展示退款详情。
- 用户提到"售后留言/沟通记录 + 退款ID"时，执行 `1688 message <refundId>` 展示留言记录。
- 若未提供必要参数，提示用户补全。

## 注意事项

1. **参数必填**：所有子命令都有严格的参数个数要求，缺失时会自动显示帮助。
2. **1688订单号**：`thirdOrderId` 即1688平台的订单号或交易号，是查询售后工单的入口参数。
3. **退款ID**：先通过 `1688` 命令获取退款ID（RefundID），再进一步查询商品、详情或留言。
