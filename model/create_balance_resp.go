package model

type CreateBalanceResp struct {
	BalanceAccountId uint64 `json:"id"`                    // id
	Balance          uint64 `json:"balance"`               // balance
	CreateTime       uint64 `json:"create_time,omitempty"` // create time
	Currency         string `json:"currency,omitempty"`    // currency
}
