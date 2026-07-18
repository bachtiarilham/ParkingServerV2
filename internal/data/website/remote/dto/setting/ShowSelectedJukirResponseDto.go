package setting

type ShowSelectedJukirResponseDto struct {
	NIK      string `json:"nik"`
	Username string `json:"username"`
	NoTelp   string `json:"notelp"`
	SaldoMin int    `json:"saldo_min"`
	IDStatus int    `json:"idstatus"`
	IDRole   int    `json:"idrole"`
	Alamat   string `json:"alamat"`
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Saldo    int    `json:"saldo"`
	Nama     string `json:"nama"`
	Foto     string `json:"foto"`
}
