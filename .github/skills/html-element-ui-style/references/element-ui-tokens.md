# Element UI 2.x 设计令牌与组件 CSS 规格

手写复刻用速查。所有值取自 element-ui 2.x 官方 `packages/theme-chalk` 源码，禁止近似替代。

## 色板（:root 变量块，直接复制）

```css
:root {
  /* 语义主色 */
  --primary: #409EFF;
  --success: #67C23A;
  --warning: #E6A23C;
  --danger:  #F56C6C;
  --info:    #909399;

  /* 文字 */
  --text-primary:   #303133;  /* 主文字 */
  --text-regular:   #606266;  /* 常规文字 */
  --text-secondary: #909399;  /* 次要文字 */

  /* 边框 */
  --border-base:    #DCDFE6;
  --border-light:   #E4E7ED;
  --border-lighter: #EBEEF5;

  /* 填充 */
  --fill-blank: #FFFFFF;
  --fill-page:  #F2F3F5;      /* 页面底 */

  /* 阴影（el-card 标准） */
  --box-shadow: 0 2px 12px 0 rgba(0, 0, 0, .1);
}
```

## 字体与排版

```css
font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
             "Microsoft YaHei", Arial, sans-serif;
font-size: 14px;        /* 基准 */
font-weight: 400;
line-height: 1.8;       /* 长文阅读 */
color: #303133;         /* 标题等重文字 */
```

## 组件 CSS 规格

### el-card

```css
.el-card {
  background-color: #FFFFFF;
  border: 1px solid #EBEEF5;
  border-radius: 4px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, .1);
  overflow: hidden;
}
.el-card__header {
  padding: 15px 20px;
  border-bottom: 1px solid #EBEEF5;   /* 与卡体分隔 */
}
.el-card__body { padding: 20px; }
```

### el-tag（同色系：10% 底 + 20% 边框）

```css
.el-tag {
  display: inline-block;
  height: 32px; line-height: 30px;    /* 留 1px 边框 */
  padding: 0 10px;
  font-size: 12px;
  border: 1px solid;
  border-radius: 4px;
  white-space: nowrap;
}
.el-tag--small { height: 24px; line-height: 22px; padding: 0 8px; }

/* 五档语义（以 primary 为例，其余换色值） */
.el-tag--primary {
  background: rgba(64, 158, 255, .1);
  border-color: rgba(64, 158, 255, .2);
  color: #409EFF;
}
/* success #67C23A / warning #E6A23C / danger #F56C6C / info #909399 同构 */
```

### el-alert（info 型）

```css
.el-alert {
  display: flex; align-items: center; gap: 8px;
  padding: 9px 16px;
  border-radius: 4px;
  font-size: 13px;
}
.el-alert--info { background-color: #F4F4F5; color: #909399; }
```

### el-button（默认 + primary）

```css
.el-button {
  display: inline-block;
  line-height: 1;
  padding: 12px 20px;
  font-size: 14px;
  border-radius: 4px;
  border: 1px solid #DCDFE6;
  background: #FFFFFF;
  color: #606266;
  cursor: pointer;
}
.el-button--primary {
  background: #409EFF; border-color: #409EFF; color: #FFFFFF;
}
.el-button--small { padding: 9px 15px; font-size: 12px; border-radius: 3px; }
```

### el-divider（分隔线）

```css
.el-divider {
  display: block;
  height: 1px;
  background-color: #EBEEF5;
  margin: 12px 0;     /* 按需调 */
}
```

### el-descriptions（描述列表，报告类常用）

```css
.el-descriptions__body {
  background: #FFFFFF;
}
.el-descriptions .item-label {
  font-weight: 700;
  color: #909399;     /* label 灰 */
  margin-right: 12px;
}
.el-descriptions .item-content { color: #606266; }
```

## 枚举 → 语义色映射惯例

| 业务语义 | class 后缀 | 色值 |
|----------|-----------|------|
| 好/高/成功 | `--success` | #67C23A |
| 中/警告 | `--warning` | #E6A23C |
| 差/低/危险 | `--danger` | #F56C6C |
| 中性/禁用 | `--info` | #909399 |
| 强调/品牌 | `--primary` | #409EFF |

映射固定在模板或渲染脚本（如 `CLASS_MAP`），不随内容漂移；优先级档（P0 最急）可用 danger、P1 warning、P2 primary、P3 info。
