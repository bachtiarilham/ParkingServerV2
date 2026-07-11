package home

import (
	dto "modulegue/internal/data/mobile/remote/dto/home"
	profilemapper "modulegue/internal/data/mobile/remote/mapper/profile"
	model "modulegue/internal/domain/mobile/model/home"
)

func ToCustomerHomeResponse(result *model.CustomerHomeModel) *dto.CustomerHomeResponseDto {
	if result == nil {
		return &dto.CustomerHomeResponseDto{
			Events: []dto.ContentsDto{},
		}
	}

	resp := &dto.CustomerHomeResponseDto{
		Events: []dto.ContentsDto{},
	}

	if result.Profile != nil {
		resp.Profile = profilemapper.ToCustomerDto(result.Profile)
	}

	if result.Contents != nil {
		for _, item := range *result.Contents {
			resp.Events = append(resp.Events, dto.ContentsDto{
				ContentId:       item.ContentId,
				ContentTypeId:   item.ContentTypeId,
				ContentTypeCode: item.ContentTypeCode,
				ContentTypeName: item.ContentTypeName,
				Title:           item.Title,
				Summary:         item.Summary,
				Body:            item.Body,
				ThumbnailUrl:    item.ThumbnailUrl,
				BannerUrl:       item.BannerUrl,
				EventLocation:   item.EventLocation,
				EventStartAt:    item.EventStartAt,
				EventEndAt:      item.EventEndAt,
				PublishAt:       item.PublishAt,
				ExpiredAt:       item.ExpiredAt,
				Priority:        item.Priority,
			})
		}
	}

	resp.UnreadNotifCount = result.UnreadNotifCount

	return resp
}
