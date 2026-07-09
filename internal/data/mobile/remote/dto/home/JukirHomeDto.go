package home

import dto "modulegue/internal/data/mobile/remote/dto/profile"

type JukirHomeResponseDto struct {
	Profile  *dto.JukirDto `json:"profile"`
	Events   []EventDto    `json:"events"`
	News     []NewsDto     `json:"news"`
	Warnings *WarningsDto  `json:"warnings"`
}
