//Token, Session

package ledger

import "time"

type Entry struct {
	ID            int64
	TransactionID string

	AccountCode AccountCode

	DebitAmount  int64
	CreditAmount int64

	CreatedAt time.Time
}
