package apis

import (
	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// ============ 通用结构 ============

// AiAnalysisQuery AI 分析通用查询参数
type AiAnalysisQuery struct {
	Limit          int    `json:"limit"`
	Page           int    `json:"page"`
	ManagerID      int    `json:"managerId,omitempty"`      // 客户经理 ID
	ManagerIDs     []int  `json:"managerIds,omitempty"`     // 客户经理 ID 列表
	Username       string `json:"username,omitempty"`       // 客户账号
	Type           string `json:"type,omitempty"`           // 处理状态/类型
	Source         string `json:"source,omitempty"`         // 来源
	RiskLevel      int    `json:"riskLevel,omitempty"`      // AI 风险等级
	SupplierLevel  int    `json:"supplierLevel,omitempty"`  // 供应商等级
	Day            string `json:"day,omitempty"`            // 日期
	SearchKey      string `json:"searchKey,omitempty"`      // 搜索关键字
	HyperskuStatus int    `json:"hyperskuStatus,omitempty"` // 异常状态
}

type AiInternationalLogisticsQuery struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	RiskLevel int    `json:"riskLevel,omitempty"` // AI 风险等级
	SearchKey string `json:"searchKey,omitempty"` // AI 风险等级
}

// AiBaseCount 统计结果
type AiBaseCount struct {
	Target string `json:"target"`
	Count  int    `json:"count"`
}

// ============ 国际物流异常分析 ============

// AiInternationalLogisticsRow 国际物流异常分析行
type AiInternationalLogisticsRow struct {
	ID                        int    `json:"id"`
	CustomerID                int    `json:"customerId"`
	CustomerUsername          string `json:"customerUsername"`
	OrderID                   string `json:"orderId"`
	CustomerOrderID           string `json:"customerOrderId"`
	PaymentTime               string `json:"paymentTime"`
	WarehouseID               int    `json:"warehouseId"`
	WarehouseName             string `json:"warehouseName"`
	TrackingNumber            string `json:"trackingNumber"`
	InitTrackingNumber        string `json:"initTrackingNumber"`
	LogisticsCompanyName      string `json:"logisticsCompanyName"`
	LatestLogisticsTrack      string `json:"latestLogisticsTrack"`
	LatestLogisticsStatusText string `json:"latestLogisticsStatusText"`
	ManagerName               string `json:"managerName"`
	HyperskuStatus            int    `json:"hyperskuStatus"`
	HyperskuSubStatus         int    `json:"hyperskuSubStatus"`
	Quantity                  int    `json:"quantity"`
	Remark                    string `json:"remark"`
	Reason                    string `json:"reason"`
	OrderStatusText           string `json:"orderStatusText"`
	AiSummary                 string `json:"aiSummary"`
	AiRiskLevel               int    `json:"aiRiskLevel"`
	AiSummaryUpdTime          string `json:"aiSummaryUpdTime"`
	OrderRemark               string `json:"orderRemark"`
}

// ============ 库存动销分析 ============

// AiInventoryRow 库存动销分析行
type AiInventoryRow struct {
	ID                    int     `json:"id"`
	AiSummaryID           int64   `json:"aiSummaryId"`
	CustomerID            int     `json:"customerId"`
	CustomerName          string  `json:"customerName"`
	GoodsID               int64   `json:"goodsId"`
	GoodsName             string  `json:"goodsName"`
	ImgURL                string  `json:"imgUrl"`
	WarehouseID           int     `json:"warehouseId"`
	WarehouseName         string  `json:"warehouseName"`
	Quantity              int     `json:"quantity"`       // 库存数量
	AlertInventory        int     `json:"alertInventory"` // 近15天消耗
	PredictUsageDays      int     `json:"predictUsageDays"`
	DealStatus            int     `json:"dealStatus"`    // AI 处理状态
	SalesQuantity         int     `json:"salesQuantity"` // 90天销量
	AvgInventory          int     `json:"avgInventory"`  // 平均库存
	InventoryTurnover     int     `json:"inventoryTurnover"`
	InventoryTurnoverDays int     `json:"inventoryTurnoverDays"`
	ShelfSalesRate        string  `json:"shelfSalesRate"` // 动销率
	SameProductSkus       int     `json:"sameProductSkus"`
	TotalPrice            float64 `json:"totalPrice"`   // 库存货值
	InventoryAge          int     `json:"inventoryAge"` // 库龄（天）
	LastOrderTime         string  `json:"lastOrderTime"`
	ManagerName           string  `json:"managerName"`
	AiSummary             string  `json:"aiSummary"`
	AiRiskLevel           int     `json:"aiRiskLevel"`
	AiSummaryUpdTime      string  `json:"aiSummaryUpdTime"`
}

// ============ 采购/售后分析 ============

// AiPurchaseRow 采购售后分析行
type AiPurchaseRow struct {
	ID             int64                `json:"id"`
	CustomerID     int                  `json:"customerId"`
	ImgURL         string               `json:"imgUrl"`
	DataID         string               `json:"dataId"`
	ThirdOrderID   string               `json:"thirdOrderId"`
	TrackingNumber string               `json:"trackingNumber"`
	Type           string               `json:"type"`
	RiskLevel      int                  `json:"riskLevel"`
	Source         string               `json:"source"`
	Summary        string               `json:"summary"`
	CrtTime        string               `json:"crtTime"`
	UpdTime        string               `json:"updTime"`
	Logistics      *AiPurchaseLogistics `json:"logistics"`
	Refund         *AiPurchaseRefund    `json:"refund"`
}

// AiPurchaseLogistics 采购物流监控信息
type AiPurchaseLogistics struct {
	ImgURL               string `json:"imgUrl"`
	TrackingNumber       string `json:"trackingNumber"`
	LogisticsCompanyName string `json:"logisticsCompanyName"`
	WarehouseName        string `json:"warehouseName"`
	HyperskuStatus       int    `json:"hyperskuStatus"`
	HyperskuSubStatus    int    `json:"hyperskuSubStatus"`
	SourceType           int    `json:"sourceType"`
	Quantity             int    `json:"quantity"`
	IdentificationCode   string `json:"identificationCode"`
}

// AiPurchaseRefund 1688 售后信息
type AiPurchaseRefund struct {
	RefundID string `json:"refundId"`
}

// ============ 供应商分析 ============

// AiSupplierRow 供应商分析行
type AiSupplierRow struct {
	AiSummaryID      int    `json:"aiSummaryId"`
	SupplierID       int    `json:"supplierId"`
	SupplierName     string `json:"supplierName"`
	SupplierURL      string `json:"supplierUrl"`
	SupplierType     int    `json:"supplierType"`
	SupplierTypeName string `json:"supplierTypeName"`
	SupplierLevel    int    `json:"supplierLevel"`
	AiSummary        string `json:"aiSummary"`
	AiSummaryUpdTime string `json:"aiSummaryUpdTime"`
}

// ============ 客户订单分析 ============

// AiCustomerOrderRow 客户订单分析行
type AiCustomerOrderRow struct {
	CustomerID       int    `json:"customerId"`
	Username         string `json:"username"`
	Tag              int    `json:"tag"`
	AiSummaryID      int64  `json:"aiSummaryId"`
	AiSummary        string `json:"aiSummary"`
	AiSummaryUpdTime string `json:"aiSummaryUpdTime"`
	AiRiskLevel      int    `json:"aiRiskLevel"`
}

// ============ API 客户端 ============

// AiAnalysisApi AI 分析查询客户端（只读）
type AiAnalysisApi struct {
	http httpclient.Client
}

func NewAiAnalysisApi() *AiAnalysisApi {
	return &AiAnalysisApi{
		http: *httpclient.DefaultClient,
	}
}

// InternationalLogisticsTable 国际物流异常分析表格数据
func (api *AiAnalysisApi) InternationalLogisticsTable(query AiInternationalLogisticsQuery) (*ApiPageResponse[AiInternationalLogisticsRow], error) {
	result := ApiPageResponse[AiInternationalLogisticsRow]{}
	if err := api.http.Post("/api/tenant/ai/analysis/table/monitor_abnormal_international_logistics", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InternationalLogisticsCountRiskLevel 国际物流异常分析按风险等级统计
func (api *AiAnalysisApi) InternationalLogisticsCountRiskLevel(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/count/riskLevel/monitor_abnormal_international_logistics", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// InternationalLogisticsCountHyperskuStatus 国际物流异常分析按处理状态统计
func (api *AiAnalysisApi) InternationalLogisticsCountHyperskuStatus(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/count/hyperskuStatus/monitor_abnormal_international_logistics", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// InventoryTable 库存动销分析表格数据
func (api *AiAnalysisApi) InventoryTable(query AiAnalysisQuery) (*ApiPageResponse[AiInventoryRow], error) {
	result := ApiPageResponse[AiInventoryRow]{}
	if err := api.http.Post("/api/tenant/ai/analysis/inventory/table", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// InventoryCountRiskLevel 库存动销分析按风险等级统计
func (api *AiAnalysisApi) InventoryCountRiskLevel(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/inventory/count/riskLevel", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// InventoryCountDealStatus 库存动销分析按处理状态统计
func (api *AiAnalysisApi) InventoryCountDealStatus(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/inventory/count/dealStatus", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// PurchaseTable 采购售后分析表格数据
func (api *AiAnalysisApi) PurchaseTable(query AiAnalysisQuery) (*ApiPageResponse[AiPurchaseRow], error) {
	result := ApiPageResponse[AiPurchaseRow]{}
	if err := api.http.Post("/api/tenant/ai/analysis/purchase/table", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// PurchaseCountRiskLevel 采购售后分析按风险等级统计
func (api *AiAnalysisApi) PurchaseCountRiskLevel(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/purchase/count/riskLevel", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// PurchaseCountType 采购售后分析按类型统计
func (api *AiAnalysisApi) PurchaseCountType(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/purchase/count/type", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// PurchaseCountSource 采购售后分析按来源统计
func (api *AiAnalysisApi) PurchaseCountSource(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/purchase/count/source", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// SupplierTable 供应商分析表格数据
func (api *AiAnalysisApi) SupplierTable(query AiAnalysisQuery) (*ApiPageResponse[AiSupplierRow], error) {
	result := ApiPageResponse[AiSupplierRow]{}
	if err := api.http.Post("/api/tenant/ai/analysis/supplier/table", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// SupplierCountSupplierLevel 供应商分析按等级统计
func (api *AiAnalysisApi) SupplierCountSupplierLevel(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/supplier/count/supplierLevel", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CustomerOrderTable 客户订单分析表格数据
func (api *AiAnalysisApi) CustomerOrderTable(query AiAnalysisQuery) (*ApiPageResponse[AiCustomerOrderRow], error) {
	result := ApiPageResponse[AiCustomerOrderRow]{}
	if err := api.http.Post("/api/tenant/ai/analysis/customer/order/table", query, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// CustomerOrderCountRiskLevel 客户订单分析按风险等级统计
func (api *AiAnalysisApi) CustomerOrderCountRiskLevel(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/customer/order/count/riskLeve", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// CustomerOrderCountType 客户订单分析按类型统计
func (api *AiAnalysisApi) CustomerOrderCountType(query AiAnalysisQuery) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Post("/api/tenant/ai/analysis/customer/order/count/type", query, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// TaskProgressCountRiskLevelByType 按类型统计 AI 任务风险等级
func (api *AiAnalysisApi) TaskProgressCountRiskLevelByType(taskType string) ([]AiBaseCount, error) {
	result := ApiResponse[[]AiBaseCount]{}
	if err := api.http.Get("/api/tenant/ai/task/progress/countRiskLevelByType?type="+taskType, &result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
