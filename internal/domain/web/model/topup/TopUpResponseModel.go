package topup

type TopUpResponseModel struct {
	TopUpTransactionID int64   `json:"topup_transaction_id"`
	TopUpCode          string  `json:"topup_code"`
	ExternalReference  string  `json:"external_reference"`
	BalanceBefore      float64 `json:"balance_before"`
	BalanceAfter       float64 `json:"balance_after"`
}
