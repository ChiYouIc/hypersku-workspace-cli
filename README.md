# HyperSKU CLI

> 基于 Go 构建的高效命令行工具，提供灵活的第三方业务 API 调用能力，涵盖售后、客户、物流、采购、供应商、第三方平台、AI 分析等场景。

## 目录结构

```
hypersku-cli/
├── .vscode/
│   └── settings.json          # VS Code 配置（Go 路径、国内代理镜像）
├── cmd/                       # cobra 命令定义（薄命令层，只做编排）
│   ├── root.go                # 根命令（全局 flags、配置加载、依赖初始化）
│   ├── after_sales.go         # 售后管理命令
│   ├── ai_analysis.go         # AI 分析查询命令
│   ├── customer.go            # 客户管理命令
│   ├── customer_order.go      # 客户订单子命令（含拦截/退件）
│   ├── customer_profile.go    # 客户画像命令
│   ├── domestic_third_trade_exception.go  # 国内第三方交易异常订单管理命令
│   ├── extra_service.go       # 增值服务查询命令
│   ├── logistics.go           # 物流管理命令
│   ├── purchase.go            # 采购订单管理命令
│   ├── purchase_after_sales.go # 采购售后查询命令
│   ├── supplier.go            # 供应商管理查询命令
│   ├── third_app_api.go       # 第三方平台数据查询命令
│   └── warehouse.go           # 仓库管理命令（含仓库/物流列表）
├── internal/                  # 内部业务逻辑（不对外暴露）
│   ├── apis/                  # 第三方 API 封装
│   │   ├── types.go           # 通用类型定义
│   │   ├── after_sales.go     # 售后 API
│   │   ├── ai_analysis.go     # AI 分析 API
│   │   ├── customer.go        # 客户 API
│   │   ├── customer_order_return.go  # 客户订单退件工单 API
│   │   ├── extra_service.go   # 增值服务 API
│   │   ├── logistics.go       # 物流 API
│   │   ├── purchase.go        # 采购 API
│   │   ├── purchase_after_sales.go  # 采购售后 API
│   │   ├── supplier.go        # 供应商 API
│   │   ├── third_app_api.go   # 第三方平台（1688/Shopify）API
│   │   └── warehouse.go       # 仓库 API
│   ├── config/                # 配置加载（config.json）
│   ├── httpclient/            # HTTP 客户端基础封装
│   └── version/               # 版本信息
├── pkg/                       # 可复用的公开包（预留）
├── scripts/                   # 辅助脚本（打包发布等）
│   └── pack.ps1               # 打包脚本：编译 + 同步 skills
├── skills/                    # Copilot Skill 使用文档
├── build/                     # 编译产物
├── main.go                    # 程序入口
├── go.mod / go.sum            # Go 模块依赖
├── Makefile                   # 构建脚本
└── README.md
```

## 环境要求

| 依赖 | 版本 |
|------|------|
| [Go](https://go.dev/dl/) | >= 1.26 |
| Git（可选） | 任意版本 |

## 快速开始

### 1. 配置 Go 代理（国内网络）

项目已内置 `.vscode/settings.json` 配置了 `https://goproxy.cn` 镜像，如未使用 VS Code，手动设置：

```powershell
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
go mod tidy
```

### 2. 编译

```powershell
# Windows
go build -o build/hypersku-cli.exe .

go build -o "$HOME\.hypersku-cli\hypersku-cli.exe" .

# 或使用 Makefile
make build

# 或使用 Powershell
powershell -ExecutionPolicy Bypass -File scripts\pack.ps1

# 上传 skill
powershell -ExecutionPolicy Bypass -File scripts\upload_skill.ps1
```

### 3. 运行

```powershell
# 查看帮助
build\hypersku-cli.exe --help

# 显示版本
build\hypersku-cli.exe --version

# 示例：查询采购订单详情
build\hypersku-cli.exe purchase info 123456
```

## 配置文件（config.json）

配置文件默认位于 `~/.hypersku-cli/config.json`，也可通过全局 flag `--config` 指定其他路径：

```json
{
  "api_base_url": "https://api.example.com",
  "api_timeout": 30,
  "api_token": "your-token-here"
}
```

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `api_base_url` | string | `""` | API 基础地址，为空时不覆盖默认地址 |
| `api_timeout` | int | `30` | API 请求超时时间（秒），小于等于 0 时使用默认值 |
| `api_token` | string | `""` | API 认证令牌，自动作为 `authorization` 请求头发送；`auth status` 亦基于此 token 远程校验登录态 |

> 配置文件不存在时自动使用默认配置；启动时会自动创建 `~/.hypersku-cli` 目录（含 `logs/`、`data/` 子目录）。

## 使用指南

### 全局 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--config` | string | `""` | 配置文件路径 |
| `-v, --version` | bool | `false` | 显示版本信息 |
| `-h, --help` | bool | `false` | 显示帮助信息 |

### 认证管理 `auth`

| 命令格式 | 说明 |
|----------|------|
| `auth login <access-token>` | 使用已有 token 登录 |
| `auth login --api-user <用户名> --api-password <密码>` | 用户名密码登录 |
| `auth status` | 查看当前登录状态 |
| `auth logout` | 退出登录（删除本地 token） |

### 售后管理 `after-sales`

| 命令格式 | 说明 |
|----------|------|
| `after-sales 1688 <thirdOrderId>` | 查询 1688 售后工单 |
| `after-sales 1688 goods <thirdOrderId> <refundId>` | 查询 1688 售后商品 |
| `after-sales 1688 detail <refundId>` | 查询 1688 售后详情 |
| `after-sales 1688 message <refundId>` | 查询 1688 售后留言 |

### 客户管理 `customer`

| 命令格式 | 说明 |
|----------|------|
| `customer detail <customerId>` | 查询客户档案 |
| `customer order info <orderId>` | 查询客户订单信息 |
| `customer order logistics <orderId>` | 查询订单物流信息 |
| `customer order address <orderId>` | 查询订单地址信息 |
| `customer order return <customerOrderId>` | 查询客户订单退件工单（按客户订单号） |
| `customer order intercept <tradeId>` | 按交易号查询包裹拦截工单 |
| `customer order return-by-trade <tradeId>` | 按交易号查询包裹退件工单 |
| `customer profile order count <customerId>` | 查询客户订单统计 |
| `customer profile order daily <customerId>` | 查询客户日订单数量 |
| `customer profile transaction count <customerId>` | 查询客户交易统计 |
| `customer profile transaction bills <customerId>` | 查询客户交易流水 |

### 物流管理 `logistics`

| 命令格式 | 说明 |
|----------|------|
| `logistics tracking <trackingNumber>` | 查询物流轨迹 |

### 采购订单管理 `purchase`

| 命令格式 | 说明 |
|----------|------|
| `purchase info <orderId>` | 查询采购订单详情 |
| `purchase info page` | 分页查询采购订单 |
| `purchase log <orderId>` | 查询采购日志 |
| `purchase logistics <orderId>` | 查询采购订单国际物流轨迹 |

### 仓库管理 `warehouse`

| 命令格式 | 说明 |
|----------|------|
| `warehouse tracking <trackingNumber>` | 查询仓库物流轨迹（快递签收、仓库签收、入库、物流轨迹、仓库操作） |
| `warehouse list --name <name>` | 按名称查询仓库列表 |
| `warehouse logistics-list --name <name>` | 按名称查询物流列表 |

### 采购售后查询 `purchase-after-sales`

| 命令格式 | 说明 |
|----------|------|
| `purchase-after-sales by-trade <tradeId>` | 按交易号查询采购售后工单 |
| `purchase-after-sales by-order <orderId>` | 按订单号查询采购售后工单 |

### 供应商管理 `supplier`

| 命令格式 | 说明 |
|----------|------|
| `supplier detail <supplierId>` | 查询供应商详情 |
| `supplier list --name <name>` | 分页查询供应商列表 |
| `supplier pur-count <loginId1,loginId2>` | 按登录 ID 查询采购次数 |
| `supplier ranking` | 供应商排行信息 |
| `supplier spu-list` | 供应商 SPU 信息 |
| `supplier after-sales` | 供应商售后信息 |
| `supplier fulfillment` | 供应商履约信息 |

### 第三方平台数据查询 `third-app-api`

| 命令格式 | 说明 |
|----------|------|
| `third-app-api ali-logistics <orderId>` | 查询 1688 订单物流信息 |
| `third-app-api ali-logistics-trace <orderId>` | 查询 1688 订单物流轨迹 |
| `third-app-api ali-order-detail <orderId>` | 查询 1688 订单详情 |
| `third-app-api ali-refund-detail <refundId>` | 查询 1688 退款单详情 |
| `third-app-api ali-refund-list <orderId>` | 查询 1688 退款单列表 |
| `third-app-api ali-refund-operations <refundId>` | 查询 1688 退款单操作记录 |
| `third-app-api ali-product --id <id>` | 查询 1688 产品信息（按 ID） |
| `third-app-api ali-product --url <url>` | 查询 1688 产品信息（按链接） |
| `third-app-api ali-supplier --login-id <id>` | 查询 1688 供应商信息 |
| `third-app-api ali-mix-config --member-id <id>` | 查询 1688 卖家混批设置 |
| `third-app-api ali-sub-accounts` | 查询 1688 子账号列表 |
| `third-app-api shop-product --store-id <id>` | 查询店铺产品 |
| `third-app-api shop-order-list --store-id <id>` | 查询店铺订单列表 |
| `third-app-api shop-inventory --store-id <id>` | 查询店铺库存 |

### AI 分析查询 `ai-analysis`

| 命令格式 | 说明 |
|----------|------|
| `ai-analysis intl-logistics` | 国际物流异常分析（分页） |
| `ai-analysis intl-logistics-count` | 国际物流异常风险等级统计 |
| `ai-analysis inventory` | 库存动销分析（分页） |
| `ai-analysis inventory-count` | 库存动销风险等级统计 |
| `ai-analysis purchase` | 采购售后分析（分页） |
| `ai-analysis purchase-count` | 采购售后风险等级统计 |
| `ai-analysis supplier` | 供应商 AI 分析（分页） |
| `ai-analysis supplier-count` | 供应商等级统计 |
| `ai-analysis customer-order` | 客户订单 AI 分析（分页） |
| `ai-analysis customer-order-count` | 客户订单风险等级统计 |
| `ai-analysis task-progress <type>` | AI 任务风险等级汇总 |

> 所有分页命令支持 `--page`/`-p`（页码）和 `--limit`/`-l`（每页条数）参数。

### 增值服务查询 `extra-service`

| 命令格式 | 说明 |
|----------|------|
| `extra-service list` | 查询增值服务名称列表 |
| `extra-service warehouses <serviceId>` | 查询增值服务关联的仓储 |
| `extra-service goods <serviceId>` | 查询增值服务关联的商品 |
| `extra-service warehouse-add <repertoryId>` | 查询仓库开通的增值服务 |

### 国内第三方交易异常订单管理 `domestic-third-trade-exception`

| 命令格式 | 说明 |
|----------|------|
| `domestic-third-trade-exception page-list` | 分页查询国内第三方交易异常订单（含物流明细） |
| `domestic-third-trade-exception message-list <monitorOrderId> <monitorLogisticsId>` | 查询异常订单留言列表 |

## Makefile 命令

| 命令 | 说明 |
|------|------|
| `make build` | 编译当前平台版本 |
| `make build-linux` | 交叉编译 Linux amd64 版本 |
| `make build-macos` | 交叉编译 macOS amd64 版本 |
| `make test` | 运行所有测试 |
| `make lint` | 代码静态检查 |
| `make run` | 直接运行查看帮助 |
| `make clean` | 清理构建产物 |
| `make tidy` | 整理 Go 模块依赖 |

## 测试

```powershell
# 运行全部测试
go test -v ./...

# 运行指定包测试
go test -v ./internal/httpclient/...
go test -v ./cmd/...
go test -v ./internal/apis/...

# 查看测试覆盖率
go test -cover ./...
```

## 技术栈

- [Go](https://go.dev/) - 编程语言
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- 标准库 `net/http` + `httptest` - HTTP 客户端与测试
