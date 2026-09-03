# auth login 输出参考

## 命令

```bash
hypersku-cli auth login
```

## 当前状态：暂未实现

`auth login` 为空操作：不执行任何逻辑，打印提示信息后正常返回（退出码 0）。

## 输出示例

```text
login 暂未实现，请手动在 ~/.hypersku-cli/config.json 中配置 api_token
```

## 手动配置登录凭证

登录凭证保存于 `~/.hypersku-cli/config.json`：

```json
{
  "api_base_url": "https://pur.hyperoms.com",
  "api_token": "your-token-here"
}
```

| 字段 | 说明 |
|------|------|
| `api_base_url` | API 基础地址，`auth status` 校验与业务查询均基于此地址 |
| `api_token` | API 认证令牌，同时用于 `auth status` 的 token 校验与业务请求的 `authorization` 请求头 |

## 退出码

| 场景 | 退出码 |
|------|--------|
| 成功（空操作） | 0 |

## 注意事项

- **幂等**：可反复执行，无任何副作用。
- **配置后校验**：手动写入 `api_token` 后，执行 `hypersku-cli auth status` 验证是否生效。
