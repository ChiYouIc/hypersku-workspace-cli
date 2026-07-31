# get-customer-order-address 输出参考

## 命令

```bash
hypersku-cli customer order address <orderId>
```

## 输出示例

```markdown
【收货地址】
收货人：John Doe
国家：United States
省份：California
城市：Los Angeles
区域：Downtown
详细地址：123 Main Street, Apt 4B
邮编：90001
欧盟税号：EU123456789
VAT：GB123456789
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 收货人 | 收货人姓 + 名 |
| 国家 | 收货国家 |
| 省份 | 收货省份/州 |
| 城市 | 收货城市 |
| 区域 | 收货区域/区县 |
| 详细地址 | 街道及门牌号详情 |
| 邮编 | 邮政编码 |
| 欧盟税号 | 欧盟税务编号（非欧盟订单可能为空） |
| VAT | 增值税号（VAT Number） |
