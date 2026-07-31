package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestAfterSalesApi_Get1688AfterSales(t *testing.T) {
	tests := []struct {
		name         string
		thirdOrderId string
		want         *[]apis.AfterSales1688Info
		wantLen      int
		wantErr      bool
	}{
		{name: "存在一个售后的工单的交易号", thirdOrderId: "3314076099641067093", wantErr: false, wantLen: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := apis.NewAfterSalesApi()
			got, gotErr := af.Get1688AfterSales(tt.thirdOrderId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get1688AfterSales() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get1688AfterSales() succeeded unexpectedly")
			}

			if len(*got) != tt.wantLen {
				t.Errorf("len %v, want %v", len(*got), tt.wantLen)
			}
		})
	}
}
