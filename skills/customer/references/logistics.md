# get-customer-order-logistics 输出参考

## 命令

```bash
hypersku-cli customer order logistics <orderId>
```

## 输出示例

```markdown
|物流单号|承运商|
|----|----|
|1Z999AA10123456784|UPS|
|LX123456789CN|DHL|
|EE987654321CN|EMS|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 物流单号 | 快递运单号（Tracking Number） |
| 承运商 | 物流承运商名称（如 UPS/DHL/FedEx/EMS） |

> **注意**：该命令仅返回物流单号和承运商，不包含具体物流轨迹。如需查询国内段物流轨迹，请使用 `hypersku-cli logistics tracking <trackingNumber>`。
