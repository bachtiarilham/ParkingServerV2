package payment

type PembayaranQrisSectionModel struct {
	Title            *string     `json:"title,omitempty"`
	QrContent        *IsiQrModel `json:"qr_content,omitempty"`
	Information      *string     `json:"masaBerlakuQr,omitempty"`
	Countdown        *int64      `json:"countdown,omitempty"`
	AlternativeLabel *string     `json:"alternative_label,omitempty"`
}
