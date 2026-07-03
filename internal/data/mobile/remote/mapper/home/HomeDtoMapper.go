package home

import (
	dto "modulegue/internal/data/mobile/remote/dto/home"
	model "modulegue/internal/domain/mobile/model/home"
)

func ToHomeResponse(result *model.HomeModel) *dto.HomeResponseDto {
	if result == nil {
		return &dto.HomeResponseDto{
			Events: []dto.EventDto{},
			News:   []dto.NewsDto{},
		}
	}

	resp := &dto.HomeResponseDto{
		Events: []dto.EventDto{},
		News:   []dto.NewsDto{},
	}

	if result.Profile != nil {
		resp.Profile = &dto.ProfileDto{
			Id:   result.Profile.ID,
			Name: result.Profile.Name,
		}
	}

	if result.CustomerSummary != nil {
		resp.CustomerSummary = &dto.CustomerSummaryDto{
			Saldo:       result.CustomerSummary.Saldo,
			ExpiredDate: result.CustomerSummary.ExpiredDate,
		}
	}

	if result.JukirSummary != nil {
		resp.JukirSummary = &dto.JukirSummaryDto{
			Pendapatan: result.JukirSummary.Pendapatan,
			Lokasi:     result.JukirSummary.Lokasi,
			Area:       result.JukirSummary.Area,
			Zona:       result.JukirSummary.Zona,
		}
	}

	if result.Events != nil {
		for _, ev := range result.Events {
			resp.Events = append(resp.Events, dto.EventDto{
				Id:          ev.ID,
				Title:       ev.Title,
				Description: ev.Description,
				Date:        ev.Date.Format("2006-01-02T15:04:05Z"),
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
				Date:        ev.Date.Format("2006-01-02T15:04:05Z"),
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
