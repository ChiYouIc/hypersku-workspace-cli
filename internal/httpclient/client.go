package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// DefaultClient 是全局默认的 HTTP 客户端单例，所有子命令均可直接使用
var DefaultClient *Client

// once 确保 Init 只执行一次
var once sync.Once

// Init 初始化全局默认的 HTTP 客户端，由根命令的 PersistentPreRun 调用
// 多次调用不会重复初始化
func Init(opts ...Option) {
	once.Do(func() {
		DefaultClient = New(opts...)
	})
}

// Reset 重置全局客户端，允许重新初始化（通常用于测试）
func Reset() {
	once = sync.Once{}
	DefaultClient = nil
}

// Client 是 HTTP 客户端的封装，用于调用第三方 API
type Client struct {
	baseURL    string
	httpClient *http.Client
	headers    map[string]string
	debug      bool
}

// New 创建一个新的 HTTP 客户端
func New(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// SetBaseURL 设置基础 URL
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// SetHeader 设置全局请求头
func (c *Client) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetTimeout 设置请求超时时间
func (c *Client) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// Get 发送 GET 请求
func (c *Client) Get(path string, result interface{}) error {
	return c.doRequest(http.MethodGet, path, nil, result)
}

// Post 发送 POST 请求
func (c *Client) Post(path string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPost, path, body, result)
}

// Put 发送 PUT 请求
func (c *Client) Put(path string, body interface{}, result interface{}) error {
	return c.doRequest(http.MethodPut, path, body, result)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(path string, result interface{}) error {
	return c.doRequest(http.MethodDelete, path, nil, result)
}

// doRequest 执行 HTTP 请求的核心方法
func (c *Client) doRequest(method, path string, body, result interface{}) error {
	url := c.buildURL(path)

	var reqBody io.Reader
	if body != nil {
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("序列化请求体失败: %w", err)
		}
		reqBody = bytes.NewReader(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置全局请求头
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// 默认设置为 JSON 格式
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	if c.debug {
		fmt.Printf("[HTTP] %s %s\n", method, url)
		if body != nil {
			fmt.Printf("[HTTP] 请求体: %+v\n", body)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应体失败: %w", err)
	}

	if c.debug {
		fmt.Printf("[HTTP] 响应状态: %s\n", resp.Status)
		fmt.Printf("[HTTP] 响应体: %s\n", string(respBody))
	}

	// 检查 HTTP 状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &HTTPError{
			StatusCode: resp.StatusCode,
			Body:       string(respBody),
		}
	}

	// 解析响应 JSON
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("解析响应 JSON 失败: %w", err)
		}
	}

	return nil
}

// buildURL 拼接完整的请求 URL，自动处理斜杠
func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}

	baseURL := strings.TrimRight(c.baseURL, "/")
	requestPath := strings.TrimLeft(path, "/")

	if requestPath == "" {
		return baseURL
	}
	return baseURL + "/" + requestPath
}
