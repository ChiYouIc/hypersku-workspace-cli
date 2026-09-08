---
name: html-element-ui-style
description: '输出独立 HTML 页面/模板/报告时统一使用 Element UI（element-ui 2.x）视觉风格。Use when 生成静态 HTML、邮件/报告/画像等单文件页面、HTML 模板改版或手写原生 CSS 版式时，需套用 el-card / el-tag / el-alert / el-button 等组件视觉规格与 Element 设计令牌。不适用于 Vue 工程内直接使用 Element 组件库的场景（那类直接引组件，无需本 skill）。'
---

# HTML 输出统一 Element UI 风格

工程中任何需要**输出独立 HTML**（单文件页面、报告、模板、原型）的场景，视觉层一律按 Element UI 2.x 规格手写实现——不引入 Vue/组件库运行时，纯 CSS 复刻组件视觉。

## When to Use

- 生成客户画像报告、数据报表、导出页等单文件 HTML
- 新建或改版 HTML 模板（如 `skills/*/templates/*.html`）
- 原本自定义样式的 HTML 需要统一到 Element 视觉体系

**不适用**：Vue 工程内直接用 Element 组件（那是引库，不是复刻）；交互页面/动态模板等非报告内容——那类走姊妹技能 `html-element-ui-cdn`（引 CDN 真组件，数据驱动）。判断标准：**打开时有没有网、要不要交互**——离线/存档走本 skill，联网/交互走 CDN skill。

## 核心原则

1. **零依赖单文件**：不引外部 CSS/JS/字体，所有样式内联在 `<style>`，产物双击可开、可离线分发。
2. **设计令牌先行**：`:root` CSS 变量必须与 Element 官方令牌一致（见 [element-ui-tokens.md](references/element-ui-tokens.md)），禁止自创近似色。
3. **组件规格复刻**：卡片、标签、按钮等按 Element 的精确规格（边框/圆角/阴影/尺寸）实现，class 命名沿用 `el-` BEM 风格（`el-card__header`、`el-tag--success`），让产物视觉与 Element 应用无法区分。
4. **语义色即业务色**：枚举值 → 语义色映射固定（成功/警告/危险/信息/主色五档），由模板或渲染脚本决定，AI 不得随意配色。

## Procedure

1. **建骨架**：HTML5 doctype + `lang="zh-CN"` + `<meta charset="UTF-8">` + viewport；页面容器背景 `#F2F3F5`（`--fill-page`），内容区白卡浮层。
2. **注入令牌**：从 [element-ui-tokens.md](references/element-ui-tokens.md) 复制 `:root` 变量块（色板 / 字体栈 / 边框 / 阴影）到 `<style>` 顶部。
3. **按需复刻组件**：查令牌参考的「组件 CSS 规格」章节，用到的组件逐个实现（el-card / el-tag / el-alert / el-button / el-divider / el-descriptions 等）；没用到的组件不写，保持轻量。
4. **布局**：优先 CSS Grid/Flex；统计数值卡用 grid `repeat(auto-fit, minmax(200px, 1fr))` 自适应多列。
5. **枚举配色**：业务枚举（如 高/中/低、P0-P3）映射到语义 class（`--success/--warning/--danger/--info/--primary`），映射关系写在模板或渲染脚本里，输出前可断言。
6. **自检**：浏览器打开产物（或截图）核对——卡片阴影是否 4px 圆角 + `0 2px 12px` 阴影；tag 是否同色系 10% 底 + 20% 边框；正文 14px / 行高 1.8 / 色值 `#606266`。
7. **模板管线场景**（如画像 HTML）：只改模板的 `<style>` 与结构 class，**不动槽位占位符、`<!--SCRIPT_ITEMS-->`、`<template>` 片段**；渲染脚本侧若 class 前缀变化需同步 `CLASS_MAP`。

## Quality Checklist

- [ ] 零外部依赖，单文件可离线打开
- [ ] `:root` 令牌与 Element 官方一致（色值逐项核对）
- [ ] class 命名为 `el-` BEM 风格
- [ ] 枚举 → 语义色映射固定且可断言
- [ ] 正文 14px、`#303133`/`#606266` 文字色、1.8 行高
- [ ] 模板管线产物：占位符零残留（渲染后 `{{` 计数为 0）

## References

- [element-ui-tokens.md](references/element-ui-tokens.md) — 设计令牌（色板/字体/边框/阴影）与组件 CSS 规格速查
