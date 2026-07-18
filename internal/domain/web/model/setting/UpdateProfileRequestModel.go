package setting

type UpdateProfilRequestModel struct {
	IDJukir int `json:"id_jukir"`
	// Field di bawah menggunakan pointer agar tidak tereksekusi jika kosong (omitempty)
	Nama     *string `json:"nama,omitempty"`
	Username *string `json:"username,omitempty"`
	Alamat   *string `json:"alamat,omitempty"`
	NoTelp   *string `json:"notelp,omitempty"`
	Password *string `json:"password,omitempty"`
	IDStatus *string `json:"id_status,omitempty"`
	Foto     *string `json:"foto,omitempty"`
}
