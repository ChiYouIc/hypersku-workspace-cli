---
name: hypersku-customer-profile-analysis
display_name: 客户转化画像
display_name_en: Customer Conversion Profile
description_zh: 对**已注册且无已支付订单**的客户生成转化画像：leadQuality 线索质量、intentHeat 意向热度、followPriority 跟进优先级三项综合评分，附信息完整度（infoCompleteness 五层基准）+ 线索速读 + 按转化阶段分流的转化策略与销售话术。
description_en: Generate conversion profiles for registered customers with no paid orders: three scores (leadQuality, intentHeat, followPriority) plus infoCompleteness (5-layer basis), reading summary, stage-specific conversion strategies, and sales scripts.
description: HyperSKU 客户转化画像。本 skill 限定"已注册且无已支付订单"的池内客户，输出三项评分（leadQuality/intentHeat/followPriority）+ 信息完整度（infoCompleteness 五层加权）+ reading + strategy + script；followPriority 由本包 scripts/follow_priority.py 查表计算、infoCompleteness 由本包 scripts/info_completeness.py 五层加权计算；要求网页形态时只产槽位 JSON 并调用本包 scripts/render_profile.py 装配固定版式 HTML；与 hypersku-customer-profile（基础画像汇总报告）严格分工——后者覆盖全部 hypersku-cli customer 子命令的档案汇总，本 skill 只关注未转化客户的销售视角评估。
version: 3.0.0
author: owen
tags:
  - hypersku
  - 客户画像
  - 转化
  - 销售
---

# 客户转化画像（customer-profile-analysis）

对**已注册且无已支付订单**的客户生成转化画像：三项综合评分 + 信息完整度 + 线索速读 + 转化策略 + 可用话术，帮助销售决定"跟不跟、什么时候跟、说什么"。

**本 skill 约束评估规则与输出内容**——调用方可要求组织为网页、markdown 或系统触发的结构化 JSON；未指定形态时按易读的结构化文本输出。**网页形态必须走渲染管线**（见下文），禁止 AI 直接手写 HTML。

## 与 hypersku-customer-profile 的边界（硬规则）

本 skill 与新建的 `hypersku-customer-profile` skill 是**两个独立 skill，互不调用、不共享模板、不互相引用**：

| 维度 | hypersku-customer-profile-analysis（本 skill） | hypersku-customer-profile |
|------|-----------------------------------------------|---------------------------|
| 用途 | 销售视角：未转化客户跟不跟、跟什么、说什么 | 客户视角：全客户的多维档案汇总报告 |
| 客户范围 | 仅**已注册且无已支付订单**（含已绑店未首单、下单未付） | **全部客户**（含已转化） |
| 数据来源 | `hypersku-cli customer detail`（+ 可选 profile order/transaction 聚合） | `hypersku-cli customer detail` + `customer order *` + `customer profile order *` + `customer profile transaction *` + `customer order return` |
| 核心输出 | 三项评分（leadQuality/intentHeat/followPriority） + infoCompleteness + reading + strategy + script | 档案概览 + 订单/交易/退件多维度聚合 |
| 渲染管线 | 本包 `scripts/render_profile.py` + `templates/profile.html` | 独立管线 `scripts/render_customer_profile.py` + `templates/customer_profile.html` |

**冲突处理**：分析同一客户时，应分别由两个 skill 各产出各的产物；不要在本 skill 内塞进订单/退件清单，也不要在 `customer-profile` 内塞三项评分。

## 分析对象与前置校验（硬规则）

1. **池资格**：仅分析"已注册且**是否有已支付订单＝否**"的客户。若**是否有已支付订单＝是**，停止分析并告知：该客户已转化出池，不在画像分析范围。
2. **下单未付仍在池内**（属催单场景），正常分析。
3. **转化阶段判定来自实时数据，不由 AI 猜**：
   - 绑定店铺为 `-`（空）→ **未绑店**
   - 绑定店铺非空（如"3 家店铺"）→ **已绑店未首单**（池内客户无已支付订单，绑店即处于此阶段）

## 数据输入

执行 `hypersku-cli customer detail <customerId>` 获取脱敏档案（出参字段口径见客户管理能力域 `customer` 的参考文档 references/customer-detail.md）。解析口径：

| 信号 | 出参字段 | 解析 |
|------|----------|------|
| 消费状态 | 是否有已支付订单 | 是/否；用于池资格校验 |
| 绑店 | 绑定店铺 | `-` 记未绑店；"N 家店铺"记已绑店 |
| 联系方式 | 是否有联系方式 | 是/否 |
| 姓名 | 是否有姓名 | 是/否 |
| DS 经验 | DS 经验 | 枚举文案（如"多于1年"）；未命中显原始值 |
| 广告预算 | 周广告预算 | 枚举文案（如"500-1000 USD"） |
| 订单量预期 | 月订单量预期 | 枚举文案（如"500+"） |
| 细分市场 | 细分市场 | 原始数字 1-10 |
| 意向服务 | 意向服务 | 多选逗号串（1 dropshipping / 2 DTC / 3 POD / 4 Merch / 5 Wholesale）；`-` 记未填 |
| 最近登录 | 最近登录时间 | 时间戳；`-` 记未填 |
| 国家/区域 | 国家 / 注册区域 | 英文名码 / 省级中文描述；`-` 记未填 |
| 业务标签 | 客户标签 | 文案（0 默认 / 1 老用户 / 2 新用户 / 3 潜在用户 / 4 流失用户） |
| 用户等级 | 用户等级 / 订单量级别 | L0.. / O0.. 标识 |
| 渠道归因 | 渠道来源 / 二级 / Medium / Campaign / 合作来源代码 / 来源链接 | `--` 记未归因 |

**空值约定**：普通业务字段空显示 `-`、渠道字段空显示 `--`，一律记"未填"，不得当作负面证据，只作"信号缺失"处理。

### 已录购买深度（可选增强输入，本 skill 不抓取）

为支持"已绑店未首单"客户在历史订单窗内的购买画像，本 skill **可选**接收以下聚合字段（调用方从 `hypersku-cli customer profile order count` + `profile transaction count` 喂入；本 skill 不爬页面）：

| 聚合字段 | 来源 | 用途 |
|----------|------|------|
| `recorded_order_total` | `customer profile order count` 的 Total | 评估已有订单体量 |
| `recorded_tran_amount` | `customer profile transaction count` 的 TranAmount | 评估已有交易规模 |
| `recorded_avg_daily_orders` | `customer profile order count` 的 Avg | 评估日常下单节奏 |
| `recorded_refund_rate` | 由 `profile order count + profile transaction count` 推算 | 退款/订单 比值，越高代表服务需求越急 |

上述字段**缺失时按 0 处理**（不构成负面证据）；本评分维度仅对 leadQuality 起增益修正，不影响 intentHeat 与 followPriority 查表。

### 信息完整度（infoCompleteness，脚本计算）

**不由 AI 判断**，由 `scripts/info_completeness.py` 按五层基准加权计算：

```
python scripts/info_completeness.py --input profile.json
# → {"score": 65.5, "level": "中", "layers": {...}}
```

| 层 | 字段数 | 权重 | 字段 |
|----|--------|------|------|
| 问卷 | 5 | 35% | DS经验/周广告预算/月订单量预期/细分市场/意向服务 |
| 联系方式 | 1 | 15% | 是否有联系方式 |
| 基础档案 | 6 | 25% | 公司/国家/注册区域/注册时间/最近登录/用户等级 |
| 渠道归因 | 6 | 15% | 渠道来源/二级/Medium/Campaign/合作来源代码/链接 |
| 绑店 | 1 | 10% | 绑定店铺 |

分档：低 [0, 40) / 中 [40, 70) / 高 [70, 100]

## 评估规则

三项评分为 **AI 综合判断**：以下信号权重与分档为参考基线，结论必须能溯源到具体输入字段，不得引入字段之外的信息。

### leadQuality（线索质量：高 / 中 / 低）

语义：该客户若转化，**业务价值与配合度**有多大。信号（权重降序）：

| 信号 | 强 | 中 | 弱/无 |
|------|-----|-----|-------|
| 绑定店铺 | 非空（+2） | — | 空（0） |
| DS 经验 | 多于 1 年（+2） | 6 个月–1 年（+1） | 新手/未填（0） |
| 周广告预算 | 1000+（+2） | 500–1000（+1） | 更低/未填（0） |
| 月订单量预期 | 500+（+2） | 100–500（+1） | 更低/未填（0） |
| 意向服务 | 含 1 且多选（+2） | 含 1 或含 2（+1） | 未填（0） |
| 客户标签 | 3 潜在用户（+1） | 0/1/2（0） | 4 流失用户（−1） |
| 已录购买深度（可选） | recorded_order_total≥500 且 refund_rate<5%（+2） | recorded_order_total 50–500 或 refund_rate 5–15%（+1） | 其余/未填（0） |

分档参考：**≥6 高**（成熟卖家画像，值得重点投入）；**3–5 中**（有潜力但证据不全）；**≤2 低**（新手或信号缺失为主）。

### intentHeat（意向热度：热 / 温 / 冷）

语义：**当下推进的紧迫度**。信号：

| 信号 | 计分 |
|------|------|
| 最近登录 | 7 天内 +3；30 天内 +2；90 天内 +1；更久/未填 0 |
| 绑定店铺 | 非空 +2（已迈出转化动作） |
| 意向服务 | 非空 +1 |
| 渠道归因 | 渠道来源非 `--` +1 |
| 客户标签 | 4 流失用户 −2 |

分档参考：**≥5 热**（正活跃，跟进窗口就在眼前）；**2–4 温**（有兴趣但不紧迫）；**≤1 冷**（沉默或流失倾向，需低成本触达）。

### followPriority（跟进优先级：P0 / P1 / P2 / P3）

**不由 AI 判断，由脚本查表得出**（口径与平台侧 FollowPriorityCalculator 完全一致）。在 leadQuality / intentHeat 定档后，调用本包脚本：

```
python scripts/follow_priority.py --lead-quality <高|中|低> --intent-heat <热|温|冷>
# 输出：P0-P3（加 --json 可得完整三元组）
```

交叉表如下，仅作展示与人工核对，不得跳过脚本自行查表填值：

| leadQuality \\ intentHeat | 热 | 温 | 冷 |
|--------------------------|-----|-----|-----|
| 高 | **P0** 立即跟进（当日） | **P1** 本周跟进 | **P2** 常规池 |
| 中 | **P1** 本周跟进 | **P2** 常规池 | **P3** 低频维护 |
| 低 | **P2** 常规池 | **P3** 低频维护 | **P3** 低频维护 |

## 输出文本写作规则

- **线索速读（reading，≤300 字）**：一段概括——客户身份（国家/区域 + 经验档 + **infoCompleteness 分档与口径**）→ 意向概况 → 所处转化阶段 → 最值得利用的 1–2 个信号 → 最主要的 1 个信息缺口。须体现评分的关键依据。**可出现"信息完整度：高/中/低"等定性表述，但不得出现具体百分比数字**。
- **转化策略（strategy，≤300 字）**：2–3 条可执行动作（切入点/沟通时机/利用哪个信号）；**必须与转化阶段一致**（见下）；**不得在策略中给出 infoCompleteness 数值**。
- **销售话术（script，2-3 条）**：每条**独立可用**、角度互补，覆盖不同跟进时机或切入点（如：①价值切入——平台能力与客户画像匹配点；②信号破冰——用问卷自报的经验/预算/品类打开话题；③行动催促——低成本明确下一步）；每条结尾给一个**低门槛明确行动指令**（绑店咨询 / 首单选品协助）。单条 ≤200 字、合计 ≤500 字。**不得出现 infoCompleteness 数值**。
- **条数下限**：至少 2 条；客户信息极少时写通用角度（不得虚构信号凑个性化），不得只给 1 条，也不得为凑第 3 条重复同一切入点。

### 策略按转化阶段硬分流（不得混淆）

- **未绑店 → 促绑店**：策略与话术围绕绑定店铺的价值（选品、订单、物流一体化管理），不涉及下单细节。
- **已绑店未首单 → 促首单**：策略与话术围绕降低首单门槛（选品建议、首单流程引导、小额试单），默认以绑店动作为信任基础；可选引入"已录购买深度"作为背书信号（仅在客户提供时可用）。

## 网页形态输出管线（硬规则）

要求输出网页/HTML 时，AI **只产槽位 JSON，不碰模板**：

```
hypersku-cli customer detail <customerId>   → 脱敏信号（+ 可选 profile 聚合）
        ↓
AI 按「评估规则」计算 leadQuality / intentHeat，调用 scripts/info_completeness.py 算分档
        ↓
产出槽位 JSON（下方契约，含 stage / leadQuality / intentHeat / infoCompletenessLevel）
        ↓
python scripts/follow_priority.py --lead-quality <高|中|低> --intent-heat <热|温|冷>
        ↓ 把脚本输出的 followPriority 写回槽位 JSON
python scripts/render_profile.py --data profile.json -o out.html
        ↓
模板 templates/profile.html + JSON → 固定版式 HTML（枚举色值、布局全由模板决定）
```

1. **职责边界**：布局、样式、枚举→色值映射固定在只读模板 `templates/profile.html`；AI 的产出物只有槽位 JSON。禁止在上下文中现编 HTML、禁止改写模板内容。
2. **先落盘再调用**：槽位 JSON 先写为文件（如 `profile.json`），再执行渲染脚本；缺 `--data` / 校验失败时修正 JSON 后重跑，不得改脚本或模板来绕过校验。
3. **信任脚本校验**：脚本对必填字段、枚举值域、字数上限、话术条数做全量断言，任何一项违规即退出码 2 拒绝渲染——出现校验失败说明 JSON 不合契约，回头改 JSON。
4. **安全**：所有槽位值经 HTML 转义后注入，话术文本不支持任何 HTML 标签。
5. **followPriority 由脚本计算**：两档评分定档后必须调用 `scripts/follow_priority.py` 取值，把脚本输出写入槽位 JSON；禁止自行查表或直接填写。渲染脚本会再次校验一致性，与交叉表不符的 followPriority 将被拒绝渲染（退出码 2）。
6. **infoCompleteness 由脚本计算**：按五层基准加权得出，必须调用 `scripts/info_completeness.py` 取值，把 `level` 字段写入槽位 JSON；禁止自行打分。

## 输出内容契约

无论何种形态，内容必须包含（可按调用方要求组织为 JSON / markdown / 网页等；网页形态以槽位 JSON 形式交付给渲染脚本）：

1. `customerId`（客户标识）
2. 转化阶段标签：**未绑店** / **已绑店未首单**
3. `leadQuality` ∈ {高, 中, 低}
4. `intentHeat` ∈ {热, 温, 冷}
5. `followPriority` ∈ {P0, P1, P2, P3}（**由 scripts/follow_priority.py 查表得出**，必须与 leadQuality × intentHeat 交叉表一致）
6. `infoCompletenessLevel` ∈ {低, 中, 高}（**由 scripts/info_completeness.py 五层加权得出**，仅输出分档，不输出百分比）
7. `reading`（≤300 字，可引用"信息完整度：高/中/低"等定性词，**不得出现具体百分比数字**）
8. `strategy`（≤300 字，与阶段一致；**不得出现 infoCompleteness 数值**）
9. `script`（2–3 条，每条独立可用、角度互补，与阶段一致；**不得出现 infoCompleteness 数值**）

## 硬约束（红线）

1. **禁止虚构意向信号**：所有判断须能溯源到输入字段；字段未填只能表述"未提供/数据有限"，不得脑补动机。
2. **禁止硬承诺**：不得出现具体价格、收益数字、时效承诺（如"月入 X 美元""7 天出单"）。
3. **枚举严格**：三项评分与 infoCompletenessLevel 只取规定枚举值，不得输出"较高/P0.5/52%"等变体；followPriority 必须来自脚本计算，不得手填偏离交叉表的值。
4. **阶段一致**：未绑店客户不得输出促首单话术，反之亦然。
5. **字数与条数**：reading ≤300、strategy ≤300、script 2–3 条且单条 ≤200、合计 ≤500（中文字符），超限必须收敛；条数不足 2 不得交付，凑不出第 3 条时给 2 条即可。
6. **infoCompleteness 表达边界**：
   - `reading` 中可使用"信息完整度为 高/中/低"等定性描述；
   - 禁止在 **任何**输出字段中出现具体百分比数字（如"65.5%"）、分数（如"score: 65.5"）；
   - 不得输出层明细（layers 对象的 ratio / weighted 等内部结构）；
   - 数值由调用方系统持有，本 skill 仅按需引用分档。
7. **已出池即止**：是否有已支付订单＝是时拒绝分析，不得输出画像。

## 冷启动处理

刚注册、问卷全空的客户**照常分析**：评分按信号缺失处理（通常 leadQuality 低、intentHeat 依赖最近登录），infoCompleteness 必然为"低"，`reading` 中声明"**数据有限，结论仅供参考**"；不设门槛、不拒答、不输出降级版画像。
