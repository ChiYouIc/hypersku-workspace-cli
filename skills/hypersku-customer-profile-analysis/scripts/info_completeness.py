#!/usr/bin/env python3
"""
info_completeness.py — 客户档案信息完整度（infoCompleteness）计算器

输入：hypersku-cli customer detail <customerId> 的脱敏档案（JSON）
输出：综合分（0-100）+ 分档（低/中/高）+ 各层明细

五层基准与权重（与 SKILL.md 输出契约一致）：

  | 层           | 字段数 | 权重 | 字段含义                                                     |
  |--------------|--------|------|--------------------------------------------------------------|
  | 问卷         |    5   |  35% | DS经验/周广告预算/月订单量预期/细分市场/意向服务             |
  | 联系方式     |    1   |  15% | 是否有联系方式                                                |
  | 基础档案     |    6   |  25% | 公司/国家/注册区域/注册时间/最近登录/客户标签                 |
  | 渠道归因     |    6   |  15% | 渠道来源/二级/Medium/Campaign/合作来源代码/链接              |
  | 绑店         |    1   |  10% | 绑定店铺                                                      |

打分规则：
  - 各层 ratio = 该层已填字段数 / 该层字段总数；空值判定与 CLI 对齐：
    * 普通业务字段空 = "-" → 记未填
    * 渠道字段空 = "--" → 记未填
  - 综合分 = Σ(ratio × weight) × 100，保留 1 位小数
  - 分档：低  [0, 40) / 中 [40, 70) / 高 [70, 100]

用法：
  python info_completeness.py --input profile.json          # → JSON
  python info_completeness.py --input profile.json --text   # → 单行可读文本（"中 (65.5%)"）
  cat profile.json | python info_completeness.py --input -  # 读 stdin

退出码：0 成功；2 输入不合法（缺字段、非 dict、读文件失败）。
"""

import argparse
import json
import sys
from pathlib import Path

# ---------------------------------------------------------------------------
# 五层基准（字段名与权重）
# ---------------------------------------------------------------------------

# 问卷五项（任一已填即得分）
QUESTIONNAIRE_FIELDS = ("EngagedTime", "WeeklyAdBudget", "OrderVolume", "Niche", "ServiceInterest")
QUESTIONNAIRE_WEIGHT = 0.35

# 联系方式
CONTACT_FIELDS = ("HasContactWay",)
CONTACT_WEIGHT = 0.15

# 基础档案 6 个核心字段
PROFILE_FIELDS = ("Company", "CountryName", "RegRegion", "SignedUpAt", "LastLoginTime", "Level")
PROFILE_WEIGHT = 0.25

# 渠道归因 6 字段（5/6 视为关键字段；缺失容忍度更低，列入下一版迭代）
CHANNEL_FIELDS = ("ChannelSource", "ChannelSourceSub", "ChannelMedium", "ChannelCampaign", "PartnerSource", "ChannelURL")
CHANNEL_WEIGHT = 0.15

# 绑店
STORE_FIELDS = ("Stores",)
STORE_WEIGHT = 0.10

# 加权和应为 1.0（自检）
LAYERS = [
    ("问卷",     QUESTIONNAIRE_FIELDS, QUESTIONNAIRE_WEIGHT),
    ("联系方式", CONTACT_FIELDS,       CONTACT_WEIGHT),
    ("基础档案", PROFILE_FIELDS,       PROFILE_WEIGHT),
    ("渠道归因", CHANNEL_FIELDS,       CHANNEL_WEIGHT),
    ("绑店",     STORE_FIELDS,         STORE_WEIGHT),
]
assert abs(sum(w for _, _, w in LAYERS) - 1.0) < 1e-9, "层权重之和必须等于 1.0"

# 空值集合（与 CLI 出参口径对齐）
EMPTY_NORMAL = ("", "-")       # 普通业务字段空显示 "-"
EMPTY_CHANNEL = ("", "--")     # 渠道字段空显示 "--"


def is_filled(value, layer_name: str) -> bool:
    """判定字段是否已填。"""
    if value is None:
        return False
    if isinstance(value, bool):
        # 联系方式类布尔字段：True 计已填
        return value is True
    s = str(value).strip()
    if s in EMPTY_NORMAL:
        return False
    if layer_name == "渠道归因" and s in EMPTY_CHANNEL:
        return False
    return True


def calc_layer(data: dict, fields, layer_name: str) -> dict:
    """计算单层完整度。"""
    total = len(fields)
    filled = 0
    for f in fields:
        if f in data and is_filled(data.get(f), layer_name):
            filled += 1
    ratio = filled / total if total else 0.0
    return {
        "filled": filled,
        "total": total,
        "ratio": round(ratio, 4),
        "weight": None,  # 由调用层补齐
        "weighted": round(ratio * 0, 4),  # 占位，外部写入
    }


def level_of(score: float) -> str:
    """综合分 → 分档。"""
    if score >= 70:
        return "高"
    if score >= 40:
        return "中"
    return "低"


def calculate(data: dict) -> dict:
    """计算信息完整度。"""
    if not isinstance(data, dict):
        raise ValueError("输入必须是 JSON 对象（dict）")

    layers = {}
    score = 0.0
    for name, fields, weight in LAYERS:
        info = calc_layer(data, fields, name)
        info["weight"] = weight
        info["weighted"] = round(info["ratio"] * weight * 100, 2)
        score += info["weighted"]
        layers[name] = info

    score = round(score, 1)
    return {
        "score": score,
        "level": level_of(score),
        "layers": layers,
    }


def fail(msg: str) -> None:
    print("[info_completeness] 计算失败：%s" % msg, file=sys.stderr)
    sys.exit(2)


def main() -> None:
    parser = argparse.ArgumentParser(description="客户档案信息完整度计算器（五层加权）")
    parser.add_argument("--input", required=True, help="输入 JSON 路径，或 - 读 stdin")
    parser.add_argument("--text", action="store_true", help="以可读单行输出（默认 JSON）")
    args = parser.parse_args()

    # 读 JSON
    try:
        raw = sys.stdin.read() if args.input == "-" else Path(args.input).read_text(encoding="utf-8-sig")
        data = json.loads(raw)
    except (OSError, json.JSONDecodeError) as e:
        fail("读取输入失败：%s" % e)

    try:
        result = calculate(data)
    except ValueError as e:
        fail(str(e))

    if args.text:
        print("%s (%.1f%%)" % (result["level"], result["score"]))
    else:
        print(json.dumps(result, ensure_ascii=False, indent=2))


if __name__ == "__main__":
    main()
