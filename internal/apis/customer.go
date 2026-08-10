package apis

import "github.com/hypersku/hypersku-cli/internal/httpclient"

var CustomerOrderStatus = map[int]string{
	1:  "待付款",
	2:  "待采购",
	3:  "待发货",
	4:  "待收货",
	5:  "待评价",
	6:  "交易关闭",
	7:  "退款中",
	8:  "待找货",
	9:  "支付中",
	10: "超期",
	11: "退件",
	14: "分开发货",
}

type CustomerOrderApi struct {
	http httpclient.Client
}

func NewCustomerOrderApi() *CustomerOrderApi {
	return &CustomerOrderApi{
		http: *httpclient.DefaultClient,
	}
}

// CustomerOrderInfo 订单信息
type CustomerOrderInfo struct {
	ID                    string                   `json:"id"`
	ActualAmount          float64                  `json:"actualAmount"`          // 实际金额
	Amount                float64                  `json:"amount"`                // 订单金额
	Freight               float64                  `json:"freight"`               // 运费
	BrandingServiceAmount float64                  `json:"brandingServiceAmount"` // 增值服务费
	TaxAmount             float64                  `json:"taxAmount"`             // 税费
	TariffAmount          float64                  `json:"tariffAmount"`          // 关税
	CurrencySymbol        string                   `json:"currencySymbol"`        // 币种符号
	CurrencyCode          string                   `json:"currencyCode"`          // 币种
	WarehouseName         string                   `json:"warehouseName"`         // 仓库
	PurchaseStatus        int                      `json:"purchaseStatus"`        // 采购状态
	Status                int                      `json:"status"`                // 客户订单状态
	PaymentTime           string                   `json:"paymentTime"`           // 支付时间
	CrtTime               string                   `json:"crtTime"`               // 订单创建时间
	GoodsList             []CustomerOrderGoodsInfo `json:"goodsList"`             // 订单商品项
	LogisticsList         []CustomerOrderLogistics `json:"logisticsList"`         // 订单物流
	OrdersAddress         CustomerOrderAddress     `json:"ordersAddress"`         // 订单地址
}

// CustomerOrderGoodsInfo 订单商品
type CustomerOrderGoodsInfo struct {
	ID             string  `json:"id"`
	ProductID      int     `json:"productId"`      // spu
	GoodsID        string  `json:"goodsId"`        // sku
	GoodsName      string  `json:"goodsName"`      // 名称
	AttrStr        string  `json:"attrStr"`        // 属性
	ImgURL         string  `json:"imgUrl"`         // 图片
	Num            int     `json:"num"`            // 数量
	SellingPrice   float64 `json:"sellingPrice"`   // 销售价
	UnitPrice      float64 `json:"unitPrice"`      // 单价
	CurrencyCode   string  `json:"currencyCode"`   // 币种
	CurrencySymbol string  `json:"currencySymbol"` // 币种符号
	Weight         float64 `json:"weight"`         // 重量
	IsInventory    int     `json:"isInventory"`    // 1 表示使用库存
}

// CustomerOrderLogistics 订单物流
type CustomerOrderLogistics struct {
	ID                    string `json:"id"`
	TrackingNumber        string `json:"trackingNumber"`        // 快递单号
	ExpressDelivery       string `json:"expressDelivery"`       // 承运商
	InitialTrackingNumber string `json:"initialTrackingNumber"` // 初始化快递单号
}

// CustomerOrderAddress 订单地址
type CustomerOrderAddress struct {
	ID               string `json:"id"`
	FirstName        string `json:"firstName"`        // 姓
	LastName         string `json:"lastName"`         // 名
	Address          string `json:"address"`          // 地址
	CountryName      string `json:"countryName"`      // 国家
	SecondRegionName string `json:"secondRegionName"` // 省份
	ThirdRegionName  string `json:"thirdRegionName"`  // 城市
	FourthRegionName string `json:"fourthRegionName"` // 区域
	TaxNo            string `json:"taxNo"`            // 欧盟税号
	VatNo            string `json:"vatNo"`            // VAT
	ZipCode          string `json:"zipCode"`          // 邮编
}

// GetOrderInfo 获取客户订单信息
func (o *CustomerOrderApi) GetOrderInfo(orderId string) (*CustomerOrderInfo, error) {
	body := struct {
		Id              string `json:"id"`
		Page            int    `json:"page"`
		Limit           int    `json:"limit"`
		ThirdSearchType int    `json:"thirdSearchType"`
		StartTime       string `json:"startTime"`
		EndTime         string `json:"endTime"`
	}{
		Id:              orderId,
		Page:            1,
		Limit:           10,
		ThirdSearchType: 1,
	}

	apiResponse := &ApiPageResponse[CustomerOrderInfo]{}
	if err := o.http.Post("/api/customer/manager/orders/list/page", body, apiResponse); err != nil {
		return nil, err
	}
	rows := apiResponse.Data.Rows
	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}
