package helper

type BaseResponse struct {
	Status     int         `json:"status"`
	Message    string      `json:"message"`
	Error      interface{} `json:"error,omitempty"`
	Data       interface{} `json:"data,omitempty"`
	Pagination interface{} `json:"pagination,omitempty"`
}

type Page struct {
	Size        int `json:"size"`
	TotalData   int `json:"totalData"`
	CurrentPage int `json:"currentPage"`
	TotalPage   int `json:"totalPage"`
}
