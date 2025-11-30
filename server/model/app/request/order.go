package request

type GetOrderList struct {
	Date string `json:"date" form:"date"`
	// UserId 用户ID
	UserId uint `json:"userId" form:"userId"`
}
