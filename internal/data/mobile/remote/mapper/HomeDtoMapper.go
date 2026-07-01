package mapper

import (
	"modulegue/internal/data/mobile/remote/dto"
	"modulegue/internal/domain/mobile/model"
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

	if result.Summary != nil {
		resp.Summary = &dto.SummaryDto{
			Saldo:       result.Summary.Saldo,
			ExpiredDate: result.Summary.ExpiredDate,
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

	for _, nw := range result.News {
		resp.News = append(resp.News, dto.NewsDto{
			Id:          nw.ID,
			Title:       nw.Title,
			Description: nw.Description,
			Date:        nw.Date.Format("2006-01-02T15:04:05Z"),
			ImageUrl:    nw.ImageURL,
			Tag:         nw.ContentType,
		})
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
