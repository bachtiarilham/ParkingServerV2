package payment

type PembayaranModel struct {
	Title               *string                     `json:"title,omitempty"`
	StatusCard          *PembayaranStatusCardModel  `json:"status_card,omitempty"`
	TotalPembayaran     *int64                      `json:"total_pembayaran,omitempty"`
	DetailLabel         *string                     `json:"detail_label,omitempty"`
	QrisSection         *PembayaranQrisSectionModel `json:"qris_section,omitempty"`
	PaymentOptionsTitle *string                     `json:"payment_options_title,omitempty"`
	PaymentOptions      []PembayaranOptionModel     `json:"payment_options,omitempty"`
	PrintButtonLabel    *string                     `json:"print_button_label,omitempty"`
}
