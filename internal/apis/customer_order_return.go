package apis

import "github.com/hypersku/hypersku-cli/internal/httpclient"

type CustomerOrderReturnApi struct {
	http *httpclient.Client
}

var (
	WorkOrderType = map[int]string{
		0:  "退件",
		3:  "海外退件",
		9:  "自取件",
		15: "进出口查验",
		19: "服务商反馈异常件",
		20: "不合规申报协同件",
		21: "派送异常件",
		22: "库内查验退件协件",
		23: "取消签名服务协同件",
		24: "巴西待交税工单件",
		99: "国内退件",
		50: "海外退件",
		51: "转库存",
	}

	WorkOrderStatus = map[int]string{
		0:  "待处理",
		1:  "放行",
		5:  "重发",
		6:  "销毁",
		12: "重发中",
		13: "已重发",
		14: "转库存",
		15: "重派异常",
	}

	ReturnStatus = map[int]string{
		0: "关闭",
		1: "待处理",
		2: "完成",
		3: "处理中",
		6: "退款",
		7: "待审核",
	}
)

type CustomerOrderReturnInfoQuery struct {
	Page            int
	Limit           int
	CustomerOrderId string
	Status          int // 状态 9:全部、7:待审核、1:待处理、3:处理中、2:完成、0:关闭、6:退款
}

type CustomerOrderReturnInfo struct {
	CountryCode         string  `json:"countryCode"`         // 国家
	CustomerOrderID     string  `json:"customerOrderId"`     // HyperSKU 订单号
	CustomerOrderNumber string  `json:"customerOrderNumber"` // 客户单号
	Describing          string  `json:"describing"`          // 说明
	InstorageCreatedOn  string  `json:"instorageCreatedOn"`  // 入库时间
	Remark              string  `json:"remark"`              // 留言
	RemarkUpdateTime    string  `json:"remarkUpdateTime"`    // 留言时间
	Status              int     `json:"status"`              // 状态
	TrackingNumber      string  `json:"trackingNumber"`      // 快递单号
	WaybillNumber       string  `json:"waybillNumber"`       // 运单号
	Weight              float64 `json:"weight"`              // 重量
	WorkOrderState      int     `json:"workOrderState"`      // 工单状态
	WorkOrderType       int     `json:"workOrderType"`       // 工单类型
	CrtTime             string  `json:"crtTime"`             // 创建时间
	CanResend           bool    `json:"canResend"`
	CountryID           int     `json:"countryId"`
	CreateOn            string  `json:"createOn"`
	CustomerDiyOrderID  string  `json:"customerDiyOrderId"`
	CustomerID          int     `json:"customerId"`
	CustomerLogisticsID int     `json:"customerLogisticsId"`
	CustomerName        string  `json:"customerName"`
	Expired             bool    `json:"expired"`
	FinancialStatus     int     `json:"financialStatus"`
	Freight             float64 `json:"freight"`
	GoodsID             int64   `json:"goodsId"`
	HandleState         int     `json:"handleState"`
	ID                  int     `json:"id"`
	ImgURL              string  `json:"imgUrl"`
	OpFormat            int     `json:"opFormat"`
	OrderCode           string  `json:"orderCode"`
	OrderCodeName       string  `json:"orderCodeName"`
	OrderState          int     `json:"orderState"`
	OrderStateStr       string  `json:"orderStateStr"`
	PackageReturnMark   bool    `json:"packageReturnMark"`
	ProductCode         string  `json:"productCode"`
	PurOrderID          string  `json:"purOrderId"`
	SaleName            string  `json:"saleName"`
	ServiceName         string  `json:"serviceName"`
	SourceID            int     `json:"sourceId"`
	StatusStr           string  `json:"statusStr"`
	UpdTime             string  `json:"updTime"`
	WarehouseCode       string  `json:"warehouseCode"`
	WorkOrderCode       string  `json:"workOrderCode"`
	WorkOrderID         int     `json:"workOrderId"`
}

func NewCustomerOrderReturnApi() *CustomerOrderReturnApi {
	return &CustomerOrderReturnApi{
		http: httpclient.DefaultClient,
	}
}

// Page 分页查询客户订单退款数据
func (api *CustomerOrderReturnApi) Page(query CustomerOrderReturnInfoQuery) (*ApiPageResponse[CustomerOrderReturnInfo], error) {

	body := map[string]any{
		"page":                  query.Page,
		"limit":                 query.Limit,
		"orderStatus":           3,
		"returnOrderSearchType": 1,
		"status":                query.Status,
		"type":                  1,
		"workOrderType":         0,
	}

	if query.CustomerOrderId != "" {
		body["customerOrderId"] = query.CustomerOrderId
		body["returnOrderSearchText"] = query.CustomerOrderId
	}

	result := ApiPageResponse[CustomerOrderReturnInfo]{}

	if err := api.http.Post("/api/tenant/return/order/list", body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetOrderReturnInfo 获取订单退款信息
func (api *CustomerOrderReturnApi) GetOrderReturnInfo(customerOrderId string) (*CustomerOrderReturnInfo, error) {
	result, err := api.Page(CustomerOrderReturnInfoQuery{
		Limit:           10,
		Page:            1,
		CustomerOrderId: customerOrderId,
		Status:          9,
	})

	if err != nil {
		return nil, err
	}

	rows := result.Data.Rows
	if len(rows) == 0 {
		return nil, nil
	}

	return &result.Data.Rows[0], nil
}
