package repository

import (
	"context"
	"database/sql"
	"fmt"
	"modulegue/internal/domain/jukir_wallet"
)

type JukirWalletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) jukir_wallet.Repository {
	return &JukirWalletRepository{db: db}
}

func (r *JukirWalletRepository) GetJukirWalletByUserID(ctx context.Context, userID int64) (*jukir_wallet.JukirWallet, error) {
	query := `
		SELECT id, jukir_user_id, current_balance, total_topup, total_withdrawn, status, created_at, updated_at
		FROM jukir_wallet
		WHERE jukir_user_id = ?
		LIMIT 1
	`
	var w jukir_wallet.JukirWallet
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&w.ID, &w.UserID, &w.CurrentBalance, &w.TotalTopup, &w.TotalWithdrawn, &w.Status, &w.CreatedAt, &w.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("jukir wallet not found for user id: %d", userID)
		}
		return nil, fmt.Errorf("get jukir wallet by user id: %w", err)
	}
	return &w, nil
}

func (r *JukirWalletRepository) UpdateJukirWalletBalance(ctx context.Context, walletID int64, newBalance int64) error {
	query := `UPDATE jukir_wallet SET current_balance = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, newBalance, walletID)
	if err != nil {
		return fmt.Errorf("update jukir wallet balance: %w", err)
	}
	return nil
}

func (r *JukirWalletRepository) CreateJukirWalletHistory(ctx context.Context, history *jukir_wallet.JukirWalletHistory) error {
	query := `
		INSERT INTO jukir_wallet_history (
			jukir_wallet_id, transaction_type, amount, previous_balance, new_balance, reference_id, description, created_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		history.WalletID, history.TransactionType, history.Amount, history.PreviousBalance, history.NewBalance, history.ReferenceID, history.Description,
	)
	if err != nil {
		return fmt.Errorf("insert jukir wallet history: %w", err)
	}
	return nil
}
