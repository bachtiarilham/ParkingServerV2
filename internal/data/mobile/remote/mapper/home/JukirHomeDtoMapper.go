package home

import (
	"modulegue/core/utils"
	dto "modulegue/internal/data/mobile/remote/dto/home"
	profilemapper "modulegue/internal/data/mobile/remote/mapper/profile"
	model "modulegue/internal/domain/mobile/model/home"
)

func ToJukirHomeResponse(result *model.JukirHomeModel) *dto.JukirHomeResponseDto {
	if result == nil {
		return &dto.JukirHomeResponseDto{
			Events: []dto.EventDto{},
			News:   []dto.NewsDto{},
		}
	}

	resp := &dto.JukirHomeResponseDto{
		Events: []dto.EventDto{},
		News:   []dto.NewsDto{},
	}

	if result.Profile != nil {
		resp.Profile = profilemapper.ToJukirDto(result.Profile)
	}

	if result.Events != nil {
		for _, ev := range result.Events {
			resp.Events = append(resp.Events, dto.EventDto{
				Id:          ev.ID,
				Title:       ev.Title,
				Description: ev.Description,
				Date:        utils.FormatIndonesianDate(ev.Date),
				ImageUrl:    ev.ImageURL,
				Tag:         ev.ContentType,
			})
		}
	}

	if result.News != nil {
		for _, ev := range result.News {
			resp.News = append(resp.News, dto.NewsDto{
				Id:          ev.ID,
				Title:       ev.Title,
				Description: ev.Description,
				Date:        utils.FormatIndonesianDate(ev.Date),
				ImageUrl:    ev.ImageURL,
				Tag:         ev.ContentType,
			})
		}
	}

	if result.Warnings != nil {
		resp.Warnings = &dto.WarningsDto{
			Profile: result.Warnings.Profile,
			Parking: result.Warnings.Parking,
			Finance: result.Warnings.Finance,
		}
	}

	return resp
}
