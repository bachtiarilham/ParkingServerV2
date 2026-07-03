package jukir_wallet

import "context"

type Repository interface {
	// FindByUserID(ctx context.Context, userID int64) (*JukirWallet, error)

	// Update(ctx context.Context, wallet *JukirWallet) error

	GetJukirWalletByUserID(ctx context.Context, userID int64) (*JukirWallet, error)
	UpdateJukirWalletBalance(ctx context.Context, walletID int64, newBalance int64) error
	CreateJukirWalletHistory(ctx context.Context, history *JukirWalletHistory) error
}
