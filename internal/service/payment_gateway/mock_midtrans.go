// internal/service/payment_gateway/mock_midtrans.go
package payment_gateway

import (
	"context"
	"fmt"
	"time"
)

// MockMidtransService adalah implementasi dummy dari interface Service
type MockMidtransService struct {
	// Bisa menambahkan field konfigurasi jika diperlukan nanti
}

// NewMockMidtransService membuat instance dari MockMidtransService
func NewMockMidtransService() *MockMidtransService {
	return &MockMidtransService{}
}

// RequestPayment mensimulasikan pembuatan pembayaran di Midtrans
func (s *MockMidtransService) RequestPayment(ctx context.Context, details PaymentDetails) (QRString string, VAString string, BankName string, ExpiredAt string, err error) {
	// Simulasi delay jaringan (opsional)
	// time.Sleep(100 * time.Millisecond)

	// Simulasi pembuatan QRIS
	qrString := fmt.Sprintf("SIMULATED_QRIS_STRING_FOR_%s", details.OrderID)
	// VA dan Bank Name kosong karena tipe pembayaran adalah QRIS
	vaString := ""
	bankName := ""
	// Simulasi expiry dalam 1 jam
	expiredAt := time.Now().Add(1 * time.Hour).Format("2006-01-02T15:04:05Z07:00")

	return qrString, vaString, bankName, expiredAt, nil
}

// VerifyPayment mensimulasikan pengecekan status pembayaran di Midtrans
func (s *MockMidtransService) VerifyPayment(ctx context.Context, orderID string) (Status string, PaidAt *string, err error) {
	// Simulasi delay jaringan (opsional)
	// time.Sleep(100 * time.Millisecond)

	// Simulasi status pembayaran sukses
	status := "settlement" // Atau "capture" untuk QRIS
	// Simulasi waktu pembayaran
	paidAtStr := time.Now().Format("2006-01-02T15:04:05Z07:00")
	paidAt := &paidAtStr

	return status, paidAt, nil
}
