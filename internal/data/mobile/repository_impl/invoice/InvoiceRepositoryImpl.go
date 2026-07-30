package invoice

import (
	"context"
	"database/sql"
	"fmt"

	model "modulegue/internal/domain/mobile/model/invoice"
	"modulegue/internal/domain/mobile/repository"
)

type InvoiceRepositoryImpl struct {
	db *sql.DB
}

func NewInvoiceRepositoryImpl(db *sql.DB) repository.InvoiceRepository {
	return &InvoiceRepositoryImpl{db: db}
}

func (r *InvoiceRepositoryImpl) GetInvoice(ctx context.Context, code string) (*model.UniversalInvoiceResponseModel, error) {
	// First check the payment_type in payment_transaction
	var paymentType string
	err := r.db.QueryRowContext(ctx, "SELECT payment_type FROM payment_transaction WHERE transaction_code = ? LIMIT 1", code).Scan(&paymentType)
	if err != nil {
		return nil, fmt.Errorf("lookup payment type: %w", err)
	}

	switch paymentType {
	case "PARKING":
		return r.getParkingInvoice(ctx, code)
	case "TOPUP":
		return r.getTopupInvoice(ctx, code)
	case "TRANSFER":
		return r.getTransferInvoice(ctx, code)
	case "MEMBERSHIP":
		return r.getMembershipInvoice(ctx, code)
	default:
		return nil, fmt.Errorf("unknown payment type: %s", paymentType)
	}
}

func (r *InvoiceRepositoryImpl) getParkingInvoice(ctx context.Context, txCode string) (*model.UniversalInvoiceResponseModel, error) {
	query := `
SELECT 
    -- 1. Main Invoice fields
    pt.transaction_code,
    pt.payment_type AS transaction_type,
    'Parkir On-Street' AS title,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'PAID'
        WHEN pt.transaction_status = 'PENDING' THEN 'PENDING'
        ELSE 'FAILED'
    END AS status,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'Lunas'
        WHEN pt.transaction_status = 'PENDING' THEN 'Menunggu Pembayaran'
        ELSE 'Gagal'
    END AS status_text,
    pt.created_at AS created_at,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS total_amount,
    
    -- 2. Parking Details fields
    lp.location_name AS parking_location_name,
    ps.plate_number AS parking_license_plate,
    mvt.vehicle_type_name AS parking_vehicle_type,
    ps.started_at AS parking_check_in_time,
    ps.finished_at AS parking_check_out_time,
    CONCAT(
        TIMESTAMPDIFF(HOUR, ps.started_at, COALESCE(ps.finished_at, NOW())), ' Jam ',
        MOD(TIMESTAMPDIFF(MINUTE, ps.started_at, COALESCE(ps.finished_at, NOW())), 60), ' Menit'
    ) AS parking_duration_text,
    officer.full_name AS parking_attendant_name,
    -- 3. Price Breakdown fields
    pt.amount AS base_price,
    pt.fee_amount AS admin_fee,
    pt.tax_amount AS tax_amount,
    0 AS discount_amount,
    '' AS discount_code,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS final_total,
    -- 4. Payment Method fields
    mpm.payment_method_code AS channel_code,
    mpm.payment_method_name AS channel_name,
    CASE 
        WHEN mpm.payment_method_code = 'WALLET' THEN ui.phone_number
        ELSE pt.external_reference
    END AS account_no,
    mpm.logo_url AS icon_url,
    -- 5. Customer Info fields
    CAST(ui.id AS CHAR) AS user_id,
    ui.full_name AS full_name,
    ui.email AS email,
    ui.phone_number AS phone
FROM payment_transaction pt
JOIN parking_session ps ON pt.reference_id = ps.id
JOIN location_parking lp ON ps.location_id = lp.id
JOIN master_vehicle_type mvt ON ps.vehicle_type_id = mvt.id
JOIN user_identity ui ON pt.user_id = ui.id
JOIN master_payment_method mpm ON pt.payment_method_id = mpm.id
LEFT JOIN user_identity officer ON ps.officer_user_id = officer.id
WHERE pt.transaction_code = ? AND pt.payment_type = 'PARKING';
`
	var (
		res           model.UniversalInvoiceResponseModel
		detail        model.ParkingInvoiceDetailModel
		price         model.PriceBreakdownModel
		payment       model.PaymentMethodModel
		customer      model.CustomerInfoModel
		checkOutTime  sql.NullTime
		attendantName sql.NullString
		email         sql.NullString
		phone         sql.NullString
		accountNo     sql.NullString
		iconUrl       sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, txCode).Scan(
		&res.TrxCode,
		&res.TransactionType,
		&res.Title,
		&res.Status,
		&res.StatusText,
		&res.CreatedAt,
		&res.TotalAmount,
		&detail.LocationName,
		&detail.LicensePlate,
		&detail.VehicleType,
		&detail.CheckInTime,
		&checkOutTime,
		&detail.DurationText,
		&attendantName,
		&price.BasePrice,
		&price.AdminFee,
		&price.TaxAmount,
		&price.DiscountAmount,
		&price.DiscountCode,
		&price.FinalTotal,
		&payment.ChannelCode,
		&payment.ChannelName,
		&accountNo,
		&iconUrl,
		&customer.UserID,
		&customer.FullName,
		&email,
		&phone,
	)
	if err != nil {
		return nil, err
	}

	if checkOutTime.Valid {
		detail.CheckOutTime = &checkOutTime.Time
	}
	detail.AttendantName = attendantName.String
	payment.AccountNo = accountNo.String
	payment.IconURL = iconUrl.String
	customer.Email = email.String
	customer.Phone = phone.String
	res.ParkingDetails = &detail
	res.PriceBreakdown = price
	res.PaymentMethod = payment
	res.CustomerInfo = customer

	return &res, nil
}

func (r *InvoiceRepositoryImpl) getTopupInvoice(ctx context.Context, txCode string) (*model.UniversalInvoiceResponseModel, error) {
	query := `
SELECT 
    -- 1. Main Invoice fields
    pt.transaction_code,
    pt.payment_type AS transaction_type,
    'Top Up Saldo' AS title,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'PAID'
        WHEN pt.transaction_status = 'PENDING' THEN 'PENDING'
        ELSE 'FAILED'
    END AS status,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'Lunas'
        WHEN pt.transaction_status = 'PENDING' THEN 'Menunggu Pembayaran'
        ELSE 'Gagal'
    END AS status_text,
    pt.created_at AS created_at,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS total_amount,
    
    -- 2. Wallet Details (Top-Up) fields
    mpm.payment_method_name AS wallet_sender_name,
    COALESCE(pt.external_reference, '-') AS wallet_sender_account,
    'Saldo LineSpot' AS wallet_recipient_name,
    ui.phone_number AS wallet_recipient_account,
    pt.external_reference AS wallet_bank_ref_no,
    -- 3. Price Breakdown fields
    pt.amount AS base_price,
    pt.fee_amount AS admin_fee,
    pt.tax_amount AS tax_amount,
    0 AS discount_amount,
    '' AS discount_code,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS final_total,
    -- 4. Payment Method fields
    mpm.payment_method_code AS channel_code,
    mpm.payment_method_name AS channel_name,
    COALESCE(pt.external_reference, '-') AS account_no,
    mpm.logo_url AS icon_url,
    -- 5. Customer Info fields
    CAST(ui.id AS CHAR) AS user_id,
    ui.full_name AS full_name,
    ui.email AS email,
    ui.phone_number AS phone
FROM payment_transaction pt
JOIN user_identity ui ON pt.user_id = ui.id
JOIN master_payment_method mpm ON pt.payment_method_id = mpm.id
WHERE pt.transaction_code = ? AND pt.payment_type = 'TOPUP';
`
	var (
		res       model.UniversalInvoiceResponseModel
		detail    model.WalletInvoiceDetailModel
		price     model.PriceBreakdownModel
		payment   model.PaymentMethodModel
		customer  model.CustomerInfoModel
		email     sql.NullString
		phone     sql.NullString
		accountNo sql.NullString
		iconUrl   sql.NullString
		bankRefNo sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, txCode).Scan(
		&res.TrxCode,
		&res.TransactionType,
		&res.Title,
		&res.Status,
		&res.StatusText,
		&res.CreatedAt,
		&res.TotalAmount,
		&detail.SenderName,
		&detail.SenderAccount,
		&detail.RecipientName,
		&detail.RecipientAccount,
		&bankRefNo,
		&price.BasePrice,
		&price.AdminFee,
		&price.TaxAmount,
		&price.DiscountAmount,
		&price.DiscountCode,
		&price.FinalTotal,
		&payment.ChannelCode,
		&payment.ChannelName,
		&accountNo,
		&iconUrl,
		&customer.UserID,
		&customer.FullName,
		&email,
		&phone,
	)
	if err != nil {
		return nil, err
	}

	detail.BankRefNo = bankRefNo.String
	payment.AccountNo = accountNo.String
	payment.IconURL = iconUrl.String
	customer.Email = email.String
	customer.Phone = phone.String
	res.WalletDetails = &detail
	res.PriceBreakdown = price
	res.PaymentMethod = payment
	res.CustomerInfo = customer

	return &res, nil
}

func (r *InvoiceRepositoryImpl) getTransferInvoice(ctx context.Context, txCode string) (*model.UniversalInvoiceResponseModel, error) {
	query := `
SELECT 
    -- 1. Main Invoice fields
    pt.transaction_code,
    pt.payment_type AS transaction_type,
    'Transfer Saldo' AS title,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'PAID'
        WHEN pt.transaction_status = 'PENDING' THEN 'PENDING'
        ELSE 'FAILED'
    END AS status,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'Lunas'
        WHEN pt.transaction_status = 'PENDING' THEN 'Menunggu Pembayaran'
        ELSE 'Gagal'
    END AS status_text,
    pt.created_at AS created_at,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS total_amount,
    
    -- 2. Wallet Details (Transfer) fields
    sender.full_name AS wallet_sender_name,
    sender.phone_number AS wallet_sender_account,
    receiver.full_name AS wallet_recipient_name,
    receiver.phone_number AS wallet_recipient_account,
    '' AS wallet_bank_ref_no,
    -- 3. Price Breakdown fields
    pt.amount AS base_price,
    pt.fee_amount AS admin_fee,
    pt.tax_amount AS tax_amount,
    0 AS discount_amount,
    '' AS discount_code,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS final_total,
    -- 4. Payment Method fields
    mpm.payment_method_code AS channel_code,
    mpm.payment_method_name AS channel_name,
    sender.phone_number AS account_no,
    mpm.logo_url AS icon_url,
    -- 5. Customer Info (Pengirim) fields
    CAST(sender.id AS CHAR) AS user_id,
    sender.full_name AS full_name,
    sender.email AS email,
    sender.phone_number AS phone
FROM payment_transaction pt
JOIN user_identity sender ON pt.user_id = sender.id
JOIN user_identity receiver ON pt.reference_id = receiver.id
JOIN master_payment_method mpm ON pt.payment_method_id = mpm.id
WHERE pt.transaction_code = ? AND pt.payment_type = 'TRANSFER';
`
	var (
		res       model.UniversalInvoiceResponseModel
		detail    model.WalletInvoiceDetailModel
		price     model.PriceBreakdownModel
		payment   model.PaymentMethodModel
		customer  model.CustomerInfoModel
		email     sql.NullString
		phone     sql.NullString
		accountNo sql.NullString
		iconUrl   sql.NullString
		bankRefNo sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, txCode).Scan(
		&res.TrxCode,
		&res.TransactionType,
		&res.Title,
		&res.Status,
		&res.StatusText,
		&res.CreatedAt,
		&res.TotalAmount,
		&detail.SenderName,
		&detail.SenderAccount,
		&detail.RecipientName,
		&detail.RecipientAccount,
		&bankRefNo,
		&price.BasePrice,
		&price.AdminFee,
		&price.TaxAmount,
		&price.DiscountAmount,
		&price.DiscountCode,
		&price.FinalTotal,
		&payment.ChannelCode,
		&payment.ChannelName,
		&accountNo,
		&iconUrl,
		&customer.UserID,
		&customer.FullName,
		&email,
		&phone,
	)
	if err != nil {
		return nil, err
	}

	detail.BankRefNo = bankRefNo.String
	payment.AccountNo = accountNo.String
	payment.IconURL = iconUrl.String
	customer.Email = email.String
	customer.Phone = phone.String
	res.WalletDetails = &detail
	res.PriceBreakdown = price
	res.PaymentMethod = payment
	res.CustomerInfo = customer

	return &res, nil
}

func (r *InvoiceRepositoryImpl) getMembershipInvoice(ctx context.Context, txCode string) (*model.UniversalInvoiceResponseModel, error) {
	query := `
SELECT 
    -- 1. Main Invoice fields
    pt.transaction_code,
    pt.payment_type AS transaction_type,
    'Membership Langganan' AS title,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'PAID'
        WHEN pt.transaction_status = 'PENDING' THEN 'PENDING'
        ELSE 'FAILED'
    END AS status,
    CASE 
        WHEN pt.transaction_status = 'SUCCESS' THEN 'Lunas'
        WHEN pt.transaction_status = 'PENDING' THEN 'Menunggu Pembayaran'
        ELSE 'Gagal'
    END AS status_text,
    pt.created_at AS created_at,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS total_amount,
    
    -- 2. Membership Details fields
    mp.package_name AS membership_package_name,
    mu.activated_at AS membership_period_start,
    mu.expired_at AS membership_period_end,
    2 AS membership_max_vehicles, -- Default / Hardcoded mapping
    FALSE AS membership_is_auto_renew,
    mu.expired_at AS membership_next_billing_date,
    -- 3. Price Breakdown fields
    pt.amount AS base_price,
    pt.fee_amount AS admin_fee,
    pt.tax_amount AS tax_amount,
    0 AS discount_amount,
    '' AS discount_code,
    (pt.amount + pt.fee_amount + pt.tax_amount) AS final_total,
    -- 4. Payment Method fields
    mpm.payment_method_code AS channel_code,
    mpm.payment_method_name AS channel_name,
    CASE 
        WHEN mpm.payment_method_code = 'WALLET' THEN ui.phone_number
        ELSE pt.external_reference
    END AS account_no,
    mpm.logo_url AS icon_url,
    -- 5. Customer Info fields
    CAST(ui.id AS CHAR) AS user_id,
    ui.full_name AS full_name,
    ui.email AS email,
    ui.phone_number AS phone
FROM payment_transaction pt
JOIN membership_user mu ON pt.reference_id = mu.id
JOIN membership_package mp ON mu.package_id = mp.id
JOIN user_identity ui ON pt.user_id = ui.id
JOIN master_payment_method mpm ON pt.payment_method_id = mpm.id
WHERE pt.transaction_code = ? AND pt.payment_type = 'MEMBERSHIP';
`
	var (
		res             model.UniversalInvoiceResponseModel
		detail          model.MembershipInvoiceDetailModel
		price           model.PriceBreakdownModel
		payment         model.PaymentMethodModel
		customer        model.CustomerInfoModel
		nextBillingDate sql.NullTime
		email           sql.NullString
		phone           sql.NullString
		accountNo       sql.NullString
		iconUrl         sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, txCode).Scan(
		&res.TrxCode,
		&res.TransactionType,
		&res.Title,
		&res.Status,
		&res.StatusText,
		&res.CreatedAt,
		&res.TotalAmount,
		&detail.PackageName,
		&detail.PeriodStart,
		&detail.PeriodEnd,
		&detail.MaxVehicles,
		&detail.IsAutoRenew,
		&nextBillingDate,
		&price.BasePrice,
		&price.AdminFee,
		&price.TaxAmount,
		&price.DiscountAmount,
		&price.DiscountCode,
		&price.FinalTotal,
		&payment.ChannelCode,
		&payment.ChannelName,
		&accountNo,
		&iconUrl,
		&customer.UserID,
		&customer.FullName,
		&email,
		&phone,
	)
	if err != nil {
		return nil, err
	}

	if nextBillingDate.Valid {
		detail.NextBillingDate = &nextBillingDate.Time
	}
	payment.AccountNo = accountNo.String
	payment.IconURL = iconUrl.String
	customer.Email = email.String
	customer.Phone = phone.String
	res.MembershipDetails = &detail
	res.PriceBreakdown = price
	res.PaymentMethod = payment
	res.CustomerInfo = customer

	return &res, nil
}
