package topup

import (
	"context"
	"database/sql"
	"fmt"

	"modulegue/core/utils"
	model "modulegue/internal/domain/web/model/topup"
	"modulegue/internal/domain/web/repository"
)

type TopUpRepositoryImpl struct {
	db *sql.DB
}

func NewTopUpRepositoryImpl(db *sql.DB) repository.TopUpRepository {
	return &TopUpRepositoryImpl{db: db}
}

func (r *TopUpRepositoryImpl) TopUp(ctx context.Context, reqModel model.TopUpRequestModel) (*model.TopUpResponseModel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin topup transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	walletID, balanceBefore, err := r.getWalletForTopUp(ctx, tx, reqModel.IDUser)
	if err != nil {
		return nil, err
	}

	const insertQuery = `
INSERT INTO payment_topup_transaction (
    topup_code,
    user_id,
    wallet_id,
    payment_method_id,
    payment_status_id,

    amount,
    admin_fee,
    total_amount,

    transaction_status,
    external_reference,
    provider_name,
    qr_string,

    paid_at,
    expired_at,
    failed_reason,
    completed_at,

    created_at,
    updated_at
)
VALUES (
    ?,
    ?,
    ?,

    (
        SELECT id
        FROM master_payment_method
        WHERE payment_method_code = 'MANUAL'
        LIMIT 1
    ),

    (
        SELECT id
        FROM master_payment_status
        WHERE status_code = 'SUCCESS'
        LIMIT 1
    ),

    ?,
    0,
    ?,

    'SUCCESS',
    ?,
    'MANUAL',
    NULL,

    NOW(),
    NULL,
    NULL,
    NOW(),

    NOW(),
    NOW()
);
`

	totalAmount := float64(reqModel.NominalTopUp)
	result, err := tx.ExecContext(
		ctx,
		insertQuery,
		reqModel.TopUpCode,
		reqModel.IDUser,
		walletID,
		reqModel.NominalTopUp,
		totalAmount,
		reqModel.ExternalReference,
	)
	if err != nil {
		return nil, fmt.Errorf("insert topup transaction: %w", err)
	}

	resultID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get last insert id: %w", err)
	}

	balanceAfter := balanceBefore + float64(reqModel.NominalTopUp)

	const updateWalletQuery = `
UPDATE wallet_account
SET
    current_balance_amount = current_balance_amount + ?,
    updated_at = NOW()
WHERE id = ?;
`
	if _, err := tx.ExecContext(ctx, updateWalletQuery, reqModel.NominalTopUp, walletID); err != nil {
		return nil, fmt.Errorf("update wallet balance: %w", err)
	}

	const insertHistoryQuery = `
INSERT INTO wallet_history (
    wallet_id,
    user_id,
    reference_type,
    reference_id,
    mutation_type,
    amount,
    balance_before,
    balance_after,
    description,
    created_at
)
VALUES (
    ?,
    ?,
    'TOPUP',
    ?,
    'CREDIT',
    ?,
    ?,
    ?,
    'Top up saldo berhasil',
    NOW()
);
`
	if _, err := tx.ExecContext(
		ctx,
		insertHistoryQuery,
		walletID,
		reqModel.IDUser,
		resultID,
		reqModel.NominalTopUp,
		balanceBefore,
		balanceAfter,
	); err != nil {
		return nil, fmt.Errorf("insert wallet history: %w", err)
	}

	const updateSummaryQuery = `
UPDATE summary_user_home
SET
    saldo = ?,
    updated_at = NOW()
WHERE user_id = ?;
`
	if _, err := tx.ExecContext(ctx, updateSummaryQuery, balanceAfter, reqModel.IDUser); err != nil {
		return nil, fmt.Errorf("update summary user home: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit topup transaction: %w", err)
	}
	tx = nil

	return &model.TopUpResponseModel{
		TopUpTransactionID: resultID,
		TopUpCode:          reqModel.TopUpCode,
		ExternalReference:  reqModel.ExternalReference,
		BalanceBefore:      balanceBefore,
		BalanceAfter:       balanceAfter,
	}, nil
}

func (r *TopUpRepositoryImpl) getWalletForTopUp(ctx context.Context, tx *sql.Tx, userID int) (int64, float64, error) {
	const query = `
SELECT
    wa.id AS walletId,
    wa.current_balance_amount AS balanceBefore
FROM wallet_account wa
JOIN master_wallet_type mwt
    ON mwt.id = wa.wallet_type_id
WHERE wa.user_id = ?
  AND mwt.wallet_type_code = 'BALANCE'
  AND wa.is_active = 1
LIMIT 1
FOR UPDATE;
`

	var (
		walletID      sql.NullInt64
		balanceBefore sql.NullFloat64
	)

	if err := tx.QueryRowContext(ctx, query, userID).Scan(&walletID, &balanceBefore); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("Top up gagal, wallet user tidak ditemukan")
		}
		return 0, 0, fmt.Errorf("get wallet for topup: %w", err)
	}

	return utils.NullInt64Value(walletID), utils.NullFloat64Value(balanceBefore), nil
}
