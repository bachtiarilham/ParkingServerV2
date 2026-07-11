package home

import (
	model "modulegue/internal/domain/mobile/model/profile"
)

type JukirHomeModel struct {
	Profile          *model.JukirModel `json:"profile"`
	Contents         *[]ContentsModel  `json:"contents"`
	UnreadNotifCount int64             `json:"warnings"`
}
