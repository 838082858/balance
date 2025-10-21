package model

// 查询请求
type GetBalanceReq struct {
	Id uint64 `json:"id" binding:"required"`
}
