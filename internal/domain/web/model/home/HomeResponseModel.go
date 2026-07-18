package home

type HomeResponseModel struct {
	NIK             string                  `json:"nik"`
	NoTelp          string                  `json:"no_telp"`
	GrafikPerDay    *[]GrafikPerDayModel    `json:"grafikperday"`
	TransaksiAkhir  *[]TransaksiAkhirModel  `json:"transaksiakhir"`
	GrafikPerMinggu *GrafikPerMingguModel   `json:"grafikperminggu"`
	ListProvinsi    *[]ProvinsiModel        `json:"listprovinsi"`
	ListShift       *[]ShiftModel           `json:"listshift"`
	Alamat          string                  `json:"alamat"`
	ListKecamatan   *[]KecamatanModel       `json:"listkecamatan"`
	PetugasLapangan *[]PetugasLapanganModel `json:"petugaslapangan"`
	ID              int                     `json:"id"`
	StatsGlobal     *StatsGlobalModel       `json:"statsglobal"`
	Token           string                  `json:"token"`
	ListKabupaten   *[]KabupatenModel       `json:"listkabupaten"`
	Nama            string                  `json:"nama"`
	Foto            string                  `json:"foto"`
	ListZona        *[]ZonaModel            `json:"listzona"`
	ListRole        *[]RoleModel            `json:"listrole"`
	ListDesa        *[]DesaModel            `json:"listdesa"`
	PasswordHash    string
	UserId          int64
	RoleId          int64
}

type GrafikPerDayModel struct {
	Tanggal  string `json:"tanggal"`
	JmlMotor int    `json:"jml_motor"`
	JmlMobil int    `json:"jml_mobil"`
	Total    int    `json:"total"`
}

type TransaksiAkhirModel struct {
	PlatNomor string `json:"plat_nomor"`
	NamaJukir string `json:"nama_jukir"`
	Waktu     string `json:"waktu"`
	Kendaraan string `json:"kendaraan"`
	Tarif     int    `json:"tarif"`
}

type GrafikPerMingguModel struct {
	Motor int `json:"motor"`
	Mobil int `json:"mobil"`
	Total int `json:"total"`
}

type ProvinsiModel struct {
	NamaProv string `json:"nama_prov"`
	ID       int    `json:"id"`
}

type ShiftModel struct {
	JamKeluar  string `json:"jam_keluar"`
	NamaShift  string `json:"nama_shift"`
	JamMasuk   string `json:"jam_masuk"`
	ID         int    `json:"id"`
	NamaParlok string `json:"nama_parlok"`
}

type KecamatanModel struct {
	NamaKec string `json:"nama_kec"`
	ID      int    `json:"id"`
}

type PetugasLapanganModel struct {
	Parlok string `json:"parlok"`
	Nama   string `json:"nama"`
	Foto   string `json:"foto"`
}

type StatsGlobalModel struct {
	JmlTransaksi int `json:"jml_transaksi"`
	TotalPAD     int `json:"total_pad"`
	JmlPetugas   int `json:"jml_petugas"`
	Saldo        int `json:"saldo"`
}

type KabupatenModel struct {
	NamaKab string `json:"nama_kab"`
	ID      int    `json:"id"`
}

type ZonaModel struct {
	Keterangan string `json:"keterangan"`
	ID         int    `json:"id"`
}

type RoleModel struct {
	NamaRole string `json:"nama_role"`
	ID       int    `json:"id"`
}

type DesaModel struct {
	NamaDesa string `json:"nama_desa"`
	ID       int    `json:"id"`
}
