package home

import dto "modulegue/internal/data/mobile/remote/dto/profile"

type CustomerHomeResponseDto struct {
	Profile *dto.CustomerDto `json:"profile"`
	Events  []ContentsDto    `json:"contents"`
}
