# HyperSKU CLI

> 基于 Go 构建的高效命令行工具，提供灵活的第三方业务 API 调用能力，涵盖售后、客户、物流、采购等场景。

## 目录结构

```
hypersku-cli/
├── .vscode/
│   └── settings.json          # VS Code 配置（Go 路径、国内代理镜像）
├── cmd/                       # cobra 命令定义（薄命令层，只做编排）
│   ├── root.go                # 根命令（全局 flags、配置加载、依赖初始化）
│   ├── after_sales.go         # 售后管理命令
│   ├── customer.go            # 客户管理命令
│   ├── logistics.go           # 物流管理命令
│   ├── purchase.go            # 采购订单管理命令
│   ├── warehouse.go           # 仓库管理命令
│   └── domestic_third_trade_exception.go  # 国内第三方交易异常订单管理命令
├── internal/                  # 内部业务逻辑（不对外暴露）
│   ├── apis/                  # 第三方 API 封装
│   │   ├── types.go           # 通用类型定义
│   │   ├── after_sales.go     # 售后 API
│   │   ├── customer.go        # 客户 API
│   │   ├── customer_order_return.go  # 客户订单退件工单 API
│   │   ├── logistics.go       # 物流 API
│   │   ├── purchase.go        # 采购 API
│   │   ├── warehouse.go       # 仓库 API
│   │   └── domestic_third_trade_exception.go  # 国内第三方交易异常订单 API
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
| `api_token` | string | `""` | API 认证令牌，自动作为 `authorization` 请求头发送 |

> 配置文件不存在时自动使用默认配置；启动时会自动创建 `~/.hypersku-cli` 目录（含 `logs/`、`data/` 子目录）。

## 使用指南

### 全局 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--config` | string | `""` | 配置文件路径 |
| `-v, --version` | bool | `false` | 显示版本信息 |
| `-h, --help` | bool | `false` | 显示帮助信息 |

### 子命令

| 命令 | 说明 |
|------|------|
| `after-sales` | 售后管理 |
| `customer` | 客户管理 |
| `logistics` | 物流管理 |
| `purchase` | 采购订单管理 |
| `warehouse` | 仓库管理 |
| `domestic-third-trade-exception` | 国内第三方交易异常订单管理 |

#### 售后管理 `after-sales`

| 命令 | 说明 |
|------|------|
| `1688 <thirdOrderId>` | 查询 1688 售后工单 |
| `goods <thirdOrderId> <refundId>` | 查询 1688 售后商品 |
| `detail <refundId>` | 查询 1688 售后详情 |
| `message <refundId>` | 查询 1688 售后留言 |

#### 客户管理 `customer`

| 命令 | 说明 |
|------|------|
| `order info <orderId>` | 查询订单信息 |
| `order logistics <orderId>` | 查询订单物流信息 |
| `order address <orderId>` | 查询订单地址信息 |
| `order return <customerOrderId>` | 查询客户订单退件工单 |

#### 物流管理 `logistics`

| 命令 | 说明 |
|------|------|
| `tracking <trackingNumber>` | 查询物流轨迹 |

#### 采购订单管理 `purchase`

| 命令 | 说明 |
|------|------|
| `info <orderId>` | 查询采购订单详情 |
| `page` | 分页查询采购订单 |
| `log <orderId>` | 查询采购日志 |
| `logistics <orderId>` | 查询采购订单国际物流轨迹 |

`purchase page` 支持以下过滤参数：

| Flag | 默认值 | 说明 |
|------|--------|------|
| `-p, --page` | `1` | 页码 |
| `-l, --limit` | `10` | 页大小 |
| `--start` | `""` | 开始时间，格式：`yyyy-MM-dd HH:mm:ss` |
| `--end` | `""` | 结束时间，格式：`yyyy-MM-dd HH:mm:ss` |
| `--thirdOrderId` | `""` | 交易号、第三方订单号 |
| `--trackingNumber` | `""` | 物流单号 |

#### 仓库管理 `warehouse`

| 命令 | 说明 |
|------|------|
| `tracking <trackingNumber>` | 查询仓库物流轨迹（快递/仓库签收、入库、物流轨迹、仓库操作） |

#### 国内第三方交易异常订单管理 `domestic-third-trade-exception`

| 命令 | 说明 |
|------|------|
| `page-list` | 分页查询国内第三方交易异常订单（含物流明细） |
| `message-list <monitorOrderId> <monitorLogisticsId>` | 查询异常订单留言列表 |

`domestic-third-trade-exception page-list` 支持以下过滤参数：

| Flag | 默认值 | 说明 |
|------|--------|------|
| `-p, --page` | `1` | 页码 |
| `-l, --limit` | `10` | 页大小 |
| `-s, --hypersku-status` | `0` | 异常主状态（1-未发货，2-假发货，3-未到货，4-假签收，5-未签收，6-退件，7-丢件，8-未入库，9-丢包裹，10-无货） |
| `-c, --hypersku-sub-status` | `[1,2]` | 异常子状态列表（1-待处理，2-处理中，3-已处理，4-已关闭，5-已拒绝） |
| `-b, --buyer-id` | `""` | 买家 ID（可选） |

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
