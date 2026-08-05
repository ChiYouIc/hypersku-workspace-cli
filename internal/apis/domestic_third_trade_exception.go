package apis

import (
	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

type DomesticThirdTradeExceptionQuery struct {
	HyperskuStatus        int    `json:"hyperskuStatus"`
	HyperskuSubStatusList []int  `json:"hyperskuSubStatusList"`
	BuyerID               string `json:"buyerId"`
	Page                  int    `json:"page"`
	Limit                 int    `json:"limit"`
}

type DomesticThirdTradeExceptionOrder struct {
	ID                 int                                    `json:"id"`               // monitorOrderId
	WarehouseName      string                                 `json:"warehouseName"`    // 仓库名称
	TotalAmount        float64                                `json:"totalAmount"`      // 订单总金额
	ThirdOrderID       string                                 `json:"thirdOrderId"`     // 第三方订单号（交易号、1688订单号）
	Quantity           int                                    `json:"quantity"`         // 商品数量
	PaymentTime        string                                 `json:"paymentTime"`      // 订单付款时间
	HasRefundRecords   bool                                   `json:"hasRefundRecords"` // 是否存在退款记录
	BuyerName          string                                 `json:"buyerName"`        // 销售员名称
	CrtTime            string                                 `json:"crtTime"`          // 创建时间
	SourceType         int                                    `json:"sourceType"`       // 采购来源
	CreatedTime        string                                 `json:"createdTime"`
	BuyerID            int                                    `json:"buyerId"`
	DeliverSeparately  bool                                   `json:"deliverSeparately"`
	Freight            float64                                `json:"freight"`
	HasCompletedOrder  bool                                   `json:"hasCompletedOrder"`
	HasCustomized      bool                                   `json:"hasCustomized"`
	IdentificationCode string                                 `json:"identificationCode"`
	ImgURL             string                                 `json:"imgUrl"`
	LogisticsList      []DomesticThirdTradeExceptionLogistics `json:"logisticsList"`
	PurchaseType       int                                    `json:"purchaseType"`
	RefundAmount       int                                    `json:"refundAmount"`
	SupplierID         int                                    `json:"supplierId"`
	SupplierLoginID    string                                 `json:"supplierLoginId"`
	TradeStatus        int                                    `json:"tradeStatus"`
	TradeStatusText    string                                 `json:"tradeStatusText"`
	UpdTime            string                                 `json:"updTime"`
	WarehouseID        int                                    `json:"warehouseId"`
}

type DomesticThirdTradeExceptionLogistics struct {
	ID                       int    `json:"id"`                       // monitorLogisticsId
	CrtTime                  string `json:"crtTime"`                  // 创建时间
	ThirdOrderID             int64  `json:"thirdOrderId"`             // 第三方订单号（交易号、1688订单号）
	TrackingNumber           string `json:"trackingNumber"`           // 物流单号
	LogisticsCompanyName     string `json:"logisticsCompanyName"`     // 物流公司
	HyperskuStatusChangeTime string `json:"hyperskuStatusChangeTime"` // 最近更新时间
	OrderTypes               string `json:"orderTypes"`               // 订单类型
	AiRiskLevel              int    `json:"aiRiskLevel"`
	AiSummary                string `json:"aiSummary"`
	AiSummaryUpdTime         string `json:"aiSummaryUpdTime"`
	DeliveredTime            string `json:"deliveredTime"`
	Customized               int    `json:"customized"`
	DeliveredType            string `json:"deliveredType"`
	ExtendHour               int    `json:"extendHour"`
	HasAiSummary             bool   `json:"hasAiSummary"`
	HyperskuStatus           int    `json:"hyperskuStatus"`
	HyperskuSubStatus        int    `json:"hyperskuSubStatus"`
	LatestLogisticsTrack     string `json:"latestLogisticsTrack"`
	LatestOperateID          int    `json:"latestOperateId"`
	LatestOperateName        string `json:"latestOperateName"`
	LatestOperateTime        string `json:"latestOperateTime"`
	LogisticsStatus          int    `json:"logisticsStatus"`
	LogisticsStatusText      string `json:"logisticsStatusText"`
	MonitorOrderID           int    `json:"monitorOrderId"`
	SignedTime               string `json:"signedTime"`
	SignedWarehouseID        int    `json:"signedWarehouseId"`
	SignedWarehouseName      string `json:"signedWarehouseName"`
	TrackingUpdateTime       string `json:"trackingUpdateTime"`
	UpdName                  string `json:"updName"`
	UpdTime                  string `json:"updTime"`
	UpdUser                  string `json:"updUser"`
	WarehouseOperator        string `json:"warehouseOperator"`
	WarehouseSignedTime      string `json:"warehouseSignedTime"`
	WarehouseStatus          int    `json:"warehouseStatus"`
	WarehouseStatusText      string `json:"warehouseStatusText"`
}

type DomesticThirdTradeExceptionMessage struct {
	ID                 int    `json:"id"`
	Remark             string `json:"remark"`             // remark
	MonitorLogisticsID int    `json:"monitorLogisticsId"` // monitorLogisticsId
	MonitorOrderID     int    `json:"monitorOrderId"`     // monitorOrderId
	Source             int    `json:"source"`             // 来源
	CrtTime            string `json:"crtTime"`            // 创建时间
	CrtHost            string `json:"crtHost"`
	CrtName            string `json:"crtName"`
	CrtUser            string `json:"crtUser"`
}

type DomesticThirdTradeExceptionApi struct {
	http httpclient.Client
}

func NewDomesticThirdTradeExceptionApi() *DomesticThirdTradeExceptionApi {
	return &DomesticThirdTradeExceptionApi{
		http: *httpclient.DefaultClient,
	}
}

// PageList 分页查询国内第三方交易异常订单（第三方平台采购单）
func (api *DomesticThirdTradeExceptionApi) PageList(query DomesticThirdTradeExceptionQuery) (*ApiPageResponse[DomesticThirdTradeExceptionOrder], error) {
	result := ApiPageResponse[DomesticThirdTradeExceptionOrder]{}

	if err := api.http.Post("/api/tenant/monitor/abnormal/purchase/logistics/table/list", query, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// MessageList 获取国内第三方交易异常订单留言列表
func (api *DomesticThirdTradeExceptionApi) MessageList(monitorOrderId, monitorLogisticsId string) (*[]DomesticThirdTradeExceptionMessage, error) {
	result := &ApiResponse[[]DomesticThirdTradeExceptionMessage]{}
	if err := api.http.Get("/api/tenant/monitor/abnormal/purchase/remark/"+monitorOrderId+"/"+monitorLogisticsId, result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
