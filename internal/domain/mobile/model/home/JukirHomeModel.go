package home

import (
	model "modulegue/internal/domain/mobile/model/profile"
)

type JukirHomeModel struct {
	Profile  *model.JukirModel `json:"profile"`
	Events   []EventsModel     `json:"events"`
	News     []NewsModel       `json:"news"`
	Warnings *WarningsModel    `json:"warnings"`
}
