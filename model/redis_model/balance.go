package redis_model

type Balance struct {
	BalanceAccountId uint64 // balance account id
	Balance          uint64 // balance
	CreateTime       uint64 // create time
	UpdateTime       uint64 // update time
	Currency         string // currency
}
