# pur-count 输出参考

## 命令

`hypersku-cli supplier pur-count <loginId1,loginId2,...>`

多个登录 ID 用英文逗号分隔，一次批量查询。

## 输出示例

```
|供应商登录ID|采购次数|
|----|----|
|xyz_supplier|128|
|yy_craft|56|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 供应商登录ID | 平台侧供应商登录账号 |
| 采购次数 | 该供应商历史被采购次数 |
