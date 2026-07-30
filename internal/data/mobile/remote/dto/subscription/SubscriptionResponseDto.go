package subscription

type SubscribeResponseDto struct {
	ActivePaket *ActivePaketDto  `json:"active_paket,omitempty"`
	Benefits    []BenefitsDto    `json:"benefits,omitempty"`
	Statistik   *StatistikDto    `json:"statistik,omitempty"`
	ListPaket   []DetailPaketDto `json:"list_paket,omitempty"`
	Faq         []FaqDto         `json:"faq,omitempty"`
}

type ActivePaketDto struct {
	ActivePackageName    *string `json:"active_package_name"`
	ActivePackageExpired *string `json:"active_package_expired"`
}

type BenefitsDto struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DetailPaketDto struct {
	Name        string   `json:"name"`
	Price       int64    `json:"price"`
	PriceLabel  string   `json:"price_label"`
	PeriodLabel string   `json:"period_label"`
	InfoLabel   string   `json:"info_label"`
	BadgeLabel  *string  `json:"badge_label"`
	Benefits    []string `json:"benefits"`
}

type StatistikDto struct {
	TotalJamParkirBulanLalu       int    `json:"total_jam_parkir_bulan_lalu"`
	TotalBiayaParkirBulanLaluText string `json:"total_biaya_parkir_bulan_lalu_text"`
	TotalPersentaseHematText      string `json:"total_persentase_hemat_text"`
}

type FaqDto struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}
