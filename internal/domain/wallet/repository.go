package wallet

import "context"

type Repository interface {
	FindByUserID(
		ctx context.Context,
		userID int64,
	) (*Wallet, error)

	Update(
		ctx context.Context,
		wallet *Wallet,
	) error
}
