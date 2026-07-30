package subscription

type SubscribeResponseModel struct {
	ActivePaket *ActivePaketModel  `json:"active_paket,omitempty"`
	Benefits    []BenefitsModel    `json:"benefits,omitempty"`
	Statistik   *StatistikModel    `json:"statistik,omitempty"`
	ListPaket   []DetailPaketModel `json:"list_paket,omitempty"`
	Faq         []FaqModel         `json:"faq,omitempty"`
}

type ActivePaketModel struct {
	ActivePackageName    *string `json:"active_package_name"`
	ActivePackageExpired *string `json:"active_package_expired"`
}

type BenefitsModel struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DetailPaketModel struct {
	Name        string   `json:"name"`
	Price       int64    `json:"price"`
	PriceLabel  string   `json:"price_label"`
	PeriodLabel string   `json:"period_label"`
	InfoLabel   string   `json:"info_label"`
	BadgeLabel  *string  `json:"badge_label"`
	Benefits    []string `json:"benefits"`
}

type StatistikModel struct {
	TotalJamParkirBulanLalu       int    `json:"total_jam_parkir_bulan_lalu"`
	TotalBiayaParkirBulanLaluText string `json:"total_biaya_parkir_bulan_lalu_text"`
	TotalPersentaseHematText      string `json:"total_persentase_hemat_text"`
}

type FaqModel struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

func GetDefaultFaqs() []FaqModel {
	return []FaqModel{
		{
			Question: "Bagaimana cara membayar membership?",
			Answer:   "Anda dapat membayar menggunakan saldo LineSpot, QRIS, Virtual Account, atau gerai retail terdekat.",
		},
		{
			Question: "Apakah saldo LineSpot dapat dicairkan?",
			Answer:   "Mohon maaf, saldo yang sudah dimasukkan tidak dapat diuangkan kembali namun tidak memiliki masa kedaluwarsa.",
		},
		{
			Question: "Bagaimana jika membership belum aktif setelah pembayaran?",
			Answer:   "Silakan segarkan halaman aplikasi. Jika status belum berubah dalam 10 menit, hubungi bantuan admin dengan melampirkan bukti transaksi.",
		},
	}
}
