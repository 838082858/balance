package model

type TransferBalanceReq struct {
	TransactionId   uint64 `json:"transaction_id" binding:"required"`
	TransactionType string `json:"transaction_type" binding:"required"`
	FromAccountId   uint64 `json:"from_account_id" binding:"required"`
	ToAccountId     uint64 `json:"to_account_id" binding:"required"`
	Amount          uint64 `json:"amount" binding:"required"`
	Currency        string `json:"currency" binding:"required"`
	ExpansionFactor uint64 `json:"expansion_factor" binding:"required"`
}
