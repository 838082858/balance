package model

// 创建请求
type CreateBalanceReq struct {
	Id       uint64 `json:"id" binding:"required"`
	Balance  uint64 `json:"balance"`
	Currency string `json:"currency"`
}
