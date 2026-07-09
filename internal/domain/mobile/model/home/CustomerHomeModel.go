package home

import (
	model "modulegue/internal/domain/mobile/model/profile"
)

type CustomerHomeModel struct {
	Profile  *model.CustomerModel `json:"profile"`
	Events   []EventsModel        `json:"events"`
	News     []NewsModel          `json:"news"`
	Warnings *WarningsModel       `json:"warnings"`
}
