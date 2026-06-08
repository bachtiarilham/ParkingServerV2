package ledger

type AccountCode string

const (
	AccountCash            AccountCode = "CASH"
	AccountWalletLiability AccountCode = "WALLET_LIABILITY"
	AccountRevenue         AccountCode = "REVENUE"
)
