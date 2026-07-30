package httpclient

import (
	"strings"
	"testing"
)

func TestHTTPError(t *testing.T) {
	err := &HTTPError{
		StatusCode: 404,
		Body:       `{"message":"not found"}`,
	}

	errStr := err.Error()

	if !strings.Contains(errStr, "404") {
		t.Errorf("期望错误信息包含状态码 404，实际得到 %q", errStr)
	}

	if !strings.Contains(errStr, "not found") {
		t.Errorf("期望错误信息包含响应体，实际得到 %q", errStr)
	}
}

func TestHTTPError_EmptyBody(t *testing.T) {
	err := &HTTPError{
		StatusCode: 500,
		Body:       "",
	}

	errStr := err.Error()

	if !strings.Contains(errStr, "500") {
		t.Errorf("期望错误信息包含状态码 500，实际得到 %q", errStr)
	}
}
