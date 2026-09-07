# auth login 输出参考

## 命令

```bash
hypersku-cli auth login <access-token>                                        # 方式一：直接使用已有 token
hypersku-cli auth login --api-user <用户名> --api-password <密码>               # 方式二：用户名密码登录
```

## 登录流程

1. **方式一（access-token）**：CLI 携带该 token 请求 `GET {api_base_url}/api/admin/user/front/info?token=xxx` 校验；校验通过后 token 写入 `~/.hypersku-cli/config.json` 的 `api_token` 字段。
2. **方式二（用户名密码）**：CLI 请求 `POST {api_base_url}/api/auth/jwt/token`（body：`username`/`password`）换取 access token，成功后写入配置文件。

两种方式均未提供凭证时，命令报错退出（退出码 1）。

## 输出示例

成功（退出码 0，stdout）：

```text
Logged in as owen.chi
```

失败（退出码 1，stderr）：

```text
Error: 登录失败: token 校验失败: 响应中无用户数据
```

```text
Error: login fail, username or password is error
```

```text
Error: 请提供 access token，或通过 --api-user 与 --api-password 指定账号密码
```

```text
Error: accepts at most 1 arg(s), received 2
```

## 凭证存储

登录成功后 token 保存于 `~/.hypersku-cli/config.json`：

```json
{
  "api_base_url": "https://pur.hyperoms.com",
  "api_token": "<登录获取的 token>"
}
```

| 字段 | 说明 |
|------|------|
| `api_base_url` | API 基础地址，首次运行 CLI 时自动创建并默认 `https://pur.hyperoms.com`，可按需修改 |
| `api_token` | API 认证令牌，登录成功后自动写入，同时用于 `auth status` 的 token 校验与业务请求的 `authorization` 请求头 |
| `api_token_updated_at` | token 最近一次更新时间（RFC3339 格式，如 `2026-09-07T15:30:00+08:00`），登录成功时自动写入，登出时清除 |

## 退出码

| 场景 | 退出码 |
|------|--------|
| 登录成功 | 0 |
| token 失效 / 账密错误 / 未提供凭证 / 多余位置参数 / 配置写入失败 | 1 |

## 注意事项

- **覆盖语义**：重复登录会覆盖配置中的旧 `api_token`。
- **登录后校验**：可执行 `hypersku-cli auth status` 确认登录态。
- **用户名展示**：方式一成功输出的用户名依次取服务端返回的 username/name/nickname，均缺省时显示 `unknown`。
