package apis

import (
	"fmt"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// ============ 供应商管理 ============

// SupplierInfo 供应商信息
type SupplierInfo struct {
	ID                int    `json:"id"`
	Type              int    `json:"type"`              // 类型：1=1688, 17=京东, 18=淘宝, 19=天猫, 25=拼多多, 99=HyperFrom
	PlatformMemberID  string `json:"platformMemberId"`  // 平台供应商会员ID
	PlatformLoginID   string `json:"platformLoginId"`   // 平台供应商的登录ID
	Name              string `json:"name"`              // 供应商名称
	URL               string `json:"url"`               // 供应商链接
	Email             string `json:"email"`             // 邮箱
	ContactsName      string `json:"contactsName"`      // 联系人姓名
	ContactsPhone     string `json:"contactsPhone"`     // 联系人手机号
	WangWangAccount   string `json:"wangWangAccount"`   // 旺旺号
	ArrivalPeriod     int    `json:"arrivalPeriod"`     // 到货周期(天)
	PaymentWay        int    `json:"paymentWay"`        // 付款方式
	ProceedsWay       int    `json:"proceedsWay"`       // 收款方式
	State             int    `json:"state"`             // 状态：启用、停用、删除
	CountryName       string `json:"countryName"`       // 国家
	SecondRegionName  string `json:"secondRegionName"`  // 省份
	ThirdRegionName   string `json:"thirdRegionName"`   // 城市
	AddressDetail     string `json:"addressDetail"`     // 具体地址
	PurchaserName     string `json:"purchaserName"`     // 采购员
	AssociateCategory string `json:"associateCategory"` // 关联类目
	Remark            string `json:"remark"`            // 备注
	Wechat            string `json:"wechat"`            // 微信
}

// SupplierPageQuery 供应商分页查询参数
type SupplierPageQuery struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Name   string `json:"name,omitempty"`
	States string `json:"states,omitempty"`
}

// SupplierPageResult 供应商分页查询结果
type SupplierPageResult struct {
	Rows  []map[string]interface{} `json:"rows"`
	Total int64                    `json:"total"`
}

// SupplierRankingQuery 供应商排行查询参数
type SupplierRankingQuery struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Name      string `json:"name,omitempty"`
	SearchDay int    `json:"searchDay,omitempty"` // 近几天
	Sort      string `json:"sort,omitempty"`      // 排序字段
	Type      string `json:"type,omitempty"`      // 类型
}

// SupplierRankingInfo 供应商排行信息
type SupplierRankingInfo struct {
	SupplyID      int     `json:"supplyId"`
	Name          string  `json:"name"`             // 供应商名称
	SecondRegion  string  `json:"secondRegionName"` // 城市
	URL           string  `json:"url"`              // 供应商链接
	PurTotal      int     `json:"purTotal"`         // 采购数量
	SpuNum        int     `json:"spuNum"`           // SPU数量
	SkuNum        int     `json:"skuNum"`           // SKU数量
	CategoryName  string  `json:"categoryName"`     // 类目
	AvgPurPrice   float64 `json:"avgPurPrice"`      // 采购平均价
	AddService    string  `json:"addService"`       // 是否为增值服务
	CustomerNum   int     `json:"customerNum"`      // 客户数
	ShippingFee   float64 `json:"shippingFee"`      // 运费
	ShipmentsNum  int     `json:"shipmentsNum"`     // 已发货个数
	UnshippedNum  int     `json:"unshippedNum"`     // 未发货个数
	AvgTime       float64 `json:"avgTime"`          // 平均时效
	OrdNum        int     `json:"ordNum"`           // 订单数量
	TotalAmount   float64 `json:"totalAmount"`      // 交易总金额
	SalesAmount   float64 `json:"salesAmount"`      // 销售总金额
	ReturnAmount  float64 `json:"returnAmount"`     // 退款总金额
	DeliveryRatio float64 `json:"deliveryRatio"`    // 发货比例
	AbnormalRatio float64 `json:"abnormalRatio"`    // 异常比例
	ProfitMargin  float64 `json:"profitMargin"`     // 利润率
	CurrencyCode  string  `json:"currencyCode"`     // 币种
	CrtTime       string  `json:"crtTime"`          // 创建时间
}

// SupplierAfterSales 供应商售后信息
type SupplierAfterSales struct {
	SupplyID         int    `json:"supplyId"`
	StatisticsDate   string `json:"statisticsDate"` // 统计时间
	SupplierName     string `json:"supplierName"`
	SuccessNum       int    `json:"successNum"`       // 成功数量
	OverTimeCloseNum int    `json:"overTimeCloseNum"` // 超时关闭数量
	CloseNum         int    `json:"closeNum"`         // 关闭数量
	TradeTotalNum    int    `json:"tradeTotalNum"`    // 交易全部数量
	ReturnSuccessNum int    `json:"returnSuccessNum"` // 退款/退货成功数量
	ReturnFailureNum int    `json:"returnFailureNum"` // 退款/退货失败数量
	ReturnTotalNum   int    `json:"returnTotalNum"`   // 退款/退货总数
}

// SupplierFulfillment 供应商履约信息
type SupplierFulfillment struct {
	SupplyID      int     `json:"supplyId"`
	SupplierName  string  `json:"supplierName"`
	ShipmentsNum  int     `json:"shipmentsNum"`  // 已发货个数
	UnshippedNum  int     `json:"unshippedNum"`  // 未发货个数
	AvgTime       float64 `json:"avgTime"`       // 平均时效
	DeliveryRatio float64 `json:"deliveryRatio"` // 发货比例
	AbnormalRatio float64 `json:"abnormalRatio"` // 异常比例
	WdNum         int     `json:"wdNum"`         // 错发个数
	LessNum       int     `json:"lessNum"`       // 少发个数
	TotalAmount   float64 `json:"totalAmount"`   // 交易总金额
}

// SupplierApi 供应商管理查询客户端（只读）
type SupplierApi struct {
	http httpclient.Client
}

func NewSupplierApi() *SupplierApi {
	return &SupplierApi{
		http: *httpclient.DefaultClient,
	}
}

// GetSupplier 根据 ID 获取供应商详情
func (api *SupplierApi) GetSupplier(id int) (*SupplierInfo, error) {
	result := &ApiResponse[SupplierInfo]{}
	if err := api.http.Get(fmt.Sprintf("/api/tenant/supplier/getSupplier/%d", id), result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// ListSupplier 分页查询供应商列表
func (api *SupplierApi) ListSupplier(query SupplierPageQuery) (*ApiPageResponse[map[string]interface{}], error) {
	result := ApiPageResponse[map[string]interface{}]{}
	params := url.Values{}
	params.Set("page", fmt.Sprintf("%d", query.Page))
	params.Set("limit", fmt.Sprintf("%d", query.Limit))
	if query.Name != "" {
		params.Set("name", query.Name)
	}
	if query.States != "" {
		params.Set("states", query.States)
	}
	if err := api.http.Get("/api/tenant/supplier/list/page?"+params.Encode(), &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SupplierPurNumCount 根据供应商登录ID获取采购次数
func (api *SupplierApi) SupplierPurNumCount(loginIds []string) (map[string]int, error) {
	result := &ApiResponse[map[string]int]{}
	if err := api.http.Post("/api/tenant/supplier/supplierPurNumCount", loginIds, result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// SupplyRankingList 查询供应商排行信息
func (api *SupplierApi) SupplyRankingList(query SupplierRankingQuery) (*ApiPageResponse[SupplierRankingInfo], error) {
	result := ApiPageResponse[SupplierRankingInfo]{}
	if err := api.http.Post("/api/tenant/supplier/list/supplyList", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SupplySPUList 查询供应商 SPU 信息
func (api *SupplierApi) SupplySPUList(query SupplierRankingQuery) (*ApiPageResponse[SupplierRankingInfo], error) {
	result := ApiPageResponse[SupplierRankingInfo]{}
	if err := api.http.Post("/api/tenant/supplier/list/getSupplySPUList", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SupplierAfterSalesList 查询供应商售后信息
func (api *SupplierApi) SupplierAfterSalesList(query SupplierRankingQuery) (*ApiPageResponse[SupplierAfterSales], error) {
	result := ApiPageResponse[SupplierAfterSales]{}
	if err := api.http.Post("/api/tenant/supplier/list/getSupplierAfterSales", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SupplierFulfillmentList 查询供应商履约信息
func (api *SupplierApi) SupplierFulfillmentList(query SupplierRankingQuery) (*ApiPageResponse[SupplierFulfillment], error) {
	result := ApiPageResponse[SupplierFulfillment]{}
	if err := api.http.Post("/api/tenant/supplier/list/getSupplierFulfillments", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
