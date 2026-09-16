package apis

import (
	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// PurchaseAfterSalesApi 采购售后查询客户端（只读）
type PurchaseAfterSalesApi struct {
	http httpclient.Client
}

func NewPurchaseAfterSalesApi() *PurchaseAfterSalesApi {
	return &PurchaseAfterSalesApi{
		http: *httpclient.DefaultClient,
	}
}

// PurchaseAfterSalesQuery 采购售后查询参数
type PurchaseAfterSalesQuery struct {
	Page         int   `json:"page"`
	Limit        int   `json:"limit"`
	ThirdOrderID int64 `json:"thirdOrderId,omitempty"` // 交易号
	TradeID      int64 `json:"tradeId,omitempty"`      // 交易ID
	OrderID      int64 `json:"orderId,omitempty"`      // 订单ID
	// 以下字段按用户要求不使用
	// Status, AbnormalType, HandlerId, ManagerId 等
}

// PurchaseAfterSalesTrade 采购售后工单
type PurchaseAfterSalesTrade struct {
	ID                 int     `json:"id"`
	ThirdOrderID       int64   `json:"thirdOrderId"`       // 交易号
	TradeID            int64   `json:"tradeId"`            // 交易ID
	OrderID            int64   `json:"orderId"`            // 订单ID
	CustomerOrderID    int64   `json:"customerOrderId"`    // 客户订单号
	WorkOrderCode      string  `json:"workOrderCode"`      // 工单号
	WorkOrderType      int     `json:"workOrderType"`      // 工单类型
	WorkOrderState     int     `json:"workOrderState"`     // 工单状态
	Status             int     `json:"status"`             // 处理状态
	StatusStr          string  `json:"statusStr"`          // 状态描述
	AbnormalType       int     `json:"abnormalType"`       // 异常类型
	AbnormalTypeStr    string  `json:"abnormalTypeStr"`    // 异常类型描述
	PurchaseSource     int     `json:"purchaseSource"`     // 采购来源
	IdentificationCode string  `json:"identificationCode"` // 识别码
	TrackingNumber     string  `json:"trackingNumber"`     // 快递单号
	WarehouseName      string  `json:"warehouseName"`      // 仓库名称
	CustomerUsername   string  `json:"customerUsername"`   // 客户账号
	HandlerName        string  `json:"handlerName"`        // 处理人
	ManagerName        string  `json:"managerName"`        // 客户经理
	Quantity           int     `json:"quantity"`           // 数量
	TotalAmount        float64 `json:"totalAmount"`        // 金额
	ImgUrl             string  `json:"imgUrl"`             // 图片
	CrtTime            string  `json:"crtTime"`            // 创建时间
	UpdTime            string  `json:"updTime"`            // 更新时间
	Remark             string  `json:"remark"`             // 备注
	Describe           string  `json:"describe"`           // 描述
	GoodsName          string  `json:"goodsName"`          // 商品名称
	GoodsSku           string  `json:"goodsSku"`           // SKU
	ThirdOrderIdStr    string  `json:"thirdOrderIdStr"`    // 交易号(字符串)
}

// PageList 按交易号查询采购售后工单列表
func (api *PurchaseAfterSalesApi) PageList(query PurchaseAfterSalesQuery) (*ApiPageResponse[PurchaseAfterSalesTrade], error) {
	result := ApiPageResponse[PurchaseAfterSalesTrade]{}
	if err := api.http.Post("/api/tenant/purchase/aftersales/list", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
