# ali-sub-accounts 输出参考

## 命令

`hypersku-cli third-app-api ali-sub-accounts`

无需参数，查询当前 1688 主账号下的所有子账号。

## 输出示例

```
主账号: main_account (b2b-1234567890)
|子账号|MemberID|
|----|----|
|sub_purchase01|b2b-1234567890-01|
|sub_service02|b2b-1234567890-02|
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 主账号 | 主账号登录 ID 与 memberId |
| 子账号 | 子账号登录 ID |
| MemberID | 子账号对应的会员 ID |

无子账号时输出 `无子账号`。
