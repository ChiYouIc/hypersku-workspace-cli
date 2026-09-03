package apis

import (
	"fmt"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// Auth 认证相关 API（基于 ~/.hypersku-cli/config.json 中的 api_base_url / api_token）
type Auth struct {
	http *httpclient.Client
}

// NewAuthApi 创建认证 API 客户端（复用全局默认客户端，baseURL 来自配置）
func NewAuthApi() *Auth {
	return &Auth{
		http: httpclient.DefaultClient,
	}
}

// UserInfo 当前登录用户信息（GET /api/admin/user/front/info 的 data 载荷）
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Proxy    bool   `json:"proxy"`
	RoleCode string `json:"roleCode"`
	Roles    string `json:"roles"`
	Username string `json:"username"`
}

// GetUserInfo 携带 token 请求服务端校验，并返回当前用户信息。
// token 缺失、失效或服务端拒绝时返回错误（调用方据此判定为未登录）。
func (a *Auth) GetUserInfo(token string) (*UserInfo, error) {
	if token == "" {
		return nil, fmt.Errorf("api_token 未配置")
	}

	result := &UserInfo{}
	path := "/api/admin/user/front/info?token=" + url.QueryEscape(token)
	if err := a.http.Get(path, result); err != nil {
		return nil, err
	}
	// token 失效时服务端返回 HTTP 200 + {"message":"...","status":40101}（无 data），
	// 因此成功必须同时满足：业务码正常（0 或 200）且返回了带 id 的用户数据
	if result.ID == "" {
		return nil, fmt.Errorf("token 校验失败: 响应中无用户数据")
	}
	return result, nil
}
