package paymentgate

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "modulegue/internal/domain/mobile/model/payment_gate"
	"modulegue/internal/domain/mobile/repository"

	"github.com/google/uuid"
)

type PaymentGateRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentGateRepositoryImpl(db *sql.DB) repository.PaymentGateRepository {
	return &PaymentGateRepositoryImpl{db: db}
}

func (r *PaymentGateRepositoryImpl) PayParking(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error) {
	txCode := fmt.Sprintf("TRX-PRK-%s", uuid.New().String())

	var sessionID int64
	if req.TargetID == nil || *req.TargetID == "" {
		return nil, fmt.Errorf("target_id (session code) tidak boleh kosong")
	}

	// 1. Fetch session ID by session code (TargetID)
	getSessionQuery := `SELECT id FROM parking_session WHERE session_code = ? LIMIT 1;`
	err := r.db.QueryRowContext(ctx, getSessionQuery, *req.TargetID).Scan(&sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sesi parkir dengan kode %s tidak ditemukan", *req.TargetID)
		}
		return nil, fmt.Errorf("gagal mendapatkan data sesi parkir: %w", err)
	}

	// 2. Insert into payment_transaction using the retrieved sessionID
	query := `
	INSERT INTO payment_transaction 
	  (transaction_code, user_id, payment_type, reference_id, payment_method_id, amount, partner_share, company_share, gov_share, midtrans_share, transaction_status, created_at, updated_at)
	VALUES 
	  (?, ?, 'PARKING', ?, COALESCE((SELECT id FROM master_payment_method WHERE payment_method_code = ? LIMIT 1), 1), ?, ?, ?,?,?, 'PENDING', NOW(), NOW());
	`
	_, err = r.db.ExecContext(ctx, query, txCode, req.UserID, sessionID, req.PaymentMethodCode, req.Amount, req.JukirShare, req.CompanyShare, req.GovShare, req.MidtransShare)
	if err != nil {
		return nil, fmt.Errorf("pay parking: %w", err)
	}

	return &model.PayResponseModel{
		OrderID:       txCode,
		GrossAmount:   req.Amount,
		PaymentMethod: req.PaymentMethodCode,
		Status:        "PENDING",
	}, nil
}

func (r *PaymentGateRepositoryImpl) PayTransfer(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error) {
	txCode := fmt.Sprintf("TRX-TRF-%s", uuid.New().String())
	var receiverID int64
	if req.TargetID != nil {
		fmt.Sscanf(*req.TargetID, "%d", &receiverID)
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transfer tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Deduct sender balance
	deductQuery := `UPDATE wallet_account SET current_balance_amount = current_balance_amount - ? WHERE user_id = ? AND wallet_type_id = 1;`
	if _, err = tx.ExecContext(ctx, deductQuery, req.Amount, req.UserID); err != nil {
		return nil, fmt.Errorf("deduct sender balance: %w", err)
	}

	// Add receiver balance
	addQuery := `UPDATE wallet_account SET current_balance_amount = current_balance_amount + ? WHERE user_id = ? AND wallet_type_id = 1;`
	if _, err = tx.ExecContext(ctx, addQuery, req.Amount, receiverID); err != nil {
		return nil, fmt.Errorf("add receiver balance: %w", err)
	}

	// Insert transaction log
	insertQuery := `
	INSERT INTO payment_transaction 
	  (transaction_code, user_id, payment_type, reference_id, payment_method_id, amount, transaction_status, paid_at, created_at, updated_at)
	VALUES 
	  (?, ?, 'TRANSFER', ?, COALESCE((SELECT id FROM master_payment_method WHERE payment_method_code = ? LIMIT 1), 3), ?, 'SUCCESS', NOW(), NOW(), NOW());
	`
	if _, err = tx.ExecContext(ctx, insertQuery, txCode, req.UserID, receiverID, req.PaymentMethodCode, req.Amount); err != nil {
		return nil, fmt.Errorf("insert transfer transaction ledger: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transfer tx: %w", err)
	}

	return &model.PayResponseModel{
		OrderID:       txCode,
		GrossAmount:   req.Amount,
		PaymentMethod: req.PaymentMethodCode,
		Status:        "SUCCESS",
	}, nil
}

func (r *PaymentGateRepositoryImpl) PayMembership(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error) {
	txCode := fmt.Sprintf("TRX-MBR-%s", uuid.New().String())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin parking callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var packageID int64
	// if req.TargetID != nil {
	// 	fmt.Sscanf(*req.TargetID, "%d", &packageID)
	if req.TargetID == nil || *req.TargetID == "" {
		return nil, fmt.Errorf("target_id (session code) tidak boleh kosong")
	}

	getPackageQuery := `SELECT id FROM membership_package WHERE package_name = ? LIMIT 1;`
	err = r.db.QueryRowContext(ctx, getPackageQuery, *req.TargetID).Scan(&packageID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("sesi parkir dengan kode %s tidak ditemukan", *req.TargetID)
		}
		return nil, fmt.Errorf("gagal mendapatkan data sesi parkir: %w", err)
	}

	query := `
	INSERT INTO payment_transaction 
	  (transaction_code, user_id, payment_type, reference_id, payment_method_id, amount, transaction_status, created_at, updated_at)
	VALUES 
	  (?, ?, 'MEMBERSHIP', ?, COALESCE((SELECT id FROM master_payment_method WHERE payment_method_code = ? LIMIT 1), 1), ?, 'PENDING', NOW(), NOW());
	`
	_, err = r.db.ExecContext(ctx, query, txCode, req.UserID, packageID, req.PaymentMethodCode, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("pay membership: %w", err)
	}

	var duration time.Duration
	durationQuery := `
	SELECT duration_days FROM membership_package WHERE id = ? LIMIT 1
	`
	err = r.db.QueryRowContext(ctx, durationQuery, packageID).Scan(&duration)
	if err != nil {
		return nil, fmt.Errorf("pay membership: %w", err)
	}

	membershipQuery := `
	INSERT INTO membership_user 
	  (user_id, package_id, activated_at, expired_at, status, created_at, updated_at)
	VALUES 
	  (?, ?, NOW(), ? , 'ACTIVE', NOW(), NOW());
	`
	_, err = r.db.ExecContext(ctx, membershipQuery, req.UserID, packageID, (time.Now().Add(duration * 24 * time.Hour)))
	if err != nil {
		return nil, fmt.Errorf("pay membership: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit parking callback tx: %w", err)
	}

	return &model.PayResponseModel{
		OrderID:       txCode,
		GrossAmount:   req.Amount,
		PaymentMethod: req.PaymentMethodCode,
		Status:        "PENDING",
	}, nil
}

func (r *PaymentGateRepositoryImpl) PayTopUp(ctx context.Context, req model.PayRequestModel) (*model.PayResponseModel, error) {
	txCode := fmt.Sprintf("TRX-TOPUP-%s", uuid.New().String())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin parking callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	query := `
	INSERT INTO payment_transaction 
	  (transaction_code, user_id, payment_type, reference_id, payment_method_id, amount, fee_amount, transaction_status, created_at, updated_at)
	VALUES 
	  (?, ?, 'TOPUP', NULL, COALESCE((SELECT id FROM master_payment_method WHERE payment_method_code = ? LIMIT 1), 1), ?, 0, 'PENDING', NOW(), NOW());
	`
	_, err = r.db.ExecContext(ctx, query, txCode, req.UserID, req.PaymentMethodCode, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("pay topup: %w", err)
	}

	// updateBalanceQuery := `
	// 	UPDATE wallet_account
	// 	SET current_balance_amount = current_balance_amount + ?
	// 	WHERE user_id = ? AND wallet_type_id = 1`

	// if _, err = tx.ExecContext(ctx, updateBalanceQuery, req.Amount, req.UserID); err != nil {
	// 	return nil, fmt.Errorf("update wallet account: %w", err)
	// }
	// // 2. Catat Riwayat Mutasi ke wallet_history
	// insertHistoryQuery := `
	// 	INSERT INTO wallet_history
	// 		(wallet_id, user_id, reference_type, reference_id, mutation_type, amount, balance_before, balance_after, description)
	// 	SELECT
	// 		id,
	// 		user_id,
	// 		'PARKING',
	// 		null,
	// 		'DEBIT',
	// 		?,
	// 		current_balance_amount + ?,
	// 		current_balance_amount,
	// 		'Pembayaran Parkir'
	// 	FROM wallet_account
	// 	WHERE user_id = ? AND wallet_type_id = 1`
	// if _, err = tx.Exec(insertHistoryQuery, req.Amount, req.Amount, req.UserID); err != nil {
	// 	return nil, fmt.Errorf("update wallet history: %w", err)

	// }

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit parking callback tx: %w", err)
	}

	return &model.PayResponseModel{
		OrderID:       txCode,
		GrossAmount:   req.Amount,
		PaymentMethod: req.PaymentMethodCode,
		Status:        "PENDING",
	}, nil
}

func (r *PaymentGateRepositoryImpl) PayCashParking(ctx context.Context, req model.PayRequestModel) (string, *model.PayResponseModel, error) {
	txCode := fmt.Sprintf("TRX-PRK-%s", uuid.New().String())

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return "", nil, fmt.Errorf("begin parking callback tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	var sessionID int64
	getSessionQuery := `SELECT id FROM parking_session WHERE session_code = ? LIMIT 1;`
	err = tx.QueryRowContext(ctx, getSessionQuery, *req.TargetID).Scan(&sessionID)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil, fmt.Errorf("sesi parkir dengan kode %s tidak ditemukan", *req.TargetID)
		}
		return "", nil, fmt.Errorf("gagal mendapatkan data sesi parkir: %w", err)
	}

	query := `
	INSERT INTO payment_transaction 
	  (transaction_code, user_id, payment_type, reference_id, payment_method_id, amount, partner_share, company_share, gov_share, midtrans_share, transaction_status, created_at, updated_at, paid_at,external_reference)
	VALUES 
	  (?, ?, 'PARKING', ?, COALESCE((SELECT id FROM master_payment_method WHERE payment_method_code = ? LIMIT 1), 1), ?, ?, ?,?,?, 'SUCCESS', NOW(), NOW(),NOW(),'cash');
	`
	_, err = r.db.ExecContext(ctx, query, txCode, req.UserID, sessionID, req.PaymentMethodCode, req.Amount, req.JukirShare, req.CompanyShare, req.GovShare, req.MidtransShare)
	if err != nil {
		return "", nil, fmt.Errorf("pay parking: %w", err)
	}

	// 2. Update parking session status to Completed and assign customer_user_id
	updateSessionQuery := `
	UPDATE parking_session 
	SET customer_user_id = ?, parking_status_id = 3
	WHERE id = ?;
	`
	if _, err = tx.ExecContext(ctx, updateSessionQuery, req.UserID, sessionID); err != nil {
		return "", nil, fmt.Errorf("update parking session: %w", err)
	}

	updateBalanceQuery := `
		UPDATE wallet_account 
		SET current_balance_amount = current_balance_amount - ? 
		WHERE user_id = ? AND wallet_type_id = 1`

	if _, err = tx.ExecContext(ctx, updateBalanceQuery, (req.CompanyShare + req.GovShare), req.UserID); err != nil {
		return "", nil, fmt.Errorf("update parking session: %w", err)
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
			'DEBIT', 
			?, 
			current_balance_amount + ?, 
			current_balance_amount, 
			'Pembayaran Parkir'
		FROM wallet_account
		WHERE user_id = ? AND wallet_type_id = 1`
	if _, err = tx.Exec(insertHistoryQuery, sessionID, (req.CompanyShare + req.GovShare), (req.CompanyShare + req.GovShare), req.UserID); err != nil {
		return "", nil, fmt.Errorf("update parking session: %w", err)

	}

	if err = tx.Commit(); err != nil {
		return "", nil, fmt.Errorf("commit parking callback tx: %w", err)
	}

	return txCode, &model.PayResponseModel{
		OrderID: txCode,
		Status:  "SUCCESS",
	}, nil
}

func (r *PaymentGateRepositoryImpl) IsMember(ctx context.Context, sessioncode string) (*model.MemberCheckResult, error) {
	query := `
	SELECT 
		-- Menghasilkan TRUE (1) jika membership aktif ditemukan, selain itu FALSE (0)
		IF(mu.id IS NOT NULL, TRUE, FALSE) AS is_member,
		mu.package_id,
		mp.package_name,
		mu.expired_at
	FROM parking_session ps
	-- 1. Hubungkan ke vehicle_registry menggunakan plat nomor dari sesi parkir
	LEFT JOIN vehicle_registry vr 
		ON vr.plate_number = ps.plate_number 
	AND vr.is_active = 1
	-- 2. Hubungkan ke membership_user berdasarkan pemilik kendaraan (vr.user_id)
	LEFT JOIN membership_user mu 
		ON mu.user_id = vr.user_id 
	AND mu.status = 'ACTIVE' 
	AND mu.expired_at > NOW()
	-- 3. Hubungkan ke detail paket untuk mendapatkan nama paketnya
	LEFT JOIN membership_package mp 
		ON mp.id = mu.package_id
	WHERE ps.session_code = ?
	LIMIT 1;
	`
	var result model.MemberCheckResult
	err := r.db.QueryRowContext(ctx, query, sessioncode).Scan(
		&result.IsMember,
		&result.PackageID,
		&result.PackageName,
		&result.ExpiredAt,
	)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
