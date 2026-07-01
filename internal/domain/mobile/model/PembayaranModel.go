package model

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

type PembayaranStatusCardModel struct {
	Title     *string `json:"title,omitempty"`
	Message   *string `json:"message,omitempty"`
	IsSuccess *bool   `json:"is_success,omitempty"`
}

type PembayaranQrisSectionModel struct {
	Title            *string     `json:"title,omitempty"`
	QrContent        *IsiQrModel `json:"qr_content,omitempty"`
	Information      *string     `json:"masaBerlakuQr,omitempty"`
	Countdown        *int64      `json:"countdown,omitempty"`
	AlternativeLabel *string     `json:"alternative_label,omitempty"`
}

type IsiQrModel struct {
	SessionID     int64  `json:"session_id"`
	PlatNomor     string `json:"plat_nomor"`
	Lokasi        string `json:"lokasi"`
	WaktuMasuk    string `json:"waktu_masuk"`
	Durasi        string `json:"durasi"`
	Nominal       int64  `json:"nominal"`
	IsPaid        bool   `json:"isPaid"`
	PaymentStatus int64  `json:"paymentStatus"`
	IsExpired     bool   `json:"isExpired"`
	StatusMessage string `json:"statusMessage"`
}

type PembayaranOptionModel struct {
	Type     *string `json:"type,omitempty"`
	Title    *string `json:"title,omitempty"`
	Subtitle *string `json:"subtitle,omitempty"`
}
