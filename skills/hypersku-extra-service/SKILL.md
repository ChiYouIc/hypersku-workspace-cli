---
name: hypersku-extra-service
display_name: 增值服务查询
display_name_en: Extra Service Query
description_zh: 通过 hypersku-cli extra-service 子命令查询增值服务名称列表、增值服务关联的仓储与商品、以及仓库开通的增值服务。
description_en: Query extra service names, associated warehouses and goods, and warehouse-enabled extra services via hypersku-cli extra-service commands.
description: HyperSKU 增值服务查询。当用户提到增值服务/增值服务列表/仓库增值服务/增值服务关联商品/增值服务关联仓储时，通过 hypersku-cli extra-service 子命令查询对应数据。
version: 1.0.0
author: owen
---

# 增值服务查询

通过 `hypersku-cli extra-service` 子命令查询增值服务相关信息，支持名称列表、关联仓储、关联商品及仓库增值服务查询。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 增值服务名称列表 | 获取所有可用增值服务名称 | `hypersku-cli extra-service list` | [list.md](references/list.md) |
| 关联仓储查询 | 根据服务 ID 查询关联的仓储 | `hypersku-cli extra-service warehouses <serviceId>` | [warehouses.md](references/warehouses.md) |
| 关联商品查询 | 根据服务 ID 查询关联的商品 | `hypersku-cli extra-service goods <serviceId>` | [goods.md](references/goods.md) |
| 仓库增值服务 | 根据仓库 ID 查询开通的增值服务 | `hypersku-cli extra-service warehouse-add <repertoryId>` | [warehouse-add.md](references/warehouse-add.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

- 用户提到"增值服务列表/有哪些增值服务"时，执行 `list` 查询所有可用增值服务名称。
- 用户提到"增值服务仓储/服务关联仓库"并提供服务 ID 时，执行 `warehouses <serviceId>` 查询关联仓储。
- 用户提到"增值服务商品/服务关联商品"并提供服务 ID 时，执行 `goods <serviceId>` 查询关联商品。
- 用户提到"仓库增值服务/仓库开通了什么增值服务"并提供仓库 ID 时，执行 `warehouse-add <repertoryId>` 查询仓库增值服务。
- 若未提供必要参数（服务 ID 或仓库 ID），提示用户提供。

## 注意事项

1. **只读查询**：所有命令均为只读，不包含写入/修改操作。
2. **ID 参数**：`warehouses` 和 `goods` 需要增值服务 ID（从 `list` 命令获取）；`warehouse-add` 需要仓库 ID。
3. **与仓库管理区别**：本命令专注增值服务维度查询；仓库基础信息请使用 `hypersku-cli warehouse list` 命令。
