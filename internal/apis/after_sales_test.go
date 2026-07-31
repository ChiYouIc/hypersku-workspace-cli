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

func TestAfterSalesApi_Get1688AfterSalesGoods(t *testing.T) {
	tests := []struct {
		name         string
		thirdOrderId string
		refundId     string
		wantErr      bool
		wantLen      int
	}{
		{
			name:         "存在售后商品的工单",
			thirdOrderId: "3314076099641067093",
			refundId:     "TQ278495904826618234",
			wantErr:      false,
			wantLen:      1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := apis.NewAfterSalesApi()
			got, gotErr := af.Get1688AfterSalesGoods(tt.thirdOrderId, tt.refundId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get1688AfterSalesGoods() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get1688AfterSalesGoods() succeeded unexpectedly")
			}
			if len(*got) != tt.wantLen {
				t.Errorf("len %v, want %v", len(*got), tt.wantLen)
			}
		})
	}
}

func TestAfterSalesApi_Get1688AfterSalesDetail(t *testing.T) {
	tests := []struct {
		name     string
		refundId string
		wantErr  bool
		wantNil  bool
	}{
		{
			name:     "存在售后详情的退款ID",
			refundId: "TQ278495904826618234",
			wantErr:  false,
			wantNil:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := apis.NewAfterSalesApi()
			got, gotErr := af.Get1688AfterSalesDetail(tt.refundId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get1688AfterSalesDetail() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get1688AfterSalesDetail() succeeded unexpectedly")
			}
			if tt.wantNil && got != nil {
				t.Errorf("expected nil, got %v", got)
			}
			if !tt.wantNil && got == nil {
				t.Error("expected non-nil result")
			}
		})
	}
}

func TestAfterSalesApi_Get1688AfterSalesMessage(t *testing.T) {
	tests := []struct {
		name     string
		refundId string
		wantErr  bool
		wantLen  int
	}{
		{
			name:     "存在留言记录的退款ID",
			refundId: "TQ278495904826618234",
			wantErr:  false,
			wantLen:  1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			af := apis.NewAfterSalesApi()
			got, gotErr := af.Get1688AfterSalesMessage(tt.refundId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("Get1688AfterSalesMessage() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Get1688AfterSalesMessage() succeeded unexpectedly")
			}
			if len(*got) != tt.wantLen {
				t.Errorf("len %v, want %v", len(*got), tt.wantLen)
			}
		})
	}
}
