package apis

import (
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

type AfterSalesApi struct {
	http httpclient.Client
}

func NewAfterSalesApi() *AfterSalesApi {
	return &AfterSalesApi{
		http: *httpclient.DefaultClient,
	}
}

type AfterSales1688Info struct {
	ID                   int               `json:"id"`
	ActionDesc           string            `json:"actionDesc"`         // 当前操作
	AllRefundAmount      float64           `json:"allRefundAmount"`    // 退款金额
	BackFreight          float64           `json:"backFreight"`        // 退货运费
	CrtTime              string            `json:"crtTime"`            // 创建时间
	DisputeRequestDesc   string            `json:"disputeRequestDesc"` // 纠纷请求类型
	DisputeType          string            `json:"disputeType"`        // 退款类型（售中/售后）
	Freight              float64           `json:"freight"`            // 退款运费
	PurchaseSource       int               `json:"purchaseSource"`     // 采购平台来源
	RefundAmount         float64           `json:"refundAmount"`       // 退款金额
	RefundID             string            `json:"refundId"`           // 退款id
	ResponsibleParty     int               `json:"responsibleParty"`   // 责任方
	TimeDesc             string            `json:"timeDesc"`           //剩余处理时间
	Status               string            `json:"status"`             // 退款状态
	ThirdOrderID         string            `json:"thirdOrderId"`       // 1688订单号
	Action               string            `json:"action"`
	DisputeRequest       string            `json:"disputeRequest"`
	HasMonitor           bool              `json:"hasMonitor"`
	ImgURL               string            `json:"imgUrl"`
	OgID                 int               `json:"ogId"`
	OrderID              int64             `json:"orderId"`
	Percentage           int               `json:"percentage"`
	RefundSource         int               `json:"refundSource"`
	RefundSourceDesc     string            `json:"refundSourceDesc"`
	RefundStatus         int               `json:"refundStatus"`
	RefundUserID         int               `json:"refundUserId"`
	RefundUserName       string            `json:"refundUserName"`
	ResponsiblePartyText string            `json:"responsiblePartyText"`
	StatusDesc           string            `json:"statusDesc"`
	StatusShowDesc       string            `json:"statusShowDesc"`
	SubOrderIds          string            `json:"subOrderIds"`
	SyncWarehouseStatus  int               `json:"syncWarehouseStatus"`
	TenantOrdersGoodsVo  []AfterSalesGoods `json:"tenantOrdersGoodsVo"`
	TriggerType          int               `json:"triggerType"`
	UpdTime              string            `json:"updTime"`
}

// 售后商品
type AfterSalesGoods struct {
	ID        string `json:"id"`
	GoodsName string `json:"goodsName"` // 商品名称
	GoodsAttr string `json:"goodsAttr"` // 属性
	GoodsID   int64  `json:"goodsId"`   // sku
	ImgURL    string `json:"imgUrl"`    // 图片
	Num       int    `json:"num"`       // 数量
	OrderID   string `json:"orderId"`   // 订单id
}

// 售后消息
type AfterSalesMessage struct {
	OperatorLoginID string `json:"operatorLoginId"` // 留言方
	OperateRemark   string `json:"operateRemark"`   // 操作留言
	Discription     string `json:"discription"`     // 留言
	GmtCreate       string `json:"gmtCreate"`       // 留言时间
	GmtModified     string `json:"gmtModified"`
	ID              int64  `json:"id"`
	MessageStatus   int    `json:"messageStatus"`
	MsgType         int    `json:"msgType"`
	OperatorRoleID  int    `json:"operatorRoleId"`
	RefundID        string `json:"refundId"`
}

// 售后详情
type AfterSaleDetail struct {
	AllRefundAmount float64 `json:"allRefundAmount"` // 退款总金额
	ApplyTime       string  `json:"applyTime"`       // 申请退款时间
	CompanyName     string  `json:"companyName"`     // 卖家
	Freight         float64 `json:"freight"`         // 退款运费
	Mobile          string  `json:"mobile"`          // 手机
	Phone           string  `json:"phone"`           // 电话
	RefundAmount    float64 `json:"refundAmount"`    // 商品退款金额
	RefundDesc      string  `json:"refundDesc"`      // 退款说明
	RefundReason    string  `json:"refundReason"`    // 退款原因
	RefundType      string  `json:"refundType"`      // 退款服务
	SellerLoginID   string  `json:"sellerLoginId"`   // 会员登录名
	StatusDesc      string  `json:"statusDesc"`      // 退款状态
	ThirdOrderID    string  `json:"thirdOrderId"`    // 交易号
	RefundID        string  `json:"refundId"`        // 退款id
	Status          string  `json:"status"`
	Email           string  `json:"email"`
	Name            string  `json:"name"`
}

// Get1688AfterSales 查询1688售后工单
func (af *AfterSalesApi) Get1688AfterSales(thirdOrderId string) (*[]AfterSales1688Info, error) {
	params := url.Values{}
	params.Set("page", "1")
	params.Set("limit", "20")
	params.Set("thirdOrderId", thirdOrderId)

	apiResponse := &ApiPageResponse[AfterSales1688Info]{}
	if err := af.http.Get("/api/tenant/thirdOrders/refund/refundList?"+params.Encode(), apiResponse); err != nil {
		return nil, err
	}

	return &apiResponse.Data.Rows, nil
}

// Get1688RefundGoods 查询1688退款商品项
func (af *AfterSalesApi) Get1688AfterSalesGoods(thirdOrderId string, refundId string) (*[]AfterSalesGoods, error) {
	params := url.Values{}
	params.Set("refundId", refundId)
	params.Set("thirdOrderId", thirdOrderId)

	refundGoodsList := &[]AfterSalesGoods{}
	if err := af.http.Post("/api/tenant/thirdAliRefund/refundList/findByreFundId?"+params.Encode(), nil, refundGoodsList); err != nil {
		return nil, err
	}

	return refundGoodsList, nil
}

// Get1688RefundDetail 查询1688退款详情
func (af *AfterSalesApi) Get1688AfterSalesDetail(refundId string) (*AfterSaleDetail, error) {
	detail := &AfterSaleDetail{}
	if err := af.http.Post("/api/tenant/thirdAliRefund/getDetailByRefundId/"+refundId, nil, detail); err != nil {
		return nil, err
	}

	if detail.RefundID == "" {
		return nil, nil
	}

	return detail, nil
}

// Get1688RefundMessage 查询1688售后留言
func (af *AfterSalesApi) Get1688AfterSalesMessage(refundId string) (*[]AfterSalesMessage, error) {

	messageList := &[]AfterSalesMessage{}
	if err := af.http.Post("/api/tenant/thirdAliRefund/refundList/findMessageByRefundId/"+refundId, nil, messageList); err != nil {
		return nil, err
	}

	return messageList, nil
}
