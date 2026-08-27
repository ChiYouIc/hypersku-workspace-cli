package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestCustomerInfoApi_GetCustomerInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerId string
		wantErr    bool
	}{
		{name: "获取客户信息", customerId: "1000249401", wantErr: false},
		{name: "客户不存在", customerId: "9999999999", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerInfoApi()
			got, gotErr := api.GetCustomerInfo(tt.customerId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerInfo() succeeded unexpectedly")
			}
			if got == nil {
				t.Errorf("GetCustomerInfo() return nil")
				return
			}
			if got.CustomerID == 0 {
				t.Errorf("GetCustomerInfo() 返回的 customerId 为 0，疑似未查到客户")
			}
			t.Logf("客户档案: customerId=%d hasOrder=%v stores=%s engagedTime=%d weeklyAdBudget=%d orderVolume=%d niche=%d serviceInterest=%s countryName=%s signedUpAt=%s level=%s orderLevel=%s tag=%d",
				got.CustomerID, got.HasOrder, got.Stores, got.EngagedTime, got.WeeklyAdBudget, got.OrderVolume, got.Niche, got.ServiceInterest, got.CountryName, got.SignedUpAt, got.Level, got.OrderLevel, got.Tag)
		})
	}
}
