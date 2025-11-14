package model

type TransferBalanceResp struct {
	TransactionId   uint64 `json:"transaction_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"`
	FromAccountId   uint64 `json:"from_account_id" binding:"required"`
	ToAccountId     uint64 `json:"to_account_id" binding:"required"`
	Amount          uint64 `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	CreateTime      uint64 `json:"create_time" binding:"required"`
	UpdateTime      uint64 `json:"update_time,omitempty" `
	ExpansionFactor uint64 `json:"expansion_factor" binding:"required"`
}
