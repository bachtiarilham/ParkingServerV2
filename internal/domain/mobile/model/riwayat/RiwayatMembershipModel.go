package riwayat

import "time"

// Response utama API Riwayat Membership
type RiwayatMembershipModel struct {
	Summary MembershipSummaryModel       `json:"summary"`
	Items   []RiwayatMembershipItemModel `json:"items"`
}

// MembershipSummaryModel merepresentasikan status langganan aktif di Top Metric Card
type MembershipSummaryModel struct {
	PackageName       string     `json:"package_name"` // e.g. "LineSpot Gold Personal"
	IsActive          bool       `json:"is_active"`
	ActiveUntil       *time.Time `json:"active_until,omitempty"`
	NextBillingAmount int64      `json:"next_billing_amount"` // e.g. 59000
	IsAutoRenew       bool       `json:"is_auto_renew"`
}

// RiwayatMembershipItemModel murni merepresentasikan bukti pembayaran/tagihan real
type RiwayatMembershipItemModel struct {
	ID          string    `json:"id"`
	InvoiceNo   string    `json:"invoice_no"`   // e.g. "INV/20260722/LSP-0091"
	PackageName string    `json:"package_name"` // e.g. "LineSpot Gold Personal - 1 Bulan"
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`
	Amount      int64     `json:"amount"`
	PaidAt      time.Time `json:"paid_at"`
	Status      string    `json:"status"` // "DIBAYAR" | "BATAL" | "GAGAL"
}
