package utils

import (
	"fmt"

	model "modulegue/internal/domain/mobile/model/invoice"

	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/line"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/props"
)

func GenerateInvoicePDF(invoice model.UniversalInvoiceResponseModel) ([]byte, error) {
	cfg := config.NewBuilder().Build()
	m := maroto.New(cfg)

	// --- 1. HEADER SECTION ---
	m.AddRows(
		row.New(20).Add(
			col.New(8).Add(
				text.New("LINESPOT", props.Text{
					Size:  16,
					Style: fontstyle.Bold,
				}),
				text.New(invoice.Title, props.Text{
					Top:  6,
					Size: 10,
				}),
			),
			col.New(4).Add(
				text.New("INVOICE", props.Text{
					Size:  18,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
				text.New(fmt.Sprintf("No: %s", invoice.TrxCode), props.Text{
					Top:   7,
					Size:  9,
					Align: align.Right,
				}),
			),
		),
		line.NewRow(2, props.Line{}),
	)

	// --- 2. TRANSACTION & CUSTOMER METADATA ---
	formattedTime := invoice.CreatedAt.Format("2006-01-02 15:04:05")
	m.AddRows(
		row.New(20).Add(
			col.New(6).Add(
				text.New("DITAGIHKAN KEPADA:", props.Text{Size: 8, Style: fontstyle.Bold}),
				text.New(invoice.CustomerInfo.FullName, props.Text{Top: 4, Size: 9, Style: fontstyle.Bold}),
				text.New(invoice.CustomerInfo.Email, props.Text{Top: 8, Size: 8}),
				text.New(invoice.CustomerInfo.Phone, props.Text{Top: 12, Size: 8}),
			),
			col.New(6).Add(
				text.New(fmt.Sprintf("Tanggal: %s", formattedTime), props.Text{Size: 8, Align: align.Right}),
				text.New(fmt.Sprintf("Metode Bayar: %s", invoice.PaymentMethod.ChannelName), props.Text{Top: 4, Size: 8, Align: align.Right}),
				text.New(fmt.Sprintf("Status: %s", invoice.StatusText), props.Text{Top: 8, Size: 9, Style: fontstyle.Bold, Align: align.Right}),
			),
		),
		line.NewRow(2, props.Line{}),
	)

	// --- 3. DYNAMIC DETAILS (PARKING / TOPUP / MEMBERSHIP / TRANSFER) ---
	switch invoice.TransactionType {
	case "PARKING":
		if invoice.ParkingDetails != nil {
			checkInStr := invoice.ParkingDetails.CheckInTime.Format("2006-01-02 15:04:05")
			checkOutStr := "-"
			if invoice.ParkingDetails.CheckOutTime != nil {
				checkOutStr = invoice.ParkingDetails.CheckOutTime.Format("2006-01-02 15:04:05")
			}
			m.AddRows(
				row.New(35).Add(
					col.New(12).Add(
						text.New("INFORMASI PARKIR", props.Text{Size: 9, Style: fontstyle.Bold}),
						text.New(fmt.Sprintf("Lokasi         : %s", invoice.ParkingDetails.LocationName), props.Text{Top: 4, Size: 8}),
						text.New(fmt.Sprintf("Nomor Plat     : %s (%s)", invoice.ParkingDetails.LicensePlate, invoice.ParkingDetails.VehicleType), props.Text{Top: 9, Size: 8}),
						text.New(fmt.Sprintf("Waktu Masuk    : %s", checkInStr), props.Text{Top: 14, Size: 8}),
						text.New(fmt.Sprintf("Waktu Keluar   : %s", checkOutStr), props.Text{Top: 19, Size: 8}),
						text.New(fmt.Sprintf("Durasi Parkir  : %s", invoice.ParkingDetails.DurationText), props.Text{Top: 24, Size: 8}),
						text.New(fmt.Sprintf("Nama Jukir     : %s", invoice.ParkingDetails.AttendantName), props.Text{Top: 29, Size: 8}),
					),
				),
				line.NewRow(2, props.Line{}),
			)
		}

	case "TOPUP", "TRANSFER":
		if invoice.WalletDetails != nil {
			m.AddRows(
				row.New(25).Add(
					col.New(12).Add(
						text.New("INFORMASI TRANSFER / MUTASI SALDO", props.Text{Size: 9, Style: fontstyle.Bold}),
						text.New(fmt.Sprintf("Pengirim       : %s (%s)", invoice.WalletDetails.SenderName, invoice.WalletDetails.SenderAccount), props.Text{Top: 4, Size: 8}),
						text.New(fmt.Sprintf("Penerima       : %s (%s)", invoice.WalletDetails.RecipientName, invoice.WalletDetails.RecipientAccount), props.Text{Top: 9, Size: 8}),
						text.New(fmt.Sprintf("No Referensi   : %s", invoice.WalletDetails.BankRefNo), props.Text{Top: 14, Size: 8}),
					),
				),
				line.NewRow(2, props.Line{}),
			)
		}

	case "MEMBERSHIP":
		if invoice.MembershipDetails != nil {
			startStr := invoice.MembershipDetails.PeriodStart.Format("2006-01-02")
			endStr := invoice.MembershipDetails.PeriodEnd.Format("2006-01-02")
			autoRenewStr := "Tidak"
			if invoice.MembershipDetails.IsAutoRenew {
				autoRenewStr = "Ya"
			}
			m.AddRows(
				row.New(25).Add(
					col.New(12).Add(
						text.New("INFORMASI PAKET BERLANGGANAN", props.Text{Size: 9, Style: fontstyle.Bold}),
						text.New(fmt.Sprintf("Paket Layanan  : %s", invoice.MembershipDetails.PackageName), props.Text{Top: 4, Size: 8}),
						text.New(fmt.Sprintf("Masa Aktif     : %s s.d %s", startStr, endStr), props.Text{Top: 9, Size: 8}),
						text.New(fmt.Sprintf("Maks Kendaraan : %d Plat Nomor", invoice.MembershipDetails.MaxVehicles), props.Text{Top: 14, Size: 8}),
						text.New(fmt.Sprintf("Perpanjangan   : %s", autoRenewStr), props.Text{Top: 19, Size: 8}),
					),
				),
				line.NewRow(2, props.Line{}),
			)
		}
	}

	// --- 4. PRICE BREAKDOWN TABLE ---
	m.AddRows(
		row.New(10).Add(
			col.New(8).Add(text.New("Rincian Biaya", props.Text{Style: fontstyle.Bold, Size: 10})),
			col.New(4).Add(text.New("Jumlah", props.Text{Style: fontstyle.Bold, Size: 10, Align: align.Right})),
		),
		line.NewRow(2, props.Line{}),
	)

	// Base Price Row
	m.AddRows(
		row.New(8).Add(
			col.New(8).Add(text.New("Tarif Dasar / Subtotal", props.Text{Size: 9})),
			col.New(4).Add(text.New(formatRupiahString(invoice.PriceBreakdown.BasePrice), props.Text{Size: 9, Align: align.Right})),
		),
	)

	// Admin Fee Row (If > 0)
	if invoice.PriceBreakdown.AdminFee > 0 {
		m.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New("Biaya Admin / Layanan", props.Text{Size: 9})),
				col.New(4).Add(text.New(formatRupiahString(invoice.PriceBreakdown.AdminFee), props.Text{Size: 9, Align: align.Right})),
			),
		)
	}

	// Tax Amount Row (If > 0)
	if invoice.PriceBreakdown.TaxAmount > 0 {
		m.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New("PPN (11%)", props.Text{Size: 9})),
				col.New(4).Add(text.New(formatRupiahString(invoice.PriceBreakdown.TaxAmount), props.Text{Size: 9, Align: align.Right})),
			),
		)
	}

	// Discount Row (If > 0)
	if invoice.PriceBreakdown.DiscountAmount > 0 {
		promoStr := "Potongan Promo"
		if invoice.PriceBreakdown.DiscountCode != "" {
			promoStr = fmt.Sprintf("Promo (%s)", invoice.PriceBreakdown.DiscountCode)
		}
		m.AddRows(
			row.New(8).Add(
				col.New(8).Add(text.New(promoStr, props.Text{Size: 9})),
				col.New(4).Add(text.New("-"+formatRupiahString(invoice.PriceBreakdown.DiscountAmount), props.Text{Size: 9, Align: align.Right})),
			),
		)
	}

	m.AddRows(line.NewRow(2, props.Line{}))

	// --- 5. TOTAL & QR CODE FOOTER ---
	m.AddRows(
		row.New(35).Add(
			// QR Code Verification
			col.New(4).Add(
				code.NewQr(fmt.Sprintf("https://linespot.com/verify/%s", invoice.TrxCode)),
			),
			// Grand Total
			col.New(8).Add(
				text.New("TOTAL AKHIR", props.Text{
					Top:   5,
					Size:  11,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
				text.New(formatRupiahString(invoice.PriceBreakdown.FinalTotal), props.Text{
					Top:   13,
					Size:  16,
					Style: fontstyle.Bold,
					Align: align.Right,
				}),
			),
		),
	)

	// Generate document bytes
	document, err := m.Generate()
	if err != nil {
		return nil, fmt.Errorf("generate pdf: %w", err)
	}

	return document.GetBytes(), nil
}

func formatRupiahString(amount int64) string {
	str := fmt.Sprintf("%d", amount)
	var result []rune
	count := 0
	for i := len(str) - 1; i >= 0; i-- {
		if count > 0 && count%3 == 0 {
			result = append([]rune{'.'}, result...)
		}
		result = append([]rune{rune(str[i])}, result...)
		count++
	}
	return "Rp " + string(result)
}
