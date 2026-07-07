package home

import "time"

type HomeModel struct {
	Profile  *ProfileModel  `json:"profile"`
	Events   []EventsModel  `json:"events"`
	News     []NewsModel    `json:"news"`
	Warnings *WarningsModel `json:"warnings"`
}

type ProfileModel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhotoUrl    string `json:"photo_url"`
	Saldo       int64  `json:"saldo"` // Asumsi dalam satuan terkecil (misalnya Rupiah)
	ExpiredDate string `json:"expiredDate"`
	Pendapatan  int64  `json:"pendapatan"`
	Lokasi      string `json:"lokasi"`
	Area        string `json:"area"`
	Zona        string `json:"zona"`
}

type EventsModel struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	ImageURL    string    `json:"image_url"`
	ContentType string    `json:"content_type"` // "event" atau "news"
}

type NewsModel struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	ImageURL    string    `json:"image_url"`
	ContentType string    `json:"content_type"` // "event" atau "news"
}

type WarningsModel struct {
	Profile string `json:"profile"`
	Parking string `json:"parking"`
	Finance string `json:"finance"`
}
