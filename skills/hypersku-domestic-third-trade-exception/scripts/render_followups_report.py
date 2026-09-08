#!/usr/bin/env python3
"""
render_followups_report.py — 国内第三方交易异常订单跟进工单汇总报告渲染器

数据流：
  调用方运行 hypersku-cli domestic-third-trade-exception page-list + 多次 message-list
        ↓
  按 SKILL.md "输出槽位契约"组装 followups.json
        ↓
  python scripts/render_followups_report.py --data followups.json -o report.html
        ↓
  模板 templates/followups_report.html + JSON → 固定版式 HTML

职责边界：
  - AI 只产槽位 JSON（聚合数据），本脚本负责校验、HTML 转义、装配。
  - 模板只读（templates/followups_report.html），本脚本不改任何样式与结构。

用法：
  python render_followups_report.py --data followups.json [-o out.html] [--stdout]
  cat followups.json | python render_followups_report.py --data -

退出码：0 成功；2 输入不合法（fail-fast）。

槽位 JSON 契约（必填）：
  title             报告标题
  generatedAt       生成时间（YYYY-MM-DD HH:MM）
  query             {hyperskuStatus:int 1-10, hyperskuSubStatusList:[int 1-5], buyerId:str}
  totals            {orderCount:int >=0, logisticsCount:int >=0, messageCount:int >=0}
  mainStatusMatrix  [{key:str, mainStatus:str,
                      pending:int, processing:int, resolved:int, closed:int, rejected:int,
                      total:int == pending+processing+resolved+closed+rejected}]

槽位 JSON 契约（可选，按数据是否喂入渲染；为空时该 section 不渲染）：
  byWarehouse       [{name:str, count:int >=0}]
  bySourceType      [{key:str, name:str, count:int >=0}]
  handlerActivity   [{name:str, messageCount:int >=0, ordersCount:int >=0,
                      lastMessageAt:str YYYY-MM-DD HH:MM:SS}]
  dailyMessages     [{date:str YYYY-MM-DD, count:int >=0}]
"""

import argparse
import html
import json
import re
import sys
from pathlib import Path

# ---------------------------------------------------------------------------
# 契约定义
# ---------------------------------------------------------------------------

# 异常主状态（与 page-list.md 一致：1-未发货…10-无货）
VALID_MAIN_STATUS = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

# 异常子状态（与 page-list.md 一致：1-待处理…5-已拒绝）
VALID_SUB_STATUS = {1, 2, 3, 4, 5}

# date / datetime 格式
GEN_TIME_RE = re.compile(r"^\d{4}-\d{2}-\d{2} \d{2}:\d{2}(:\d{2})?$")
DATE_RE = re.compile(r"^\d{4}-\d{2}-\d{2}$")
DATETIME_RE = re.compile(r"^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$")

REQUIRED_FIELDS = ["title", "generatedAt", "query", "totals", "mainStatusMatrix"]


# ---------------------------------------------------------------------------
# 错误处理
# ---------------------------------------------------------------------------

def fail(msg: str) -> None:
    """校验失败：退出码 2。"""
    print("[render_followups_report] 校验失败：%s" % msg, file=sys.stderr)
    sys.exit(2)


# ---------------------------------------------------------------------------
# 校验：必填齐全 + 枚举合法 + 数值合法 + matrix 一致
# ---------------------------------------------------------------------------

def _ensure_int(v, where: str, lo: int = 0) -> int:
    if isinstance(v, bool) or not isinstance(v, int) or v < lo:
        fail("%s 必须是 ≥%d 的整数，当前 %r" % (where, lo, v))
    return v


def validate(data: dict) -> None:
    if not isinstance(data, dict):
        fail("输入必须是 JSON 对象（dict）")

    # 1) 必填齐全
    missing = [f for f in REQUIRED_FIELDS if f not in data]
    if missing:
        fail("缺少必填字段: %s" % ", ".join(missing))

    # 2) 基础字段
    title = data.get("title", "")
    if not isinstance(title, str) or not title.strip():
        fail("title 必须为非空字符串")

    generated_at = data.get("generatedAt", "")
    if not isinstance(generated_at, str) or not GEN_TIME_RE.match(generated_at):
        fail("generatedAt 必须匹配 YYYY-MM-DD HH:MM[:SS]，当前 %r" % generated_at)

    # 3) query
    query = data.get("query")
    if not isinstance(query, dict):
        fail("query 必须为 dict")
    hs_status = _ensure_int(query.get("hyperskuStatus"), "query.hyperskuStatus")
    if hs_status not in VALID_MAIN_STATUS:
        fail("query.hyperskuStatus=%d 不在主状态枚举 1-10 内" % hs_status)

    sub_list = query.get("hyperskuSubStatusList")
    if not isinstance(sub_list, list) or len(sub_list) == 0:
        fail("query.hyperskuSubStatusList 必须为非空列表")
    for s in sub_list:
        _ensure_int(s, "query.hyperskuSubStatusList[]")
        if s not in VALID_SUB_STATUS:
            fail("query.hyperskuSubStatusList 内 %d 不在子状态枚举 1-5 内" % s)

    buyer_id = query.get("buyerId", "")
    if not isinstance(buyer_id, str):
        fail("query.buyerId 必须为字符串")

    # 4) totals
    totals = data.get("totals")
    if not isinstance(totals, dict):
        fail("totals 必须为 dict")
    order_count = _ensure_int(totals.get("orderCount"),    "totals.orderCount")
    logistics_count = _ensure_int(totals.get("logisticsCount"), "totals.logisticsCount")
    message_count = _ensure_int(totals.get("messageCount"),  "totals.messageCount")
    if order_count == 0:
        fail("totals.orderCount 不能为 0（批次元数据缺失，请回头排查 page-list 数据）")
    if message_count == 0 and "handlerActivity" in data:
        # 若有 handlerActivity 数据则 messageCount 必须 > 0
        fail("totals.messageCount=0 与 handlerActivity 互相矛盾")

    # 5) mainStatusMatrix 完整 + 每行 total = 子状态加总
    matrix = data.get("mainStatusMatrix")
    if not isinstance(matrix, list) or len(matrix) == 0:
        fail("mainStatusMatrix 必须为非空列表")
    row_total = 0
    for i, row in enumerate(matrix, 1):
        if not isinstance(row, dict):
            fail("mainStatusMatrix 第 %d 行必须为 dict" % i)
        if "key" not in row or "mainStatus" not in row:
            fail("mainStatusMatrix 第 %d 行缺 key 或 mainStatus" % i)
        try:
            row_key_int = int(row["key"])
        except (TypeError, ValueError):
            fail("mainStatusMatrix 第 %d 行 key 必须是整数（字符串形式也可），当前 %r" % (i, row["key"]))
        if row_key_int not in VALID_MAIN_STATUS:
            fail("mainStatusMatrix 第 %d 行 key=%d 不在主状态枚举 1-10 内" % (i, row_key_int))
        if not isinstance(row["mainStatus"], str) or not row["mainStatus"].strip():
            fail("mainStatusMatrix 第 %d 行 mainStatus 必须为非空字符串" % i)

        children = ["pending", "processing", "resolved", "closed", "rejected"]
        vals = []
        for c in children:
            vals.append(_ensure_int(row.get(c), "mainStatusMatrix[%d].%s" % (i, c)))
        declared_total = _ensure_int(row.get("total"), "mainStatusMatrix[%d].total" % i)
        if declared_total != sum(vals):
            fail("mainStatusMatrix 第 %d 行 total=%d 与子状态加总 %d 不一致" % (
                i, declared_total, sum(vals)))
        row_total += declared_total

    # matrix 总体 total 与 totals.orderCount 不强制相等（matrix 是异常状态聚合；
    # 同一订单在部分系统中可能贡献多种主状态——这里仅做软警告，留给调用方对齐）

    # 6) 可选 section 形状
    for sec, fmt_check in [
        ("byWarehouse",     lambda x: isinstance(x, list)),
        ("bySourceType",    lambda x: isinstance(x, list)),
        ("handlerActivity", lambda x: isinstance(x, list)),
        ("dailyMessages",   lambda x: isinstance(x, list)),
    ]:
        if sec not in data:
            continue
        if not fmt_check(data[sec]):
            fail("%s 必须为列表" % sec)
        for j, item in enumerate(data[sec], 1):
            if not isinstance(item, dict):
                fail("%s 第 %d 项必须为 dict" % (sec, j))

    # byWarehouse 详情
    if "byWarehouse" in data:
        for j, item in enumerate(data["byWarehouse"], 1):
            if not isinstance(item.get("name"), str) or not item["name"].strip():
                fail("byWarehouse[%d].name 必须非空" % j)
            _ensure_int(item.get("count"), "byWarehouse[%d].count" % j)

    # bySourceType 详情
    if "bySourceType" in data:
        for j, item in enumerate(data["bySourceType"], 1):
            if not isinstance(item.get("key"), str):
                fail("bySourceType[%d].key 必须为字符串" % j)
            if not isinstance(item.get("name"), str) or not item["name"].strip():
                fail("bySourceType[%d].name 必须非空" % j)
            _ensure_int(item.get("count"), "bySourceType[%d].count" % j)

    # handlerActivity 详情
    if "handlerActivity" in data:
        for j, item in enumerate(data["handlerActivity"], 1):
            if not isinstance(item.get("name"), str) or not item["name"].strip():
                fail("handlerActivity[%d].name 必须非空" % j)
            _ensure_int(item.get("messageCount"),  "handlerActivity[%d].messageCount" % j)
            _ensure_int(item.get("ordersCount"),   "handlerActivity[%d].ordersCount" % j)
            la = item.get("lastMessageAt", "")
            if not isinstance(la, str) or not DATETIME_RE.match(la):
                fail("handlerActivity[%d].lastMessageAt 必须匹配 YYYY-MM-DD HH:MM:SS，当前 %r" % (j, la))

    # dailyMessages 详情
    if "dailyMessages" in data:
        for j, item in enumerate(data["dailyMessages"], 1):
            d = item.get("date", "")
            if not isinstance(d, str) or not DATE_RE.match(d):
                fail("dailyMessages[%d].date 必须匹配 YYYY-MM-DD，当前 %r" % (j, d))
            _ensure_int(item.get("count"), "dailyMessages[%d].count" % j)


# ---------------------------------------------------------------------------
# HTML 转义
# ---------------------------------------------------------------------------

def esc(value) -> str:
    if value is None:
        return ""
    return html.escape(str(value), quote=True)


# ---------------------------------------------------------------------------
# 装配：从模板提取片段，逐条克隆填充
# ---------------------------------------------------------------------------

PLACEHOLDER_RE = re.compile(r"\{\{\s*([a-zA-Z_]+)\s*\}\}")
TEMPLATE_BLOCK_RE = re.compile(
    r'<template\s+id="(tpl-[a-z0-9-]+)"[^>]*>(.*?)</template>', re.DOTALL)


def _render_rows(template_html: str, tpl_id: str, items, placeholder: str, row_func) -> str:
    """查找 tpl-id 片段，对 items 逐条克隆填充，删除原片段并把渲染结果
    替换到 placeholder 处；返回更新后的 HTML。"""
    m = re.search(
        r'<template\s+id="%s"[^>]*>(.*?)</template>' % tpl_id, template_html, re.DOTALL)
    if not m:
        fail("模板缺少 %s 片段" % tpl_id)
    item_tpl = m.group(1)
    rendered = "".join(
        row_func(item_tpl, i, item)
        for i, item in enumerate(items, 1)
    )
    new_html = template_html.replace(m.group(0), "", 1)  # 移除模板片段本身（避免 {{rowXxx}} 残留）
    return new_html.replace(placeholder, rendered)


def _fill_tpl(item_tpl: str, replacements: dict) -> str:
    out = item_tpl
    for k, v in replacements.items():
        out = out.replace("{{%s}}" % k, v)
    return out


def render(template_html: str, data: dict) -> str:
    """把 JSON 装配到模板，生成成品 HTML。"""
    out = template_html

    # ---- 标量槽位 ----
    scalar_slots = {
        "title":              esc(data["title"]),
        "generatedAt":        esc(data["generatedAt"]),
        "queryStatusLabel":   _format_hypersku_status(data["query"]["hyperskuStatus"]),
        "querySubStatusLabel": _format_hypersku_sub_statuses(data["query"]["hyperskuSubStatusList"]),
        "queryBuyerLabel":    esc(data["query"].get("buyerId") or "-"),
        "totalOrderCount":    esc(data["totals"]["orderCount"]),
        "totalLogisticsCount": esc(data["totals"]["logisticsCount"]),
        "totalMessageCount":  esc(data["totals"]["messageCount"]),
    }

    # 矩阵子状态加总（按列聚合）
    matrix = data["mainStatusMatrix"]
    matrix_total = {"pending": 0, "processing": 0, "resolved": 0, "closed": 0, "rejected": 0, "total": 0}
    for row in matrix:
        for k in matrix_total:
            matrix_total[k] += row.get(k, 0)

    scalar_slots.update({
        "matrixTotalPending":     esc(matrix_total["pending"]),
        "matrixTotalProcessing":  esc(matrix_total["processing"]),
        "matrixTotalResolved":    esc(matrix_total["resolved"]),
        "matrixTotalClosed":      esc(matrix_total["closed"]),
        "matrixTotalRejected":    esc(matrix_total["rejected"]),
        "matrixTotalTotal":       esc(matrix_total["total"]),
    })

    for key, val in scalar_slots.items():
        out = out.replace("{{%s}}" % key, val)

    # ---- 列表槽位 ----
    # 矩阵主体
    out = _render_rows(
        out, "tpl-matrix-row", matrix, "<!--MATRIX_ROWS-->",
        lambda tpl, i, row: _fill_tpl(tpl, {
            "rowMainKey":    esc(row["key"]),
            "rowMainStatus": esc(row["mainStatus"]),
            "rowPending":    esc(row.get("pending", 0)),
            "rowProcessing": esc(row.get("processing", 0)),
            "rowResolved":   esc(row.get("resolved", 0)),
            "rowClosed":     esc(row.get("closed", 0)),
            "rowRejected":   esc(row.get("rejected", 0)),
            "rowTotal":      esc(row.get("total", 0)),
        })
    )

    # 可选 list section：数据为空时仅插入空字符串（依赖 _toggle_section 把整 section 隐藏）
    for tpl_id, placeholder, key, row_func in [
        ("tpl-warehouse-row", "<!--WAREHOUSE_ROWS-->", "byWarehouse",
         lambda tpl, i, item: _fill_tpl(tpl, {
             "rowIdx":   esc(i),
             "rowName":  esc(item["name"]),
             "rowCount": esc(item["count"]),
         })),
        ("tpl-source-row", "<!--SOURCE_ROWS-->", "bySourceType",
         lambda tpl, i, item: _fill_tpl(tpl, {
             "rowKey":   esc(item["key"]),
             "rowName":  esc(item["name"]),
             "rowCount": esc(item["count"]),
         })),
        ("tpl-handler-row", "<!--HANDLER_ROWS-->", "handlerActivity",
         lambda tpl, i, item: _fill_tpl(tpl, {
             "rowIdx":           esc(i),
             "rowName":          esc(item["name"]),
             "rowMessageCount":  esc(item["messageCount"]),
             "rowOrdersCount":   esc(item["ordersCount"]),
             "rowLastMessageAt": esc(item["lastMessageAt"]),
         })),
        ("tpl-daily-row", "<!--DAILY_ROWS-->", "dailyMessages",
         lambda tpl, i, item: _fill_tpl(tpl, {
             "rowDate":  esc(item["date"]),
             "rowCount": esc(item["count"]),
         })),
    ]:
        items = data.get(key, [])
        if items:
            out = _render_rows(out, tpl_id, items, placeholder, row_func)
        else:
            # 数据为空：删除模板片段并插入空字符串（section 自身会被 _toggle_section 隐藏）
            m = re.search(r'<template id="%s">.*?</template>' % tpl_id, out, re.DOTALL)
            if m:
                out = out.replace(m.group(0), "", 1)
            out = out.replace(placeholder, "")

    # ---- section 显隐：可选 section 数据空时，整段替换为占位 + 由 CSS 隐藏 ----
    out = _toggle_section(out, "section-warehouse",     "WAREHOUSE_DATA_FLAG",     "byWarehouse",     data)
    out = _toggle_section(out, "section-source",        "SOURCE_DATA_FLAG",        "bySourceType",    data)
    out = _toggle_section(out, "section-handler",       "HANDLER_DATA_FLAG",       "handlerActivity", data)
    out = _toggle_section(out, "section-daily",         "DAILY_DATA_FLAG",         "dailyMessages",   data)

    # ---- 终检：残留占位符 ----
    leftover = PLACEHOLDER_RE.findall(out)
    if leftover:
        fail("渲染后残留占位符: %s" % ", ".join(sorted(set(leftover))))
    return out


def _toggle_section(html_text: str, section_class: str, flag_id: str, key: str, data: dict) -> str:
    """在 section 起始标签加 data-flag 占位，并在 CSS 路径上控制显隐（CSS 文件不引入
    时，将带 data-empty="true" 标记）。最简实现：在 section 的 div 上加 data-empty
    属性，模板 CSS 用属性选择器 [data-empty="true"] { display:none } 控制。
    为避免硬编 CSS，本脚本在 data 缺失时直接给整个 section 加 display:none style。"""
    has_data = bool(data.get(key))
    # 找首个含 section_class 的 <section ...> 起首标签，加 style="display:none | none"
    style = "display:none" if not has_data else "display:block"
    return re.sub(
        r'(<section[^>]*class="[^"]*\b%s\b[^"]*"[^>]*)>' % section_class,
        lambda m: "%s style=\"%s\" data-empty=\"%s\">" % (m.group(1), style, "true" if not has_data else "false"),
        html_text, count=1)


# ---------------------------------------------------------------------------
# 枚举→文案
# ---------------------------------------------------------------------------

MAIN_STATUS_TEXT = {
    1: "未发货", 2: "假发货", 3: "未到货", 4: "假签收", 5: "未签收",
    6: "退件",   7: "丢件",   8: "未入库", 9: "丢包裹", 10: "无货",
}
SUB_STATUS_TEXT = {
    1: "待处理", 2: "处理中", 3: "已处理", 4: "已关闭", 5: "已拒绝",
}


def _format_hypersku_status(s: int) -> str:
    return "%d-%s" % (s, MAIN_STATUS_TEXT.get(s, "未知"))


def _format_hypersku_sub_statuses(lst) -> str:
    parts = []
    for s in lst:
        parts.append("%d-%s" % (s, SUB_STATUS_TEXT.get(s, "未知")))
    return "+".join(parts)


# ---------------------------------------------------------------------------
# 入口
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(description="国内第三方交易异常订单跟进工单汇总 HTML 渲染器")
    parser.add_argument("--data", required=True, help="槽位 JSON 路径，或 - 读 stdin")
    parser.add_argument("-o", "--output", default=None, help="输出 HTML 路径，缺省 followups_report.html")
    parser.add_argument("--stdout", action="store_true", help="输出到 stdout 而非文件")
    args = parser.parse_args()

    # 读 JSON（容忍 UTF-8 BOM，Windows PowerShell 写入习惯）
    try:
        raw = sys.stdin.read() if args.data == "-" else Path(args.data).read_text(encoding="utf-8-sig")
        data = json.loads(raw)
    except (OSError, json.JSONDecodeError) as e:
        fail("读取槽位 JSON 失败：%s" % e)

    # 默认模板：同包 templates/followups_report.html
    template_path = Path(__file__).resolve().parent.parent / "templates" / "followups_report.html"
    try:
        template_html = template_path.read_text(encoding="utf-8")
    except OSError as e:
        fail("读取模板失败：%s" % e)

    validate(data)
    result = render(template_html, data)

    if args.stdout:
        sys.stdout.write(result)
    else:
        out_path = Path(args.output) if args.output else Path("followups_report.html")
        out_path.write_text(result, encoding="utf-8")
        print("[render_followups_report] 已生成: %s" % out_path)


if __name__ == "__main__":
    main()
