package home

import "time"

type ContentsModel struct {
	ContentId       int64
	ContentTypeId   int64
	ContentTypeCode string
	ContentTypeName string
	Title           string
	Summary         string
	Body            string
	ThumbnailUrl    string
	BannerUrl       string
	EventLocation   string
	EventStartAt    time.Time
	EventEndAt      time.Time
	PublishAt       time.Time
	ExpiredAt       time.Time
	Priority        int64
}
