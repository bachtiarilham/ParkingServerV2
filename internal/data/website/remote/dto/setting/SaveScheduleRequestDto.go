package setting

type SaveScheduleRequestDto struct {
	ID        int    `json:"ID"`
	IDUser    int    `json:"IDUser"`
	IDLokasi  int    `json:"IDLokasi"`
	IDZona    int    `json:"IDZona"`
	IDArea    int    `json:"IDArea"`
	IDShift   int    `json:"IDShift"`
	DateAwal  string `json:"DateAwal"`
	DateAkhir string `json:"DateAkhir"`
}

// Gunakan AsyncTaskResponse yang sudah dideklarasikan di endpoint Register
