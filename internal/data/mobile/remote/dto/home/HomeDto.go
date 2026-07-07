package home

type HomeResponseDto struct {
	Profile  *ProfileDto  `json:"profile"`
	Events   []EventDto   `json:"events"`
	News     []NewsDto    `json:"news"`
	Warnings *WarningsDto `json:"warnings"`
}

type ProfileDto struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Saldo       int64  `json:"saldo"`
	ExpiredDate string `json:"expiredDate"`
	Pendapatan  int64  `json:"pendapatan"`
	Lokasi      string `json:"lokasi"`
	Area        string `json:"area"`
	Zona        string `json:"zona"`
}

type EventDto struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}

type NewsDto struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}

type WarningsDto struct {
	Profile string `json:"profile"`
	Parking string `json:"parking"`
	Finance string `json:"finance"`
}
