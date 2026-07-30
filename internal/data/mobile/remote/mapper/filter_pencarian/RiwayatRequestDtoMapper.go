package filterpencarian

import (
	"modulegue/core/utils"
	dto "modulegue/internal/data/mobile/remote/dto/filter_pencarian"
	model "modulegue/internal/domain/mobile/model/filter_pencarian"
)

func ToRiwayatRequestModel(src dto.FilterPencarianDto) *model.FilterPencarianModel {
	startDate, err := utils.ParseISODate(src.StartDate)
	if err != nil {
		return nil
	}

	endDate, err := utils.ParseISODate(src.EndDate)
	if err != nil {
		return nil
	}

	return &model.FilterPencarianModel{
		SearchTypeCode: src.SearchTypeCode,
		StartDate:      startDate,
		EndDate:        endDate,
	}
}
