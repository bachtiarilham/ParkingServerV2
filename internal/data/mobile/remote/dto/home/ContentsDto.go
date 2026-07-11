package home

import "time"

type ContentsDto struct {
	ContentId       int64     `json:"content_id"`
	ContentTypeId   int64     `json:"content_type_id"`
	ContentTypeCode string    `json:"content_type_code"`
	ContentTypeName string    `json:"content_type_name"`
	Title           string    `json:"title"`
	Summary         string    `json:"summary"`
	Body            string    `json:"body"`
	ThumbnailUrl    string    `json:"thumbnail_url"`
	BannerUrl       string    `json:"banner_url"`
	EventLocation   string    `json:"event_location"`
	EventStartAt    time.Time `json:"event_startat"`
	EventEndAt      time.Time `json:"event_endat"`
	PublishAt       time.Time `json:"publishat"`
	ExpiredAt       time.Time `json:"expiredat"`
	Priority        int64     `json:"priority"`
}
