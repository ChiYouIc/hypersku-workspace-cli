# auth logout 输出参考

## 命令

```bash
hypersku-cli auth logout
```

## 当前状态：暂未实现

`auth logout` 为空操作：不清理任何凭证，打印提示信息后正常返回（退出码 0）。

## 输出示例

```text
logout 暂未实现，请手动删除 ~/.hypersku-cli/config.json 中的 api_token
```

## 手动退出登录

编辑 `~/.hypersku-cli/config.json`，将 `api_token` 置空或删除该字段：

```json
{
  "api_base_url": "https://pur.hyperoms.com",
  "api_token": ""
}
```

## 退出码

| 场景 | 退出码 |
|------|--------|
| 成功（空操作） | 0 |

## 注意事项

- **幂等**：可反复执行，无任何副作用。
- **后续影响**：手动清除 `api_token` 后，所有业务查询命令（purchase/customer/logistics 等）与 `auth status` 均会判定为未登录。
