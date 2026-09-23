# shop-inventory 输出参考

## 命令

`hypersku-cli third-app-api shop-inventory --store-id <店铺ID>`

`--store-id` 必填。固定查询第 1 页、每页 50 条。

## 输出示例

```
|位置ID|库存项ID|可用数量|
|----|----|----|
|1234567|8264983111|1200|
|1234568|8264983112|350|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 位置ID | 库存位置 ID（location） |
| 库存项ID | 库存商品项 ID（inventory item） |
| 可用数量 | 该位置的可用库存数 |

查询失败时输出 `查询失败: <message>`。
