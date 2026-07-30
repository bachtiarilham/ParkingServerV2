package paymentgate

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"
)

type PaymentCallbackRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentCallbackRepositoryImpl(db *sql.DB) repository.PaymentCallbackRepository {
	return &PaymentCallbackRepositoryImpl{db: db}
}

func (r *PaymentCallbackRepositoryImpl) GetPaymentTransaction(ctx context.Context, txCode string) (*model.PaymentTransactionModel, error) {
	query := `
	SELECT 
		payment_type, 
		user_id, 
		COALESCE(reference_id, 0), 
		partner_share,
		Amount
	FROM payment_transaction 
	WHERE transaction_code = ? 
	LIMIT 1;
	`
	var txModel model.PaymentTransactionModel
	err := r.db.QueryRowContext(ctx, query, txCode).Scan(
		&txModel.PaymentType,
		&txModel.UserID,
		&txModel.ReferenceID,
		&txModel.PartnerShare,
		&txModel.Amount,
	)
	if err != nil {
		return nil, err
	}
	return &txModel, nil
}

func (r *PaymentCallbackRepositoryImpl) GetPaymentStatus(ctx context.Context, txCode string) (status string, err error) {
	query := `SELECT transaction_status FROM payment_transaction WHERE transaction_code = ? LIMIT 1;`
	err = r.db.QueryRowContext(ctx, query, txCode).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *PaymentCallbackRepositoryImpl) GetPaymentStatusCash(ctx context.Context, sessionCode string) (status string, err error) {
	query := `SELECT parking_status_id FROM parking_session WHERE session_code = ? LIMIT 1;`
	err = r.db.QueryRowContext(ctx, query, sessionCode).Scan(&status)
	if err != nil {
		return "", err
	}
	return status, nil
}

func (r *PaymentCallbackRepositoryImpl) ProcessParkingCallback(ctx context.Context, txCode string, extRef string, sessionID int64, req model.PaymentTransactionModel) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin parking callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Update payment transaction
	updateTxQuery := `
	UPDATE payment_transaction 
	SET transaction_status = 'SUCCESS', paid_at = NOW(), external_reference = ?
	WHERE transaction_code = ?;
	`
	if _, err = tx.ExecContext(ctx, updateTxQuery, extRef, txCode); err != nil {
		return fmt.Errorf("update payment transaction: %w", err)
	}

	// 2. Update parking session status to Completed and assign customer_user_id
	updateSessionQuery := `
	UPDATE parking_session 
	SET customer_user_id = ?, parking_status_id = 3
	WHERE id = ?;
	`
	if _, err = tx.ExecContext(ctx, updateSessionQuery, req.UserID, sessionID); err != nil {
		return fmt.Errorf("update parking session: %w", err)
	}

	updateBalanceQuery := `
		UPDATE wallet_account 
		SET current_balance_amount = current_balance_amount + ? 
		WHERE user_id = ? AND wallet_type_id = 1`

	if _, err = tx.ExecContext(ctx, updateBalanceQuery, req.PartnerShare, req.UserID); err != nil {
		return fmt.Errorf("update parking session: %w", err)
	}
	// 2. Catat Riwayat Mutasi ke wallet_history
	insertHistoryQuery := `
		INSERT INTO wallet_history 
			(wallet_id, user_id, reference_type, reference_id, mutation_type, amount, balance_before, balance_after, description)
		SELECT 
			id, 
			user_id, 
			'PARKING', 
			?, 
			'CREDIT', 
			?, 
			current_balance_amount - ?, 
			current_balance_amount, 
			'Pembayaran Parkir'
		FROM wallet_account
		WHERE user_id = ? AND wallet_type_id = 1`
	if _, err = tx.ExecContext(ctx, insertHistoryQuery, sessionID, req.PartnerShare, req.PartnerShare, req.UserID); err != nil {
		return fmt.Errorf("update parking session: %w", err)

	}

	// var lokasiid int64
	// ambilokasi :=
	// 	`
	// SELECT
	// 	location_id
	// 	FROM parking_session
	// 	WHERE id = ?
	// LIMIT 1;
	// `
	// if err := tx.QueryRowContext(ctx, ambilokasi, sessionID).Scan(
	// 	&lokasiid,
	// ); err != nil {
	// 	return fmt.Errorf("ambil lokasi id: %w", err)
	// }

	// summarydailyquery :=
	// 	`
	// INSERT INTO summary_officer_daily_report
	// 	(report_date, officer_user_id, location_id, total_transaction, total_jukir_share, motor_count)
	// VALUES
	// 	(CURRENT_DATE(), ?, ?, 1, ?, 1)
	// ON DUPLICATE KEY UPDATE
	// 	total_transaction = total_transaction + 1,
	// 	total_jukir_share = total_jukir_share + ?,
	// 	motor_count = motor_count + 1
	// `
	// if _, err = tx.ExecContext(ctx, summarydailyquery, req.UserID, lokasiid, req.PartnerShare, req.PartnerShare); err != nil {
	// 	return fmt.Errorf("update summary: %w", err)

	// }

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit parking callback tx: %w", err)
	}

	return nil
}

func (r *PaymentCallbackRepositoryImpl) ProcessTopupCallback(ctx context.Context, txCode string, extRef string, userID int64, amount int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin topup callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Update payment transaction
	updateTxQuery := `
	UPDATE payment_transaction 
	SET transaction_status = 'SUCCESS', paid_at = NOW(), external_reference = ?
	WHERE transaction_code = ?;
	`
	if _, err = tx.ExecContext(ctx, updateTxQuery, extRef, txCode); err != nil {
		return fmt.Errorf("update payment transaction: %w", err)
	}

	// 2. Add balance to user wallet
	updateWalletQuery := `
	UPDATE wallet_account 
	SET current_balance_amount = current_balance_amount + ? 
	WHERE user_id = ? AND wallet_type_id = 1;
	`
	if _, err = tx.ExecContext(ctx, updateWalletQuery, amount, userID); err != nil {
		return fmt.Errorf("update wallet balance: %w", err)
	}

	// 3. Log mutation to wallet_history
	historyQuery := `
	INSERT INTO wallet_history 
	  (wallet_id, user_id, reference_type, reference_id, mutation_type, amount, balance_before, balance_after, description)
	SELECT id, user_id, 'TOPUP', 1, 'CREDIT', ?, current_balance_amount - ?, current_balance_amount, 'Isi Ulang Saldo via Midtrans'
	FROM wallet_account
	WHERE user_id = ? AND wallet_type_id = 1;
	`
	if _, err = tx.ExecContext(ctx, historyQuery, amount, amount, userID); err != nil {
		return fmt.Errorf("insert wallet history log: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit topup callback tx: %w", err)
	}

	return nil
}

func (r *PaymentCallbackRepositoryImpl) ProcessMembershipCallback(ctx context.Context, txCode string, extRef string, userID int64, packageID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin membership callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Update payment transaction
	updateTxQuery := `
	UPDATE payment_transaction 
	SET transaction_status = 'SUCCESS', paid_at = NOW(), external_reference = ?
	WHERE transaction_code = ?;
	`
	if _, err = tx.ExecContext(ctx, updateTxQuery, extRef, txCode); err != nil {
		return fmt.Errorf("update payment transaction: %w", err)
	}

	// 2. Insert/update membership status
	membershipQuery := `
	INSERT INTO membership_user 
	  (user_id, package_id, activated_at, expired_at, status)
	VALUES 
	  (?, ?, NOW(), DATE_ADD(NOW(), INTERVAL 30 DAY), 'ACTIVE')
	ON DUPLICATE KEY UPDATE 
	  expired_at = DATE_ADD(expired_at, INTERVAL 30 DAY), status = 'ACTIVE';
	`
	if _, err = tx.ExecContext(ctx, membershipQuery, userID, packageID); err != nil {
		return fmt.Errorf("insert or update membership: %w", err)
	}

	// buatschquery := `
	// 	SELECT
	// 		mp.package_name AS membership_package_name,
	// 		mu.expired_at AS membership_expired_at
	// 	FROM membership_user mu
	// 	JOIN membership_package mp
	// 		ON mu.package_id = mp.id
	// 	WHERE mu.user_id = ?
	// 	  AND mu.package_id = ?
	// 	ORDER BY mu.expired_at DESC
	// 	LIMIT 1
	// `
	// var packageName string
	// var expiredAt time.Time
	// err = r.db.QueryRowContext(ctx, buatschquery, userID, packageID).Scan(&packageName, &expiredAt)
	// if err != nil {
	// 	return err
	// }

	// updateSummary := `
	// UPDATE summary_customer_home
	// SET membership_package_name = ?, membership_expired_at = ?
	// WHERE user_id = ?;
	// `
	// if _, err = tx.ExecContext(ctx, updateSummary, packageName, expiredAt, userID); err != nil {
	// 	return fmt.Errorf("update payment transaction: %w", err)
	// }

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit membership callback tx: %w", err)
	}

	return nil
}
