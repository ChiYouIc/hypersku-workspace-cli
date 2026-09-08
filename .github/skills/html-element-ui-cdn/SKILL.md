---
name: html-element-ui-cdn
description: '非报告类 HTML（交互页面、动态模板、内嵌工具原型）直接引用 element-ui CDN + Vue 2 实现数据驱动渲染。Use when 生成带交互的 HTML 原型/工具页（筛选、表单、表格分页、弹窗、tab 切换），或需要运行时注入数据（fetch/postMessage/URL 参数）的动态 HTML 模板。报告/画像等离线分发单文件不适用（走 html-element-ui-style 纯 CSS 复刻）。'
---

# HTML 交互/动态模板：Element UI CDN 方案

**非报告内容**的 HTML 输出用本 skill：直接引 element-ui 2.x CDN + Vue 2，用真组件渲染、数据驱动，不手写复刻 CSS。

## 与 html-element-ui-style 的分工（硬边界）

| 场景 | 用哪个 | 原因 |
|------|--------|------|
| 报告/画像/导出页（离线分发、存档、邮件附件） | `html-element-ui-style`（纯 CSS 复刻） | 零依赖可离线打开，产物自包含 |
| 交互页面/动态模板/工具原型（本 skill） | CDN 引组件 | 需要真组件行为（表单校验、弹窗、分页），数据可运行时注入 |

判断标准一句话：**打开时有没有网、要不要交互**——离线/存档走静态复刻，联网/交互走 CDN。

## When to Use

- 带交互的 HTML 原型：筛选、搜索、tab 切换、弹窗确认
- 需要丰富组件的工具页：el-table 分页排序、el-form 校验、el-date-picker、el-select 级联
- 运行时注入数据的动态模板：fetch 接口 / postMessage / URL 参数 → 页面不重新生成
- 内嵌 webview 页面（VS Code webview、管理后台静态壳）

**不适用**：客户画像报告等离线单文件产物；Vue 工程内开发（那是引依赖包，不是 CDN）。

## Procedure

1. **骨架与 CDN 引入**（顺序与版本锁定是硬规则）：

   ```html
   <link rel="stylesheet" href="https://unpkg.com/element-ui@2.15.14/lib/theme-chalk/index.css">
   <script src="https://unpkg.com/vue@2.7.16/dist/vue.min.js"></script>
   <script src="https://unpkg.com/element-ui@2.15.14/lib/index.js"></script>
   ```

   - **Vue 必须先于 element-ui 加载**（element-ui 是 Vue 2 插件，Vue 3 不兼容）
   - **版本锁定**：写死 `vue@2.7.16` + `element-ui@2.15.14`，禁止 `@latest`（element-ui 已停止维护，latest 无收益只有风险）
   - 内网/私有部署场景把 unpkg 换成镜像源，版本号同样锁定

2. **数据驱动**：业务数据收敛到 `window.XXX` 全局对象（或 created 里 fetch），模板区禁止硬编码业务内容；数据结构尽量与既有槽位 JSON 同构（如画像的 `PROFILE`），便于同一份数据两种形态复用。

3. **挂载边界**：Vue 实例只挂载 `#app` 容器；`<title>` 等 head 内容**不在 Vue 域内**，`{{}}` 不会被替换——标题要么写死，要么 JS 里 `document.title = ...` 赋值。

4. **枚举 → 组件属性映射**：业务枚举映射到 Element 的 type 体系并固定（如 el-tag：`success / warning / danger / info`，默认态是空串不是 `"primary"`）；映射写在 JS 常量里（等价于静态方案的 `CLASS_MAP`），不随内容漂移。

5. **组件选型参考**（动态场景常用的）：
   - 明细展示：`el-descriptions`（border + size=small）
   - 列表分页：`el-table` + `el-pagination`
   - 时间轴/步骤：`el-timeline` / `el-steps`
   - 反馈：`el-alert`（show-icon）/ `el-message` / `el-dialog`
   - 布局：`el-row`/`el-col` 栅格，页面底色 `#F2F3F5`

6. **自检**：断网打开应给出可降级提示或至少不白屏报错（可在数据加载失败时 el-alert 提示）；有交互时逐个验证事件；数据注入路径（fetch/URL 参数）实际跑一遍。

7. **场景自检**：本 skill 的产物必须有交互或运行时注数需求——纯展示内容是 `html-element-ui-style` 的领地，引 CDN 即违反路由规则。

## Quality Checklist

- [ ] CDN 版本锁定且 Vue 先于 element-ui 加载
- [ ] 业务数据与模板分离（window 全局或 fetch），模板区无硬编码内容
- [ ] Vue 只挂载 `#app`；`<title>` 未使用 `{{}}` 占位
- [ ] 枚举 → type 映射为固定 JS 常量
- [ ] 数据加载失败有 el-alert 降级提示
- [ ] 交互（筛选/分页/弹窗）逐项验证通过

## CDN 组合可行性（已验证）

`vue@2.7.16` + `element-ui@2.15.14`（unpkg CDN，`file://` 直开）实测可用：el-tag / el-descriptions / el-timeline / el-alert 渲染正常，数据驱动与枚举映射生效。验证结论即下方规则，无需示例文件。
