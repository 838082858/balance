package mysql_model

// 余额数据库，表模型
type Balance struct {
	Id               uint64 `gorm:"column:id;primary_key;AUTO_INCREMENT"` // id
	BalanceAccountId uint64 `gorm:"column:balance_account_id"`            // balance account id
	Balance          uint64 `gorm:"column:balance"`                       // balance
	CreateTime       uint64 `gorm:"column:create_time;autoCreateTime"`    // create time
	UpdateTime       uint64 `gorm:"column:update_time;autoUpdateTime"`    // update time
	Currency         string `gorm:"column:currency"`                      // currency
	Version          uint64 `gorm:"column:version"`
}

func (m *Balance) TableName() string {
	return "balance"
}
