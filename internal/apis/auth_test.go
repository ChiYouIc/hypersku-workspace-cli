package apis

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

func TestMain(m *testing.M) {
	// 初始化 HTTP 客户端（测试环境下没有配置文件，使用默认配置）
	httpclient.Reset()
	httpclient.Init(
		httpclient.WithBaseURL("https://pur.hyperoms.com"),
		httpclient.WithHeader("authorization", "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJvd2VuLmNoaUBldGFpbGVyaHViLmNvbSIsInVzZXJJZCI6IjE4MCIsIm5hbWUiOiLmsaDlj4siLCJjaWQiOiIxMDAiLCJ0eXBlIjoiMSIsImV4cCI6MTc4OTAzMzQ2OX0.o4SJ5indFYZCnx4Nq976-97BjsJym7Xo0FPhJChDC32y19cczklmoHsPtHjDsjmF5vZwBtF2IrKvqzsV-AP0EyIHKJC3v8ufqqyOqNCnBjsfmWVNo1LpZx61GJA1sGDPOXcjcBMczKzlgByITt2BRc7rpPcYlC1bn0xcldS87ZY"),
	)

	code := m.Run()

	os.Exit(code)
}

func TestGetUserInfo(t *testing.T) {
	auth := NewAuthApi()
	token := "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJvd2VuLmNoaUBldGFpbGVyaHViLmNvbSIsInVzZXJJZCI6IjE4MCIsIm5hbWUiOiLmsaDlj4siLCJjaWQiOiIxMDAiLCJ0eXBlIjoiMSIsImV4cCI6MTc4OTAzMzQ2OX0.o4SJ5indFYZCnx4Nq976-97BjsJym7Xo0FPhJChDC32y19cczklmoHsPtHjDsjmF5vZwBtF2IrKvqzsV-AP0EyIHKJC3v8ufqqyOqNCnBjsfmWVNo1LpZx61GJA1sGDPOXcjcBMczKzlgByITt2BRc7rpPcYlC1bn0xcldS87ZY"
	info, err := auth.GetUserInfo(token)
	if err != nil {
		t.Errorf("获取用户信息失败: %v", err)
		return
	}
	s, _ := json.Marshal(info)
	t.Logf("用户信息: %s", string(s))
}

func TestGetUserInfoEmptyToken(t *testing.T) {
	auth := NewAuthApi()
	if _, err := auth.GetUserInfo(""); err == nil {
		t.Error("空 token 应返回错误")
	}
}
