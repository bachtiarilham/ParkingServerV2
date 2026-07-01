package home

// import (
// 	"context"
// 	// "errors"
// 	// "fmt"
// 	"modulegue/internal/domain/mobile/model"
// 	"modulegue/internal/domain/mobile/repository"
// )

// type GetHomeInput struct {
// 	UserID int64
// 	Limit  int // Untuk pagination event/news
// 	Offset int // Untuk pagination event/news
// }

// type GetHomeOutput struct {
// 	Profile      *model.Profile      `json:"profile"`
// 	Summary      *model.Summary      `json:"summary"`
// 	JukirSummary *model.JukirSummary `json:"jukirSummary"`
// 	Events       []model.Events      `json:"events"`
// 	News         []model.News        `json:"news"`
// 	Warnings     *model.Warnings     `json:"warnings"`
// }

// type GetHomeUseCase struct {
// 	homeRepo repository.HomeRepository
// }

// func NewGetHomeUseCase(homeRepo repository.HomeRepository) *GetHomeUseCase {
// 	return &GetHomeUseCase{
// 		homeRepo: homeRepo,
// 	}
// }

// func (uc *GetHomeUseCase) Execute(ctx context.Context, input GetHomeInput) (GetHomeOutput, error) {
// 	var output GetHomeOutput
// 	var err error

// 	// 1. Ambil Profil
// 	output.Profile, err = uc.homeRepo.GetProfile(ctx, input.UserID)
// 	if err != nil {
// 		// Log error jika perlu
// 		// Bisa return error atau kosongkan profil
// 		output.Profile = &model.Profile{} // Atau return error
// 	}

// 	// 2. Ambil Summary
// 	output.Summary, err = uc.homeRepo.GetSummary(ctx, input.UserID)
// 	if err != nil {
// 		// Log error
// 		output.Summary = &model.Summary{Saldo: 0} // Default jika error
// 	}

// 	// 3. Ambil Events & News (gabungkan dulu, lalu pisahkan di usecase)
// 	allContent, err := uc.homeRepo.GetRecentEventsAndNews(ctx, input.Limit, input.Offset)
// 	if err != nil {
// 		// Log error
// 		// Biarkan slice kosong jika error
// 		allContent = []model.News{}
// 	}

// 	// Pisahkan Events dan News
// 	for _, item := range allContent {
// 		if item.ContentType == "event" {
// 			output.Events = append(output.Events, item)
// 		} else if item.ContentType == "news" {
// 			output.News = append(output.News, item)
// 		}
// 		// Jika ada tipe lain, abaikan
// 	}

// 	// 4. Ambil Warnings
// 	output.Warnings, err = uc.homeRepo.GetWarnings(ctx, input.UserID)
// 	if err != nil {
// 		// Log error
// 		// Kembalikan struct kosong jika error
// 		output.Warnings = &model.Warnings{}
// 	}

// 	return output, nil
// }
