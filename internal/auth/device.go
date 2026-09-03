package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 默认认证服务端点（占位域名，服务端就绪后替换，
// 也可通过 config.json 的 auth 节或环境变量覆盖）。
const (
	defaultAuthBase    = "https://auth.hypersku.com"
	defaultDeviceScope = ""
)

// Endpoints 描述 Device Code Flow 依赖的服务端点。
type Endpoints struct {
	DeviceAuthorizationURL string // 设备码申请端点（RFC 8628 §3.1）
	TokenURL               string // token 端点（RFC 8628 §3.4）
	ClientID               string
	Scopes                 []string
}

// DefaultEndpoints 返回内置默认端点。
func DefaultEndpoints() Endpoints {
	return Endpoints{
		DeviceAuthorizationURL: defaultAuthBase + "/oauth/device/code",
		TokenURL:               defaultAuthBase + "/oauth/token",
		ClientID:               "hypersku-cli",
		Scopes:                 nil,
	}
}

// DeviceAuthorizationResponse 是设备码申请响应（RFC 8628 §3.2）。
type DeviceAuthorizationResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"` // 秒
	Interval                int    `json:"interval"`   // 最小轮询间隔（秒）
}

// BestVerificationURI 优先返回带用户码的完整验证 URL。
func (d *DeviceAuthorizationResponse) BestVerificationURI() string {
	if d == nil {
		return ""
	}
	if d.VerificationURIComplete != "" {
		return d.VerificationURIComplete
	}
	return d.VerificationURI
}

// TokenResponse 是 token 端点响应；Error 字段非空表示 RFC 8628 §3.5 定义的错误。
type TokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type,omitempty"`
	RefreshToken     string `json:"refresh_token,omitempty"`
	Scope            string `json:"scope,omitempty"`
	ExpiresIn        int    `json:"expires_in,omitempty"`
	Account          string `json:"account,omitempty"` // 服务端可选下发的账号标识
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// TokenError 表示 token 端点返回的协议错误（authorization_pending / slow_down /
// access_denied / expired_token 等）。
type TokenError struct {
	Code        string
	Description string
	StatusCode  int
}

func (e *TokenError) Error() string {
	if e.Description != "" {
		return fmt.Sprintf("%s: %s", e.Code, e.Description)
	}
	return e.Code
}

// Client 是 Device Code Flow 客户端。
type Client struct {
	hc *http.Client
	ep Endpoints
}

// NewClient 基于给定端点创建客户端。
func NewClient(ep Endpoints) *Client {
	return &Client{
		hc: &http.Client{Timeout: 10 * time.Second},
		ep: ep,
	}
}

// StartDeviceFlow 向服务端申请设备码与验证 URL。
func (c *Client) StartDeviceFlow(ctx context.Context) (*DeviceAuthorizationResponse, error) {
	form := url.Values{}
	form.Set("client_id", c.ep.ClientID)
	if len(c.ep.Scopes) > 0 {
		form.Set("scope", strings.Join(c.ep.Scopes, " "))
	} else if defaultDeviceScope != "" {
		form.Set("scope", defaultDeviceScope)
	}

	data, err := c.postForm(ctx, c.ep.DeviceAuthorizationURL, form)
	if err != nil {
		return nil, fmt.Errorf("请求设备码失败: %w", err)
	}

	var resp DeviceAuthorizationResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析设备码响应失败: %w", err)
	}
	if resp.DeviceCode == "" || resp.VerificationURI == "" {
		return nil, fmt.Errorf("服务端响应缺少 device_code 或 verification_uri")
	}
	if resp.ExpiresIn <= 0 {
		resp.ExpiresIn = 600
	}
	return &resp, nil
}

// ExchangeToken 用设备码换取 token（单次尝试，不做轮询）。
// 协议错误（如 authorization_pending）以 *TokenError 返回。
func (c *Client) ExchangeToken(ctx context.Context, deviceCode string) (*TokenResponse, error) {
	form := url.Values{}
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Set("device_code", deviceCode)
	form.Set("client_id", c.ep.ClientID)

	data, statusCode, err := c.postFormStatus(ctx, c.ep.TokenURL, form)
	if err != nil {
		return nil, fmt.Errorf("请求 token 失败: %w", err)
	}

	var resp TokenResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("解析 token 响应失败: %w", err)
	}
	if resp.Error != "" {
		return nil, &TokenError{Code: resp.Error, Description: resp.ErrorDescription, StatusCode: statusCode}
	}
	if resp.AccessToken == "" {
		return nil, fmt.Errorf("服务端响应缺少 access_token")
	}
	return &resp, nil
}

// AwaitToken 以给定间隔轮询 token 端点直到授权完成、被拒绝或设备码过期。
// slow_down 错误会自动放宽轮询间隔。
func (c *Client) AwaitToken(ctx context.Context, d *DeviceAuthorizationResponse, interval time.Duration) (*Credentials, error) {
	if interval <= 0 {
		interval = 5 * time.Second
	}
	expiresIn := time.Duration(d.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = 10 * time.Minute
	}
	deadline := time.Now().Add(expiresIn)

	for {
		resp, err := c.ExchangeToken(ctx, d.DeviceCode)
		if err == nil {
			return TokenToCredentials(resp, time.Now()), nil
		}
		var te *TokenError
		if !errors.As(err, &te) {
			return nil, err
		}
		switch te.Code {
		case "authorization_pending":
			// 继续等待
		case "slow_down":
			interval += 5 * time.Second
		default:
			// access_denied / expired_token 等，不可恢复
			return nil, fmt.Errorf("授权失败: %w", te)
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("设备码已过期，请重新执行 auth login")
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
		}
	}
}

// TryAdvancePendingDevice 推进一次暂存的设备授权（供 status 轮询场景使用）：
//   - 服务端已发放 token：保存凭证、清理中间态，返回 true
//   - 仍在 pending：返回 false（中间态保留，等待下一次轮询）
//   - 设备码过期/被拒绝：清理中间态，返回 false
func TryAdvancePendingDevice(ctx context.Context, c *Client, s *Store, now time.Time) bool {
	p, err := s.LoadPendingDevice()
	if err != nil || p == nil {
		return false
	}
	if !p.ExpiresAt.IsZero() && now.After(p.ExpiresAt) {
		_ = s.ClearPendingDevice()
		return false
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.ExchangeToken(ctx, p.DeviceCode)
	if err == nil {
		if saveErr := s.SaveCredentials(TokenToCredentials(resp, now)); saveErr != nil {
			return false
		}
		_ = s.ClearPendingDevice()
		return true
	}

	var te *TokenError
	if errors.As(err, &te) {
		switch te.Code {
		case "expired_token", "access_denied":
			_ = s.ClearPendingDevice()
		}
	}
	return false
}

// TokenToCredentials 将 token 响应转换为可持久化的凭证。
func TokenToCredentials(tr *TokenResponse, now time.Time) *Credentials {
	c := &Credentials{
		AccessToken:  tr.AccessToken,
		TokenType:    tr.TokenType,
		RefreshToken: tr.RefreshToken,
		Scope:        tr.Scope,
		Account:      tr.Account,
		ObtainedAt:   now,
	}
	if tr.ExpiresIn > 0 {
		c.ExpiresAt = now.Add(time.Duration(tr.ExpiresIn) * time.Second)
	}
	return c
}

// ---------- 内部 ----------

func (c *Client) postForm(ctx context.Context, u string, form url.Values) ([]byte, error) {
	data, _, err := c.postFormStatus(ctx, u, form)
	return data, err
}

func (c *Client) postFormStatus(ctx context.Context, u string, form url.Values) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "HyperSKU-CLI/auth")

	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, resp.StatusCode, err
	}
	return body, resp.StatusCode, nil
}
