package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestCustomerOrderReturnApi_Page(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		query   apis.CustomerOrderReturnInfoQuery
		want    *apis.ApiPageResponse[apis.CustomerOrderReturnInfo]
		wantErr bool
	}{
		{
			name: "普通分页查询",
			query: apis.CustomerOrderReturnInfoQuery{
				Page:  1,
				Limit: 10,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerOrderReturnApi()
			got, gotErr := api.Page(tt.query)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Page() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Page() succeeded unexpectedly")
			}

			if got == nil {
				t.Errorf("Page() return is nil")
			}

			if len(got.Data.Rows) != tt.query.Limit {
				t.Errorf("Page() rows = 10")
			}

		})
	}
}

func TestCustomerOrderReturnApi_GetOrderReturnInfo(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		customerOrderId string
		want            *apis.CustomerOrderReturnInfo
		wantErr         bool
	}{
		{
			name:            "根据客户单号查询",
			customerOrderId: "1152251173282013185",
			want:            nil,
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := apis.NewCustomerOrderReturnApi()
			got, gotErr := api.GetOrderReturnInfo(tt.customerOrderId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetOrderReturnInfo() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("GetOrderReturnInfo() succeeded unexpectedly")
			}

			if got == nil {
				t.Errorf("GetOrderReturnInfo() is nil")
			}
		})
	}
}
