package payment

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	model "modulegue/internal/domain/mobile/model/payment"
	"modulegue/internal/domain/mobile/repository"
)

type PaymentRepositoryImpl struct {
	db *sql.DB
}

func NewPaymentRepositoryImpl(db *sql.DB) repository.PaymentRepository {
	return &PaymentRepositoryImpl{db: db}
}

func (r *PaymentRepositoryImpl) PostParking(ctx context.Context, req model.PostParkingRequestModel) (*model.PostParkingResponseModel, error) {
	locationID, locationName, err := r.findLocationByName(ctx, req.LokasiParkir)
	if err != nil {
		return nil, err
	}

	vehicleTypeID, _, err := r.findVehicleTypeByName(ctx, req.JenisKendaraan)
	if err != nil {
		return nil, err
	}

	startedAt, err := parseDateTime(req.WaktuMasuk)
	if err != nil {
		return nil, fmt.Errorf("parse waktu_masuk: %w", err)
	}

	sessionCode := fmt.Sprintf("PS-%d", time.Now().UnixNano())
	plateNumber := strconv.FormatInt(req.NomorPolisi, 10)

	result, err := r.db.ExecContext(
		ctx,
		`
		INSERT INTO parking_session (
			session_code,
			location_id,
			vehicle_type_id,
			plate_number,
			session_status,
			started_at,
			created_at,
			updated_at
		) VALUES (?, ?, ?, ?, 'active', ?, NOW(), NOW())
		`,
		sessionCode,
		locationID,
		vehicleTypeID,
		plateNumber,
		startedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert parking_session: %w", err)
	}

	sessionID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("get parking_session id: %w", err)
	}

	title := "Pembayaran Parkir"
	statusTitle := "Menunggu Pembayaran"
	statusMessage := "Silakan lanjutkan pembayaran parkir"
	isSuccess := false
	totalPembayaran := int64(0)
	detailLabel := "Detail Parkir"
	qrisTitle := "QRIS"
	qrisInfo := "QR akan dibuat saat pembayaran dikonfirmasi"
	countdown := int64(0)
	alternativeLabel := "Metode pembayaran lain"
	paymentOptionsTitle := "Pilih Metode Pembayaran"
	printButtonLabel := "Cetak"
	optionQrisType := "QRIS"
	optionQrisTitle := "QRIS"
	optionQrisSubtitle := "Bayar dengan scan QR"
	optionCashType := "CASH"
	optionCashTitle := "Tunai"
	optionCashSubtitle := "Bayar langsung ke petugas"

	return &model.PostParkingResponseModel{
		Title: &title,
		StatusCard: &model.PembayaranStatusCardModel{
			Title:     &statusTitle,
			Message:   &statusMessage,
			IsSuccess: &isSuccess,
		},
		TotalPembayaran: &totalPembayaran,
		DetailLabel:     &detailLabel,
		QrisSection: &model.PembayaranQrisSectionModel{
			Title:       &qrisTitle,
			Information: &qrisInfo,
			Countdown:   &countdown,
			QrContent: &model.IsiQrModel{
				SessionID:     sessionID,
				PlatNomor:     plateNumber,
				Lokasi:        locationName,
				WaktuMasuk:    req.WaktuMasuk,
				Durasi:        "",
				Nominal:       totalPembayaran,
				IsPaid:        false,
				PaymentStatus: 0,
				IsExpired:     false,
				StatusMessage: statusMessage,
			},
			AlternativeLabel: &alternativeLabel,
		},
		PaymentOptionsTitle: &paymentOptionsTitle,
		PaymentOptions: []model.PembayaranOptionModel{
			{Type: &optionQrisType, Title: &optionQrisTitle, Subtitle: &optionQrisSubtitle},
			{Type: &optionCashType, Title: &optionCashTitle, Subtitle: &optionCashSubtitle},
		},
		PrintButtonLabel: &printButtonLabel,
	}, nil
}

func (r *PaymentRepositoryImpl) PostPaymentParking(ctx context.Context, req model.PostPaymentParkingRequestModel) (*model.PostPaymentParkingResponseModel, error) {
	sessionID, err := strconv.ParseInt(req.SessionID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid session_id: %w", err)
	}

	endedAt := time.Now()
	durationMinutes := 0
	if startedAt, err := parseDateTime(req.WaktuMasuk); err == nil {
		durationMinutes = int(endedAt.Sub(startedAt).Minutes())
		if durationMinutes < 0 {
			durationMinutes = 0
		}
	}

	txCode := fmt.Sprintf("TRX-%d", time.Now().UnixNano())
	paymentMethodID, paymentMethodName, err := r.findPaymentMethod(ctx, "QRIS")
	if err != nil {
		return nil, err
	}

	_, err = r.db.ExecContext(
		ctx,
		`
		INSERT INTO financial_parking_transaction (
			transaction_code,
			operation_type,
			transaction_source,
			payment_method,
			session_id,
			location_id,
			vehicle_type_id,
			plate_number,
			payment_method_id,
			subtotal_amount,
			final_amount,
			transaction_status,
			paid_at,
			occurred_at,
			created_at
		)
		SELECT ?, 'PAYMENT', 'MOBILE', ?, ps.id, ps.location_id, ps.vehicle_type_id, ps.plate_number, ?, ?, ?, 'paid', ?, ?, NOW()
		FROM parking_session ps
		WHERE ps.id = ?
		`,
		txCode,
		paymentMethodName,
		paymentMethodID,
		req.Nominal,
		req.Nominal,
		endedAt,
		endedAt,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("insert financial_parking_transaction: %w", err)
	}

	_, err = r.db.ExecContext(
		ctx,
		`
		UPDATE parking_session
		SET session_status = 'ended',
			ended_at = ?,
			duration_minutes = ?,
			updated_at = NOW()
		WHERE id = ?
		`,
		endedAt,
		durationMinutes,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("update parking_session: %w", err)
	}

	title := "Pembayaran Berhasil"
	successTitle := "Pembayaran Berhasil"
	successDescription := "Pembayaran parkir berhasil diproses"
	totalAmount := strconv.FormatInt(req.Nominal, 10)
	paymentStatus := "paid"
	referenceNumber := txCode
	verificationMessage := "Transaksi telah terverifikasi"
	thankYouTitle := "Terima Kasih"
	thankYouDescription := "Terima kasih telah menggunakan layanan parkir"
	downloadLabel := "Unduh Bukti"
	backToHomeLabel := "Kembali ke Home"

	sessionLabel := "Session ID"
	plateLabel := "Plat Nomor"
	locationLabel := "Lokasi"
	methodLabel := "Metode Pembayaran"

	return &model.PostPaymentParkingResponseModel{
		Title:               &title,
		SuccessTitle:        &successTitle,
		SuccessDescription:  &successDescription,
		TotalAmount:         &totalAmount,
		PaymentStatus:       &paymentStatus,
		ReferenceNumber:     &referenceNumber,
		VerificationMessage: &verificationMessage,
		ThankYouTitle:       &thankYouTitle,
		ThankYouDescription: &thankYouDescription,
		DownloadLabel:       &downloadLabel,
		BackToHomeLabel:     &backToHomeLabel,
		Details: []model.PostPaymentParkingDetailItemModel{
			{Label: &sessionLabel, Value: &req.SessionID},
			{Label: &plateLabel, Value: &req.PlatNomor},
			{Label: &locationLabel, Value: &req.Lokasi},
			{Label: &methodLabel, Value: &paymentMethodName},
		},
	}, nil
}

func (r *PaymentRepositoryImpl) GetPembayaranStatus(ctx context.Context, sessionId string) (*model.PostPaymentParkingResponseModel, error) {
	var (
		txCode        sql.NullString
		finalAmount   sql.NullInt64
		status        sql.NullString
		plateNumber   sql.NullString
		locationName  sql.NullString
		paymentMethod sql.NullString
	)

	err := r.db.QueryRowContext(
		ctx,
		`
		SELECT
			fpt.transaction_code,
			fpt.final_amount,
			fpt.transaction_status,
			COALESCE(fpt.plate_number, ''),
			COALESCE(pl.location_name, ''),
			COALESCE(pm.payment_method_name, COALESCE(fpt.payment_method, ''))
		FROM financial_parking_transaction fpt
		LEFT JOIN parking_session ps ON ps.id = fpt.session_id
		LEFT JOIN parking_location pl ON pl.id = fpt.location_id
		LEFT JOIN payment_method pm ON pm.id = fpt.payment_method_id
		WHERE fpt.session_id = ?
		ORDER BY fpt.id DESC
		LIMIT 1
		`,
		sessionId,
	).Scan(
		&txCode,
		&finalAmount,
		&status,
		&plateNumber,
		&locationName,
		&paymentMethod,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("get pembayaran status: %w", err)
	}

	title := "Status Pembayaran"
	successTitle := "Pembayaran Berhasil"
	successDescription := "Pembayaran parkir sudah diterima"
	totalAmount := strconv.FormatInt(finalAmount.Int64, 10)
	paymentStatus := status.String
	referenceNumber := txCode.String
	verificationMessage := "Status pembayaran terkonfirmasi"
	thankYouTitle := "Terima Kasih"
	thankYouDescription := "Transaksi parkir selesai"
	downloadLabel := "Unduh Bukti"
	backToHomeLabel := "Kembali ke Home"
	plateLabel := "Plat Nomor"
	locationLabel := "Lokasi"
	methodLabel := "Metode Pembayaran"

	return &model.PostPaymentParkingResponseModel{
		Title:               &title,
		SuccessTitle:        &successTitle,
		SuccessDescription:  &successDescription,
		TotalAmount:         &totalAmount,
		PaymentStatus:       &paymentStatus,
		ReferenceNumber:     &referenceNumber,
		VerificationMessage: &verificationMessage,
		ThankYouTitle:       &thankYouTitle,
		ThankYouDescription: &thankYouDescription,
		DownloadLabel:       &downloadLabel,
		BackToHomeLabel:     &backToHomeLabel,
		Details: []model.PostPaymentParkingDetailItemModel{
			{Label: &plateLabel, Value: stringPtr(plateNumber.String)},
			{Label: &locationLabel, Value: stringPtr(locationName.String)},
			{Label: &methodLabel, Value: stringPtr(paymentMethod.String)},
		},
	}, nil
}

func (r *PaymentRepositoryImpl) findLocationByName(ctx context.Context, name string) (int64, string, error) {
	var id int64
	var locationName string
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, location_name FROM parking_location WHERE location_name = ? AND is_active = 1 LIMIT 1`,
		strings.TrimSpace(name),
	).Scan(&id, &locationName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("location not found: %s", name)
		}
		return 0, "", fmt.Errorf("find location: %w", err)
	}
	return id, locationName, nil
}

func (r *PaymentRepositoryImpl) findVehicleTypeByName(ctx context.Context, name string) (int64, string, error) {
	normalized := strings.TrimSpace(name)
	query := `
		SELECT id, vehicle_type_name
		FROM vehicle_type
		WHERE LOWER(vehicle_type_name) = LOWER(?) OR LOWER(vehicle_type_code) = LOWER(?)
		LIMIT 1
	`
	var id int64
	var vehicleTypeName string
	err := r.db.QueryRowContext(ctx, query, normalized, normalized).Scan(&id, &vehicleTypeName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("vehicle type not found: %s", name)
		}
		return 0, "", fmt.Errorf("find vehicle type: %w", err)
	}
	return id, vehicleTypeName, nil
}

func (r *PaymentRepositoryImpl) findPaymentMethod(ctx context.Context, name string) (int64, string, error) {
	query := `
		SELECT id, payment_method_name
		FROM payment_method
		WHERE UPPER(payment_method_name) = UPPER(?) OR UPPER(payment_method_code) = UPPER(?)
		LIMIT 1
	`
	var id int64
	var paymentMethodName string
	err := r.db.QueryRowContext(ctx, query, name, name).Scan(&id, &paymentMethodName)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("payment method not found: %s", name)
		}
		return 0, "", fmt.Errorf("find payment method: %w", err)
	}
	return id, paymentMethodName, nil
}

func parseDateTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	layouts := []string{
		"2006-01-02 15:04:05",
		time.RFC3339,
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, value); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unsupported datetime format: %s", value)
}

func stringPtr(v string) *string {
	return &v
}
