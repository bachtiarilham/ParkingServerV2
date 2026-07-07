package repository

// import (
// 	"context"
// 	"database/sql"
// 	"fmt"
// 	"modulegue/internal/domain/payment"
// 	"time"
// )

// type PaymentRepository struct {
// 	db *sql.DB
// }

// func NewPaymentRepository(db *sql.DB) payment.Repository {
// 	return &PaymentRepository{db: db}
// }

// func (r *PaymentRepository) GetActiveSessionByCode(ctx context.Context, code string) (*payment.ParkingSession, error) {
// 	query := `
// 		SELECT id, session_code, customer_user_id, location_id, started_at, duration_minutes, session_status, plate_number, vehicle_type_id
// 		FROM parking_session
// 		WHERE session_code = ? AND session_status = 'active'
// 		LIMIT 1
// 	`
// 	var s payment.ParkingSession
// 	var durationMinutes *int
// 	err := r.db.QueryRowContext(ctx, query, code).Scan(
// 		&s.ID, &s.Code, &s.CustomerID, &s.LocationID, &s.StartedAt, &durationMinutes, &s.Status, &s.PlateNumber, &s.VehicleTypeID,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("active parking session not found for code: %s", code)
// 		}
// 		return nil, fmt.Errorf("get active session by code: %w", err)
// 	}
// 	s.DurationMinutes = durationMinutes // Assign the scanned duration if available
// 	return &s, nil
// }

// func (r *PaymentRepository) GetTariffForLocationAndVehicle(ctx context.Context, locationID int64, vehicleTypeID int64) (int64, error) {
// 	// Asumsi tariff ada di admin_location_setting, dan kita hanya ambil tariff_car_amount untuk semua kendaraan
// 	// (bisa disesuaikan dengan logika sebenarnya)
// 	query := `
// 		SELECT CASE WHEN vt.vehicle_type_code = 'motorcycle' THEN als.tariff_motor_amount ELSE als.tariff_car_amount END
// 		FROM admin_location_setting als
// 		JOIN vehicle_type vt ON 1=1 -- Simplified, adjust based on how tariff is linked
// 		WHERE als.location_id = ? AND vt.id = ?
// 		LIMIT 1
// 	`
// 	var tariff int64
// 	err := r.db.QueryRowContext(ctx, query, locationID, vehicleTypeID).Scan(&tariff)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			// Default tariff jika tidak ditemukan setting
// 			tariff = 2000 // Contoh default
// 		} else {
// 			return 0, fmt.Errorf("get tariff: %w", err)
// 		}
// 	}
// 	return tariff, nil
// }

// func (r *PaymentRepository) CreateFinancialTransaction(ctx context.Context, tx *payment.FinancialTransaction) error {
// 	// Hitung subtotal jika belum dihitung di usecase (disarankan dihitung di usecase)
// 	// subtotal = tariff * duration_hours (dibulatkan ke atas)
// 	// final_amount = subtotal - discount + penalty
// 	// Untuk sementara, kita gunakan subtotal sebagai final_amount
// 	query := `
// 		INSERT INTO financial_parking_transaction (
// 			transaction_code, operation_type, transaction_source, session_id, location_id,
// 			customer_user_id, subtotal_amount, final_amount, currency_code, transaction_status, occurred_at, created_at
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
// 	`
// 	result, err := r.db.ExecContext(ctx, query,
// 		tx.Code, tx.OperationType, tx.TransactionSource, tx.SessionID, tx.LocationID,
// 		tx.CustomerID, tx.SubtotalAmount, tx.FinalAmount, tx.CurrencyCode, tx.Status, tx.OccurredAt,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert financial transaction: %w", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return fmt.Errorf("get last insert id: %w", err)
// 	}
// 	tx.ID = id // Update ID di struct
// 	return nil
// }

// func (r *PaymentRepository) GetFinancialTransactionByCode(ctx context.Context, code string) (*payment.FinancialTransaction, error) {
// 	query := `
// 		SELECT id, transaction_code, operation_type, transaction_source, session_id, location_id,
// 		       customer_user_id, jukir_user_id, subtotal_amount, final_amount, currency_code, transaction_status, paid_at, occurred_at, created_at, successful_payment_event_id,
// 		       gov_share, company_share, jukir_share, payment_method
// 		FROM financial_parking_transaction
// 		WHERE transaction_code = ?
// 		LIMIT 1
// 	`
// 	var tx payment.FinancialTransaction
// 	var sessionID, customerID, jukirID, successfulPaymentEventID *int64
// 	var paidAt *time.Time
// 	var govShare, companyShare, jukirShare int64
// 	var paymentMethod *string
// 	err := r.db.QueryRowContext(ctx, query, code).Scan(
// 		&tx.ID, &tx.Code, &tx.OperationType, &tx.TransactionSource, &sessionID, &tx.LocationID,
// 		&customerID, &jukirID, &tx.SubtotalAmount, &tx.FinalAmount, &tx.CurrencyCode, &tx.Status, &paidAt, &tx.OccurredAt, &tx.CreatedAt, &successfulPaymentEventID,
// 		&govShare, &companyShare, &jukirShare, &paymentMethod,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("financial transaction not found for code: %s", code)
// 		}
// 		return nil, fmt.Errorf("get financial transaction by code: %w", err)
// 	}
// 	tx.SessionID = sessionID
// 	tx.CustomerID = customerID
// 	tx.JukirID = jukirID
// 	tx.PaidAt = paidAt
// 	tx.SuccessfulPaymentEventID = successfulPaymentEventID
// 	tx.GovShare = govShare
// 	tx.CompanyShare = companyShare
// 	tx.JukirShare = jukirShare
// 	if paymentMethod != nil {
// 		tx.PaymentMethod = *paymentMethod
// 	}
// 	return &tx, nil
// }

// func (r *PaymentRepository) UpdateFinancialTransactionStatus(ctx context.Context, transactionID int64, status string, paidAt *time.Time) error {
// 	query := `UPDATE financial_parking_transaction SET transaction_status = ?, paid_at = ? WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, status, paidAt, transactionID)
// 	if err != nil {
// 		return fmt.Errorf("update financial transaction status: %w", err)
// 	}
// 	// Tambahkan logika untuk membuat history jika diperlukan
// 	return nil
// }

// func (r *PaymentRepository) CreatePaymentEvent(ctx context.Context, event *payment.PaymentEvent) error {
// 	query := `
// 		INSERT INTO financial_payment_event (
// 			payment_event_code, payment_context_type, reference_entity_type, reference_entity_id,
// 			gross_amount, net_amount, currency_code, payment_status, created_at, expired_at,
// 			payment_channel_name, provider_reference, channel_code
// 		)
// 		VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW(), ?, ?, ?, ?)
// 	`
// 	result, err := r.db.ExecContext(ctx, query,
// 		event.Code, event.ContextType, event.ReferenceEntityType, event.ReferenceEntityID,
// 		event.GrossAmount, event.NetAmount, event.CurrencyCode, event.Status, event.ExpiredAt,
// 		event.PaymentChannelName, event.ProviderReference, event.ChannelCode,
// 	)
// 	if err != nil {
// 		return fmt.Errorf("insert payment event: %w", err)
// 	}
// 	id, err := result.LastInsertId()
// 	if err != nil {
// 		return fmt.Errorf("get last insert id: %w", err)
// 	}
// 	event.ID = id
// 	return nil
// }

// func (r *PaymentRepository) GetPaymentEventByCode(ctx context.Context, code string) (*payment.PaymentEvent, error) {
// 	query := `
// 		SELECT id, payment_event_code, payment_context_type, reference_entity_type, reference_entity_id,
// 		       gross_amount, net_amount, currency_code, payment_status, created_at, expired_at, settled_at, failed_at,
// 		       payment_channel_name, provider_reference, channel_code
// 		FROM financial_payment_event
// 		WHERE payment_event_code = ?
// 		LIMIT 1
// 	`
// 	var event payment.PaymentEvent
// 	var expiredAt, settledAt, failedAt *time.Time
// 	err := r.db.QueryRowContext(ctx, query, code).Scan(
// 		&event.ID, &event.Code, &event.ContextType, &event.ReferenceEntityType, &event.ReferenceEntityID,
// 		&event.GrossAmount, &event.NetAmount, &event.CurrencyCode, &event.Status, &event.CreatedAt, &expiredAt, &settledAt, &failedAt,
// 		&event.PaymentChannelName, &event.ProviderReference, &event.ChannelCode,
// 	)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("payment event not found for code: %s", code)
// 		}
// 		return nil, fmt.Errorf("get payment event by code: %w", err)
// 	}
// 	event.ExpiredAt = expiredAt
// 	event.SettledAt = settledAt
// 	event.FailedAt = failedAt
// 	return &event, nil
// }

// func (r *PaymentRepository) LinkPaymentEventToTransaction(ctx context.Context, paymentEventID int64, transactionID int64) error {
// 	query := `UPDATE financial_parking_transaction SET successful_payment_event_id = ? WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, paymentEventID, transactionID)
// 	if err != nil {
// 		return fmt.Errorf("link payment event to transaction: %w", err)
// 	}
// 	// Juga update payment_event untuk mereferensikan transaction (jika belum)
// 	return nil
// }

// func (r *PaymentRepository) UpdatePaymentEventStatus(ctx context.Context, eventID int64, status string, settledAt, failedAt *time.Time) error {
// 	query := `UPDATE financial_payment_event SET payment_status = ?, settled_at = ?, failed_at = ? WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, status, settledAt, failedAt, eventID)
// 	if err != nil {
// 		return fmt.Errorf("update payment event status: %w", err)
// 	}
// 	return nil
// }

// func (r *PaymentRepository) EndParkingSession(ctx context.Context, sessionID int64, endedAt time.Time) error {
// 	query := `UPDATE parking_session SET session_status = 'ended', ended_at = ?, updated_at = NOW() WHERE id = ?`
// 	_, err := r.db.ExecContext(ctx, query, endedAt, sessionID)
// 	if err != nil {
// 		return fmt.Errorf("end parking session: %w", err)
// 	}
// 	return nil
// }
