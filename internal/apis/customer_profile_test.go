package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestCustomerProfileApi_GetCustomerProfileOrderCount(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerId string
		startDate  string
		endDate    string
		wantErr    bool
		nonNil     bool
	}{
		{name: "不为nil", customerId: "1000249401", startDate: "2026-07-28 00:00:00", endDate: "2026-08-27 00:00:00", wantErr: false, nonNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerProfileApi()
			got, gotErr := api.GetCustomerProfileOrderCount(tt.customerId, tt.startDate, tt.endDate)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerProfileOrderCount() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerProfileOrderCount() succeeded unexpectedly")
			}
			if tt.nonNil && got == nil {
				t.Errorf("GetCustomerProfileOrderCount() return nil")
			}
			if tt.nonNil && got != nil {
				t.Logf("订单统计: total=%d avg=%d max=%d min=%d", got.Total, got.Avg, got.Max, got.Min)
			}
		})
	}
}

func TestCustomerProfileApi_GetCustomerProfileReceiveVisitorsInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerId string
		startDate  string
		endDate    string
		page       int
		limit      int
		wantErr    bool
		nonNil     bool
	}{
		{name: "不为nil", customerId: "1000249401", startDate: "2026-07-28 00:00:00", endDate: "2026-08-27 00:00:00", page: 1, limit: 90, wantErr: false, nonNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerProfileApi()
			got, gotErr := api.GetCustomerProfileReceiveVisitorsInfo(tt.customerId, tt.startDate, tt.endDate, tt.page, tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerProfileReceiveVisitorsInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerProfileReceiveVisitorsInfo() succeeded unexpectedly")
			}
			if tt.nonNil && (got == nil || got.Total == 0) {
				t.Errorf("GetCustomerProfileReceiveVisitorsInfo() return nil")
			}
			if tt.nonNil && got != nil {
				t.Logf("日订单数量: total=%d rows=%d", got.Total, len(got.Rows))
			}
		})
	}
}

func TestCustomerProfileApi_GetCustomerProfileTransactionCount(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerId string
		startDate  string
		endDate    string
		wantErr    bool
		nonNil     bool
	}{
		{name: "不为nil", customerId: "1000249401", startDate: "2026-07-28 00:00:00", endDate: "2026-08-27 00:00:00", wantErr: false, nonNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerProfileApi()
			got, gotErr := api.GetCustomerProfileTransactionCount(tt.customerId, tt.startDate, tt.endDate)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerProfileTransactionCount() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerProfileTransactionCount() succeeded unexpectedly")
			}
			if tt.nonNil && got == nil {
				t.Errorf("GetCustomerProfileTransactionCount() return nil")
			}
			if tt.nonNil && got != nil {
				t.Logf("交易统计: tranAmount=%.2f allOrderNum=%d customerPrice=%.2f", got.TranAmount, got.AllOrderNum, got.CustomerPrice)
			}
		})
	}
}

func TestCustomerProfileApi_GetCustomerProfileTransactionBillRecords(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerId string
		startDate  string
		endDate    string
		page       int
		limit      int
		wantErr    bool
		nonNil     bool
	}{
		{name: "不为nil", customerId: "1000249401", startDate: "2026-07-28 00:00:00", endDate: "2026-08-27 00:00:00", page: 1, limit: 10, wantErr: false, nonNil: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerProfileApi()
			got, gotErr := api.GetCustomerProfileTransactionBillRecords(tt.customerId, tt.startDate, tt.endDate, tt.page, tt.limit)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetCustomerProfileTransactionBillRecords() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetCustomerProfileTransactionBillRecords() succeeded unexpectedly")
			}
			if tt.nonNil && got == nil {
				t.Errorf("GetCustomerProfileTransactionBillRecords() return nil")
			}
			if tt.nonNil && got != nil {
				t.Logf("交易流水: total=%d rows=%d", got.Total, len(got.Rows))
			}
		})
	}
}
