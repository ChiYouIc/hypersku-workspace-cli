# auth status 输出参考

## 命令

```bash
hypersku-cli auth status            # 文本输出（默认）
hypersku-cli auth status --json     # JSON 输出（程序化集成）
```

## 输出示例

已登录（退出码 0）：

```text
Logged in as user@example.com
```

未登录（退出码 1）：

```text
Logged out
```

JSON 格式（`--json`）：

```json
{"logged_in":true,"account":"user@example.com"}
{"logged_in":false}
```

## 字段说明

| 字段/输出 | 说明 |
|-----------|------|
| `Logged in as <account>` | 已登录且凭证有效，account 为账号标识（缺省 `unknown`） |
| `Logged out` | 未登录或凭证已过期 |
| `logged_in` | JSON 模式布尔字段，登录态 |
| `account` | JSON 模式账号标识，仅登录时返回 |

## 退出码

| 场景 | 退出码 |
|------|--------|
| 已登录 | 0 |
| 未登录 / 凭证过期 | 1 |
| 读取凭证失败 | 1 |

## 注意事项

- **幂等、无副作用**：可反复执行，不会重置登录态。
- **轮询推进**：若存在进行中的设备授权（`auth login` 后未完成），status 会自动尝试推进一次——服务端发放 token 后自动保存凭证并转为已登录。因此连接器/AI 场景的登录轮询就是循环执行本命令。
- **凭证有效期**：凭证过期后视为未登录，需重新执行 `auth login`。
