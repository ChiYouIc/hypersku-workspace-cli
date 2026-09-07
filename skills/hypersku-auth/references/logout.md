# auth logout 输出参考

## 命令

```bash
hypersku-cli auth logout
```

## 登出流程

1. 读取本地配置，将 `api_token` 与 `api_token_updated_at` 置空并写回（保留 `api_base_url` 等其他配置）。
2. 本地凭证清除后即视为未登录，无需调用服务端接口。

## 输出示例

成功登出（退出码 0，stdout）：

```text
Logged out
```

本来就没有 token（退出码 0，stdout）：

```text
Already logged out
```

## 退出码

| 场景 | 退出码 |
|------|--------|
| 登出成功 | 0 |
| 本地无 token（幂等） | 0 |
| 配置写入失败 | 1 |

## 注意事项

- **幂等**：本地无 token 时重复执行输出 `Already logged out`，正常返回退出码 0。
- **只清 token**：`api_base_url` 等其他配置保留，重新登录时无需重新配置基础地址；`api_token_updated_at` 随 token 一并清除。
- **后续影响**：登出后所有业务查询命令（purchase/customer/logistics 等）与 `auth status` 均判定为未登录；需要时重新执行 `hypersku-cli auth login`。
