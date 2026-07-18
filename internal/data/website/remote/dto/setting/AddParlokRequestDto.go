package setting

type AddParlokRequestDto struct {
	NamaParlok   string  `json:"NamaParlok"`
	JalanParlok  string  `json:"JalanParlok"`
	IDZona       int     `json:"IDZona"`
	IDArea       int     `json:"IDArea"`
	IDDes        int     `json:"IDDes"`
	IDKec        int     `json:"IDKec"`
	IDKab        int     `json:"IDKab"`
	IDProv       int     `json:"IDProv"`
	LatMinArea   float64 `json:"LatMinArea"`
	LatMaxArea   float64 `json:"LatMaxArea"`
	LngMinArea   float64 `json:"LngMinArea"`
	LngMaxArea   float64 `json:"LngMaxArea"`
	AltitudeArea float64 `json:"AltitudeArea"`
	CenterAreaX  float64 `json:"CenterAreaX"`
	CenterAreaY  float64 `json:"CenterAreaY"`
}

// Gunakan AsyncTaskResponse yang sudah dideklarasikan di endpoint Register
