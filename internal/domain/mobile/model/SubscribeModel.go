package model

type SubscribeModel struct {
	StatusCard  *StatusCardModel   `json:"status_card,omitempty"`
	PackageCard []PackageCardModel `json:"package_card,omitempty"`
	Promo       []PromoModel       `json:"promo,omitempty"`
}

type StatusCardModel struct {
	PaketAktif *string `json:"paket_aktif,omitempty"`
	Kadaluarsa *string `json:"kadaluarsa,omitempty"`
	Benefit    *string `json:"benefit,omitempty"`
}

type PackageCardModel struct {
	NamaPaket    *string  `json:"nama_paket,omitempty"`
	Harga        *int64   `json:"harga,omitempty"`
	MasaBerlaku  *string  `json:"masa_berlaku,omitempty"`
	JumlahDiskon *int64   `json:"jumlah_diskon,omitempty"`
	Deskripsi    *string  `json:"deskripsi,omitempty"`
	Benefit      []string `json:"benefit,omitempty"`
}

type PromoModel struct {
	SNk   []string             `json:"snk,omitempty"`
	Promo []PromoTerpilihModel `json:"promo,omitempty"`
}

type PromoTerpilihModel struct {
	NamaPromo    *string `json:"nama_promo,omitempty"`
	Deskripsi    *string `json:"deskripsi,omitempty"`
	JumlahDiskon *int64  `json:"jumlah_diskon,omitempty"`
}
