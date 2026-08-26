package apis

import (
	"fmt"
	"net/url"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

// CustomerProfileApi 客户画像
type CustomerProfileApi struct {
	http *httpclient.Client
}

func NewCustomerProfileApi() *CustomerProfileApi {
	return &CustomerProfileApi{
		http: httpclient.DefaultClient,
	}
}

// ---------- 订单 ----------

type CustomerProfileOrderCount struct {
	StartDate string `json:"startDate"` // 开始日期
	EndDate   string `json:"endDate"`   // 结束日期
	Avg       int    `json:"avg"`       // 平均订单数
	Max       int    `json:"max"`       // 日最大订单数
	Min       int    `json:"min"`       // 日最小订单数
	Total     int    `json:"total"`     // 总订单数
}

type CustomerProfileReceiveVisitorsInfo struct {
	CrtDate               string  `json:"crtDate"`      // 日期
	DaysNum               int     `json:"daysNum"`      // 订单数量
	FkOrderNum            int     `json:"fkOrderNum"`   // 付款订单数量
	FulfilledNum          int     `json:"fulfilledNum"` // 履约订单数量
	OverdueNum            int     `json:"overdueNum"`   // 超期订单数量
	TkOrderNum            int     `json:"tkOrderNum"`   // 退款订单数量
	AssociateRadio        float64 `json:"associateRadio"`
	FkOrderNumRatio       float64 `json:"fkOrderNumRatio"`
	HyperskuRate          float64 `json:"hyperskuRate"`
	HyperskuRatio         float64 `json:"hyperskuRatio"`
	NotLinkAssociateRadio float64 `json:"notLinkAssociateRadio"`
	OrderNumRatio         float64 `json:"orderNumRatio"`
	OthersRate            float64 `json:"othersRate"`
	OverdueRatio          float64 `json:"overdueRatio"`
	Rate                  float64 `json:"rate"`
	RefundRatio           float64 `json:"refundRatio"`
}

// GetCustomerProfileOrderCount 客户画像-订单
func (api *CustomerProfileApi) GetCustomerProfileOrderCount(customerId, startDate, endDate string) (*CustomerProfileOrderCount, error) {
	params := url.Values{}
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)

	result := &ApiResponse[CustomerProfileOrderCount]{}
	if err := api.http.Get("/api/customer/order/statistics/customer/thirty/days/avg/orders/"+customerId+"?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCustomerProfileReceiveVisitorsInfo 客户画像-日订单数量
func (api *CustomerProfileApi) GetCustomerProfileReceiveVisitorsInfo(customerId, startDate, endDate string, page, limit int) (*PageData[CustomerProfileReceiveVisitorsInfo], error) {
	params := url.Values{}
	params.Add("customerId", customerId)
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	params.Add("daterType", "day")
	params.Add("page", fmt.Sprint(page))
	params.Add("limit", fmt.Sprint(limit))

	result := &ApiPageResponse[CustomerProfileReceiveVisitorsInfo]{}
	if err := api.http.Get("/api/customer/manager/customer/receiveVisitorsInfoByDays?"+params.Encode(), result); err != nil {
		return nil, err
	}
	return result.Data, nil
}

// ---------- 交易 ----------

type CustomerProfileTransactionData struct {
	ActualAmount   float64 `json:"actualAmount"`   // 实际成交额
	AllOrderNum    int     `json:"allOrderNum"`    // 总订单数
	CustomerPrice  float64 `json:"customerPrice"`  // 客单价
	RefundAmount   float64 `json:"refundAmount"`   // 退款金额
	RefundOrderNum int     `json:"refundOrderNum"` // 退款订单数
	TranAmount     float64 `json:"tranAmount"`     // 总交易额
}

type CustomerProfileTransactionChart struct {
	ActualAmountList []string `json:"actualAmountList"`
	Days             []string `json:"days"`
}

type CustomerProfileTransactionCountResult struct {
	Chart CustomerProfileTransactionChart `json:"chart"`
	Data  CustomerProfileTransactionData  `json:"data"`
}

type CustomerProfileTransactionBillRecordItem struct {
	ActualTranAmountCny float64 `json:"actualTranAmountCny"`
	ActualTranAmountEur float64 `json:"actualTranAmountEur"`
	ActualTranAmountUsd float64 `json:"actualTranAmountUsd"`
	Day                 string  `json:"day"`
	OrderNum            float64 `json:"orderNum"`
	Quarter             string  `json:"quarter"`
	RefundOrderNum      float64 `json:"refundOrderNum"`
	RefundTranAmountCny float64 `json:"refundTranAmountCny"`
	RefundTranAmountEur float64 `json:"refundTranAmountEur"`
	RefundTranAmountUsd float64 `json:"refundTranAmountUsd"`
	TranAmountCny       float64 `json:"tranAmountCny"`
	TranAmountEur       float64 `json:"tranAmountEur"`
	TranAmountUsd       float64 `json:"tranAmountUsd"`
	TranFeeCny          float64 `json:"tranFeeCny"`
	TranFeeEur          float64 `json:"tranFeeEur"`
	TranFeeUsd          float64 `json:"tranFeeUsd"`
	Year                int     `json:"year"`
	YearMonth           string  `json:"yearMonth"`
}

// GetCustomerProfileTransactionCount 客户画像-交易
func (api *CustomerProfileApi) GetCustomerProfileTransactionCount(customerId, startDate, endDate string) (*CustomerProfileTransactionData, error) {
	params := url.Values{}
	params.Add("customerId", customerId)
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	params.Add("daterType", "day")

	result := &ApiResponse[CustomerProfileTransactionCountResult]{}
	if err := api.http.Get("/api/customer/manager/orders/transaction/data?"+params.Encode(), result); err != nil {
		return nil, err
	}

	return &result.Data.Data, nil
}

// GetCustomerProfileTransactionBillRecords 客户画像-交易流水
func (api *CustomerProfileApi) GetCustomerProfileTransactionBillRecords(customerId, startDate, endDate string, page, limit int) (*PageData[CustomerProfileTransactionBillRecordItem], error) {

	params := url.Values{}
	params.Add("customerId", customerId)
	params.Add("startDate", startDate)
	params.Add("endDate", endDate)
	params.Add("daterType", "day")
	params.Add("page", fmt.Sprint(page))
	params.Add("limit", fmt.Sprint(limit))

	result := &ApiPageResponse[CustomerProfileTransactionBillRecordItem]{}
	if err := api.http.Get("/api/customer/manager/orders/billing/Analysis/data?"+params.Encode(), result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
