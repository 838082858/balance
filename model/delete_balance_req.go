package model

// 删除请求
type DeleteBalanceReq struct {
	Id uint64 `json:"id" binding:"required"`
}
