//PaymentInfo, PaymentType

package payment

import (
	"time"
)

// ParkingSession mewakili sesi parkir yang sedang aktif
type ParkingSession struct {
	ID              int64     `json:"id"`
	Code            string    `json:"code"`
	CustomerID      int64     `json:"customer_id"`
	LocationID      int64     `json:"location_id"`
	StartedAt       time.Time `json:"started_at"`
	DurationMinutes *int      `json:"duration_minutes,omitempty"` // Bisa null jika masih aktif
	Status          string    `json:"status"`                     // active, ended, voided
	PlateNumber     string    `json:"plate_number"`
	VehicleTypeID   int64     `json:"vehicle_type_id"`
}

// FinancialTransaction mewakili transaksi keuangan (parkir, topup, dll)
type FinancialTransaction struct {
	ID                       int64      `json:"id"`
	Code                     string     `json:"code"`
	OperationType            string     `json:"operation_type"`       // parking_fee, topup, refund
	TransactionSource        string     `json:"transaction_source"`   // qr_scan, manual, etc
	SessionID                *int64     `json:"session_id,omitempty"` // Reference ke parking_session
	LocationID               int64      `json:"location_id"`
	CustomerID               *int64     `json:"customer_id,omitempty"`
	SubtotalAmount           int64      `json:"subtotal_amount"`
	FinalAmount              int64      `json:"final_amount"` // Sudah termasuk diskon/penalty
	CurrencyCode             string     `json:"currency_code"`
	Status                   string     `json:"status"` // unpaid, paid, voided, disputed
	PaidAt                   *time.Time `json:"paid_at,omitempty"`
	OccurredAt               time.Time  `json:"occurred_at"`
	CreatedAt                time.Time  `json:"created_at"`
	SuccessfulPaymentEventID *int64     `json:"successful_payment_event_id,omitempty"` // Reference ke payment_event
}

// PaymentEvent mewakili event pembayaran dari gateway
type PaymentEvent struct {
	ID                  int64      `json:"id"`
	Code                string     `json:"code"`
	ContextType         string     `json:"context_type"`          // transaction, topup
	ReferenceEntityType string     `json:"reference_entity_type"` // financial_parking_transaction
	ReferenceEntityID   int64      `json:"reference_entity_id"`
	GrossAmount         int64      `json:"gross_amount"`
	NetAmount           int64      `json:"net_amount"`
	CurrencyCode        string     `json:"currency_code"`
	Status              string     `json:"status"` // pending, settled, failed
	CreatedAt           time.Time  `json:"created_at"`
	ExpiredAt           *time.Time `json:"expired_at,omitempty"`
	SettledAt           *time.Time `json:"settled_at,omitempty"`
	FailedAt            *time.Time `json:"failed_at,omitempty"`
	PaymentChannelName  string     `json:"payment_channel_name"`         // qris, virtual_account
	ProviderReference   string     `json:"provider_reference,omitempty"` // VA number, QRIS string
	ChannelCode         string     `json:"channel_code,omitempty"`       // bca_va, gopay_qris
}
