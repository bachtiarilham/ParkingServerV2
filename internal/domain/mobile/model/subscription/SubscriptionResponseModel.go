package subscription

import "time"

type SubscriptionResponseModel struct {
	ActivePackageName    string        `json:"package_name"`
	ActivePackageExpired time.Time     `json:"package_expired"`
	ActivePackageBenefit []string      `json:"benefit_package"`
	ListPaket            ListPaket     `json:"list_paket"`
	PromoTersedia        PromoTersedia `json:"promo_tersedia"`
}

type ListPaket struct {
	Bulanan   []DetailPaket `json:"bulanan"`
	EnamBulan []DetailPaket `json:"enam_bulan"`
	Tahunan   []DetailPaket `json:"tahunan"`
}

type DetailPaket struct {
	NamaPaket      string   `json:"nama_paket"`
	Harga          int64    `json:"harga"`
	CoverageLokasi []string `json:"coverage_lokasi"`
	BenefitPackage []string `json:"benefit_package"`
}

type PromoTersedia struct {
	SyaratDanKetentuan []string      `json:"syarat_dan_ketentuan"`
	EachPromo          []DetailPromo `json:"each_promo"`
}

type DetailPromo struct {
	NamaPromo   string  `json:"nama_promo"`
	BesarDiskon float64 `json:"besar_diskon"`
}
