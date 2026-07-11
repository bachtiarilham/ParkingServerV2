package home

import dto "modulegue/internal/data/mobile/remote/dto/profile"

type JukirHomeResponseDto struct {
	Profile          *dto.JukirDto `json:"profile"`
	Events           []ContentsDto `json:"contents"`
	UnreadNotifCount int64         `json:"warnings"`
}
