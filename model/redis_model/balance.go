package redis_model

type Balance struct {
	BalanceAccountId uint64 `json:"id"`                    // balance account id
	Balance          uint64 `json:"balance"`               // balance
	CreateTime       uint64 `json:"create_time,omitempty"` // create time
	UpdateTime       uint64 `json:"update_time,omitempty"` // update time
	Currency         string `json:"currency"`              // currency
}
