package apis

import (
	"fmt"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// Warehouse 仓库 API
type Warehouse struct {
	http *httpclient.Client
}

func NewWarehouseApi() *Warehouse {
	return &Warehouse{
		http: httpclient.DefaultClient,
	}
}

// WarehouseTrackingInfo 仓库物流轨迹信息
type WarehouseTrackingInfo struct {
	Abnormal         int                     `json:"abnormal"`
	ActionList       []WarehouseAction       `json:"actionList"`       // 仓库操作列表
	ExpressSignInfo  string                  `json:"expressSignInfo"`  // 签收信息
	ExpressSignTime  string                  `json:"expressSignTime"`  // 签收时间
	FinalNo          string                  `json:"finalNo"`          // 物流单号
	InputInfo        string                  `json:"inputInfo"`        // 入库信息
	InputStatus      int                     `json:"inputStatus"`      // 入库状态
	InstoreStatus    int                     `json:"instoreStatus"`    // 入库状态
	InstoreTime      string                  `json:"instoreTime"`      // 入库时间
	IsSend           int                     `json:"isSend"`           // 是否发货，0：未发货，1：已发货
	LogisticsSteps   []WarehouseTrackingStep `json:"logisticsSteps"`   // 物流轨迹步骤
	SignInfo         string                  `json:"signInfo"`         // 签收信息
	SignStatus       int                     `json:"signStatus"`       // 签收状态
	SignTime         string                  `json:"signTime"`         // 签收时间
	StoreAddressName string                  `json:"storeAddressName"` // 仓库地址名称
	Name             string                  `json:"name"`
	PackageList      []WarehousePackage      `json:"packageList"`
	StoreAddressID   int                     `json:"storeAddressId"`
}

// WarehouseAction 仓库物流轨迹动作信息
type WarehouseAction struct {
	ActionTime string `json:"actionTime"` // 操作时间
	Content    string `json:"content"`    // 内容
}

// WarehouseTrackingStep 仓库物流轨迹步骤信息
type WarehouseTrackingStep struct {
	AcceptTime string `json:"acceptTime"` // 接收时间
	Remark     string `json:"remark"`     // 备注
}

// WarehousePackage 仓库包裹信息
type WarehousePackage struct {
	AdvanceWarehouse int    `json:"advanceWarehouse"`
	OtherID          int    `json:"otherId"`
	Status           int    `json:"status"`
	StorageInTime    string `json:"storageInTime"`
	StoreNo          string `json:"storeNo"`
}

// GetWarehouseTracking 获取仓库物流轨迹
func (api *Warehouse) GetWarehouseTracking(trackingNumber string) ([]*WarehouseTrackingInfo, error) {

	body := map[string]string{
		"expressNumber": trackingNumber,
	}

	result := []*WarehouseTrackingInfo{}
	if err := api.http.Post("/api/tenant/ordersLogistics/queryTrackPackStatusListByPackage", body, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// WarehouseInfo 仓库信息
type WarehouseInfo struct {
	Name                string  `json:"name"`           // 仓库名称
	Code                string  `json:"code"`           // 仓库代码
	Attribute           string  `json:"attribute"`      // 仓库属性
	AreaName            string  `json:"areaName"`       // 区域
	CanStock            bool    `json:"can_stock"`      // 是否支持备货
	CityName            string  `json:"cityName"`       // 城市
	ProvinceName        string  `json:"provinceName"`   // 省
	StatusStr           string  `json:"status_str"`     // 状态描述
	StreetAddress       string  `json:"streetAddress"`  // 详细地址
	BeijingCutTime      string  `json:"beijingCutTime"` // 截单时间 北京时间
	LocalCutTime        string  `json:"localCutTime"`   // 截单时间 本地时间
	AreaID              int     `json:"areaId"`
	CityID              int     `json:"cityId"`
	ContactsName        string  `json:"contactsName"`
	ContactsPhone       string  `json:"contactsPhone"`
	CountryID           int     `json:"countryId"`
	DefaultRepertory    bool    `json:"default_repertory"`
	ExitFactoryPriceSum float64 `json:"exitFactoryPriceSum"`
	ID                  int     `json:"id"`
	ProductKindCount    int     `json:"productKindCount"`
	ProvinceID          int     `json:"provinceId"`
	Status              int     `json:"status"`
	StorageID           int     `json:"storageId"`
	TotalStoreNum       string  `json:"totalStoreNum"`
	UpdHost             string  `json:"upd_host"`
	UpdName             string  `json:"upd_name"`
	UpdTime             string  `json:"upd_time"`
	UpdUser             string  `json:"upd_user"`
}

// GetWarehousePage 获取仓库信息
func (api *Warehouse) GetWarehousePage(warehouseName string, page, limit int) (*PageData[WarehouseInfo], error) {

	params := url.Values{}
	params.Add("page", fmt.Sprintf("%d", page))
	params.Add("limit", fmt.Sprintf("%d", limit))
	params.Add("status", "1")
	if warehouseName != "" {
		params.Add("name", warehouseName)
	}

	result := &ApiPageResponse[WarehouseInfo]{}
	if err := api.http.Get("/api/tenant/repertory/repertory-manager/find/all?"+params.Encode(), result); err != nil {
		return nil, err
	}

	return result.Data, nil
}

type WarehouseLogisticsInfo struct {
	ID           int    `json:"id"`           // 数据id
	Alias        string `json:"alias"`        // 别名
	Name         string `json:"name"`         // 名称
	Contain      int    `json:"contain"`      // 物流属性
	IsShow       bool   `json:"isShow"`       // 支持批量
	PriceType    int    `json:"priceType"`    // 报价方式
	ShippingType int    `json:"shippingType"` // 物流类型
	CrtHost      string `json:"crtHost"`
	CrtName      string `json:"crtName"`
	CrtTime      string `json:"crtTime"`
	CrtUser      string `json:"crtUser"`
	IsApply      bool   `json:"isApply"`
	Status       int    `json:"status"`
	Type         bool   `json:"type"`
	UpdHost      string `json:"updHost"`
	UpdName      string `json:"updName"`
	UpdTime      string `json:"updTime"`
	UpdUser      string `json:"updUser"`
}

// GetLogisticsPage 分页查询物流
func (api *Warehouse) GetLogisticsPage(logisticsName string, countryId, page, limit int) (*PageData[WarehouseLogisticsInfo], error) {
	params := url.Values{}
	params.Add("page", fmt.Sprint(page))
	params.Add("limit", fmt.Sprint(limit))
	params.Add("status", "1")
	params.Add("shippingType", "0")
	if logisticsName != "" {
		params.Add("name", logisticsName)
	}

	if countryId != 0 {
		params.Add("countryId", fmt.Sprint(countryId))
	}

	result := &ApiPageResponse[WarehouseLogisticsInfo]{}
	if err := api.http.Get("/api/tenant/logistics/list/page?"+params.Encode(), result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
