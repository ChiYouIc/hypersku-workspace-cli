package apis

import (
	"fmt"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

type Purchase struct {
	HTTP *httpclient.Client
}

func NewPurchaseApi() *Purchase {
	return &Purchase{
		HTTP: httpclient.DefaultClient,
	}
}

// ---------- 请求参数 ----------

type QueryPage struct {
	Page           int    `json:"page"`
	Limit          int    `json:"limit"`
	StartDate      string `json:"startDate"`      // 结束时间
	EndDate        string `json:"endDate"`        // 开始时间
	Id             string `json:"id"`             // 订单号
	TransactionNo  string `json:"transactionNo"`  // 第三方交易号
	TrackingNumber string `json:"trackingNumber"` // 物流单号
}

// ---------- 响应体 ----------

// Order 订单信息
type Order struct {
	ID                    string        `json:"id"`
	OrderId               string        `json:"orderId"`            // 订单号
	CustomerOrderId       string        `json:"customerOrderId"`    // 客户订单号
	AmountActuallyPaid    float64       `json:"amountActuallyPaid"` // 金额
	CrtTime               string        `json:"crt_time"`           // 下单时间
	Freight               float64       `json:"freight"`            // 运费
	Type                  int           `json:"type"`               // 类型
	Warehouse             string        `json:"warehouse"`          // 仓库
	Status                int           `json:"status"`             // 状态
	CanOrderReplenishment bool          `json:"canOrderReplenishment"`
	Consumable            bool          `json:"consumable"`
	CustomerId            int64         `json:"customerId"`
	CustomerLogisticsId   int           `json:"customerLogisticsId"`
	CustomerLogisticsName string        `json:"customerLogisticsName"`
	CustomerTag           int           `json:"customerTag"`
	CustomerUsername      string        `json:"customerUsername"`
	Discount              float64       `json:"discount"`
	GoodsList             []GoodsItem   `json:"goodsList"`
	GoodsNum              int           `json:"goodsNum"`
	HasModifyPrice        bool          `json:"hasModifyPrice"`
	HasRefOrd             bool          `json:"hasRefOrd"`
	InStorageNum          int           `json:"inStorageNum"`
	IntLogisticsId        int           `json:"intLogisticsId"`
	IntLogisticsName      string        `json:"intLogisticsName"`
	Level                 string        `json:"level"`
	LevelRangeStr         string        `json:"levelRangeStr"`
	Opt                   int           `json:"opt"`
	OrderLevel            string        `json:"orderLevel"`
	OrderLevelRangeStr    string        `json:"orderLevelRangeStr"`
	OrdersAddress         OrdersAddress `json:"ordersAddress"`
	OriginId              int           `json:"originId"`
	OriginName            string        `json:"originName"`
	PurchasingDays        int           `json:"purchasingDays"`
	Remarks               string        `json:"remarks"`
	ShowMoreAddress       bool          `json:"showMoreAddress"`
	Source                int           `json:"source"`
	Split                 int           `json:"split"`
	StatusStr             string        `json:"statusStr"`
	StockHasValueAdded    bool          `json:"stockHasValueAdded"`
	StoreId               int           `json:"storeId"`
	StoreName             string        `json:"storeName"`
	Tax                   float64       `json:"tax"`
	ThirdName             string        `json:"thirdName"`
	TimeZoneId            int           `json:"timeZoneId"`
	TransferInventory     bool          `json:"transferInventory"`
	Uname                 string        `json:"uname"`
	UpdTime               string        `json:"upd_time"`
	WarehouseId           int           `json:"warehouseId"`
}

// GoodsItem 商品明细
type GoodsItem struct {
	ID                     string  `json:"id"`
	CstOrdGoodsId          string  `json:"cstOrdGoodsId"`
	CategoryEnName         string  `json:"categoryEnName"`
	CategoryId             int     `json:"categoryId"`
	CategoryName           string  `json:"categoryName"`
	CrtTime                string  `json:"crt_time"`
	DefaultRepertoryId     int     `json:"defaultRepertoryId"`
	GoodsAttr              string  `json:"goodsAttr"`
	GoodsAttrEn            string  `json:"goodsAttrEn"`
	GoodsId                int64   `json:"goodsId"`
	GoodsName              string  `json:"goodsName"`
	GoodsSource            int     `json:"goodsSource"` // 来源
	GoodsSourceUrl         string  `json:"goodsSourceUrl"`
	HasReplenishment       bool    `json:"hasReplenishment"`
	ImgUrl                 string  `json:"imgUrl"`
	IsRefund               int     `json:"isRefund"`
	Num                    int     `json:"num"`
	OrderGoodsStatus       int     `json:"orderGoodsStatus"`
	OrderGoodsStatusStr    string  `json:"orderGoodsStatusStr"`
	ProductId              int     `json:"productId"`
	PurchaseAfterSale      bool    `json:"purchaseAfterSale"`
	PurchaseRecordNum      int     `json:"purchaseRecordNum"`
	PurchaseSource         int     `json:"purchaseSource"`
	RefundStatus           int     `json:"refundStatus"`
	RepurchaseStatus       int     `json:"repurchaseStatus"`
	SubCategoryEnName      string  `json:"subCategoryEnName"`
	SubCategoryId          int     `json:"subCategoryId"`
	SubCategoryName        string  `json:"subCategoryName"`
	SupplierId             int     `json:"supplierId"`
	ThirdLineItemId        string  `json:"thirdLineItemId"`
	Type                   int     `json:"type"`
	UnitPrice              float64 `json:"unitPrice"`
	UpdTime                string  `json:"upd_time"`
	Virtual                bool    `json:"virtual"`
	Weight                 float64 `json:"weight"`
	InventoryCode          string  `json:"inventoryCode"`          // 库存代码
	InventoryWarehouseName string  `json:"inventoryWarehouseName"` // 库存仓库
	TrackingNumber         string  `json:"trackingNumber"`         // 快递单号
	ThirdOrderId           string  `json:"thirdOrderId"`           // 第三方订单号
}

// OrdersAddress 收货地址
type OrdersAddress struct {
	ID               string `json:"id"`
	OrderId          string `json:"orderId"`
	Address          string `json:"address"`
	AddressContinued string `json:"addressContinued,omitempty"`
	CountryId        int    `json:"countryId"`
	CountryImg       string `json:"countryImg"`
	CountryName      string `json:"countryName"`
	FirstName        string `json:"firstName"`
	Flag             string `json:"flag"`
	LastName         string `json:"lastName"`
	Phone            string `json:"phone"`
	SecondRegionName string `json:"secondRegionName"`
	ThirdRegionName  string `json:"thirdRegionName"`
	FourthRegionName string `json:"fourthRegionName"`
	ZipCode          string `json:"zipCode"`
	TaxNo            string `json:"taxNo"`
	VatNo            string `json:"vatNo"`
}

// InternationalLogistics 国际物流
type InternationalLogistics struct {
	ExpressDelivery string        `json:"expressDelivery"` // 物流商
	TrackingNumber  string        `json:"trackingNumber"`  // 物流单号
	TruckInfo       DeliveryTrack `json:"truckInfo"`       // 轨迹
	Type            int           `json:"type"`            // 类型
	Warehouse       string        `json:"warehouse"`       // 仓库
}

// DeliveryTrack 物流轨迹
type DeliveryTrack struct {
	TrackingNumber string              `json:"trackingNumber"` // 物流单号
	List           []DeliveryTrackItem `json:"list"`           // 轨迹记录
	Status         string              `json:"status"`         // 状态
}

// DeliveryTrackItem 物流轨迹记录
type DeliveryTrackItem struct {
	ProcessDate     string `json:"processDate"`     // 时间
	ProcessContent  string `json:"processContent"`  // 内容
	ProcessLocation string `json:"processLocation"` // 位置
}

// PurchaseLog 采购日志
type PurchaseLog struct {
	ID      int    `json:"id"`
	OrderId int    `json:"orderId"` // 订单号
	Status  int    `json:"status"`  // 状态
	Uname   string `json:"uname"`   // 操作人
	UpdTime string `json:"updTime"` // 操作时间
}

// ---------- API 方法 ----------

// PageList 查询订单列表
func (p *Purchase) PageList(query QueryPage) (*ApiPageResponse[Order], error) {
	// 构建查询参数
	params := url.Values{}
	if query.Page > 0 {
		params.Set("page", fmt.Sprintf("%d", query.Page))
	}
	if query.Limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", query.Limit))
	}
	if query.StartDate != "" {
		params.Set("startDate", query.StartDate)
	}
	if query.EndDate != "" {
		params.Set("endDate", query.EndDate)
	}
	if query.Id != "" {
		params.Set("id", query.Id)
		params.Set("orderSearchText", query.Id)
	}
	if query.TransactionNo != "" {
		params.Set("transactionNo", query.TransactionNo)
	}
	if query.TrackingNumber != "" {
		params.Set("trackingNumber", query.TrackingNumber)
		params.Set("logisticsSearchText", query.TrackingNumber)
	}

	// 拼接查询字符串
	path := "/api/tenant/orders/query/manage/list/page"
	if len(params) > 0 {
		path = path + "?" + params.Encode()
	}

	result := &ApiPageResponse[Order]{}
	if err := p.HTTP.Get(path, result); err != nil {
		return nil, fmt.Errorf("查询订单列表失败: %w", err)
	}

	return result, nil
}

// GetOrderInfo 查询订单信息
func (p *Purchase) GetOrderInfo(orderId string) (*Order, error) {
	// 构建查询参数
	params := url.Values{}
	params.Set("id", orderId)
	params.Set("page", "1")
	params.Set("limit", "1")
	params.Set("orderSearchText", orderId)

	path := "/api/tenant/orders/query/manage/list/page" + "?" + params.Encode()
	result := &ApiPageResponse[Order]{}
	if err := p.HTTP.Get(path, result); err != nil {
		return nil, fmt.Errorf("查询订单列表失败: %w", err)
	}

	if result.Data.Total == 0 {
		return nil, nil
	}

	return &result.Data.Rows[0], nil

}

// GetPurchaseLog 查询采购日志
func (p *Purchase) GetPurchaseLog(orderId string) (*[]PurchaseLog, error) {

	path := "/api/tenant/orders/log/" + orderId
	result := make([]PurchaseLog, 0)
	if err := p.HTTP.Get(path, &result); err != nil {
		return nil, fmt.Errorf("查询采购日志失败: %w", err)
	}

	return &result, nil
}

// GetInternationalLogistics 查询订单国际物流轨迹
func (p *Purchase) GetInternationalLogistics(orderId string) (*[]InternationalLogistics, error) {

	path := "/api/tenant/orders/getLogisticsInformation/" + orderId + "?isBack=true"
	result := &ApiPageResponse[InternationalLogistics]{}

	if err := p.HTTP.Get(path, &result); err != nil {
		return nil, fmt.Errorf("查询采购日志失败: %w", err)
	}

	return &result.Data.Rows, nil
}
