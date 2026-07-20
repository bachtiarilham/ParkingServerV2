package repository

import (
	"context"
	model "modulegue/internal/domain/web/model/setting"
)

type SettingRepository interface {
	AddParlok(context.Context, model.AddParlokRequestModel) error
	RegisterJukir(context.Context, model.RegisterRequestModel) error
	SaveSchedule(context.Context, model.SaveScheduleRequestModel) error
	SaveTarif(context.Context, model.SaveTarifRequestModel) error
	ShowSelectedJukir(context.Context, model.ShowSelectedJukirRequestModel) (*model.ShowSelectedJukirResponseModel, error)
	UpdateProfil(context.Context, model.UpdateProfilRequestModel) error
}
