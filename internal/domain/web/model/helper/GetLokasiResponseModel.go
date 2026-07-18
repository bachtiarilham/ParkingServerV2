package helper

type GetLokasiResponseModel struct {
	LokasiItem *[]LokasiItemModel `json:"lokasi_list"`
}

type LokasiItemModel struct {
	ID          int    `json:"ID"`
	NamaParlok  string `json:"NamaParlok"`
	JalanParlok string `json:"JalanParlok"`
}
