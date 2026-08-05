package apis

import (
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
