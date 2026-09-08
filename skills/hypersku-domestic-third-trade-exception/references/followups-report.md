# 跟进工单汇总报告（followups-report）

由调用方（AI）跑一次 `hypersku-cli domestic-third-trade-exception page-list` 取一批异常订单，再对每条 logistics 跑 `message-list` 取留言，按 SKILL.md 中的"输出槽位契约"组装 JSON 并落盘，然后调用本包 `scripts/render_followups_report.py` 装配为 HTML 报告。

## 数据流

```
hypersku-cli domestic-third-trade-exception page-list --hypersku-status N --hypersku-sub-status X,Y
        ↓ 每行 logistics 记录带 MonitorOrderId + MonitorLogisticsId
对每条 logistics：
  hypersku-cli domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>
        ↓ 收集 CrtTime / CrtName / Remark
按 SKILL.md "输出槽位契约"组装 followups.json 并落盘
        ↓
python scripts/render_followups_report.py --data followups.json -o report.html
        ↓
templates/followups_report.html + JSON → 固定版式 HTML
```

## 数据维度

报告核心信息：

1. **批次总览**：orderCount / logisticsCount / messageCount
2. **主状态 × 子状态 矩阵**：每行 = 一个主状态；每列 = 一个子状态；行尾 total = 5 列加总（渲染前自检）
3. **Top 仓库分布**（可选）
4. **Top 采购来源分布**（可选）
5. **处理人活跃度**（可选）：来自 message-list 的 CrtName 聚合
6. **每日留言量**（可选）

可选 section 在对应数据缺失时被自动隐藏（section 节点加 `style="display:none"`）。

## 数据示例与契约校验

完整 JSON 示例与校验规则由 `scripts/render_followups_report.py` 内 contract + `references/page-list.md` 一并覆盖；调用方按 SKILL.md 输出槽位契约组装即可。校验失败由渲染脚本退出码 2 兜底（AI 修改 JSON，不要改脚本或模板绕过校验）。

## 注意

- `handlerActivity` 来自 message-list 字段 `CrtName`；不与 `page-list` 直接关联，亦不构成"工单系统在办"语义。
- `dailyMessages` 来自 message-list 字段 `CrtTime` 的日期切分（小时分舍去）。
- 单批次建议 `orderCount` ≤ 500；超出请按 `--hypersku-sub-status` 拆批。
