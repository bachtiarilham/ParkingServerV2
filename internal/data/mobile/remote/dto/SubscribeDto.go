package dto

type SubscribeDto struct {
	StatusCard  *StatusCardDto   `json:"status_card,omitempty"`
	PackageCard []PackageCardDto `json:"package_card,omitempty"`
	Promo       []PromoDto       `json:"promo,omitempty"`
}

type StatusCardDto struct {
	PaketAktif *string `json:"paket_aktif,omitempty"`
	Kadaluarsa *string `json:"kadaluarsa,omitempty"`
	Benefit    *string `json:"benefit,omitempty"`
}

type PackageCardDto struct {
	NamaPaket    *string  `json:"nama_paket,omitempty"`
	Harga        *int64   `json:"harga,omitempty"`
	MasaBerlaku  *string  `json:"masa_berlaku,omitempty"`
	JumlahDiskon *int64   `json:"jumlah_diskon,omitempty"`
	Deskripsi    *string  `json:"deskripsi,omitempty"`
	Benefit      []string `json:"benefit,omitempty"`
}

type PromoDto struct {
	SNk   []string           `json:"snk,omitempty"`
	Promo []PromoTerpilihDto `json:"promo,omitempty"`
}

type PromoTerpilihDto struct {
	NamaPromo    *string `json:"nama_promo,omitempty"`
	Deskripsi    *string `json:"deskripsi,omitempty"`
	JumlahDiskon *int64  `json:"jumlah_diskon,omitempty"`
}
