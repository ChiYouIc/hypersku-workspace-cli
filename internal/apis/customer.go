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

// CustomerInfo 客户档案
type CustomerInfo struct {
	// ===== 身份标识 =====
	CustomerID int `json:"customerId"` // 客户ID（平台唯一标识，卡片 Header 由宿主渲染，AI 链路仅传 ID）

	// ===== 消费状态与转化阶段 =====
	HasOrder bool   `json:"hasOrder"` // 是否有已支付订单（true = 已出池，拒绝触发）
	Stores   string `json:"stores"`   // 绑定店铺集合（建议裁剪为「绑定状态 + 店铺数」形态；非空 = 已绑店）

	// ===== 问卷五项（完整度第 1 层 ~60% 权重 + 基础画像展示 + AI 输入）=====
	EngagedTime     int    `json:"engagedTime"`     // DS 经验：1 还没开始 / 2 少于6个月 / 3 6个月-1年 / 4 多于1年
	WeeklyAdBudget  int    `json:"weeklyAdBudget"`  // 周广告预算：1 <200USD / 2 200-500USD / 3 500-1000USD / 4 >1000USD / 5 200-1000USD
	OrderVolume     int    `json:"orderVolume"`     // 月订单量预期：0 刚刚开始 / 1 1-100 / 2 100-500 / 3 500+
	Niche           int    `json:"niche"`           // 细分市场（1-10 枚举）
	ServiceInterest string `json:"serviceInterest"` // 意向服务（多选逗号串：1 dropshipping / 2 DTC / 3 POD / 4 Merch / 5 Wholesale）

	// ===== 脱敏布尔信号（下游按原始字段计算回传，替代 PII 原值）=====
	HasContactWay bool `json:"hasContactWay"` // 是否有可用联系方式（chat_account 或 reply_to_email 非空）
	HasFirstName  bool `json:"hasFirstName"`  // 是否填写了名（first_name 非空）

	// ===== 基础档案（保留非 PII 业务属性）=====
	Company     string `json:"company"`     // 公司（业务属性）
	CountryName string `json:"countryName"` // 国家
	RegRegion   string `json:"regRegion"`   // 注册区域

	// ===== 时间特征 =====
	SignedUpAt    string `json:"signedUpAt"`    // 注册时间（刚注册 + 问卷全空 → 冷启动声明）
	LastLoginTime string `json:"lastLoginTime"` // 最近登录（两表冗余，取数口径归数据属主）

	// ===== 渠道归因（卡片 Header Meta + AI 输入特征）=====
	PartnerSource    string `json:"partnerSource"`    // 合作来源代码
	ChannelSource    string `json:"channelSource"`    // 渠道来源
	ChannelSourceSub string `json:"channelSourceSub"` // 渠道来源二级
	ChannelMedium    string `json:"channelMedium"`    // 渠道 Medium
	ChannelCampaign  string `json:"channelCampaign"`  // 渠道 Campaign
	ChannelURL       string `json:"channelUrl"`       // 渠道来源链接

	// ===== 业务标签 =====
	Tag        int    `json:"tag"`        // 客户标签：0 默认 / 1 老用户 / 2 新用户 / 3 潜在用户 / 4 流失用户
	Level      string `json:"level"`      // 用户等级（L0 等）
	OrderLevel string `json:"orderLevel"` // 订单量级别（O0）
}

// GetCustomerInfo 获取客户档案
func (api *CustomerInfoApi) GetCustomerInfo(customerId string) (*CustomerInfo, error) {
	result := &ApiResponse[CustomerInfo]{}
	if err := api.http.Get("/api/customer/outer/mcp/customer/profile/detail/"+customerId, result); err != nil {
		return nil, err
	}

	return &result.Data, nil
}
