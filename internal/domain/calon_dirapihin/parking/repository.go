package parking

import "context"

type Repository interface {
	Create(ctx context.Context, session *ParkingSession) error
	FindByID(ctx context.Context, id int64) (*ParkingSession, error)
	Update(ctx context.Context, session *ParkingSession) error
}
