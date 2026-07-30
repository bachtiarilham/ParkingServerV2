package riwayat

import "time"

// Response utama API Riwayat Top Up & Transfer
type RiwayatTransaksiDto struct {
	Summary WalletSummaryDto          `json:"summary"`
	Items   []RiwayatTransaksiItemDto `json:"items"`
}

// WalletSummaryDto merepresentasikan statistik arus kas saldo di Top Metric Card
type WalletSummaryDto struct {
	TotalIncome    int64 `json:"total_income"`    // e.g. 250000 (+Rp250.000)
	TotalExpense   int64 `json:"total_expense"`   // e.g. 85000 (-Rp85.000)
	CurrentBalance int64 `json:"current_balance"` // e.g. 165000
}

// RiwayatTransaksiItemDto merepresentasikan transaksi saldo (Kredit / Debit)
type RiwayatTransaksiItemDto struct {
	ID              string    `json:"id"`
	ReferenceNo     string    `json:"reference_no"`     // e.g. "TRX-882190"
	Title           string    `json:"title"`            // e.g. "Top Up via BCA VA" atau "Transfer ke Budi"
	TransactionType string    `json:"transaction_type"` // "TOP_UP" | "TRANSFER" | "PAYMENT"
	Flow            string    `json:"flow"`             // "IN" (Kredit/+) | "OUT" (Debit/-)
	Amount          int64     `json:"amount"`
	Status          string    `json:"status"` // "BERHASIL" | "PENDING" | "GAGAL"
	CreatedAt       time.Time `json:"created_at"`
}
