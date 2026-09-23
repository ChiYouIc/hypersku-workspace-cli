# ali-supplier 输出参考

## 命令

```bash
hypersku-cli third-app-api ali-supplier --login-id <1688登录ID>
hypersku-cli third-app-api ali-supplier --domain <店铺地址>
```

`--login-id` 与 `--domain` 二选一。

## 输出示例

```
登录ID: xyz_supplier
供应商名称: 义乌市XX电子有限公司
公司名称: 义乌市XX电子有限公司
类目: 电子元器件
店铺地址: https://xyz.1688.com
跨境宝: 是
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 登录ID | 1688 卖家登录账号 |
| 供应商名称 | HyperSKU 侧供应商名称 |
| 公司名称 | 1688 店铺公司主体名称 |
| 类目 | 主营类目 |
| 店铺地址 | 1688 店铺链接 |
| 跨境宝 | 仅开通跨境宝的卖家显示该行 |
