package home

import "time"

// type NewsModel struct {
// 	News []NewsItemModel `json:"news"`
// }

type NewsModel struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	ImageURL    string    `json:"image_url"`
	ContentType string    `json:"content_type"` // "event" atau "news"
}
