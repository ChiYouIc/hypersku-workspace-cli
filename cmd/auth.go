package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/hypersku/hypersku-cli/internal/config"
	"github.com/spf13/cobra"
)

// authCmd 表示认证管理命令组
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "管理登录状态",
	Long: `管理 HyperSKU CLI 的登录状态。

子命令：
  login   登录并保存 token（支持 access-token 参数或用户名密码）
  status  校验 token 并查看当前登录状态
  logout  退出登录（删除本地配置中的 token）`,
}

// authLoginCmd 登录命令：支持两种方式——
//   - login <access-token>：直接使用已有 token 登录（服务端校验后保存）
//   - login --api-user xxx --api-password yyy：用户名密码换取 token 后保存
var authLoginCmd = func() *cobra.Command {
	var username string
	var password string

	cmd := &cobra.Command{
		Use:   "login [access-token]",
		Short: "登录 HyperSKU 账号",
		Long: `登录 HyperSKU 账号。

两种方式：
  hypersku-cli auth login <access-token>                     # 直接使用已有 token
  hypersku-cli auth login --api-user <用户名> --api-password <密码>  # 用户名密码登录`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			authApi := apis.NewAuthApi()

			// 方式一：直接使用 access token
			if len(args) > 0 && args[0] != "" {
				return loginWithToken(cmd, authApi, args[0])
			}

			// 方式二：用户名密码登录
			return loginWithPassword(cmd, authApi, username, password)
		},
	}

	cmd.Flags().StringVar(&username, "api-user", "", "登录用户名")
	cmd.Flags().StringVar(&password, "api-password", "", "登录密码")

	return cmd
}

// loginWithToken 使用已有 access token 登录：服务端校验通过后保存。
func loginWithToken(cmd *cobra.Command, authApi *apis.Auth, token string) error {
	info, err := authApi.GetUserInfo(token)
	if err != nil {
		return &exitCodeError{code: 1, msg: fmt.Sprintf("登录失败: %v", err)}
	}

	if err := saveToken(token); err != nil {
		return &exitCodeError{code: 1, msg: err.Error()}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s\n", displayName(info))
	return nil
}

// loginWithPassword 用户名密码登录：换取 token 后保存。
func loginWithPassword(cmd *cobra.Command, authApi *apis.Auth, username, password string) error {
	if username == "" || password == "" {
		return &exitCodeError{code: 1, msg: "请提供 access token，或通过 --api-user 与 --api-password 指定账号密码"}
	}

	token, err := authApi.Login(username, password)
	if err != nil {
		return &exitCodeError{code: 1, msg: err.Error()}
	}

	if err := saveToken(token); err != nil {
		return &exitCodeError{code: 1, msg: err.Error()}
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s\n", username)
	return nil
}

// authStatusCmd 校验登录状态：读取 config.json 中的 api_base_url/api_token，
// 调用 /api/admin/user/front/info?token=xxx 远程校验 token。
//
// 输出约定（供 statusMatch 匹配）：
//   - 已登录：退出码 0，stdout 首行 "Logged in as <username>"
//   - 已登录（JSON）：{"logged_in": true, "account": "..."}
//   - 未登录：退出码 1，输出 "Logged out"
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "校验当前登录状态",
	Long: `校验当前登录状态。

远程校验本地保存的 token：
token 有效时输出 "Logged in as <用户名>"（退出码 0）；
未配置或 token 失效时输出 "Logged out"（退出码 1）。`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, baseURL := authTokenFromConfig()

		if token == "" || baseURL == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
			return &exitCodeError{code: 1, msg: ""}
		}

		client := apis.NewAuthApi()
		info, err := client.GetUserInfo(token)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
			return &exitCodeError{code: 1, msg: ""}
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s\n", displayName(info))
		return nil
	},
}

// authLogoutCmd 退出登录：删除本地配置文件中的 api_token 并落盘。
// 本地 token 清除后即视为未登录（无需调用服务端接口）。
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录",
	Long: `退出登录。

删除本地保存的登录凭证（保留其他配置），清除后即视为未登录。`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := authTokenFromConfig()
		if token == "" {
			fmt.Fprintln(cmd.OutOrStdout(), "Already logged out")
			return nil
		}

		loadedConfig.APIToken = ""
		loadedConfig.APITokenUpdatedAt = ""
		if err := saveConfig(); err != nil {
			return &exitCodeError{code: 1, msg: err.Error()}
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
		return nil
	},
}

// exitCodeError 用于让 cobra 以指定退出码结束且不打印重复错误。
type exitCodeError struct {
	code int
	msg  string
}

func (e *exitCodeError) Error() string { return e.msg }

// authTokenFromConfig 从全局已加载配置中取出 api_token 与 api_base_url。
func authTokenFromConfig() (token, baseURL string) {
	if loadedConfig == nil {
		return "", ""
	}
	return strings.TrimSpace(loadedConfig.APIToken), strings.TrimSpace(loadedConfig.APIBaseURL)
}

// displayName 返回用户展示名：依次取 username/name/nickname，均缺省时返回 unknown。
func displayName(info *apis.UserInfo) string {
	if info == nil {
		return "unknown"
	}
	switch {
	case info.Username != "":
		return info.Username
	case info.Name != "":
		return info.Name
	case info.Nickname != "":
		return info.Nickname
	default:
		return "unknown"
	}
}

// saveConfig 将当前 loadedConfig 写回配置文件（优先使用 --config 指定的路径）。
func saveConfig() error {
	savePath := cfgFile
	if savePath == "" {
		defaultPath, err := config.DefaultConfigPath()
		if err != nil {
			return fmt.Errorf("获取默认配置路径失败: %w", err)
		}
		savePath = defaultPath
	}
	return config.Save(savePath, loadedConfig)
}

// saveToken 将 token 写入内存配置并落盘，同时记录更新时间（RFC3339）。
func saveToken(token string) error {
	loadedConfig.APIToken = token
	loadedConfig.APITokenUpdatedAt = time.Now().Format(time.RFC3339)
	return saveConfig()
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd())
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}
