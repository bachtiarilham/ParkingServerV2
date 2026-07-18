package setting

type RegisterRequestModel struct {
	NIK      string `json:"nik"`
	Nama     string `json:"nama"`
	NoTelp   string `json:"notelp"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	IDRole   int    `json:"idrole"`
	Alamat   string `json:"alamat"`
	Foto     string `json:"foto"`
}
