package home

type NewsDto struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Date        string `json:"date"` // String untuk format ISO8601
	ImageUrl    string `json:"imageurl"`
	Tag         string `json:"tag"`
}
