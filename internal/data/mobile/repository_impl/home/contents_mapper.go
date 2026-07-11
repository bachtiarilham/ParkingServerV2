package home

import (
	"database/sql"

	"modulegue/core/utils"
	model "modulegue/internal/domain/mobile/model/home"
)

type contentsRow struct {
	ContentId       sql.NullInt64
	ContentTypeId   sql.NullInt64
	ContentTypeCode sql.NullString
	ContentTypeName sql.NullString
	Title           sql.NullString
	Summary         sql.NullString
	Body            sql.NullString
	ThumbnailUrl    sql.NullString
	BannerUrl       sql.NullString
	EventLocation   sql.NullString
	EventStartAt    sql.NullTime
	EventEndAt      sql.NullTime
	PublishAt       sql.NullTime
	ExpiredAt       sql.NullTime
	Priority        sql.NullInt64
}

func mapContentsRowToModel(row contentsRow) model.ContentsModel {
	item := model.ContentsModel{
		ContentId:       utils.NullInt64Value(row.ContentId),
		ContentTypeId:   utils.NullInt64Value(row.ContentTypeId),
		ContentTypeCode: utils.NullStringValue(row.ContentTypeCode),
		ContentTypeName: utils.NullStringValue(row.ContentTypeName),
		Title:           utils.NullStringValue(row.Title),
		Summary:         utils.NullStringValue(row.Summary),
		Body:            utils.NullStringValue(row.Body),
		ThumbnailUrl:    utils.NullStringValue(row.ThumbnailUrl),
		BannerUrl:       utils.NullStringValue(row.BannerUrl),
		EventLocation:   utils.NullStringValue(row.EventLocation),
		EventStartAt:    utils.NullTimeValue(row.EventStartAt),
		EventEndAt:      utils.NullTimeValue(row.EventEndAt),
		PublishAt:       utils.NullTimeValue(row.PublishAt),
		ExpiredAt:       utils.NullTimeValue(row.ExpiredAt),
		Priority:        utils.NullInt64Value(row.Priority),
	}

	return item
}
