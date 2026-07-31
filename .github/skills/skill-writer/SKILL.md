---
name: skill-writer
description: '编写 Hypersku CLI 能力 Skill。当用户新增了 cmd 包下的命令，需要为其生成对应的 SKILL.md 及 references 参考文档时使用。适用场景：新增子命令后编写 skill、完善已有 skill、批量生成 skills。'
version: 1.0.0
author: owen
tags:
  - skill
  - 文档生成
  - cli
---

# Skill 编写器

根据 `cmd/` 包下的 cobra 命令定义，自动生成结构化的 SKILL.md 及其 references 参考文件到 `skills/<name>/` 目录。

## 何时使用

- 用户在 `cmd/` 下新增了子命令，需要配套编写 skill 文档
- 用户要求"为这个命令写个 skill"、"完善 SKILL.md"、"生成 skill 文档"
- 批量扫描 `cmd/` 目录，为缺少 skill 的命令补全文档

## 工作流程

### 第一步：探索命令定义

1. 读取 `cmd/` 目录，列出所有命令文件（排除 `root.go`）
2. 对每个命令文件，提取以下信息：
   - 命令变量名（如 `purchaseCmd`）
   - `Use` 字段（命令名称，如 `purchase`）
   - `Short` / `Long` 字段（用途描述）
   - `init()` 中注册的子命令列表（`AddCommand` 调用）
3. 对每个子命令，提取：
   - 子命令变量名（如 `getOrderInfoCmd`）
   - `Use` 字段（子命令名称 + 参数占位符）
   - `Short` / `Long` 字段
   - `Args` 约束（必填参数）
   - `Run` 函数中的输出格式（`fmt.Sprintf` 模板、字段列表）

### 第二步：分析数据结构

1. 找到命令文件 `Run` 中调用的 API 方法（如 `apis.NewPurchaseApi().GetOrderInfo(...)`）
2. 读取对应的 `internal/apis/` 文件，提取：
   - 请求参数类型及字段
   - 响应结构体及字段
   - 枚举映射表（如 `orderType`、`orderStatus` map）
3. 如有枚举映射，判断是否需要纳入 SKILL.md（仅在用户需要对照参考时保留，默认移除）

### 第三步：确定 skill 名称和目录

1. skill 名称取命令的 `Use` 值（如 `purchase`）
2. 目标目录：`skills/<name>/`
3. 如果目录已存在，以完善模式进行，询问用户需要更新哪些部分

### 第四步：编写 SKILL.md

按以下结构编写，参考 [purchase SKILL.md](../purchase/SKILL.md) 作为模板：

```markdown
---
name: <skill-name>
description: <一句话能力概括 + 触发关键词>
version: 1.0.0
author: owen
tags:
  - <标签1>
  - <标签2>
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

#### description 编写要点

- 必须包含"当用户需要..."的触发句式
- 关键字覆盖所有子命令的能力域（如"商品/物流/地址/日志"）
- 不超过 1024 字符

#### 能力总览编写要点

- 每行一个子命令，参考文件指向 `references/` 下的同名 `.md`

#### 意图判断编写要点

- 每个子命令用一条自然语句描述：`用户提到"关键词1/关键词2/..."时，执行 \`子命令\` 展示xxx。`
- 从子命令的 `Short`/`Long` 描述以及 `Run` 中的输出字段反推用户问法关键词
- 覆盖常见口语化表达（如"买了什么"→订单详情、"到哪了"→物流）
- 末条兜底：`若未提供 <必填参数>，提示用户提供xxx。`

### 第五步：编写 references 参考文件

为每个子命令在 `references/` 下创建对应的 `.md` 文件，结构如下：

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

- 从 `Run` 函数的 `fmt.Sprintf` / `cmd.Print` 中还原输出模板
- 用中文模拟数据填充，覆盖正常场景和多条记录场景
- 对于列表型输出（如物流轨迹、商品明细），展示 2~3 条记录

#### 字段说明生成规则

- 从响应结构体的 json tag 和上下文推断字段含义
- 中文命名，简洁描述

### 第六步：自检

完成后检查以下项目：

- [ ] frontmatter 中 `name` 与文件夹名一致
- [ ] `description` 包含触发关键词
- [ ] 能力总览表每行对应一个子命令
- [ ] 参考文件路径使用相对路径 `references/xxx.md`
- [ ] 每个子命令都有对应的 reference 文件
- [ ] SKILL.md 正文不超过 500 行（保持渐进加载友好）
- [ ] 无孤立无用的对照表（如不需要的枚举表已移除）

## 约束

1. **命名一致**：skill 文件夹名、`name` 字段、CLI 命令的 `Use` 值三者必需一致。
2. **引用路径**：始终使用 `./` 开头的相对路径引用 references。
3. **渐进加载**：输出示例等详细内容放 references，SKILL.md 只保留概要和决策逻辑。
4. **中文优先**：面向中文用户，正文使用中文，代码/命令使用英文。
5. **基于代码**：所有输出格式必须与 `cmd/*.go` 中的实际代码一致，不得凭空编造。
