//Zone, Tariff, Session

package parking

import "time"

type ParkingSession struct {
	ID          int64
	UserID      int64
	VehiclePlat string
	EntryTime   time.Time
	ExitTime    *time.Time
	Status      string
}
