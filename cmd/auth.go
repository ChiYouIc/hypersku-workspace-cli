package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hypersku/hypersku-cli/internal/apis"
	"github.com/spf13/cobra"
)

// authCmd 表示认证管理命令组
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "管理登录状态",
	Long: `管理 HyperSKU CLI 的登录状态。

登录凭证保存于 ~/.hypersku-cli/config.json（api_base_url 与 api_token）。

子命令：
  login   登录（暂未实现，正常返回）
  status  校验 token 并查看当前登录状态
  logout  退出登录（暂未实现，正常返回）`,
}

// authLoginCmd 登录命令（暂未实现：不执行任何操作，正常返回退出码 0）。
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "登录 HyperSKU 账号（暂未实现）",
	Long: `登录 HyperSKU 账号。

暂未实现：当前版本不执行任何操作，直接正常返回。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("login 暂未实现，请手动在 ~/.hypersku-cli/config.json 中配置 api_token")
		return nil
	},
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

读取 ~/.hypersku-cli/config.json 中的 api_base_url 与 api_token，
请求 /api/admin/user/front/info 远程校验 token：
token 有效时输出 "Logged in as <用户名>"（退出码 0）；
未配置或 token 失效时输出 "Logged out"（退出码 1）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token, baseURL := authTokenFromConfig()

		var info *apis.UserInfo
		if token != "" && baseURL != "" {
			client := apis.NewAuthApi()
			info, _ = client.GetUserInfo(token) // 校验失败视为未登录
		}

		if statusJSON {
			out := map[string]any{"logged_in": false}
			if info != nil {
				out["logged_in"] = true
				out["account"] = statusAccount(info)
			}
			b, _ := json.Marshal(out)
			fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b)
			if info != nil {
				return nil
			}
			return &exitCodeError{1, ""}
		}

		if info != nil {
			fmt.Fprintf(cmd.OutOrStdout(), "Logged in as %s\n", statusAccount(info))
			return nil
		}
		fmt.Fprintln(cmd.OutOrStdout(), "Logged out")
		return &exitCodeError{1, ""}
	},
}

// authLogoutCmd 退出登录（暂未实现：不执行任何操作，正常返回退出码 0）。
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录（暂未实现）",
	Long: `退出登录。

暂未实现：当前版本不执行任何操作，直接正常返回。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("logout 暂未实现，请手动删除 ~/.hypersku-cli/config.json 中的 api_token")
		return nil
	},
}

var statusJSON bool

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

// statusAccount 返回用于 status 展示的账号标识（缺省 "unknown"）。
func statusAccount(info *apis.UserInfo) string {
	if info.Username != "" {
		return info.Username
	}
	if info.Name != "" {
		return info.Name
	}
	if info.Nickname != "" {
		return info.Nickname
	}
	return "unknown"
}

func init() {
	authStatusCmd.Flags().BoolVar(&statusJSON, "json", false, "以 JSON 输出状态（statusMatchJson 场景）")
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}
