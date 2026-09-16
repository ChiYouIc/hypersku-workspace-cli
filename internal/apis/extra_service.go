package apis

import (
	"fmt"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// ExtraServiceApi 增值服务查询客户端（只读）
type ExtraServiceApi struct {
	http httpclient.Client
}

func NewExtraServiceApi() *ExtraServiceApi {
	return &ExtraServiceApi{
		http: *httpclient.DefaultClient,
	}
}

// ExtraServiceName 增值服务名称
type ExtraServiceName struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type int    `json:"type"`
}

// ExtraServiceValueAdded 增值服务关联的仓储增值服务
type ExtraServiceValueAdded struct {
	ID             int    `json:"id"`
	ExtraServiceID int    `json:"extraServiceId"`
	RepertoryID    int    `json:"repertoryId"`
	RepertoryName  string `json:"repertoryName"`
	Name           string `json:"name"`
	Status         bool   `json:"status"`
}

// ExtraServiceGoods 增值服务关联的商品
type ExtraServiceGoods struct {
	ID             int    `json:"id"`
	ExtraServiceID int    `json:"extraServiceId"`
	GoodsID        int64  `json:"goodsId"`
	GoodsName      string `json:"goodsName"`
	ImgUrl         string `json:"imgUrl"`
}

// GetExtraServiceNames 获取增值服务名称列表
func (api *ExtraServiceApi) GetExtraServiceNames() ([]ExtraServiceName, error) {
	result := &ApiResponse[[]ExtraServiceName]{}
	if err := api.http.Get("/api/tenant/extra-service/getExtraServicesNameByType", result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListValueAddedByService 查询增值服务关联的仓储增值服务
func (api *ExtraServiceApi) ListValueAddedByService(extraServiceId int) ([]ExtraServiceValueAdded, error) {
	result := &ApiResponse[[]ExtraServiceValueAdded]{}
	if err := api.http.Get(fmt.Sprintf("/api/tenant/extra-service/value-added/list/%d", extraServiceId), result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListGoodsByService 查询增值服务关联的商品
func (api *ExtraServiceApi) ListGoodsByService(extraServiceId int) ([]ExtraServiceGoods, error) {
	result := &ApiResponse[[]ExtraServiceGoods]{}
	if err := api.http.Get(fmt.Sprintf("/api/tenant/extra-service/goods/matched/list/%d", extraServiceId), result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ListWarehouseAddValue 查询仓库的增值服务列表
func (api *ExtraServiceApi) ListWarehouseAddValue(repertoryId int) ([]string, error) {
	result := &ApiResponse[[]string]{}
	if err := api.http.Get(fmt.Sprintf("/api/tenant/extra-service/query/warehouse/add-value/%d", repertoryId), result); err != nil {
		return nil, err
	}
	return result.Data, nil
}
