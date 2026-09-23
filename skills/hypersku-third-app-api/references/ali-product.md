# ali-product 输出参考

## 命令

```bash
hypersku-cli third-app-api ali-product --id <产品ID>
hypersku-cli third-app-api ali-product --url "https://detail.1688.com/offer/xxx.html"
```

`--id` 与 `--url` 二选一，均未指定时报错。

## 输出示例

```
产品ID: 728451930
名称: 硅胶手机壳 全面防摔
英文名: Silicone Phone Case
链接: https://detail.1688.com/offer/728451930.html
原价: 5.80 ¥
销售价: 4.20 ¥
重量: 0.05 kg
净重: 0.04 kg
发货地: 浙江省金华市
库存: 120000
类目: 手机保护壳 / Phone Cases
供应商: 义乌市XX电子有限公司 (156)
支持组合SKU: 是
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 产品ID | 1688 产品数字 ID |
| 名称/英文名 | 产品中英文名称 |
| 链接 | 1688 产品详情页地址 |
| 原价/销售价 | 含币种符号的价格 |
| 重量/净重 | 商品重量（kg） |
| 发货地 | 卖家发货地址 |
| 库存 | 可售库存数量 |
| 类目 | 中文 / 英文类目 |
| 供应商 | 供应商名称与 ID |
| 支持组合SKU | 仅组合商品显示该行 |
