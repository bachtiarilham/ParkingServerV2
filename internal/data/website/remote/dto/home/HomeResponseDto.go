package home

type HomeResponseData struct {
	NIK             string                `json:"nik"`
	NoTelp          string                `json:"no_telp"`
	GrafikPerDay    *[]GrafikPerDayDto    `json:"grafikperday"`
	TransaksiAkhir  *[]TransaksiAkhirDto  `json:"transaksiakhir"`
	GrafikPerMinggu *GrafikPerMingguDto   `json:"grafikperminggu"`
	ListProvinsi    *[]ProvinsiDto        `json:"listprovinsi"`
	ListShift       *[]ShiftDto           `json:"listshift"`
	Alamat          string                `json:"alamat"`
	ListKecamatan   *[]KecamatanDto       `json:"listkecamatan"`
	PetugasLapangan *[]PetugasLapanganDto `json:"petugaslapangan"`
	ID              int                   `json:"id"`
	StatsGlobal     *StatsGlobalDto       `json:"statsglobal"`
	Token           string                `json:"token"`
	ListKabupaten   *[]KabupatenDto       `json:"listkabupaten"`
	Nama            string                `json:"nama"`
	Foto            string                `json:"foto"`
	ListZona        *[]ZonaDto            `json:"listzona"`
	ListRole        *[]RoleDto            `json:"listrole"`
	ListDesa        *[]DesaDto            `json:"listdesa"`
}

// --- Sub-structs untuk LoginData ---
type GrafikPerDayDto struct {
	Tanggal  string `json:"tanggal"`
	JmlMotor int    `json:"jml_motor"`
	JmlMobil int    `json:"jml_mobil"`
	Total    int    `json:"total"`
}

type TransaksiAkhirDto struct {
	PlatNomor string `json:"plat_nomor"`
	NamaJukir string `json:"nama_jukir"`
	Waktu     string `json:"waktu"`
	Kendaraan string `json:"kendaraan"`
	Tarif     int    `json:"tarif"`
}

type GrafikPerMingguDto struct {
	Motor int `json:"motor"`
	Mobil int `json:"mobil"`
	Total int `json:"total"`
}

type ProvinsiDto struct {
	NamaProv string `json:"nama_prov"`
	ID       int    `json:"id"`
}

type ShiftDto struct {
	JamKeluar  string `json:"jam_keluar"`
	NamaShift  string `json:"nama_shift"`
	JamMasuk   string `json:"jam_masuk"`
	ID         int    `json:"id"`
	NamaParlok string `json:"nama_parlok"`
}

type KecamatanDto struct {
	NamaKec string `json:"nama_kec"`
	ID      int    `json:"id"`
}

type PetugasLapanganDto struct {
	Parlok string `json:"parlok"`
	Nama   string `json:"nama"`
	Foto   string `json:"foto"`
}

type StatsGlobalDto struct {
	JmlTransaksi int `json:"jml_transaksi"`
	TotalPAD     int `json:"total_pad"`
	JmlPetugas   int `json:"jml_petugas"`
	Saldo        int `json:"saldo"`
}

type KabupatenDto struct {
	NamaKab string `json:"nama_kab"`
	ID      int    `json:"id"`
}

type ZonaDto struct {
	Keterangan string `json:"keterangan"`
	ID         int    `json:"id"`
}

type RoleDto struct {
	NamaRole string `json:"nama_role"`
	ID       int    `json:"id"`
}

type DesaDto struct {
	NamaDesa string `json:"nama_desa"`
	ID       int    `json:"id"`
}
