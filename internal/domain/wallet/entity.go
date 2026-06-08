//Token, Session

package wallet

import "time"

type Wallet struct {
	ID     int64
	UserID int64

	BalanceAmount   int64
	AvailableAmount int64
	FrozenAmount    int64

	CreatedAt time.Time
	UpdatedAt time.Time
}
