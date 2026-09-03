# auth status 输出参考

## 命令

```bash
hypersku-cli auth status            # 文本输出（默认）
hypersku-cli auth status --json     # JSON 输出（程序化集成）
```

## 校验流程

1. 读取 `~/.hypersku-cli/config.json` 中的 `api_base_url` 与 `api_token`。
2. 请求 `GET {api_base_url}/api/admin/user/front/info?token=xxx` 远程校验 token。
3. token 有效时服务端返回用户信息（`data` 载荷）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 用户 ID |
| `name` | string | 姓名 |
| `nickname` | string | 昵称 |
| `proxy` | bool | 是否代理 |
| `roleCode` | string | 角色代码 |
| `roles` | string | 角色集合 |
| `username` | string | 用户名（status 展示的首选标识） |

## 输出示例

已登录（退出码 0，stdout）：

```text
Logged in as owen.chi
```

未登录 / token 失效 / 未配置（退出码 1，stdout）：

```text
Logged out
```

JSON 格式（`--json`）：

```json
{"logged_in":true,"account":"owen.chi"}
{"logged_in":false}
```

## 字段说明

| 字段/输出 | 说明 |
|-----------|------|
| `Logged in as <account>` | token 有效，account 依次取 username/name/nickname，缺省 `unknown` |
| `Logged out` | 未配置 api_token/api_base_url、token 失效或校验请求失败 |
| `logged_in` | JSON 模式布尔字段，登录态 |
| `account` | JSON 模式账号标识，仅登录时返回 |

## 退出码

| 场景 | 退出码 |
|------|--------|
| 已登录 | 0 |
| 未登录 / 未配置 / token 失效 / 校验失败 | 1 |

## 注意事项

- **远程校验**：本命令每次执行都会请求服务端接口实时校验 token，token 失效立即判定为未登录（服务端返回 `{"message":"User Token Forbidden or Expired!","status":40101}` 时视为未登录）。
- **无副作用**：可反复执行，不会修改配置或登录态。
- **凭证来源**：token 由 `~/.hypersku-cli/config.json` 的 `api_token` 字段提供，当前需用户手动配置（`auth login` 暂未实现）。
