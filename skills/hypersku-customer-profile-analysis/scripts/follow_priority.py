#!/usr/bin/env python3
"""
follow_priority.py — 跟进优先级（followPriority）计算器

口径来源：平台侧 com.etailerhub.ehub.ai.workbench.analysis.service.FollowPriorityCalculator，
交叉表须两侧同步修改。由 leadQuality（线索质量）× intentHeat（意向热度）查表得出：

| leadQuality \\ intentHeat | 热 | 温 | 冷 |
|--------------------------|----|----|----|
| 高                       | P0 | P1 | P2 |
| 中                       | P1 | P2 | P3 |
| 低                       | P2 | P3 | P3 |

用法：
  python follow_priority.py --lead-quality 高 --intent-heat 温        # → P1
  python follow_priority.py --lead-quality 高 --intent-heat 温 --json # → 完整三元组 JSON

退出码：0 成功；2 入参不合法（fail-fast，对齐 Java 侧语义）。
"""

import argparse
import json
import sys

# 枚举值域（与 SKILL.md 输出内容契约一致）
ENUM_LEAD_QUALITY = ("高", "中", "低")
ENUM_INTENT_HEAT = ("热", "温", "冷")

# 交叉查找表：leadQuality × intentHeat → followPriority
# （结构对齐 Java 侧 CROSS_TABLE，行序 = 线索质量降序）
CROSS_TABLE = {
    "高": {"热": "P0", "温": "P1", "冷": "P2"},
    "中": {"热": "P1", "温": "P2", "冷": "P3"},
    "低": {"热": "P2", "温": "P3", "冷": "P3"},
}


def calculate(lead_quality: str, intent_heat: str) -> str:
    """计算跟进优先级；入参不在枚举内时抛 ValueError（对齐 Java 侧 fail-fast 语义）。"""
    if lead_quality not in CROSS_TABLE:
        raise ValueError("leadQuality 取值 %r 不在枚举 %s 内" % (lead_quality, list(CROSS_TABLE)))
    row = CROSS_TABLE[lead_quality]
    if intent_heat not in row:
        raise ValueError("intentHeat 取值 %r 不在枚举 %s 内" % (intent_heat, list(row)))
    return row[intent_heat]


def fail(msg: str) -> None:
    print("[follow_priority] 计算失败：%s" % msg, file=sys.stderr)
    sys.exit(2)


def main() -> None:
    parser = argparse.ArgumentParser(description="跟进优先级计算器（leadQuality × intentHeat 查表）")
    parser.add_argument("--lead-quality", required=True, choices=ENUM_LEAD_QUALITY, help="线索质量：高/中/低")
    parser.add_argument("--intent-heat", required=True, choices=ENUM_INTENT_HEAT, help="意向热度：热/温/冷")
    parser.add_argument("--json", action="store_true", help="以 JSON 输出完整三元组")
    args = parser.parse_args()

    try:
        priority = calculate(args.lead_quality, args.intent_heat)
    except ValueError as e:
        fail(str(e))

    if args.json:
        print(json.dumps({
            "leadQuality": args.lead_quality,
            "intentHeat": args.intent_heat,
            "followPriority": priority,
        }, ensure_ascii=False))
    else:
        print(priority)


if __name__ == "__main__":
    main()
