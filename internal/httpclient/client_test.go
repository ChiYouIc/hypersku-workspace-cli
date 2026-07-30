package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// testPayload 测试用的请求/响应结构体
type testPayload struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

// TestNew 测试创建客户端
func TestNew(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("New() 返回了 nil")
	}

	if c.httpClient.Timeout != 30*time.Second {
		t.Errorf("期望默认超时 30s，实际得到 %v", c.httpClient.Timeout)
	}
}

// TestNewWithOptions 测试带选项创建客户端
func TestNewWithOptions(t *testing.T) {
	c := New(
		WithBaseURL("https://api.example.com"),
		WithTimeout(10*time.Second),
		WithHeader("X-Custom", "test"),
		WithBearerToken("my-token"),
	)

	if c.baseURL != "https://api.example.com" {
		t.Errorf("期望 baseURL = \"https://api.example.com\"，实际得到 %q", c.baseURL)
	}

	if c.httpClient.Timeout != 10*time.Second {
		t.Errorf("期望超时 10s，实际得到 %v", c.httpClient.Timeout)
	}

	if c.headers["X-Custom"] != "test" {
		t.Errorf("期望 X-Custom = \"test\"，实际得到 %q", c.headers["X-Custom"])
	}

	if c.headers["Authorization"] != "Bearer my-token" {
		t.Errorf("期望 Authorization = \"Bearer my-token\"，实际得到 %q", c.headers["Authorization"])
	}
}

// TestSetBaseURL 测试 SetBaseURL
func TestSetBaseURL(t *testing.T) {
	c := New()
	c.SetBaseURL("https://api.example.com")

	if c.baseURL != "https://api.example.com" {
		t.Errorf("期望 baseURL = \"https://api.example.com\"，实际得到 %q", c.baseURL)
	}
}

// TestSetHeader 测试 SetHeader
func TestSetHeader(t *testing.T) {
	c := New()
	c.SetHeader("X-Api-Key", "secret-123")

	if c.headers["X-Api-Key"] != "secret-123" {
		t.Errorf("期望 X-Api-Key = \"secret-123\"，实际得到 %q", c.headers["X-Api-Key"])
	}
}

// TestSetTimeout 测试 SetTimeout
func TestSetTimeout(t *testing.T) {
	c := New()
	c.SetTimeout(5 * time.Second)

	if c.httpClient.Timeout != 5*time.Second {
		t.Errorf("期望超时 5s，实际得到 %v", c.httpClient.Timeout)
	}
}

// TestGet 测试 GET 请求
func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("期望 GET 请求，实际得到 %s", r.Method)
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("期望 Accept = \"application/json\"，实际得到 %q", r.Header.Get("Accept"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"test","value":42}`))
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	var result testPayload
	if err := c.Get("/", &result); err != nil {
		t.Fatalf("Get 请求失败: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("期望 Name = \"test\"，实际得到 %q", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("期望 Value = 42，实际得到 %d", result.Value)
	}
}

// TestPost 测试 POST 请求
func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("期望 POST 请求，实际得到 %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("期望 Content-Type = \"application/json\"，实际得到 %q", r.Header.Get("Content-Type"))
		}

		var req testPayload
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("解析请求体失败: %v", err)
		}

		// 返回处理后的数据
		resp := testPayload{
			Name:  req.Name,
			Value: req.Value * 2,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	body := testPayload{Name: "hello", Value: 21}
	var result testPayload
	if err := c.Post("/", body, &result); err != nil {
		t.Fatalf("Post 请求失败: %v", err)
	}

	if result.Name != "hello" {
		t.Errorf("期望 Name = \"hello\"，实际得到 %q", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("期望 Value = 42，实际得到 %d", result.Value)
	}
}

// TestPut 测试 PUT 请求
func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("期望 PUT 请求，实际得到 %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name":"updated","value":99}`))
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	body := testPayload{Name: "original", Value: 1}
	var result testPayload
	if err := c.Put("/", body, &result); err != nil {
		t.Fatalf("Put 请求失败: %v", err)
	}

	if result.Name != "updated" {
		t.Errorf("期望 Name = \"updated\"，实际得到 %q", result.Name)
	}
}

// TestDelete 测试 DELETE 请求
func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("期望 DELETE 请求，实际得到 %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	if err := c.Delete("/", nil); err != nil {
		t.Fatalf("Delete 请求失败: %v", err)
	}
}

// TestHTTPError 测试 HTTP 错误响应
func TestClientHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"资源不存在"}`))
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	var result interface{}
	err := c.Get("/not-found", &result)

	if err == nil {
		t.Fatal("期望返回错误，但没有")
	}

	httpErr, ok := err.(*HTTPError)
	if !ok {
		t.Fatalf("期望 *HTTPError 类型，实际得到 %T", err)
	}

	if httpErr.StatusCode != http.StatusNotFound {
		t.Errorf("期望状态码 404，实际得到 %d", httpErr.StatusCode)
	}
}

// TestInitAndReset 测试全局单例的 Init 和 Reset
func TestInitAndReset(t *testing.T) {
	// 确保初始状态为 nil
	if DefaultClient != nil {
		t.Fatal("测试前期望 DefaultClient 为 nil")
	}

	Init(WithBaseURL("https://test.com"))
	if DefaultClient == nil {
		t.Fatal("Init 后期望 DefaultClient 不为 nil")
	}

	// 多次调用 Init 不会重新初始化
	oldClient := DefaultClient
	Init(WithBaseURL("https://other.com"))
	if DefaultClient != oldClient {
		t.Error("多次 Init 应返回同一个实例")
	}

	// 重置
	Reset()
	if DefaultClient != nil {
		t.Error("Reset 后期望 DefaultClient 为 nil")
	}

	// 重置后可重新初始化
	Init(WithBaseURL("https://new.com"))
	if DefaultClient == nil {
		t.Error("Reset 后 Init 应成功")
	}

	Reset()
}

// TestBuildURL 测试 URL 拼接
func TestBuildURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		path     string
		expected string
	}{
		{"无 baseURL", "", "/api/users", "/api/users"},
		{"有 baseURL", "https://api.example.com", "users", "https://api.example.com/users"},
		{"baseURL 带路径", "https://api.example.com/v1", "users", "https://api.example.com/v1/users"},
		{"路径带前导斜杠", "https://api.example.com", "/v1/users", "https://api.example.com/v1/users"},
		{"baseURL 带尾斜杠", "https://api.example.com/", "users", "https://api.example.com/users"},
		{"两边都带斜杠", "https://api.example.com/", "/v1/users", "https://api.example.com/v1/users"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New(WithBaseURL(tt.baseURL))
			got := c.buildURL(tt.path)
			if got != tt.expected {
				t.Errorf("buildURL(%q) = %q，期望 %q", tt.path, got, tt.expected)
			}
		})
	}
}

// TestGlobalHeaders 测试全局请求头被发送
func TestGlobalHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-App") != "hypersku" {
			t.Errorf("期望 X-App = \"hypersku\"，实际得到 %q", r.Header.Get("X-App"))
		}
		if r.Header.Get("Authorization") != "Bearer token123" {
			t.Errorf("期望 Authorization = \"Bearer token123\"，实际得到 %q", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	c := New(
		WithBaseURL(server.URL),
		WithHeader("X-App", "hypersku"),
		WithBearerToken("token123"),
	)

	var result map[string]interface{}
	if err := c.Get("/", &result); err != nil {
		t.Fatalf("请求失败: %v", err)
	}
}

// TestNilResult 测试不解析响应体的场景
func TestNilResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data":"should be ignored"}`))
	}))
	defer server.Close()

	c := New(WithBaseURL(server.URL))

	// result 为 nil 时不应报错
	if err := c.Get("/", nil); err != nil {
		t.Fatalf("result=nil 的请求失败: %v", err)
	}
}
