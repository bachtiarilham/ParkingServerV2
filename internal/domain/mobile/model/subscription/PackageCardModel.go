package subscription

type PackageCardModel struct {
	NamaPaket    *string  `json:"nama_paket,omitempty"`
	Harga        *int64   `json:"harga,omitempty"`
	MasaBerlaku  *string  `json:"masa_berlaku,omitempty"`
	JumlahDiskon *int64   `json:"jumlah_diskon,omitempty"`
	Deskripsi    *string  `json:"deskripsi,omitempty"`
	Benefit      []string `json:"benefit,omitempty"`
}
