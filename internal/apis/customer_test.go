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
			t.Logf("客户信息: level=%s orderNum=%d storeNum=%d durationType=%s weeklyAdBudget=%s",
				got.Level, got.OrderNum, got.StoreNum,
				apis.DurationTypeMap[got.DurationType], apis.WeeklyAdBudgetMap[got.WeeklyAdBudget])
		})
	}
}
