# get-1688-after-sales-goods 输出参考

## 命令

```bash
hypersku-cli after-sales 1688 goods <thirdOrderId> <refundId>
```

## 输出示例

```markdown
|商品名称|SKU|属性|数量|图片|
|----|----|----|----|----|
|无线蓝牙耳机|2001|颜色:黑色;版本:Pro|1|https://img.example.com/pic1.jpg|
|手机充电器|2002|颜色:白色;功率:65W|2|https://img.example.com/pic2.jpg|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 商品名称 | 售后商品名称 |
| SKU | 商品 SKU ID |
| 属性 | 商品规格属性 |
| 数量 | 售后商品数量 |
| 图片 | 商品图片 URL |
