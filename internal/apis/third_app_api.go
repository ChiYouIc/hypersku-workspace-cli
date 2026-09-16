package apis

import (
	"encoding/json"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// ============ 1688 查询（透传 ehub-alipur） ============

// AliRefundDetailQuery 1688 退款单详情查询参数
type AliRefundDetailQuery struct {
	RefundID                 string `json:"refundId"`                 // 退款单业务主键 TQ+ID
	NeedTimeOutInfo          bool   `json:"needTimeOutInfo"`          // 是否需要退款单的超时信息
	NeedOrderRefundOperation bool   `json:"needOrderRefundOperation"` // 是否需要退款单伴随的所有退款操作信息
}

// AliRefundListQuery 1688 退款单列表查询参数（按交易号）
type AliRefundListQuery struct {
	OrderID   string `json:"orderId"`   // 交易号
	QueryType string `json:"queryType"` // 1：活动；3:退款成功（只支持退款中和退款成功）
}

// AliRefundOperationQuery 1688 退款操作记录查询参数
type AliRefundOperationQuery struct {
	RefundID string `json:"refundId"` // 退款单业务主键 TQ+ID
	PageNo   string `json:"pageNo"`   // 当前页号
	PageSize string `json:"pageSize"` // 页大小
}

// AliLogisticsInfo 1688 订单物流信息
type AliLogisticsInfoResult struct {
	Success      bool                `json:"success"`
	ErrorMessage string              `json:"errorMessage"`
	ErrorCode    string              `json:"errorCode"`
	Result       []AliLogisticsOrder `json:"result"`
}

// AliLogisticsOrder 1688 物流单
type AliLogisticsOrder struct {
	LogisticsBillNo      string                 `json:"logisticsBillNo"`      // 物流单号（运单号）
	LogisticsCompanyName string                 `json:"logisticsCompanyName"` // 物流公司
	LogisticsCompanyID   string                 `json:"logisticsCompanyId"`   // 物流公司ID
	OrderEntryIds        string                 `json:"orderEntryIds"`        // 订单号列表
	Status               string                 `json:"status"`               // 当前订单发货状态
	SendGoods            []AliLogisticsSendGood `json:"sendGoods"`            // 商品信息
}

// AliLogisticsSendGood 1688 物流单商品
type AliLogisticsSendGood struct {
	SendGoodsAmount string `json:"sendGoodsAmount"`
	SendGoodsWeight string `json:"sendGoodsWeight"`
	ItemID          string `json:"itemId"`
	EntryID         string `json:"entryId"`
	Name            string `json:"name"`
}

// AliLogisticsTraceInfoResult 1688 物流轨迹信息
type AliLogisticsTraceInfoResult struct {
	Success      bool                    `json:"success"`
	ErrorMessage string                  `json:"errorMessage"`
	ErrorCode    string                  `json:"errorCode"`
	Result       []AliLogisticsTraceInfo `json:"result"`
}

// AliLogisticsTraceInfo 1688 物流轨迹
type AliLogisticsTraceInfo struct {
	LogisticsID             string               `json:"logisticsId"`      // 物流编号
	OrderID                 int64                `json:"orderId"`          // 订单编号
	LogisticsBillNo         string               `json:"logisticsBillNo"`  // 物流单编号
	LogisticsCompany        string               `json:"logisticsCompany"` // 物流公司
	OrderEntryIds           string               `json:"orderEntryIds"`    // 订单号列表
	LogisticsSteps          []AliLogisticsStep   `json:"logisticsSteps"`   // 物流跟踪步骤
	WarehouseStatusStr      string               `json:"warehouseStatusStr"`
	AliStatusStr            string               `json:"aliStatusStr"`
	RefundRemark            string               `json:"refundRemark"`
	WarehouseLogisticsSteps []AliLogisticsStep   `json:"warehouseLogisticsSteps"` // 仓库物流跟踪
	ActionList              []AliWarehouseAction `json:"actionList"`              // 仓库物流轨迹
	IsSend                  int                  `json:"isSend"`                  // 仓库是否发走了 0 未发货 1 已发货 2 部分发货
	InstoreStatus           int                  `json:"instoreStatus"`           // 是否入库
	SignStatus              int                  `json:"signStatus"`              // 签收状态 1：已签收 0：未签收
	WarehouseName           string               `json:"warehouseName"`           // 仓库名称
	AdvanceWarehouse        int                  `json:"advanceWarehouse"`        // 是否提前入库
}

// AliLogisticsStep 物流轨迹步骤
type AliLogisticsStep struct {
	AcceptTime string `json:"acceptTime"` // 时间
	Remark     string `json:"remark"`     // 备注
}

// AliWarehouseAction 仓库操作
type AliWarehouseAction struct {
	ActionID   string `json:"actionId"`
	ActionTime string `json:"actionTime"`
	Content    string `json:"content"`
	Operator   string `json:"operator"`
}

// AliOrderDetailResult 1688 订单详情
type AliOrderDetailResult struct {
	Success      bool            `json:"success"`
	ErrorMessage string          `json:"errorMessage"`
	ErrorCode    string          `json:"errorCode"`
	Result       *AliOrderDetail `json:"result"`
}

// AliOrderDetail 1688 订单详情
type AliOrderDetail struct {
	BaseInfo     *AliOrderBaseInfo     `json:"baseInfo"`
	ProductItems []AliOrderProductItem `json:"productItems"`
}

// AliOrderBaseInfo 1688 订单基本信息
type AliOrderBaseInfo struct {
	ID                int64   `json:"id"`                // 订单ID
	TotalAmount       float64 `json:"totalAmount"`       // 总金额
	SumProductPayment float64 `json:"sumProductPayment"` // 产品总付款金额
	CouponFee         float64 `json:"couponFee"`         // 红包金额
	ShippingFee       float64 `json:"shippingFee"`       // 运费
	Refund            float64 `json:"refund"`            // 退款金额
	RefundID          string  `json:"refundId"`
	SellerLoginID     string  `json:"sellerLoginId"` // 卖家登录ID（旺旺ID）
	SellerID          string  `json:"sellerId"`      // 卖家主账号ID
	CreateTime        string  `json:"createTime"`    // 创建时间
	Status            string  `json:"status"`        // 交易状态
	RefundStatus      string  `json:"refundStatus"`  // 售中退款状态
	PayTime           string  `json:"payTime"`       // 付款时间
	BuyerLoginID      string  `json:"buyerLoginId"`
	ReceivingTime     string  `json:"receivingTime"`    // 收货时间
	AllDeliveredTime  string  `json:"allDeliveredTime"` // 完全发货时间
	CompleteTime      string  `json:"completeTime"`     // 完成时间
}

// AliOrderProductItem 1688 订单商品项
type AliOrderProductItem struct {
	Name            string       `json:"name"`       // 商品名称
	ProductID       int64        `json:"productID"`  // 产品ID
	ItemAmount      float64      `json:"itemAmount"` // 实付金额（元）
	Price           float64      `json:"price"`      // 原始单价（元）
	Quantity        float64      `json:"quantity"`   // 数量
	SkuID           int64        `json:"skuID"`
	SubItemIDString string       `json:"subItemIDString"` // 子订单号
	Unit            string       `json:"unit"`            // 单位
	ProductImgUrl   string       `json:"productImgUrl"`   // 商品图片
	Status          string       `json:"status"`          // 子订单状态
	StatusStr       string       `json:"statusStr"`       // 子订单状态描述
	RefundID        string       `json:"refundId"`        // 退款ID
	RefundStatus    string       `json:"refundStatus"`    // 退款状态
	SkuInfos        []AliSkuInfo `json:"skuInfos"`        // 销售属性
}

// AliSkuInfo 1688 SKU 属性
type AliSkuInfo struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// AliRefundDetailResult 1688 退款单详情
type AliRefundDetailResult struct {
	ErrorMessage    string                `json:"errorMessage"`
	ErrorCode       string                `json:"errorCode"`
	ExtErrorMessage string                `json:"extErrorMessage"`
	Result          *AliRefundDetailModel `json:"result"`
}

// AliRefundDetailModel 1688 退款单详情模型
type AliRefundDetailModel struct {
	OpOrderRefundModelDetail *AliRefundInfo `json:"opOrderRefundModelDetail"`
}

// AliRefundInfo 1688 退款单信息（金额单位：分）
type AliRefundInfo struct {
	ApplyCarriage    int64  `json:"applyCarriage"`    // 运费的申请退款金额（分）
	ApplyPayment     int64  `json:"applyPayment"`     // 买家申请退款金额（分）
	ApplyReason      string `json:"applyReason"`      // 申请原因
	CanRefundPayment int64  `json:"canRefundPayment"` // 最大能够退款金额（分）
	FreightBill      string `json:"freightBill"`      // 运单号
	GmtApply         string `json:"gmtApply"`         // 申请退款时间
	GmtCompleted     string `json:"gmtCompleted"`     // 完成时间
	GmtCreate        string `json:"gmtCreate"`        // 创建时间
	GmtModified      string `json:"gmtModified"`      // 修改时间
	GmtTimeOut       string `json:"gmtTimeOut"`       // 超时完成时间期限
	GoodsReceived    bool   `json:"goodsReceived"`    // 买家是否已收到货
	ID               int64  `json:"id"`               // 退款单编号
	OrderID          int64  `json:"orderId"`          // 退款单对应的订单编号
	ProductName      string `json:"productName"`      // 产品名称
	RefundID         string `json:"refundId"`         // 退款单逻辑主键
	RefundCarriage   int64  `json:"refundCarriage"`   // 运费的实际退款金额（分）
	RefundPayment    int64  `json:"refundPayment"`    // 实际退款金额（分）
	RejectReason     string `json:"rejectReason"`     // 卖家拒绝原因
	RejectTimes      int    `json:"rejectTimes"`      // 拒绝次数
	SellerLoginID    string `json:"sellerLoginId"`    // 卖家LoginId
	SellerRealName   string `json:"sellerRealName"`   // 卖家名称
	SellerMobile     string `json:"sellerMobile"`     // 卖家手机
	Status           string `json:"status"`           // 退款状态
	BuyerLoginID     string `json:"buyerLoginId"`     // 买家LoginId
	OnlyRefund       bool   `json:"onlyRefund"`       // 是否仅退款
	Percentage       int    `json:"percentage"`       // 百分比
}

// AliRefundListResult 1688 退款单列表
type AliRefundListResult struct {
	ErrorMessage    string              `json:"errorMessage"`
	ErrorCode       string              `json:"errorCode"`
	ExtErrorMessage string              `json:"extErrorMessage"`
	Result          *AliRefundListModel `json:"result"`
}

// AliRefundListModel 1688 退款单列表模型
type AliRefundListModel struct {
	OpOrderRefundModels []AliRefundInfo `json:"opOrderRefundModels"`
}

// AliRefundOperationListResult 1688 退款操作记录列表
type AliRefundOperationListResult struct {
	ErrorMessage    string                   `json:"errorMessage"`
	ErrorCode       string                   `json:"errorCode"`
	ExtErrorMessage string                   `json:"extErrorMessage"`
	Result          *AliRefundOperationModel `json:"result"`
}

// AliRefundOperationModel 1688 退款操作记录模型
type AliRefundOperationModel struct {
	OpOrderRefundOperationModels []AliRefundOperation `json:"opOrderRefundOperationModels"`
}

// AliRefundOperation 1688 退款操作记录
type AliRefundOperation struct {
	ID                  int64  `json:"id"`                  // 退款操作记录流水号
	AfterOperateStatus  string `json:"afterOperateStatus"`  // 操作后的退款状态
	BeforeOperateStatus string `json:"beforeOperateStatus"` // 操作前的退款状态
	Discription         string `json:"discription"`         // 描述、说明
	GmtCreate           string `json:"gmtCreate"`           // 创建时间
	MsgType             int    `json:"msgType"`             // 留言类型 3:小二留言 4:给买家的留言 5:给卖家的留言 7:普通留言
	OperateRemark       string `json:"operateRemark"`       // 操作备注
	OperatorLoginID     string `json:"operatorLoginId"`     // 操作者-loginID
	OperatorRoleId      int    `json:"operatorRoleId"`      // 操作者角色
	RejectReason        string `json:"rejectReason"`        // 卖家拒绝退款原因
	RefundAddress       string `json:"refundAddress"`       // 退货地址
	RefundID            string `json:"refundId"`            // 退款记录ID
}

// AliSellerMixConfig 1688 卖家混批设置
type AliSellerMixConfig struct {
	GeneralHunpi bool   `json:"generalHunpi"` // 是否普通混批
	GmtCreate    string `json:"gmtCreate"`    // 创建时间
	GmtModified  string `json:"gmtModified"`  // 修改时间
	MemberID     string `json:"memberId"`     // 卖家 memberId
	MixAmount    int    `json:"mixAmount"`    // 混批金额
	MixNumber    int    `json:"mixNumber"`    // 混批数量
}

// AliSubAccountList 1688 子账号列表
type AliSubAccountList struct {
	MainUserID     string             `json:"mainUserId"`
	MainLoginID    string             `json:"mainLoginId"`
	MainMemberID   string             `json:"mainMemberId"`
	SubAccountList []AliSimpleAccount `json:"subAccountList"`
}

// AliSimpleAccount 1688 子账号信息
type AliSimpleAccount struct {
	UserID   string `json:"userId"`
	LoginID  string `json:"loginId"`
	MemberID string `json:"memberId"`
}

// AliSupplierInfo 1688 供应商信息
type AliSupplierInfo struct {
	LoginID      string `json:"loginId"`      // 供应商登录 ID
	CategoryName string `json:"categoryName"` // 类目
	CompanyName  string `json:"companyName"`  // 公司名
	ShopURL      string `json:"shopUrl"`      // 店铺地址
	SupplierName string `json:"supplierName"` // 供应商名称
	KuaJingBao   bool   `json:"kuaJingBao"`   // 是否跨境宝
}

// AliGoodsBaseInfo 1688 产品基础信息（节选）
type AliGoodsBaseInfo struct {
	ID                     int64   `json:"id"`
	Name                   string  `json:"name"`
	EnName                 string  `json:"enName"`
	URL                    string  `json:"url"`
	SourceURL              string  `json:"sourceUrl"`
	OriginalPrice          string  `json:"originalPrice"`
	OriginalCurrencySymbol string  `json:"originalCurrencySymbol"`
	SalePrice              float64 `json:"salePrice"`
	CurrencySymbol         string  `json:"currencySymbol"`
	Weight                 float64 `json:"weight"`
	SuttleWeight           float64 `json:"suttleWeight"`
	SendGoodsAddress       string  `json:"sendGoodsAddress"`
	TotalStore             int     `json:"totalStore"`
	LimitNumber            int     `json:"limitNumber"`
	CategoryName           string  `json:"categoryName"`
	CategoryEnName         string  `json:"categoryEnName"`
	SupplierID             string  `json:"supplierId"`
	SupplierName           string  `json:"supplierName"`
	IsCombinedSku          bool    `json:"isCombinedSku"`
}

// ============ 店铺查询（Shopify 等第三方店铺） ============

// AppApiShopParam 店铺 API 查询参数
type AppApiShopParam struct {
	StoreID         int    `json:"storeId"` // 店铺 ID（必填）
	ProductID       int64  `json:"productId,omitempty"`
	VariantID       int64  `json:"variantId,omitempty"`
	OrderID         int64  `json:"orderId,omitempty"`
	Status          string `json:"status,omitempty"`          // 订单状态 open/closed/cancelled/any
	LocationID      int64  `json:"locationId,omitempty"`      // 库存位置
	InventoryItemID int64  `json:"inventoryItemId,omitempty"` // sku 库存项 ID
	Title           string `json:"title,omitempty"`
	PageNum         int    `json:"pageNum"`
	PageSize        int    `json:"pageSize"`
}

// ShopInventoryLevel 店铺库存位置关联
type ShopInventoryLevel struct {
	LocationID      int64 `json:"location_id"`       // 库存位置 ID
	InventoryItemID int64 `json:"inventory_item_id"` // 库存项 ID
	Available       int   `json:"available"`         // 可用库存数量
}

// ThirdAppApi 第三方 APP API 查询客户端（只读）
type ThirdAppApi struct {
	http httpclient.Client
}

func NewThirdAppApi() *ThirdAppApi {
	return &ThirdAppApi{
		http: *httpclient.DefaultClient,
	}
}

// GetAliLogisticsInfo 获取 1688 订单物流信息
func (api *ThirdAppApi) GetAliLogisticsInfo(orderId string) (*AliLogisticsInfoResult, error) {
	result := &AliLogisticsInfoResult{}
	params := url.Values{}
	params.Set("orderId", orderId)
	if err := api.http.Get("/api/tenant/third/app/api/1688/order/logistics/info?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliLogisticsTraceInfo 获取 1688 订单物流轨迹信息
func (api *ThirdAppApi) GetAliLogisticsTraceInfo(orderId string) (*AliLogisticsTraceInfoResult, error) {
	result := &AliLogisticsTraceInfoResult{}
	params := url.Values{}
	params.Set("orderId", orderId)
	if err := api.http.Get("/api/tenant/third/app/api/1688/order/logistics/trace/info?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliOrderDetail 获取 1688 订单详情（请求时间长）
func (api *ThirdAppApi) GetAliOrderDetail(orderId string) (*AliOrderDetailResult, error) {
	result := &AliOrderDetailResult{}
	params := url.Values{}
	params.Set("orderId", orderId)
	if err := api.http.Get("/api/tenant/third/app/api/1688/order/detail/info?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// QueryAliRefundDetail 根据退款单ID查询 1688 退款单详情
func (api *ThirdAppApi) QueryAliRefundDetail(query AliRefundDetailQuery) (*AliRefundDetailResult, error) {
	result := &AliRefundDetailResult{}
	if err := api.http.Post("/api/tenant/third/app/api/1688/refund/detail", query, result); err != nil {
		return nil, err
	}
	return result, nil
}

// QueryAliRefundList 根据交易号查询 1688 退款单列表
func (api *ThirdAppApi) QueryAliRefundList(query AliRefundListQuery) (*AliRefundListResult, error) {
	result := &AliRefundListResult{}
	if err := api.http.Post("/api/tenant/third/app/api/1688/refund/list", query, result); err != nil {
		return nil, err
	}
	return result, nil
}

// QueryAliRefundOperationList 查询 1688 退款单操作记录列表
func (api *ThirdAppApi) QueryAliRefundOperationList(query AliRefundOperationQuery) (*AliRefundOperationListResult, error) {
	result := &AliRefundOperationListResult{}
	if err := api.http.Post("/api/tenant/third/app/api/1688/refund/operation/list", query, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliSellerMixConfig 查询 1688 卖家混批设置
func (api *ThirdAppApi) GetAliSellerMixConfig(memberId, loginId string) (*ApiResponse[AliSellerMixConfig], error) {
	result := &ApiResponse[AliSellerMixConfig]{}
	params := url.Values{}
	if memberId != "" {
		params.Set("memberId", memberId)
	}
	if loginId != "" {
		params.Set("loginId", loginId)
	}
	if err := api.http.Get("/api/tenant/third/app/api/1688/seller/config/marketing/mix?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliSubAccountList 查询 1688 当前账号下的所有子账号列表
func (api *ThirdAppApi) GetAliSubAccountList() (*AliSubAccountList, error) {
	result := &AliSubAccountList{}
	if err := api.http.Get("/api/tenant/third/app/api/1688/auth/account/sub/list", result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliSupplierInfo 根据登录ID或店铺地址获取 1688 供应商信息
func (api *ThirdAppApi) GetAliSupplierInfo(loginId, domain string) (*AliSupplierInfo, error) {
	result := &AliSupplierInfo{}
	params := url.Values{}
	if loginId != "" {
		params.Set("loginId", loginId)
	}
	if domain != "" {
		params.Set("domain", domain)
	}
	if err := api.http.Get("/api/tenant/third/app/api/1688/supplier/info/account?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliProductInfoById 根据 1688 产品 ID 获取产品信息
func (api *ThirdAppApi) GetAliProductInfoById(productId string) (*ApiResponse[AliGoodsBaseInfo], error) {
	result := &ApiResponse[AliGoodsBaseInfo]{}
	params := url.Values{}
	params.Set("productId", productId)
	if err := api.http.Get("/api/tenant/third/app/api/1688/product/info/id?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetAliProductInfoByUrl 根据 1688 产品链接获取产品信息
func (api *ThirdAppApi) GetAliProductInfoByUrl(productUrl string, fetchSupplier bool) (*ApiResponse[AliGoodsBaseInfo], error) {
	result := &ApiResponse[AliGoodsBaseInfo]{}
	params := url.Values{}
	params.Set("url", productUrl)
	if fetchSupplier {
		params.Set("fetchSupplier", "true")
	}
	if err := api.http.Get("/api/tenant/third/app/api/1688/product/info/url?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopProduct 查询店铺产品
func (api *ThirdAppApi) GetShopProduct(param AppApiShopParam) (*ApiResponse[json.RawMessage], error) {
	result := &ApiResponse[json.RawMessage]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/product", param, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopProductImages 查询店铺产品图片
func (api *ThirdAppApi) GetShopProductImages(param AppApiShopParam) (*ApiResponse[json.RawMessage], error) {
	result := &ApiResponse[json.RawMessage]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/product/images", param, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopProductVariants 查询店铺产品变体
func (api *ThirdAppApi) GetShopProductVariants(param AppApiShopParam) (*ApiResponse[json.RawMessage], error) {
	result := &ApiResponse[json.RawMessage]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/product/variants", param, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopOrderList 查询店铺订单列表
func (api *ThirdAppApi) GetShopOrderList(param AppApiShopParam) (*ApiResponse[json.RawMessage], error) {
	result := &ApiResponse[json.RawMessage]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/order/list", param, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopLocations 查询店铺库存位置列表
func (api *ThirdAppApi) GetShopLocations(param AppApiShopParam) (*ApiResponse[json.RawMessage], error) {
	result := &ApiResponse[json.RawMessage]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/locations", param, result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetShopInventoryLevels 查询店铺库存位置和库存商品项的关联
func (api *ThirdAppApi) GetShopInventoryLevels(param AppApiShopParam) (*ApiResponse[[]ShopInventoryLevel], error) {
	result := &ApiResponse[[]ShopInventoryLevel]{}
	if err := api.http.Post("/api/tenant/third/app/api/shop/inventory/levels", param, result); err != nil {
		return nil, err
	}
	return result, nil
}
