package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestDomesticThirdTradeExceptionApi_PageList(t *testing.T) {
	tests := []struct {
		name    string
		query   apis.DomesticThirdTradeExceptionQuery
		wantErr bool
	}{
		{name: "获取国内第三方交易丢包裹异常订单列表", query: apis.DomesticThirdTradeExceptionQuery{
			Page:                  1,
			Limit:                 10,
			HyperskuStatus:        9,
			HyperskuSubStatusList: []int{1, 2},
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewDomesticThirdTradeExceptionApi()
			got, gotErr := api.PageList(tt.query)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("PageList() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("PageList() succeeded unexpectedly")
			}

			if got == nil || len(got.Data.Rows) == 0 {
				t.Errorf("PageList() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}

func TestDomesticThirdTradeExceptionApi_MessageList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		monitorOrderId     string
		monitorLogisticsId string
		wantErr            bool
	}{
		{name: "获取国内第三方交易异常订单留言列表", monitorOrderId: "194998", monitorLogisticsId: "89535", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewDomesticThirdTradeExceptionApi()
			got, gotErr := api.MessageList(tt.monitorOrderId, tt.monitorLogisticsId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("MessageList() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("MessageList() succeeded unexpectedly")
			}

			if got == nil || len(*got) == 0 {
				t.Errorf("MessageList() = %v, want %v", got, tt.wantErr)
			}
		})
	}
}
