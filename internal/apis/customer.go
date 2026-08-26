package apis

import (
	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

var WeeklyAdBudgetMap = map[int]string{
	1: "< 200 USD",
	2: "200-500 USD",
	3: "500-1000 USD",
	4: "> 1000 USD",
	5: "200-1000 USD",
	6: ">1000USD",
}

var OrderVolumeMap = map[int]string{
	0: "justStart",
	1: "1-100",
	2: "100-500",
	3: "500+",
}

var DurationTypeMap = map[int]string{
	1: "还没开始",
	2: "少于6个月",
	3: "6个月-1年",
	4: "多于1年",
	5: "< 1 years",
	6: "1 - 3 years",
	7: "> 3 years",
}

type CustomerInfoApi struct {
	http httpclient.Client
}

func NewCustomerInfoApi() *CustomerInfoApi {
	return &CustomerInfoApi{
		http: *httpclient.DefaultClient,
	}
}

type CustomerExtendInfo struct {
	AllocatedTime  string `json:"allocatedTime"`  // 资料更新时间
	AllocatedType  int    `json:"allocatedType"`  // 类型
	CustomerSource string `json:"customerSource"` // 客户来源
	DurationType   int    `json:"durationType"`   // 从事 Dropshipping 时长类型
	Level          string `json:"level"`          // 等级
	OrderLevel     string `json:"orderLevel"`     // 订单等级
	OrderNum       int    `json:"orderNum"`       // 最近30天订单数量
	StoreNum       int    `json:"storeNum"`       // 店铺数量
	WeeklyAdBudget int    `json:"weeklyAdBudget"` // 周广告预算
	OrderVolume    int    `json:"orderVolume"`    // 每个月订单量
}

type CustomerInfo struct {
	ID              int     `json:"id"`
	RefID           int     `json:"refId"`          // 关联ID
	Cid             int     `json:"cid"`            // 租户ID
	Username        string  `json:"username"`       // 用户名
	ManagerID       int     `json:"managerId"`      // 客户经理ID
	ManagerName     string  `json:"managerName"`    // 客户经理用户名
	ManagerAssTime  string  `json:"managerAssTime"` // 客户经理分配时间
	PlatformType    int     `json:"platformType"`   // 平台类型：1-facebook
	ThirdID         string  `json:"thirdId"`        // 第三方平台ID
	Password        string  `json:"password"`       // 密码
	Gender          int     `json:"gender"`         // 1-Ms., 2-Mr.
	FirstName       string  `json:"firstName"`
	LastName        string  `json:"lastName"`
	Email           string  `json:"email"` // 联系邮箱
	Describes       int     `json:"describes"`
	Nickname        string  `json:"nickname"` // 昵称
	Company         string  `json:"company"`
	RecommenderID   int     `json:"recommenderId"`
	HeaderPortrait  string  `json:"headerPortrait"` // 头像
	LastLoginTime   string  `json:"lastLoginTime"`  // 最近登录时间
	LastLoginHost   string  `json:"lastLoginHost"`  // 最近登录ip
	LastLoginType   int     `json:"lastLoginType"`  // 最近登录类型：0.前台 1.后台
	RegHost         string  `json:"regHost"`        // 注册ip
	RecognitionCode string  `json:"recognitionCode"`
	RegRegion       string  `json:"regRegion"`    // 注册区域
	RegCountryID    int     `json:"regCountryId"` // 注册国家ID
	RegTime         string  `json:"regTime"`      // 注册时间
	RegSource       int     `json:"regSource"`    // 注册来源: 0-hypersku, 1-facebook
	TotalOrder      int64   `json:"totalOrder"`   // 总订单数
	TotalAmount     float64 `json:"totalAmount"`  // 总购买金额
	Status          int     `json:"status"`       // 状态：1：正常, 2:停用
	ReplyToEmail    string  `json:"replyToEmail"` // 回复邮箱地址
	TimeZoneID      int     `json:"timeZoneId"`   // 时区设置（默认是0）
}

// GetCustomerExtendInfo 获取客户扩展信息
func (api *CustomerInfoApi) GetCustomerExtendInfo(customerId string) (*CustomerExtendInfo, error) {
	result := &ApiResponse[CustomerExtendInfo]{}
	if err := api.http.Get("/api/customer/manager/customer/other/info/"+customerId, result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}

// GetCustomerInfo 获取客户信息
func (api *CustomerInfoApi) GetCustomerInfo(customerId string) (*CustomerInfo, error) {

	result := &CustomerInfo{}
	if err := api.http.Get("/api/customer/api/customer/info/"+customerId, result); err != nil {
		return nil, err
	}

	return result, nil

}
