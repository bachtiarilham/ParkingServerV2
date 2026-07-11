package home

import (
	model "modulegue/internal/domain/mobile/model/profile"
)

type CustomerHomeModel struct {
	Profile          *model.CustomerModel `json:"profile"`
	Contents         *[]ContentsModel     `json:"contents"`
	UnreadNotifCount int64                `json:"warnings"`
}
