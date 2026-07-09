package home

import dto "modulegue/internal/data/mobile/remote/dto/profile"

type CustomerHomeResponseDto struct {
	Profile  *dto.CustomerDto `json:"profile"`
	Events   []EventDto       `json:"events"`
	News     []NewsDto        `json:"news"`
	Warnings *WarningsDto     `json:"warnings"`
}
