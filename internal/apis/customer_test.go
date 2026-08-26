package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestCustomerInfoApi_GetCustomerExtendInfo(t *testing.T) {
	tests := []struct {
		name       string
		customerId string
		wantErr    bool
	}{
		{name: "获取客户扩展信息", customerId: "1000249401", wantErr: false},
		{name: "客户不存在", customerId: "9999999999", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerInfoApi()
			got, gotErr := api.GetCustomerExtendInfo(tt.customerId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerExtendInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerExtendInfo() succeeded unexpectedly")
			}
			if got == nil {
				t.Errorf("GetCustomerExtendInfo() return nil")
				return
			}
			t.Logf("客户扩展信息: level=%s orderLevel=%s orderNum=%d storeNum=%d durationType=%s weeklyAdBudget=%s orderVolume=%s customerSource=%s allocatedTime=%s",
				got.Level, got.OrderLevel, got.OrderNum, got.StoreNum,
				apis.DurationTypeMap[got.DurationType], apis.WeeklyAdBudgetMap[got.WeeklyAdBudget],
				apis.OrderVolumeMap[got.OrderVolume], got.CustomerSource, got.AllocatedTime)
		})
	}
}

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
			if got.ID == 0 {
				t.Errorf("GetCustomerInfo() 返回的 id 为 0，疑似未查到客户")
			}
			t.Logf("客户档案: id=%d username=%s email=%s totalOrder=%d totalAmount=%.2f status=%d",
				got.ID, got.Username, got.Email, got.TotalOrder, got.TotalAmount, got.Status)
		})
	}
}
