package mysql_model

// 流水表
type Transaction struct {
	Id              uint64 `gorm:"column:id;AUTO_INCREMENT;primary_key" json:"id"`
	TransactionId   uint64 `gorm:"column:transaction_id" json:"transaction_id"`
	TransactionType string `gorm:"column:transaction_type" json:"transaction_type"`
	FromAccountId   uint64 `gorm:"column:from_account_id" json:"from_account_id"`
	ToAccountId     uint64 `gorm:"column:to_account_id" json:"to_account_id"`
	Amount          uint64 `gorm:"column:amount" json:"amount"`
	Currency        string `gorm:"column:currency" json:"currency"`

	//Currency        *string `gorm:"column:currency" json:"currency"`
	CreateTime      uint64 `gorm:"column:create_time;autoCreateTime" json:"create_time"`
	UpdateTime      uint64 `gorm:"column:update_time;autoUpdateTime" json:"update_time"`
	ExpansionFactor uint64 `gorm:"column:expansion_factor" json:"expansion_factor"`
}

func (m *Transaction) TableName() string {
	return "transaction"
}
