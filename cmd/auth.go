package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hypersku/hypersku-cli/internal/auth"
	"github.com/spf13/cobra"
)

// authCmd 表示认证管理命令组
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "管理登录状态",
	Long: `管理 HyperSKU CLI 的登录状态。

子命令：
  login   发起设备授权登录（Device Code Flow）
  status  查看当前登录状态（幂等、无副作用）
  logout  退出登录并清理本地凭证`,
}

// authLoginCmd 发起 Device Code Flow 登录。
//
// workbuddy 连接器场景：本命令在 stdout 立即输出认证 URL 后即退出，
// 登录态由后续轮询 `auth status` 推进（服务端发放 token 后自动保存凭证）。
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "发起设备授权登录",
	Long: `发起设备授权登录（OAuth2 Device Code Flow，RFC 8628）。

流程：
  1. 向认证服务申请设备码，输出验证 URL
  2. 在浏览器中打开 URL 并确认授权
  3. CLI（前台等待模式）或后续 auth status 轮询（连接器模式）换取 token`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := auth.DefaultStore()
		if err != nil {
			return err
		}
		client := newAuthClient()

		// 已登录则直接提示（幂等：重复 login 不报错）
		if cred, _ := store.LoadCredentials(); cred != nil && cred.Valid(time.Now()) {
			cmd.Printf("已登录%s\n", accountSuffix(cred.Account))
			cmd.Println("如需切换账号，请先执行 hypersku-cli auth logout")
			return nil
		}

		resp, err := client.StartDeviceFlow(context.Background())
		if err != nil {
			return fmt.Errorf("发起设备授权失败: %w", err)
		}

		// 暂存设备授权中间态，供 auth status / 前台轮询推进
		pending := &auth.PendingDevice{
			DeviceCode:          resp.DeviceCode,
			UserCode:            resp.UserCode,
			VerificationURI:     resp.VerificationURI,
			VerificationURIComp: resp.VerificationURIComplete,
			Interval:            resp.Interval,
			ExpiresAt:           time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second),
		}
		if err := store.SavePendingDevice(pending); err != nil {
			return fmt.Errorf("保存设备授权状态失败: %w", err)
		}

		uri := resp.BestVerificationURI()
		// 关键输出：完整 https:// URL，前后空白分隔、无引号包裹（连接器提取约定）
		fmt.Fprintf(cmd.OutOrStdout(), "%s\n", uri)
		fmt.Fprintf(cmd.ErrOrStderr(), "用户码: %s\n", resp.UserCode)
		fmt.Fprintf(cmd.ErrOrStderr(), "在浏览器完成授权后，CLI 将自动检测登录状态。\n")

		// 前台等待模式（--wait）：就地轮询直到授权完成
		if waitLogin {
			interval := time.Duration(resp.Interval) * time.Second
			if interval <= 0 {
				interval = 5 * time.Second
			}
			if interval > 30*time.Second {
				interval = 30 * time.Second
			}
			cred, err := client.AwaitToken(context.Background(), resp, interval)
			if err != nil {
				return fmt.Errorf("等待授权失败: %w", err)
			}
			if err := store.SaveCredentials(cred); err != nil {
				return fmt.Errorf("保存凭证失败: %w", err)
			}
			_ = store.ClearPendingDevice()
			cmd.Printf("登录成功%s\n", accountSuffix(cred.Account))
		}
		return nil
	},
}

// authStatusCmd 查看登录状态：幂等、无副作用（除推进 pending 设备授权外）。
//
// 输出约定（供 statusMatch / statusMatchJson 匹配）：
//   - 已登录：退出码 0，stdout 首行 "Logged in as <account>"
//   - 已登录（JSON）：{"logged_in": true, "account": "..."}
//   - 未登录：退出码 1，输出 "Logged out"
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "查看当前登录状态",
	Long: `查看当前登录状态（幂等、无副作用）。

已登录时输出 "Logged in as <账号>"（退出码 0）；未登录时输出 "Logged out"（退出码 1）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := auth.DefaultStore()
		if err != nil {
			return err
		}

		now := time.Now()
		cred, err := store.LoadCredentials()
		if err != nil {
			return fmt.Errorf("读取凭证失败: %w", err)
		}

		// 连接器轮询场景：若存在进行中的设备授权，尝试推进一次
		if cred == nil || !cred.Valid(now) {
			client := newAuthClient()
			auth.TryAdvancePendingDevice(context.Background(), client, store, now)
			cred, err = store.LoadCredentials()
			if err != nil {
				return fmt.Errorf("读取凭证失败: %w", err)
			}
		}

		if statusJSON {
			out := map[string]any{"logged_in": false}
			if cred != nil && cred.Valid(now) {
				out["logged_in"] = true
				out["account"] = statusAccount(cred)
			}
			b, _ := json.Marshal(out)
			cmd.Println(string(b))
			if cred != nil && cred.Valid(now) {
				return nil
			}
			return &exitCodeError{1, ""}
		}

		if cred != nil && cred.Valid(now) {
			cmd.Printf("Logged in as %s\n", statusAccount(cred))
			return nil
		}
		cmd.Println("Logged out")
		return &exitCodeError{1, ""}
	},
}

// authLogoutCmd 清理本地凭证（未登录时也正常返回 0）。
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "退出登录并清理本地凭证",
	Long: `退出登录，删除本地凭证与设备授权中间态。

未登录时执行本命令同样正常返回（幂等）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := auth.DefaultStore()
		if err != nil {
			return err
		}
		if err := store.ClearCredentials(); err != nil {
			return fmt.Errorf("清理凭证失败: %w", err)
		}
		if err := store.ClearPendingDevice(); err != nil {
			return fmt.Errorf("清理设备授权状态失败: %w", err)
		}
		cmd.Println("Logged out")
		return nil
	},
}

var (
	waitLogin   bool
	statusJSON  bool
	authBaseURL string
)

// exitCodeError 用于让 cobra 以指定退出码结束且不打印重复错误。
type exitCodeError struct {
	code int
	msg  string
}

func (e *exitCodeError) Error() string { return e.msg }

// newAuthClient 基于配置构建 Device Flow 客户端。
func newAuthClient() *auth.Client {
	ep := auth.DefaultEndpoints()
	if loadedConfig != nil && loadedConfig.Auth != nil {
		if base := strings.TrimRight(loadedConfig.Auth.BaseURL, "/"); base != "" {
			ep.DeviceAuthorizationURL = base + "/oauth/device/code"
			ep.TokenURL = base + "/oauth/token"
		}
		if id := loadedConfig.Auth.ClientID; id != "" {
			ep.ClientID = id
		}
	}
	return auth.NewClient(ep)
}

// accountSuffix 生成 " (账号)" 后缀。
func accountSuffix(account string) string {
	if account == "" {
		return ""
	}
	return " (" + account + ")"
}

// statusAccount 返回用于 status 展示的账号标识（缺省 "unknown"）。
func statusAccount(cred *auth.Credentials) string {
	if cred.Account != "" {
		return cred.Account
	}
	return "unknown"
}

func init() {
	authLoginCmd.Flags().BoolVar(&waitLogin, "wait", false, "发起授权后就地等待授权完成（默认只输出 URL 即退出）")
	authStatusCmd.Flags().BoolVar(&statusJSON, "json", false, "以 JSON 输出状态（statusMatchJson 场景）")
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}
