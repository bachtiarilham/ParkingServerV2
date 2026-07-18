package topup

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/topup"
	"modulegue/internal/domain/mobile/repository"
	"modulegue/internal/middleware"
)

type TopUpRepositoryImpl struct {
	db *sql.DB
}

func NewTopUpRepositoryImpl(db *sql.DB) repository.TopUpRepository {
	return &TopUpRepositoryImpl{db: db}
}

func (r *TopUpRepositoryImpl) TopUpCreate(ctx context.Context, reqModel model.TopupCreateRequestModel) (*model.TopupCreateResponseModel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin topup transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	walletID, _, err := r.getWalletForTopUp(ctx, tx, reqModel.UserID)
	if err != nil {
		return nil, err
	}

	adminFee := reqModel.AdminFee
	totalAmount := reqModel.Amount + adminFee

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
			WHERE payment_method_code = ?
			LIMIT 1
		),

		(
			SELECT id
			FROM master_payment_status
			WHERE status_code = 'PENDING'
			LIMIT 1
		),

		?,
		?,
		?,

		'PENDING',
		?,
		?,
		?,

		NULL,
		?,
		NULL,
		NULL,

		NOW(),
		NOW()
	);
	`

	if _, err := tx.ExecContext(
		ctx,
		insertQuery,
		reqModel.TopupCode,
		reqModel.UserID,
		walletID,
		reqModel.PaymentMethodCode,
		reqModel.Amount,
		adminFee,
		totalAmount,
		reqModel.ExternalReference,
		reqModel.ProviderName,
		reqModel.QRString,
		reqModel.ExpiredAt,
	); err != nil {
		return nil, fmt.Errorf("insert topup transaction: %w", err)
	}

	const returnQuery = `
	SELECT
		ptt.id AS topupTransactionId,
		ptt.topup_code AS topupCode,

		ptt.amount AS amount,
		ptt.admin_fee AS adminFee,
		ptt.total_amount AS totalAmount,

		mpm.payment_method_code AS paymentMethodCode,
		mpm.payment_method_name AS paymentMethodName,

		mps.status_code AS paymentStatusCode,
		mps.status_name AS paymentStatusName,

		ptt.qr_string AS qrString,
		ptt.expired_at AS expiredAt,
		ptt.created_at AS createdAt

	FROM payment_topup_transaction ptt

	JOIN master_payment_method mpm
		ON mpm.id = ptt.payment_method_id

	JOIN master_payment_status mps
		ON mps.id = ptt.payment_status_id

	WHERE ptt.topup_code = ?
	AND ptt.user_id = ?

	LIMIT 1;
	`

	var result model.TopupCreateResponseModel
	if err := tx.QueryRowContext(ctx, returnQuery, reqModel.TopupCode, reqModel.UserID).Scan(
		&result.TopupTransactionID,
		&result.TopupCode,
		&result.Amount,
		&result.AdminFee,
		&result.TotalAmount,
		&result.PaymentMethodCode,
		&result.PaymentMethodName,
		&result.PaymentStatusCode,
		&result.PaymentStatusName,
		&result.QRString,
		&result.ExpiredAt,
		&result.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("get topup transaction: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit topup transaction: %w", err)
	}
	tx = nil

	return &result, nil
}

func (r *TopUpRepositoryImpl) TopUpStatus(ctx context.Context, req string) (*model.TopupStatusResponseModel, error) {
	userID, ok := middleware.UserIDFromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("user id not found in context")
	}

	const query = `
	SELECT
		ptt.id AS topupTransactionId,
		ptt.topup_code AS topupCode,

		ptt.amount AS amount,
		ptt.admin_fee AS adminFee,
		ptt.total_amount AS totalAmount,

		mpm.payment_method_code AS paymentMethodCode,
		mps.status_code AS paymentStatusCode,
		mps.status_name AS paymentStatusName,

		ptt.qr_string AS qrString,
		ptt.paid_at AS paidAt,
		ptt.expired_at AS expiredAt,
		ptt.failed_reason AS failedReason,
		ptt.completed_at AS completedAt,

		wa.current_balance_amount AS currentBalance

	FROM payment_topup_transaction ptt

	JOIN master_payment_method mpm
		ON mpm.id = ptt.payment_method_id

	JOIN master_payment_status mps
		ON mps.id = ptt.payment_status_id

	JOIN wallet_account wa
		ON wa.id = ptt.wallet_id

	WHERE ptt.user_id = ?
	AND ptt.topup_code = ?

	LIMIT 1;
	`

	var result model.TopupStatusResponseModel
	var paidAt sql.NullTime
	var completedAt sql.NullTime
	var paymentStatusName sql.NullString
	var failedReason sql.NullString
	if err := r.db.QueryRowContext(ctx, query, userID, req).Scan(
		&result.TopupTransactionID,
		&result.TopupCode,
		&result.Amount,
		&result.AdminFee,
		&result.TotalAmount,
		&result.PaymentMethodCode,
		&result.PaymentStatusCode,
		&paymentStatusName,
		&result.QRString,
		&paidAt,
		&result.ExpiredAt,
		&failedReason,
		&completedAt,
		&result.CurrentBalance,
	); err != nil {
		return nil, fmt.Errorf("get topup status: %w", err)
	}

	if paidAt.Valid {
		value := paidAt.Time
		result.PaidAt = &value
	}
	if completedAt.Valid {
		value := completedAt.Time
		result.CompletedAt = &value
	}
	_ = paymentStatusName
	result.FailedReason = utils.NullStringValue(failedReason)

	return &result, nil
}

func (r *TopUpRepositoryImpl) TopUpCallback(ctx context.Context, reqModel model.QrisCallbackRequestModel) (*model.QrisCallbackResponseModel, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin topup callback transaction: %w", err)
	}
	defer func() {
		if tx != nil {
			_ = tx.Rollback()
		}
	}()

	const selectQuery = `
	SELECT
		ptt.id AS topupTransactionId,
		ptt.user_id AS userId,
		ptt.wallet_id AS walletId,
		ptt.amount AS amount,
		ptt.payment_status_id,
		mps.status_code AS paymentStatusCode
	FROM payment_topup_transaction ptt
	JOIN master_payment_status mps
		ON mps.id = ptt.payment_status_id
	WHERE ptt.topup_code = ?
	LIMIT 1
	FOR UPDATE;
	`

	var (
		topupTransactionID int64
		userID             int64
		walletID           int64
		amount             int64
		paymentStatusID    int64
		paymentStatusCode  string
	)
	if err := tx.QueryRowContext(ctx, selectQuery, reqModel.TopupCode).Scan(
		&topupTransactionID,
		&userID,
		&walletID,
		&amount,
		&paymentStatusID,
		&paymentStatusCode,
	); err != nil {
		return nil, fmt.Errorf("get topup callback row: %w", err)
	}

	result := &model.QrisCallbackResponseModel{
		TopupTransactionID: topupTransactionID,
		UserID:             userID,
		WalletID:           walletID,
		Amount:             amount,
		PaymentStatusID:    paymentStatusID,
		PaymentStatusCode:  paymentStatusCode,
	}

	incomingStatus := strings.ToLower(strings.TrimSpace(reqModel.TransactionStatus))
	if incomingStatus == "" {
		incomingStatus = "pending"
	}

	if paymentStatusCode == "SUCCESS" || paymentStatusCode == "FAILED" {
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit topup callback: %w", err)
		}
		tx = nil
		return result, nil
	}

	if incomingStatus == "pending" {
		result.PaymentStatusCode = "PENDING"
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit topup pending callback: %w", err)
		}
		tx = nil
		return result, nil
	}

	if incomingStatus != "capture" && incomingStatus != "settlement" {
		failedReason := reqModel.FailedReason
		if failedReason == "" {
			failedReason = reqModel.StatusMessage
		}
		if failedReason == "" {
			failedReason = incomingStatus
		}

		const failQuery = `
	UPDATE payment_topup_transaction ptt
	JOIN master_payment_status mps
		ON mps.status_code = 'FAILED'
	SET
		ptt.payment_status_id = mps.id,
		ptt.transaction_status = 'FAILED',
		ptt.failed_reason = ?,
		ptt.completed_at = NOW(),
		ptt.updated_at = NOW()
	WHERE ptt.topup_code = ?
	AND ptt.transaction_status = 'PENDING';
`
		if _, err := tx.ExecContext(ctx, failQuery, failedReason, reqModel.TopupCode); err != nil {
			return nil, fmt.Errorf("update topup failed callback: %w", err)
		}
		if err := tx.Commit(); err != nil {
			return nil, fmt.Errorf("commit topup failed callback: %w", err)
		}
		tx = nil

		result.PaymentStatusCode = "FAILED"
		return result, nil
	}

	currentBalance, err := r.getWalletBalanceForUpdate(ctx, tx, walletID)
	if err != nil {
		return nil, err
	}

	const successQuery = `
	UPDATE payment_topup_transaction ptt
	JOIN master_payment_status mps
		ON mps.status_code = 'SUCCESS'
	SET
		ptt.payment_status_id = mps.id,
		ptt.transaction_status = 'SUCCESS',
		ptt.paid_at = NOW(),
		ptt.completed_at = NOW(),
		ptt.failed_reason = NULL,
		ptt.updated_at = NOW()
	WHERE ptt.id = ?;
	`
	if _, err := tx.ExecContext(ctx, successQuery, topupTransactionID); err != nil {
		return nil, fmt.Errorf("update topup success callback: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`UPDATE wallet_account SET current_balance_amount = current_balance_amount + ?, updated_at = NOW() WHERE id = ?;`,
		amount,
		walletID,
	); err != nil {
		return nil, fmt.Errorf("update wallet balance: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`INSERT INTO wallet_history (
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
	);`,
		walletID,
		userID,
		topupTransactionID,
		amount,
		currentBalance,
		currentBalance+float64(amount),
	); err != nil {
		return nil, fmt.Errorf("insert wallet history: %w", err)
	}

	if _, err := tx.ExecContext(
		ctx,
		`
		UPDATE summary_user_home suh
		JOIN wallet_account wa
			ON wa.user_id = suh.user_id
		JOIN master_wallet_type mwt
			ON mwt.id = wa.wallet_type_id
		SET
			suh.saldo = wa.current_balance_amount,
			suh.updated_at = NOW()
		WHERE suh.user_id = ?
		AND mwt.wallet_type_code = 'BALANCE';
		`,
		userID,
	); err != nil {
		return nil, fmt.Errorf("update summary user home: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit topup success callback: %w", err)
	}
	tx = nil

	// Success status id is updated by SQL; response keeps the current known code.
	result.PaymentStatusCode = "SUCCESS"
	return result, nil
}

func (r *TopUpRepositoryImpl) getWalletForTopUp(ctx context.Context, tx *sql.Tx, userID int64) (int64, float64, error) {
	const query = `
	SELECT
		wa.id AS walletId,
		wa.current_balance_amount AS currentBalance
	FROM wallet_account wa
	JOIN master_wallet_type mwt
		ON mwt.id = wa.wallet_type_id
	WHERE wa.user_id = ?
	AND mwt.wallet_type_code = 'BALANCE'
	AND wa.is_active = 1
	LIMIT 1;
	`

	var (
		walletID       int64
		currentBalance float64
	)

	if err := tx.QueryRowContext(ctx, query, userID).Scan(&walletID, &currentBalance); err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, fmt.Errorf("wallet balance account not found")
		}
		return 0, 0, fmt.Errorf("get wallet topup: %w", err)
	}

	return walletID, currentBalance, nil
}

func (r *TopUpRepositoryImpl) getWalletBalanceForUpdate(ctx context.Context, tx *sql.Tx, walletID int64) (float64, error) {
	var currentBalance sql.NullFloat64
	if err := tx.QueryRowContext(
		ctx,
		`SELECT current_balance_amount FROM wallet_account WHERE id = ? FOR UPDATE;`,
		walletID,
	).Scan(&currentBalance); err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("wallet not found")
		}
		return 0, fmt.Errorf("get wallet balance for update: %w", err)
	}

	return utils.NullFloat64Value(currentBalance), nil
}
