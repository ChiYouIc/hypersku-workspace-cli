# auth logout 输出参考

## 命令

```bash
hypersku-cli auth logout
```

## 输出示例

```text
Logged out
```

未登录时执行结果相同（幂等，退出码 0）。

## 执行内容

1. 删除本地凭证（`~/.hypersku-cli/credentials.json`）
2. 清理设备授权中间态（`~/.hypersku-cli/device-pending.json`）

## 退出码

| 场景 | 退出码 |
|------|--------|
| 成功（含未登录时执行） | 0 |
| 清理凭证/状态失败 | 1 |

## 注意事项

- **幂等**：未登录时执行同样正常返回，不报错。
- **切换账号**：必须先 `logout` 清理旧凭证，再 `auth login`，否则 login 会提示已登录。
- **后续影响**：登出后所有业务查询命令（purchase/customer/logistics 等）将失败，需重新登录。
