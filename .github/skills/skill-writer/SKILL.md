---
name: skill-writer
display_name: Skill 编写器
display_name_en: Skill Writer
description_zh: 编写 Hypersku CLI 技能包：为 cmd 包下的命令生成能力域 SKILL.md（含双语展示名/描述等必填 frontmatter 字段）及 references 参考文档。
description_en: Write Hypersku CLI skill packages: generate per-domain SKILL.md with required bilingual frontmatter fields and reference docs for commands under the cmd package.
description: '编写 Hypersku CLI 技能文档。当用户新增了 cmd 包下的命令，需要为其生成对应的能力域技能包（skills/<name>/SKILL.md 及 references/ 参考文档）时使用。适用场景：新增子命令后编写文档、完善已有文档、批量生成文档。'
category: developer-tools
version: 2.1.1
author: owen
---

# Skill 编写器

根据 `cmd/` 包下的 cobra 命令定义，编写结构化技能包，并在对应能力域目录 `skills/<name>/` 下维护 `SKILL.md` 与 `references/` 参考文件。

技能包标准采用 workbuddy 技能规范（https://open.workbuddy.cn/docs/skill）。目录与文件命名风格（`{skill-name}/SKILL.md` + `references/` 等）固定不变，后续扩展其它平台时沿用同一套结构。

## 技能包规范（硬约束）

### 目录结构（仅两级）

```
skills/
└── {skill-name}/
    ├── SKILL.md              # ★ 技能定义（必须）
    ├── references/           # 参考资料（可选）
    │   ├── api-spec.md       #   API 规范、字段类型等
    │   └── examples.md       #   示例数据
    ├── scripts/              # 可执行脚本（可选）
    │   └── fetch-data.js     #   数据获取脚本
    └── templates/            # 模板文件（可选）
        └── report.sh         #   报告生成模板
```

1. **一个技能包 = 一个 `{skill-name}/` 目录**：包内只允许 `SKILL.md` + `references/`、`scripts/`、`templates/` 这一级子目录，**不得出现更深层嵌套**（如 `references/` 内再建子目录）。
2. **SKILL.md 必须存在**：每个技能包根目录必须有 `SKILL.md`。
3. **上传打包**（`scripts/upload_skill.ps1`）：每个能力域目录单独打 zip，**zip 根 = SKILL.md + references/**，不带外层 `{skill-name}/` 文件夹。

### SKILL.md frontmatter 字段（YAML）

| 字段 | 必填 | 说明 |
|------|------|------|
| `name` | 否 | 技能标识（与目录名一致，kebab-case） |
| `display_name` | **是** | 中文展示名（市场列表展示用） |
| `display_name_en` | **是** | 英文展示名 |
| `description` | **是** | 一句话描述技能能力，写清用途和触发词 |
| `description_zh` | **是** | 简短中文介绍 |
| `description_en` | **是** | 简短英文介绍 |
| `category` | 否 | 平台分类之一（如 `developer-tools`、`writing`） |
| `allowed-tools` | 否 | 工具白名单（逗号分隔） |
| `version` | **是** | 版本号（语义化，如 `2.0.0`） |
| `disable-model-invocation` | 否 | `true` 则 AI 不会自动触发，只能用户手动调用 |
| `user-invocable` | 否 | `false` 则隐藏菜单，仅供 AI 内部使用 |
| `author` | **是** | 合作方名称 |

**必填字段缺失会导致平台上传解析失败**（如"缺少 Skill 中文展示名（display_name）"）。

### frontmatter 模板

```yaml
---
name: <kebab-case-技能标识，与目录名一致>
display_name: <中文展示名>
display_name_en: <English Display Name>
description: <一句话描述用途与触发词，含用户会提到的关键词>
description_zh: <简短中文介绍：能力 + 主要命令/场景>
description_en: <A brief English introduction: capability + main commands/scenarios>
category: developer-tools
version: <语义化版本号>
author: owen
---
```

### references 引用方式

参考资料在 SKILL.md 中通过相对链接 `[xxx](references/xxx.md)` 引用，AI 执行技能时读取为上下文。

## 何时使用

- 用户在 `cmd/` 下新增了子命令，需要配套编写 skill 文档
- 用户要求"为这个命令写个 skill"、"完善引导文档/SKILL.md"、"生成 skill 文档"
- 批量扫描 `cmd/` 目录，为缺少技能包的命令补全文档

## 工作流程

### 第一步：探索命令定义

1. 读取 `cmd/` 目录，列出所有命令文件（排除 `root.go`）
2. 对每个命令文件，提取以下信息：
   - 命令变量名（如 `purchaseCmd`）
   - `Use` 字段（命令名称，如 `purchase`）
   - `Short` / `Long` 字段（用途描述）
   - `init()` 中注册的子命令列表（`AddCommand` 调用）
3. 对每个子命令，提取以下信息：
   - 子命令变量名（如 `getOrderInfoCmd`）
   - `Use` 字段（子命令名称 + 参数占位符）
   - `Short` / `Long` 字段
   - `Args` 约束（必填参数）
   - `Run` 函数中的输出格式（`fmt.Sprintf` 模板、字段列表）

### 第二步：分析数据结构

1. 找到命令文件 `Run` 中调用的 API 方法（如 `apis.NewPurchaseApi().GetOrderInfo(...)`）
2. 读取对应的 `internal/apis/` 文件，提取以下信息：
   - 请求参数类型及字段
   - 响应结构体及字段
   - 枚举映射表（如 `orderType`、`orderStatus` map）
3. 如有枚举映射，判断是否需要纳入 SKILL.md（仅在用户需要对照参考时保留，默认移除）

### 第三步：确定 skill 名称和目录

1. skill 名称取命令的 `Use` 值（如 `purchase`），kebab-case
2. 目标目录：`skills/<name>/`（一个独立技能包）
3. 如果目录已存在，以完善模式进行，询问用户需要更新哪些部分

### 第四步：编写能力域技能定义（skills/<name>/SKILL.md）

按以下结构编写，参考 [skills/hypersku-purchase/SKILL.md](../../../skills/hypersku-purchase/SKILL.md) 作为模板：

```markdown
---
<按上文 frontmatter 模板填写，display_name / display_name_en /
 description / description_zh / description_en / version / author 均必填>
---

# <中文标题>

<一句话简介，说明通过哪个 CLI 命令实现什么能力>

## 能力总览

| 能力 | 用途 | CLI 命令 | 参考文件 |
|------|------|----------|----------|
| ... | ... | `hypersku-cli <cmd> <subcmd> <args>` | [ref.md](references/ref.md) |

## 意图判断

当用户输入包含以下关键词或意图时，使用对应的子命令：

<以关键词 → 子命令映射，每条用一句话描述>

例如：用户提到"订单详情/订单信息/商品/买了什么"时，执行 `info <orderId>` 展示订单概要及商品明细。

## 注意事项

<列出参数约束、API 依赖、特殊场景处理等>
```

#### frontmatter 编写要点

- `display_name`：能力域的简洁中文名（如「采购订单管理」）
- `display_name_en`：对应英文名（如 `Purchase Order Management`）
- `description`：一句话写清**用途 + 触发词**，覆盖该能力域所有子命令的关键词（与意图判断呼应）
- `description_zh` / `description_en`：中英文简短介绍，内容对应
- `version`：修改文档时递增（小改 patch、新增子命令 minor、结构变更 major）

#### 能力总览编写要点

- 每行一个子命令，参考文件指向 `references/` 下的同名 `.md`

#### 意图判断编写要点

- 每个子命令用一条自然语句描述：`用户提到"关键词1/关键词2/..."时，执行 \`子命令\` 展示xxx。`
- 从子命令的 `Short`/`Long` 描述以及 `Run` 中的输出字段反推用户问法关键词
- 覆盖常见口语化表达（如"买了什么"→订单详情、"到哪了"→物流）
- 末条兜底：`若未提供 <必填参数>，提示用户提供xxx。`

### 第五步：编写 references 参考文件

为每个子命令在 `references/` 下创建对应的 `.md` 文件（**flat 结构，不再嵌套子目录**），结构如下：

```markdown
# <子命令名> 输出参考

## 命令

`hypersku-cli <cmd> <subcmd> <args>`

## 输出示例

<根据 Run 函数中的 fmt.Sprintf 模板，填充模拟数据>

## 字段说明

| 字段 | 说明 |
|------|------|
| ... | ... |
```

#### 输出示例生成规则

- 从 `Run` 函数中的 `fmt.Sprintf` / `cmd.Print` 还原输出模板
- 用中文模拟数据填充，覆盖正常场景和多条记录场景
- 对于列表型输出（如物流轨迹、商品明细），展示 2~3 条记录

#### 字段说明生成规则

- 从响应结构体的 json tag 和上下文推断字段含义
- 中文命名，简洁描述

### 第六步：自检

完成后检查以下项目：

- [ ] 目录为两级结构：`skills/<name>/SKILL.md` + `skills/<name>/references/*.md`，无更深嵌套
- [ ] frontmatter 必填字段齐全：`display_name`、`display_name_en`、`description`、`description_zh`、`description_en`、`version`、`author`
- [ ] `name` 与目录名一致（kebab-case）
- [ ] 能力总览表每行对应一个子命令
- [ ] 参考文件路径使用相对路径 `references/xxx.md`
- [ ] 每个子命令都有对应的 reference 文件
- [ ] SKILL.md 正文不超过 500 行（保持渐进加载友好）
- [ ] 无孤立无用的对照表（如不需要的枚举表已移除）
- [ ] 跨能力域引用改为文字说明（不同技能包独立上传，相对链接跨包会断链）

## 约束

1. **命名一致**：技能包文件夹名与 CLI 命令的 `Use` 值一致；frontmatter 的 `name` 同名。
2. **包内自洽**：一个技能包独立上传、独立生效，**不依赖其他技能包的文件**；跨域能力用文字提示（如"参见 `customer` 技能包"），不使用跨包相对链接。
3. **引用路径**：包内引用统一用相对路径 `references/xxx.md`、`../SKILL.md`（references 内回指本包 SKILL.md）。
4. **渐进加载**：输出示例等详细内容放 references，SKILL.md 只保留概要和决策逻辑。
5. **中文优先**：面向中文用户，正文使用中文，代码/命令/英文字段使用英文。
6. **基于代码**：所有输出格式必须与 `cmd/*.go` 中的实际代码一致，不得凭空编造。
