package subscription

type PromoTerpilihModel struct {
	NamaPromo    *string `json:"nama_promo,omitempty"`
	Deskripsi    *string `json:"deskripsi,omitempty"`
	JumlahDiskon *int64  `json:"jumlah_diskon,omitempty"`
}
