# shop-order-list 输出参考

## 命令

`hypersku-cli third-app-api shop-order-list --store-id <店铺ID> [--status <状态>]`

`--store-id` 必填；`--status` 可选，取值 `open`/`closed`/`cancelled`/`any`，默认 `any`。固定查询第 1 页、每页 20 条。

## 输出示例

```json
[{"id":54321,"order_number":"1001","financial_status":"paid","fulfillment_status":"fulfilled","total_price":"58.00"}]
```

## 字段说明

| 字段 | 说明 |
|------|------|
| (整体) | 原样输出平台返回的 JSON 数组（订单对象列表），不做二次格式化 |
| 查询失败 | 输出 `查询失败: <message>` |
