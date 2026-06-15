package dto

type HomeResponse struct {
	Profile  *ProfileDto  `json:"profile"`
	Summary  *SummaryDto  `json:"summary"`
	Events   []EventDto   `json:"events"`
	News     []NewsDto    `json:"news"`
	Warnings *WarningsDto `json:"warnings"`
}

type ProfileDto struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type SummaryDto struct {
	Saldo      int64 `json:"saldo"`
	Keuntungan int64 `json:"keuntungan"`
}

type EventDto struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}

type NewsDto struct {
	Id          string `json:"id"`
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
