package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/hypersku/hypersku-cli/internal/config"
	httpclient "github.com/hypersku/hypersku-cli/internal/httpclient"
	"github.com/hypersku/hypersku-cli/internal/version"
	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	showVersion bool
)

// rootCmd 表示根命令，当不调用任何子命令时执行
var rootCmd = &cobra.Command{
	Use:   "hypersku-cli",
	Short: "HyperSKU CLI - 高效的命令行工具",
	Long: `HyperSKU CLI 是一个功能强大的命令行工具。

它提供了丰富的子命令来帮助您完成各种任务。
使用 --help 查看所有可用命令。`,
	// PersistentPreRun 会在每个子命令执行前调用，用于初始化全局依赖
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initGlobalDependencies()
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showVersion {
			cmd.Println(version.Get().String())
			return
		}
		cmd.Help()
	},
}

// Execute 将所有子命令添加到根命令并设置 flags。
// 由 main.main() 调用，只需执行一次。
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// initGlobalDependencies 初始化全局依赖，由 PersistentPreRun 调用
func initGlobalDependencies() {
	// 1. 确定配置文件路径
	configPath := cfgFile
	if configPath == "" {
		defaultPath, err := config.DefaultConfigPath()
		if err == nil {
			configPath = defaultPath
		}
	}

	// 2. 从配置文件加载配置
	cfg, err := config.Load(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 加载配置文件失败: %v\n", err)
	}

	// 3. 构建 HTTP 客户端选项
	opts := []httpclient.Option{
		httpclient.WithTimeout(time.Duration(cfg.APITimeout) * time.Second),
		httpclient.WithHeader("User-Agent", "HyperSKU-CLI/0.1.0"),
	}

	if cfg.APIBaseURL != "" {
		opts = append(opts, httpclient.WithBaseURL(cfg.APIBaseURL))
	}

	if cfg.APIToken != "" {
		opts = append(opts, httpclient.WithHeader("authorization", cfg.APIToken))
	}

	// 4. 确保配置目录存在（后续日志、数据等使用）
	if _, err := config.EnsureConfigDir(); err != nil {
		fmt.Fprintf(os.Stderr, "警告: 创建配置目录失败: %v\n", err)
	}

	httpclient.Init(opts...)
}

func init() {
	// 全局 flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径 (默认: ~/.hypersku-cli/config.json)")
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "显示版本信息")
}
