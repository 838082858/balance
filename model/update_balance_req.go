package model

// 更新请求
type UpdateBalanceReq struct {
	Id      uint64 `json:"id" binding:"required"`
	Balance uint64 `json:"balance" binding:"required"`
}
