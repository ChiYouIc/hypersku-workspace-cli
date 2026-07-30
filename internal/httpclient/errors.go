package httpclient

import "fmt"

// HTTPError 表示 HTTP 请求返回的非成功状态码错误
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTP 请求失败，状态码: %d，响应: %s", e.StatusCode, e.Body)
}
