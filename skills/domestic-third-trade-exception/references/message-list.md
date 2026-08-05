# message-list 输出参考

## 命令

```bash
hypersku-cli domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>
```

## 参数

| 参数 | 必填 | 说明 |
|------|------|------|
| `monitorOrderId` | 是 | 监控订单 ID |
| `monitorLogisticsId` | 是 | 监控物流 ID |

> 两个参数均为位置参数，顺序不可颠倒；可从 `page-list` 输出中获取。

## 使用示例

```bash
hypersku-cli domestic-third-trade-exception message-list 123 456
```

## 输出示例

```markdown
|留言时间|留言人|留言|
|----|----|----|
|2024-01-03 09:30:00|张三|已联系快递公司核实丢件原因|
|2024-01-04 14:00:00|李四|快递公司确认包裹丢失，已发起理赔|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 留言时间 | 留言创建时间 |
| 留言人 | 留言操作人 |
| 留言 | 留言内容（异常跟进记录） |
