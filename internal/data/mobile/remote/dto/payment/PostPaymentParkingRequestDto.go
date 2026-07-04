package payment

type PostPaymentParkingRequestDto struct {
	SessionID     string `json:"session_id"`
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
