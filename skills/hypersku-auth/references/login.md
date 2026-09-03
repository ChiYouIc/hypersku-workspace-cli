# auth login 输出参考

## 命令

```bash
hypersku-cli auth login            # 输出认证 URL 后即退出（默认，连接器/AI 场景）
hypersku-cli auth login --wait     # 就地轮询等待授权完成（交互式场景）
```

## 流程说明

Device Code Flow（RFC 8628）：

1. CLI 向认证服务申请设备码，并将中间态保存到本地（`~/.hypersku-cli/device-pending.json`）
2. **stdout 第一行输出完整认证 URL**（前后空白分隔、无引号包裹），随后命令立即退出
3. 用户在浏览器打开 URL 确认授权
4. 授权进度由后续 `auth status` 轮询推进（或 `--wait` 模式就地等待）

## 输出示例

默认模式（stdout + stderr 分开展示）：

```text
[stdout]
https://auth.hypersku.com/activate?user_code=ABCD-1234

[stderr]
用户码: ABCD-1234
在浏览器完成授权后，CLI 将自动检测登录状态。
```

`--wait` 模式（授权完成后追加）：

```text
登录成功 (user@example.com)
```

已登录时重复执行（幂等）：

```text
已登录 (user@example.com)
如需切换账号，请先执行 hypersku-cli auth logout
```

## 字段说明

| 输出 | 通道 | 说明 |
|------|------|------|
| 认证 URL | stdout | 完整 `https://` 链接，第一行输出、无引号，自动化集成从此提取 |
| 用户码 | stderr | 供人工核对，通常已包含在 URL 中 |
| `已登录 (账号)` | stdout | 已有有效凭证时的幂等提示，不会再发起授权 |
| `登录成功 (账号)` | stdout | `--wait` 模式授权完成后的确认输出 |

## 退出码

| 场景 | 退出码 |
|------|--------|
| 成功输出 URL / 已登录 / `--wait` 登录成功 | 0 |
| 发起设备授权失败 / 保存状态失败 | 1 |

## 注意事项

- 默认模式**不等待授权**，输出 URL 即退出；登录态需靠 `auth status` 轮询推进。
- `--wait` 模式按服务端约定的间隔轮询（默认 5s，上限 30s），设备码过期会返回"等待授权失败"。
