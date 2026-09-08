#!/usr/bin/env python3
"""
render_profile.py — 客户画像 HTML 渲染器（模板 + 槽位 JSON → 成品 HTML）

职责边界：
  - AI 只产槽位 JSON（评分/策略/话术/完整度），本脚本负责装配成固定版式的 HTML。
  - 模板只读（templates/profile.html），本脚本不改任何样式与结构。

用法：
  python render_profile.py --data profile.json [-o out.html] [--stdout]

槽位 JSON 契约（与 SKILL.md 输出内容契约一致，v3.0.0 增 infoCompletenessLevel）：
  {
    "customerId":            "12345",          // 必填
    "stage":                 "未绑店",          // 必填，枚举：未绑店 | 已绑店未首单
    "leadQuality":           "高",             // 必填，枚举：高 | 中 | 低
    "intentHeat":            "热",             // 必填，枚举：热 | 温 | 冷
    "followPriority":        "P0",             // 必填，枚举：P0 | P1 | P2 | P3
    "infoCompletenessLevel": "高",             // 必填（v3.0.0 新增），枚举：低 | 中 | 高
    "reading":               "...",            // 必填，≤300 字
    "strategy":              "...",            // 必填，≤300 字
    "script":                ["...", "..."]    // 必填，2-3 条，单条 ≤200 字
  }
"""

import argparse
import html
import json
import re
import sys
from datetime import datetime
from pathlib import Path

# 复用同包计算器：followPriority 必须与 leadQuality × intentHeat 查表结果一致
from follow_priority import calculate as calc_follow_priority

# ---------------------------------------------------------------------------
# 契约定义
# ---------------------------------------------------------------------------

ENUM_STAGE       = {"未绑店", "已绑店未首单"}
ENUM_QUALITY     = {"高", "中", "低"}
ENUM_HEAT        = {"热", "温", "冷"}
ENUM_PRIORITY    = {"P0", "P1", "P2", "P3"}
ENUM_COMPLETENESS = {"低", "中", "高"}

REQUIRED_FIELDS = [
    "customerId", "stage", "leadQuality", "intentHeat",
    "followPriority", "infoCompletenessLevel",
    "reading", "strategy", "script",
]

LIMITS = {"reading": 300, "strategy": 300, "script_item": 200}

# 枚举值 → tag/stat 色值 class（与模板 CSS 一一对应，模板不认识之外的值）
CLASS_MAP = {
    "stage":          {"未绑店": "primary", "已绑店未首单": "success"},
    "leadQuality":    {"高": "success", "中": "warning", "低": "danger"},
    "intentHeat":     {"热": "danger", "温": "warning", "冷": "info"},
    "followPriority": {"P0": "danger", "P1": "warning", "P2": "primary", "P3": "info"},
    "infoCompleteness": {"低": "danger", "中": "warning", "高": "success"},
}

PLACEHOLDER_RE = re.compile(r"\{\{\s*([a-zA-Z]+)\s*\}\}")


# ---------------------------------------------------------------------------
# 校验（渲染前全量断言，带病不渲染）
# ---------------------------------------------------------------------------

def validate(data: dict) -> None:
    """校验槽位 JSON；任何违规抛 SystemExit 并给出可读原因。"""
    # 1) 必填齐全
    missing = [f for f in REQUIRED_FIELDS if f not in data]
    if missing:
        fail("缺少必填字段: %s" % ", ".join(missing))

    # 2) 枚举合法
    for field, allowed in [
        ("stage", ENUM_STAGE), ("leadQuality", ENUM_QUALITY),
        ("intentHeat", ENUM_HEAT), ("followPriority", ENUM_PRIORITY),
        ("infoCompletenessLevel", ENUM_COMPLETENESS),
    ]:
        if data[field] not in allowed:
            fail("字段 %s 取值 %r 不在枚举 %s 内" % (field, data[field], sorted(allowed)))

    # followPriority 一致性校验：必须等于交叉表计算结果（AI 不得手填偏离值）
    expected = calc_follow_priority(data["leadQuality"], data["intentHeat"])
    if data["followPriority"] != expected:
        fail("followPriority=%s 与交叉表计算结果 %s 不一致（leadQuality=%s × intentHeat=%s），请用 scripts/follow_priority.py 计算" % (
            data["followPriority"], expected, data["leadQuality"], data["intentHeat"]))

    # 3) 字数上限（中文字符计数）
    for field in ("reading", "strategy"):
        n = len(str(data[field]))
        if n > LIMITS[field]:
            fail("字段 %s 超长：%d 字 > 上限 %d 字" % (field, n, LIMITS[field]))

    # 4) 话术条数与单条长度
    scripts = data["script"]
    if not isinstance(scripts, list) or not (2 <= len(scripts) <= 3):
        fail("script 必须为 2-3 条的列表，当前 %d 条" % (len(scripts) if isinstance(scripts, list) else -1))
    for i, item in enumerate(scripts, 1):
        if not isinstance(item, str) or not item.strip():
            fail("script 第 %d 条为空" % i)
        if len(item) > LIMITS["script_item"]:
            fail("script 第 %d 条超长：%d 字 > 上限 %d 字" % (i, len(item), LIMITS["script_item"]))


def fail(msg: str) -> None:
    print("[render_profile] 校验失败：%s" % msg, file=sys.stderr)
    sys.exit(2)


# ---------------------------------------------------------------------------
# 渲染
# ---------------------------------------------------------------------------

def esc(value: str) -> str:
    """槽位值 HTML 转义，防注入、防破版。"""
    return html.escape(str(value), quote=True)


def render(template_html: str, data: dict) -> str:
    slots = {
        "customerId":   esc(data["customerId"]),
        "stage":        esc(data["stage"]),
        "leadQuality":  esc(data["leadQuality"]),
        "intentHeat":   esc(data["intentHeat"]),
        "followPriority": esc(data["followPriority"]),
        "infoCompleteness": esc(data["infoCompletenessLevel"]),
        "reading":      esc(data["reading"]),
        "strategy":     esc(data["strategy"]),
        "stageClass":              CLASS_MAP["stage"][data["stage"]],
        "leadQualityClass":        CLASS_MAP["leadQuality"][data["leadQuality"]],
        "intentHeatClass":         CLASS_MAP["intentHeat"][data["intentHeat"]],
        "followPriorityClass":     CLASS_MAP["followPriority"][data["followPriority"]],
        "infoCompletenessClass":   CLASS_MAP["infoCompleteness"][data["infoCompletenessLevel"]],
    }

    out = template_html.replace("{{generatedAt}}", datetime.now().strftime("%Y-%m-%d %H:%M"))

    # 列表槽位：从 <template id="tpl-script-item"> 提取片段，逐条克隆填充
    m = re.search(
        r'<template id="tpl-script-item">(.*?)</template>', out, re.DOTALL)
    if not m:
        fail("模板缺少 tpl-script-item 片段")
    item_tpl = m.group(1)
    items_html = "".join(
        item_tpl
        .replace("{{scriptIndex}}", str(i))
        .replace("{{scriptText}}", esc(text))
        for i, text in enumerate(data["script"], 1)
    )
    out = out.replace("<!--SCRIPT_ITEMS-->", items_html)
    out = out.replace(m.group(0), "")  # 移除模板片段本身

    # 标量槽位统一替换
    for key, value in slots.items():
        out = out.replace("{{%s}}" % key, value)

    # 终检：不允许残留任何未填充占位符
    leftover = PLACEHOLDER_RE.findall(out)
    if leftover:
        fail("渲染后仍残留占位符: %s" % ", ".join(sorted(set(leftover))))
    return out


# ---------------------------------------------------------------------------
# 入口
# ---------------------------------------------------------------------------

def main() -> None:
    parser = argparse.ArgumentParser(description="客户画像 HTML 渲染器")
    parser.add_argument("--data", required=True, help="槽位 JSON 文件路径（或 - 读 stdin）")
    parser.add_argument("-o", "--output", default=None, help="输出 HTML 路径，缺省 <customerId>.html")
    parser.add_argument("--stdout", action="store_true", help="输出到 stdout 而非文件")
    parser.add_argument("--template", default=None, help="覆盖模板路径（缺省用同包 templates/profile.html）")
    args = parser.parse_args()

    # 读槽位 JSON（utf-8-sig 容忍 Windows PowerShell 写入的 UTF-8 BOM）
    try:
        raw = sys.stdin.read() if args.data == "-" else Path(args.data).read_text(encoding="utf-8-sig")
        data = json.loads(raw)
    except (OSError, json.JSONDecodeError) as e:
        fail("读取槽位 JSON 失败：%s" % e)

    # 定位模板：默认与本脚本同包（上一级）的 templates/profile.html
    template_path = Path(args.template) if args.template else Path(__file__).resolve().parent.parent / "templates" / "profile.html"
    try:
        template_html = template_path.read_text(encoding="utf-8")
    except OSError as e:
        fail("读取模板失败：%s" % e)

    validate(data)
    result = render(template_html, data)

    if args.stdout:
        sys.stdout.write(result)
    else:
        out_path = Path(args.output) if args.output else Path("%s.html" % data["customerId"])
        out_path.write_text(result, encoding="utf-8")
        print("[render_profile] 已生成: %s" % out_path)


if __name__ == "__main__":
    main()
