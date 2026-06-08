//PaymentInfo, PaymentType

package payment

import "time"

type Payment struct {
	ID               int64
	UserID           int64
	ParkingSessionID int64
	Amount           int64
	Status           string
	CreatedAt        time.Time
}
