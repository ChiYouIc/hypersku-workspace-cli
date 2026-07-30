package apis

// ApiPageResponse api 分页响应
type ApiPageResponse[T any] struct {
	Status int          `json:"status"`
	Data   *PageData[T] `json:"data"`
}

type PageData[T any] struct {
	Rows  []T   `json:"rows"`
	Total int64 `json:"total"`
}

type ApiResponse[T any] struct {
	Code   int  `json:"code"`
	Rel    bool `json:"rel"`
	Status int  `json:"status"`
	Data   T    `json:"data"`
}
