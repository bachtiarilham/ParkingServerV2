package repository

import (
	"context"
)

type SyncRepo interface {
	SyncAfterParkir(ctx context.Context, userId int64, refId int64, amount int64) error
	SyncAfterMembership(ctx context.Context, userId int64, packageID int64) error
}
