package home_customer

import "time"

// EventOrNews mewakili entitas Event atau News
type EventOrNews struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	ImageURL    string    `json:"image_url"`
	ContentType string    `json:"content_type"` // "event" atau "news"
}

// Profile mewakili data profil pengguna
type Profile struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Summary mewakili ringkasan data pengguna
type Summary struct {
	Saldo       int64      `json:"saldo"`        // Asumsi dalam satuan terkecil (misalnya Rupiah)
	ExpiredDate *time.Time `json:"expired_date"` // Tidak digunakan sementara

}

// Warnings mewakili peringatan-peringatan
type Warnings struct {
	Profile string `json:"profile"`
	Parking string `json:"parking"`
	Finance string `json:"finance"`
}
