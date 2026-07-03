package payment

import (
	"context"
	"time"
)

type Repository interface {
	// GetActiveSessionByCode mencari sesi parkir aktif berdasarkan code (dari QR)
	GetActiveSessionByCode(ctx context.Context, code string) (*ParkingSession, error)

	// GetTariffForLocationAndVehicle mengambil tariff dari setting lokasi
	GetTariffForLocationAndVehicle(ctx context.Context, locationID int64, vehicleTypeID int64) (int64, error)

	// CreateFinancialTransaction membuat entitas transaksi keuangan
	CreateFinancialTransaction(ctx context.Context, tx *FinancialTransaction) error

	// GetFinancialTransactionByCode mendapatkan transaksi berdasarkan code
	GetFinancialTransactionByCode(ctx context.Context, code string) (*FinancialTransaction, error)

	// UpdateFinancialTransactionStatus memperbarui status transaksi (digunakan oleh callback)
	UpdateFinancialTransactionStatus(ctx context.Context, transactionID int64, status string, paidAt *time.Time) error

	// CreatePaymentEvent membuat event pembayaran
	CreatePaymentEvent(ctx context.Context, event *PaymentEvent) error

	// GetPaymentEventByCode mendapatkan event pembayaran berdasarkan code
	GetPaymentEventByCode(ctx context.Context, code string) (*PaymentEvent, error)

	// LinkPaymentEventToTransaction menghubungkan event pembayaran ke transaksi
	LinkPaymentEventToTransaction(ctx context.Context, paymentEventID int64, transactionID int64) error

	// EndParkingSession mengakhiri sesi parkir
	EndParkingSession(ctx context.Context, sessionID int64, endedAt time.Time) error

	UpdatePaymentEventStatus(ctx context.Context, eventID int64, status string, settledAt, failedAt *time.Time) error
}
