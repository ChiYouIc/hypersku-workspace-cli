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
		httpclient.WithHeader("authorization", "eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJvd2VuLmNoaUBldGFpbGVyaHViLmNvbSIsInVzZXJJZCI6IjE4MCIsIm5hbWUiOiLmsaDlj4siLCJjaWQiOiIxMDAiLCJ0eXBlIjoiMSIsImV4cCI6MTc4ODI1NDg3MX0.jzOOntZZCiUFAlSwyZdBlbUqRYbzDUzJlL3bBcf6DpwBmTu_w370DhEyv7gt88uaVt8Fb7COgTVs90UsbBKO4fl60AWh5A_gNVSe5v2WOyFRXfAXDMH2Rt83FTnuSO-jTQEOXhk3pFDlRScVVUnCKz66eh-CgTb0SiGFeAHCYmM"),
	)

	code := m.Run()

	os.Exit(code)
}

func TestPageList(t *testing.T) {
	purchase := NewPurchaseApi()
	res, err := purchase.PageList(QueryPage{
		Page:          1,
		Limit:         10,
		StartDate:     "2026-04-29",
		EndDate:       "2026-07-30",
		TransactionNo: "3314226687880021084",
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
