package home

import (
	"context"
	// "errors"
	// "fmt"
	"modulegue/internal/domain/home"
)

type GetDashboardInput struct {
	UserID int64
	Limit  int // Untuk pagination event/news
	Offset int // Untuk pagination event/news
}

type GetDashboardOutput struct {
	Profile  *home.Profile      `json:"profile"`
	Summary  *home.Summary      `json:"summary"`
	Events   []home.EventOrNews `json:"events"`
	News     []home.EventOrNews `json:"news"`
	Warnings *home.Warnings     `json:"warnings"`
}

type GetDashboardUseCase struct {
	homeRepo home.Repository
}

func NewGetDashboardUseCase(homeRepo home.Repository) *GetDashboardUseCase {
	return &GetDashboardUseCase{
		homeRepo: homeRepo,
	}
}

func (uc *GetDashboardUseCase) Execute(ctx context.Context, input GetDashboardInput) (GetDashboardOutput, error) {
	var output GetDashboardOutput
	var err error

	// 1. Ambil Profil
	output.Profile, err = uc.homeRepo.GetProfile(ctx, input.UserID)
	if err != nil {
		// Log error jika perlu
		// Bisa return error atau kosongkan profil
		output.Profile = &home.Profile{} // Atau return error
	}

	// 2. Ambil Summary
	output.Summary, err = uc.homeRepo.GetSummary(ctx, input.UserID)
	if err != nil {
		// Log error
		output.Summary = &home.Summary{Saldo: 0} // Default jika error
	}

	// 3. Ambil Events & News (gabungkan dulu, lalu pisahkan di usecase)
	allContent, err := uc.homeRepo.GetRecentEventsAndNews(ctx, input.Limit, input.Offset)
	if err != nil {
		// Log error
		// Biarkan slice kosong jika error
		allContent = []home.EventOrNews{}
	}

	// Pisahkan Events dan News
	for _, item := range allContent {
		if item.ContentType == "event" {
			output.Events = append(output.Events, item)
		} else if item.ContentType == "news" {
			output.News = append(output.News, item)
		}
		// Jika ada tipe lain, abaikan
	}

	// 4. Ambil Warnings
	output.Warnings, err = uc.homeRepo.GetWarnings(ctx, input.UserID)
	if err != nil {
		// Log error
		// Kembalikan struct kosong jika error
		output.Warnings = &home.Warnings{}
	}

	return output, nil
}
