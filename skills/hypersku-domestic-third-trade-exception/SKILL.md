---
name: hypersku-domestic-third-trade-exception
display_name: 国内第三方交易异常订单管理
display_name_en: Domestic Third-Party Trade Exception
description_zh: 通过 hypersku-cli domestic-third-trade-exception 子命令查询国内第三方交易异常订单、留言，并按批次生成跟进工单汇总 HTML 报告（主状态×子状态矩阵 + Top 仓库/采购来源分布 + 处理人活跃度 + 每日留言量）。
description_en: Query domestic third-party trade exception orders & messages via hypersku-cli, and render follow-up batch reports (main/sub-status matrix, top warehouses, source types, handler activity, daily messages).
description: HyperSKU 国内第三方交易异常订单管理 + 跟进工单汇总报告。当用户提到异常订单/丢包裹/丢件/未签收/假发货/假签收/退件/未入库/无货，或需要查询异常订单留言/跟进记录，或要做按批次的跟进工单汇总时使用：通过 hypersku-cli domestic-third-trade-exception 子命令查询数据；要求报告形态时只产槽位 JSON 并调用本包 scripts/render_followups_report.py 装配固定版式 HTML。
version: 3.0.0
author: owen
tags:
  - hypersku
  - cli
  - 异常订单
  - 跟进报告
---

# 国内第三方交易异常订单管理

通过 `hypersku-cli domestic-third-trade-exception` 子命令管理 Hypersku 国内第三方交易异常订单（HyperSKU 在第三方平台如 1688/淘宝 采购产生的采购单，其国内段物流异常在此监控），支持按异常主/子状态分页查询异常订单及物流明细、查看异常订单的留言记录。

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| 分页查询异常订单 | 按异常主/子状态分页查询国内第三方交易异常订单（含监控单号、交易号、物流单号、物流公司、状态等） | `hypersku-cli domestic-third-trade-exception page-list [flags]` | [page-list.md](references/page-list.md) |
| 查询异常订单留言 | 根据监控订单ID与监控物流ID查询异常订单的留言记录 | `hypersku-cli domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>` | [message-list.md](references/message-list.md) |
| 跟进工单汇总 | 对一批异常订单 + 留言按多维度聚合输出 HTML 报告 | 由调用方先跑 page-list/message-list，再调用本包 `scripts/render_followups_report.py` | [followups-report.md](references/followups-report.md) |

## 跟进工单汇总 reports（v3.0.0 新增）

### 适用场景

运维 / 客服每日对**国内第三方交易异常订单**做跟进工单盘点：已知一批 `MonitorOrderId`，已批量拉好各物流段的 `message-list`，需要一个 **单批次 HTML 报告**聚合展示：

- 批次总览（订单数 / 物流数 / 留言数）
- 主状态 × 子状态矩阵
- Top 仓库分布
- Top 采购来源分布
- 处理人活跃度（按留言聚合）
- 每日留言量（按天聚合）

### 数据-渲染职责边界（硬规则）

| 主体 | 职责 |
|------|------|
| 调用方（AI） | 跑 `page-list` 取订单 + 对每条 `MonitorLogisticsId` 跑 `message-list` 取留言；按"输出槽位契约"组装 JSON |
| `scripts/render_followups_report.py` | 校验 JSON、HTML 转义、装配固定版式 HTML |
| `templates/followups_report.html` | 只读模板，所有样式 / 色值 / 段落固定 |

**禁止**手写 HTML、禁止修改模板、禁止修改渲染脚本绕过校验——校验失败就回头改 JSON。

### 输出槽位契约

```jsonc
{
  // 报告基础信息（必填）
  "title":       "异常订单跟进工单汇总",        // 报告标题
  "generatedAt": "2026-09-08 14:58",            // 生成时间（YYYY-MM-DD HH:MM），必填
  "query": {                                    // 本批次查询条件（必填）
    "hyperskuStatus":     9,                    // 主状态（1-10，参见 page-list.md）
    "hyperskuSubStatusList": [1, 2],            // 子状态列表（1-5）
    "buyerId":            ""                    // 可选，为空时记 "-"
  },

  // 批次总览（必填，缺一不可）
  "totals": {
    "orderCount":     87,                       // 本批次订单数（Distinct MonitorOrderId 数）
    "logisticsCount": 95,                       // 本批次物流记录数
    "messageCount":   412                       // 本批次留言总条数
  },

  // 主状态 × 子状态 矩阵（必填，至少 1 行）
  "mainStatusMatrix": [
    {
      "key":              "9",
      "mainStatus":       "丢包裹",
      "pending":          0,
      "processing":       87,
      "resolved":         0,
      "closed":           0,
      "rejected":         0,
      "total":            87
    }
  ],

  // Top 仓库分布（可选，省略时该 section 不渲染；条目按 count 降序）
  "byWarehouse": [
    {"name": "广州仓库", "count": 23}
  ],

  // Top 采购来源分布（可选；key 为 sourceType 数字，name 为来源中文）
  "bySourceType": [
    {"key": "1",  "name": "1688",   "count": 70},
    {"key": "18", "name": "淘宝",   "count": 17}
  ],

  // 处理人活跃度（可选；CrtName 视为留言人；messageCount/ordersCount 由 message-list 聚合得出）
  "handlerActivity": [
    {"name": "张三", "messageCount": 87, "ordersCount": 23, "lastMessageAt": "2026-09-08 14:00"}
  ],

  // 每日留言量（可选；按 CrtTime 日期维聚合）
  "dailyMessages": [
    {"date": "2026-09-01", "count": 56},
    {"date": "2026-09-02", "count": 89}
  ]
}
```

### 渲染管线（硬规则）

```
hypersku-cli domestic-third-trade-exception page-list --hypersku-status X --hypersku-sub-status X,Y
        ↓
对每条 logistics 跑 message-list（聚合到内存）
        ↓
AI 按输出槽位契约组装 followups.json 并落盘
        ↓
python scripts/render_followups_report.py --data followups.json -o report.html
        ↓
模板 templates/followups_report.html + JSON → 固定版式 HTML
```

1. **职责边界**：布局/样式/枚举→色值映射固定在只读模板 `templates/followups_report.html`；AI 的产出物只有槽位 JSON。禁止在上下文中现编 HTML、禁止改写模板内容。
2. **先落盘再调用**：JSON 先落盘再执行渲染；缺 `--data` 或校验失败时修正 JSON 后重跑，不得改脚本或模板绕过校验。
3. **信任脚本校验**：脚本对必填字段、枚举值域、矩阵完整性、每行 `total` = 子状态加总做全量断言，任何一项违规即退出码 2 拒绝渲染。
4. **安全**：所有槽位值经 HTML 转义后注入，不支持任何 HTML 标签。

## 注意事项

1. **业务语义**：这里的"异常订单"是 HyperSKU 在第三方平台（1688/淘宝等）的采购单，其国内段物流发生异常（丢包裹/丢件/未签收等）在此监控，**不是** HyperSKU 客户订单；输出中的"交易号/1688订单号"即第三方平台的交易号。
2. **page-list 必填参数**：`--hypersku-status`（异常主状态）与 `--hypersku-sub-status`（异常子状态列表）必填，缺失时会打印可选值并显示帮助。
3. **物流维度输出**：`page-list` 输出按物流记录展开，每个订单的每条物流占一行；订单无物流记录时整单不展示。
4. **状态枚举**：主状态 `--hypersku-status` 可选 1-未发货/2-假发货/3-未到货/4-假签收/5-未签收/6-退件/7-丢件/8-未入库/9-丢包裹/10-无货；子状态 `--hypersku-sub-status` 可选 1-待处理/2-处理中/3-已处理/4-已关闭/5-已拒绝（完整对照见 [page-list.md](references/page-list.md)）。
5. **message-list 参数**：两个位置参数顺序为监控订单ID、监控物流ID，可从 `page-list` 输出中获取。
6. **reports 数据体量**：单批次建议 orderCount ≤ 500；如超出请缩小 `--hypersku-sub-status` 过滤范围或先按仓库拆分批次，避免 message-list 调用次数过多。
7. **handlerActivity 来源**：仅按 `message-list` 出参的 `CrtName` 聚合留言条数；不与 page-list 字段直接挂钩，亦不构成"工单系统在办"语义。
8. **matrix 行 total 自检**：每行 `total` 必须等于 `pending + processing + resolved + closed + rejected`，脚本会拒绝渲染，避免聚合错误。
