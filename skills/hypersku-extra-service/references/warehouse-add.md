# warehouse-add 输出参考

## 命令

`hypersku-cli extra-service warehouse-add <repertoryId>`

`repertoryId` 为仓库 ID（可从 `warehouse list` 命令获取）。

## 输出示例

```
仓库 101 开通的增值服务:

  1. 质检贴标
  2. 换包装
  3. 代贴条码
```

## 字段说明

| 字段 | 说明 |
|------|------|
| 仓库 | 查询的仓库 ID |
| 增值服务列表 | 该仓库已开通的增值服务名称（有序号） |

未开通任何服务时输出 `仓库 <N> 暂无增值服务`。
