# HyperSKU CLI

> 基于 Go 构建的高效命令行工具，提供灵活的第三方 API 调用能力。

## 目录结构

```
hypersku-cli/
├── .vscode/
│   └── settings.json          # VS Code 配置（Go 路径、国内代理镜像）
├── cmd/                       # cobra 命令定义（薄命令层，只做编排）
│   ├── root.go                # 根命令（全局 flags、依赖初始化）
│   ├── hello.go               # 示例子命令
│   └── weather.go             # 天气查询命令（调用 internal/apis/weather）
├── internal/                  # 内部业务逻辑（不对外暴露）
│   ├── apis/                  # 第三方 API 封装
│   │   └── weather/           # 天气 API 封装
│   │       ├── client.go      # 天气 API 客户端（请求/响应/业务逻辑）
│   │       └── client_test.go # 天气 API 测试
│   ├── httpclient/            # HTTP 客户端基础封装
│   │   ├── client.go          # 核心客户端 + 全局 DefaultClient 单例
│   │   ├── options.go         # 配置选项（WithBaseURL, WithTimeout 等）
│   │   ├── errors.go          # 自定义 HTTP 错误类型
│   │   ├── client_test.go     # 客户端测试
│   │   └── errors_test.go     # 错误类型测试
│   └── version/               # 版本信息
│       ├── version.go         # 版本注入（构建时注入 commit/date）
│       └── version_test.go    # 版本测试
├── build/                     # 编译产物
├── pkg/                       # 可复用的公开包（预留）
├── main.go                    # 程序入口
├── go.mod / go.sum            # Go 模块依赖
├── Makefile                   # 构建脚本
├── .gitignore
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
```

### 2. 编译

```powershell
# Windows
go build -o build/hypersku-cli.exe .

# 或使用 Makefile
make build
```

### 3. 运行

```powershell
# 查看帮助
build\hypersku-cli.exe --help

# 显示版本
build\hypersku-cli.exe --version

# 示例子命令
build\hypersku-cli.exe hello 世界
```

## 使用指南

### 全局 Flags

| Flag | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--api-base-url` | string | `""` | 第三方 API 的基础 URL |
| `--api-timeout` | int | `30` | HTTP 请求超时时间（秒） |
| `--api-token` | string | `""` | Bearer Token 认证 |
| `--config` | string | `""` | 配置文件路径 |
| `-v, --version` | bool | `false` | 显示版本信息 |
| `-h, --help` | bool | `false` | 显示帮助信息 |

### 子命令

| 命令 | 说明 |
|------|------|
| `hello [name]` | 打印欢迎信息 |
| `weather --city <名称>` | 查询天气（第三方 API 封装示例） |

### 调用第三方 API

全局 HTTP 客户端通过 `--api-base-url` 传入基础地址，所有子命令自动获取：

```powershell
# 设置 API 基础地址后调用
hypersku-cli weather --city "北京" --api-base-url https://api.weather.com
```

在 `internal/apis/` 中封装第三方 API：

```go
// internal/apis/example/client.go —— API 封装层
package example

import (
    httpclient "github.com/hypersku/hypersku-cli/internal/httpclient"
)

type Response struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

type Client struct {
    HTTP *httpclient.Client
}

func New(hc *httpclient.Client) *Client {
    return &Client{HTTP: hc}
}

func (c *Client) GetData() (*Response, error) {
    var result Response
    err := c.HTTP.Get("/api/data", &result)
    return &result, err
}
```

在 `cmd/` 中编写薄命令层：

```go
// cmd/example.go —— 命令层，只做编排
package cmd

import (
    httpclient "github.com/hypersku/hypersku-cli/internal/httpclient"
    exampleapi "github.com/hypersku/hypersku-cli/internal/apis/example"
    "github.com/spf13/cobra"
)

var exampleCmd = &cobra.Command{
    Use: "example",
    RunE: func(cmd *cobra.Command, args []string) error {
        client := exampleapi.New(httpclient.DefaultClient)
        result, err := client.GetData()
        if err != nil {
            return err
        }
        cmd.Printf("结果: %+v\n", result)
        return nil
    },
}

func init() {
    rootCmd.AddCommand(exampleCmd)
}
```

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

# 查看测试覆盖率
go test -cover ./...
```

## 添加新子命令

在 `cmd/` 目录下新建文件，按以下模板编写：

```go
package cmd

import (
    "github.com/spf13/cobra"
)

var myCmd = &cobra.Command{
    Use:   "mycmd <参数>",
    Short: "简短描述",
    Long:  `详细的命令说明`,
    Args:  cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        // 使用 httpclient.DefaultClient 调用 API
        // 业务逻辑...
        return nil
    },
}

func init() {
    myCmd.Flags().String("flag-name", "默认值", "flag 说明")
    rootCmd.AddCommand(myCmd)
}
```

## 技术栈

- [Go](https://go.dev/) - 编程语言
- [Cobra](https://github.com/spf13/cobra) - CLI 框架
- 标准库 `net/http` + `httptest` - HTTP 客户端与测试
