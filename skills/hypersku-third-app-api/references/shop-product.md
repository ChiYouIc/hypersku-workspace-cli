# shop-product 输出参考

## 命令

`hypersku-cli third-app-api shop-product --store-id <店铺ID>`

`--store-id` 必填。固定查询第 1 页、每页 20 条。

## 输出示例

```json
[{"id":8264983111,"title":"Silicone Phone Case for iPhone 15","status":"active","variants":[{"id":456123","title":"Black","price":"4.20","inventory_quantity":1200}]}]
```

## 字段说明

| 字段 | 说明 |
|------|------|
| (整体) | 原样输出平台返回的 JSON 数组（产品对象列表），不做二次格式化 |
| 查询失败 | 输出 `查询失败: <message>` |
