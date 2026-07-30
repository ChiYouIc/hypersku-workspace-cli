package httpclient

import "time"

// Option 是 HTTP 客户端的配置选项函数
type Option func(*Client)

// WithBaseURL 设置基础 URL
func WithBaseURL(url string) Option {
	return func(c *Client) {
		c.baseURL = url
	}
}

// WithTimeout 设置请求超时时间
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHeader 设置全局请求头
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.headers[key] = value
	}
}

// WithDebug 开启调试模式，会打印请求和响应信息
func WithDebug(enable bool) Option {
	return func(c *Client) {
		c.debug = enable
	}
}

// WithBearerToken 设置 Bearer Token 认证
func WithBearerToken(token string) Option {
	return func(c *Client) {
		c.headers["Authorization"] = "Bearer " + token
	}
}
