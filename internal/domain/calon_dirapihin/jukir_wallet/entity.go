//Token, Session

package jukir_wallet

import "time"

type JukirWallet struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	CurrentBalance int64     `json:"current_balance"`
	TotalTopup     int64     `json:"total_topup"`
	TotalWithdrawn int64     `json:"total_withdrawn"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type JukirWalletHistory struct {
	ID              int64     `json:"id"`
	WalletID        int64     `json:"wallet_id"`
	TransactionType string    `json:"transaction_type"`
	Amount          int64     `json:"amount"`
	PreviousBalance int64     `json:"previous_balance"`
	NewBalance      int64     `json:"new_balance"`
	ReferenceID     string    `json:"reference_id"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
}
