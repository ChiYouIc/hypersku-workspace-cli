package apis_test

import (
	"testing"

	"github.com/hypersku/hypersku-cli/internal/apis"
)

func TestCustomerOrderApi_GetOrderInfo(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name    string                  // 测试用例名称
		orderId string                  // 测试参数
		want    *apis.CustomerOrderInfo // 期望值
		wantErr bool                    // 期望异常
	}{
		{
			name:    "获取订单信息",
			orderId: "1152751001053192198",
			wantErr: false,
		},
		{
			name:    "获取订单不存在",
			orderId: "11527510010531921982",
			wantErr: true,
		},
	}
	// 遍历测试用例表
	for _, tt := range tests {

		// 执行测试用例
		t.Run(tt.name, func(t *testing.T) {
			o := apis.NewCustomerOrderApi()
			got, gotErr := o.GetOrderInfo(tt.orderId)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("GetOrderInfo() 失败: %v", gotErr)
				}
				return
			}

			// 期望错误
			if tt.wantErr {
				t.Fatal("GetOrderInfo() succeeded unexpectedly")
			}

			// 期望值比较
			if got.ID != tt.orderId {
				t.Errorf("GetOrderInfo() = %v, 期望 %v", got, tt.want)
			}
		})
	}
}
