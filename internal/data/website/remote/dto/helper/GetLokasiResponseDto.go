package helper

type GetLokasiResponseDto struct {
	LokasiItem *[]LokasiItemDto `json:"lokasi_list"`
}

type LokasiItemDto struct {
	ID          int    `json:"ID"`
	NamaParlok  string `json:"NamaParlok"`
	JalanParlok string `json:"JalanParlok"`
}
