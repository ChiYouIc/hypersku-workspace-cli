package apis

import (
	"encoding/json"
	"testing"
)

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
