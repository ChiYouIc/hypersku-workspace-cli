package apis

import "github.com/hypersku/hypersku-cli/internal/httpclient"

var WeeklyAdBudgetMap = map[int]string{
	1: "< 200 USD",
	2: "200-500 USD",
	3: "500-1000 USD",
	4: "> 1000 USD",
	5: "200-1000 USD",
	6: ">1000USD",
}

var OrderVolumeMap = map[int]string{
	0: "justStart",
	1: "1-100",
	2: "100-500",
	3: "500+",
}

var DurationTypeMap = map[int]string{
	1: "还没开始",
	2: "少于6个月",
	3: "6个月-1年",
	4: "多于1年",
	5: "< 1 years",
	6: "1 - 3 years",
	7: "> 3 years",
}

type CustomerInfoApi struct {
	http httpclient.Client
}

func NewCustomerInfoApi() *CustomerInfoApi {
	return &CustomerInfoApi{
		http: *httpclient.DefaultClient,
	}
}

type CustomerInfo struct {
	AllocatedTime  string `json:"allocatedTime"`  // 资料更新时间
	AllocatedType  int    `json:"allocatedType"`  // 类型
	CustomerSource string `json:"customerSource"` // 客户来源
	DurationType   int    `json:"durationType"`   // 从事 Dropshipping 时长类型
	Level          string `json:"level"`          // 等级
	OrderLevel     string `json:"orderLevel"`     // 订单等级
	OrderNum       int    `json:"orderNum"`       // 最近30天订单数量
	StoreNum       int    `json:"storeNum"`       // 店铺数量
	WeeklyAdBudget int    `json:"weeklyAdBudget"` // 周广告预算
	OrderVolume    int    `json:"orderVolume"`    // 每个月订单量
}

// GetCustomerInfo 获取客户信息
func (api *CustomerInfoApi) GetCustomerInfo(customerId string) (*CustomerInfo, error) {
	result := &ApiResponse[CustomerInfo]{}
	if err := api.http.Get("/api/customer/manager/customer/other/info/"+customerId, result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
