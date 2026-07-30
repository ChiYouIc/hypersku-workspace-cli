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
		httpclient.WithHeader("authorization", "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJvd2VuLmNoaUBldGFpbGVyaHViLmNvbSIsInVzZXJJZCI6IjE4MCIsIm5hbWUiOiLmsaDlj4siLCJjaWQiOiIxMDAiLCJ0eXBlIjoiMSIsImV4cCI6MTc4NTkyNDI4Nn0.Ns-vkEYRWSVpxEbE0--9JizIXnh7C9fJ_t94ZkwQyZPPszM3bN6yfdzw9UmPaniVGzmZGCsBD3m_Vwi90tMGVDT7XejF7j2ZqHzYZci9zqL8mfcZZ12TkN6cGukvgucR5Jj3VSx9nO-YCthB0rxS3Og2wa9u3TGZ7BGY-C1z1rw"),
	)

	code := m.Run()

	os.Exit(code)
}

func TestPageList(t *testing.T) {
	purchase := NewPurchaseApi()
	res, err := purchase.PageList(QueryPage{
		Page:      1,
		Limit:     10,
		StartDate: "2026-04-29",
		EndDate:   "2026-07-30",
	})
	if err != nil {
		t.Errorf("查询订单列表失败: %v", err)
		return
	}

	s, _ := json.Marshal(res)
	t.Logf("查询结果: %s", string(s))
}

func TestGetOrderInfo(t *testing.T) {
	purchase := NewPurchaseApi()
	res, err := purchase.GetOrderInfo("1153391630170980352")
	if err != nil {
		t.Errorf("查询订单列表失败: %v", err)
		return
	}

	s, _ := json.Marshal(res)
	t.Logf("查询结果: %s", string(s))
}

func TestGetPurchaseLog(t *testing.T) {

	purchase := NewPurchaseApi()
	res, err := purchase.GetPurchaseLog("1153686099957121024")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	s, _ := json.Marshal(res)
	t.Logf("查询结果: %s", string(s))
}

func TestGetInternationalLogistics(t *testing.T) {
	purchase := NewPurchaseApi()
	res, err := purchase.GetInternationalLogistics("1153686099957121024")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	s, _ := json.Marshal(res)
	t.Logf("查询结果: %s", string(s))
}
