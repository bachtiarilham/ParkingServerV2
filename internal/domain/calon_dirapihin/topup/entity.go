// Token, Session
package topup

import "time"

type Request struct {
	ID int64

	UserID int64

	Amount int64

	ReferenceNo string

	Status Status

	CreatedAt time.Time
}
