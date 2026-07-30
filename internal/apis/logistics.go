package apis

import (
	"time"

	"github.com/hypersku/hypersku-cli/internal/httpclient"
)

type Logistics struct {
	http *httpclient.Client
}

func NewLogisticsApi() *Logistics {

	opts := []httpclient.Option{
		httpclient.WithTimeout(time.Duration(30) * time.Second),
		httpclient.WithHeader("User-Agent", "HyperSKU-CLI/0.1.0"),
		httpclient.WithHeader("X-Requested-With", "XMLHttpRequest"),
		httpclient.WithBaseURL("https://track.hypersku.com"),
	}

	return &Logistics{
		http: httpclient.New(opts...),
	}
}

type LogisticsTrackingResult struct {
	Data    string                  `json:"data"`
	RetCode string                  `json:"retCode"`
	RetInfo map[string]TrackingInfo `json:"retInfo"`
}

type TrackingInfo struct {
	Items []TrackingInfoItem `json:"flow_0"`
}

type TrackingInfoItem struct {
	City    string `json:"eventCity"`    // 城市
	Country string `json:"eventCountry"` // 国家
	Detail  string `json:"eventDetail"`  // 详情
	State   string `json:"eventState"`   // 状态
	Thing   string `json:"eventThing"`   // 轨迹说明
	Time    string `json:"eventTime"`    // 时间
	ZipCode string `json:"eventZIPCode"` // 邮编
}

// GetTrackList 获取物流轨迹
func (l *Logistics) GetTracking(trackingNumber string) (*[]TrackingInfoItem, error) {

	body := map[string]string{
		"finalNo":     trackingNumber,
		"doFinalType": "0",
		"finalId":     "0",
	}

	trackResult := &LogisticsTrackingResult{}
	err := l.http.PostForm("/Index/track_all.html", body, trackResult)
	if err != nil {
		return nil, err
	}
	retInfo := trackResult.RetInfo[trackingNumber].Items
	return &retInfo, nil
}
