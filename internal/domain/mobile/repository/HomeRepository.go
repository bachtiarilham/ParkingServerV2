package repository

import (
	"context"
	model "modulegue/internal/domain/mobile/model/home"
)

type HomeRepository interface {
	GetHome(ctx context.Context, reqModel model.GetHomeReqModel) (*model.HomeModel, error)
	// GetProfile(ctx context.Context) (*model.Profile, error)
	// GetSummary(ctx context.Context) (*model.Summary, error)
	// GetJukirSummary(ctx context.Context) (*model.JukirSummary, error)
	// GetRecentEventsAndNews(ctx context.Context, limit int, offset int) ([]model.EventOrNews, error)
	// GetWarnings(ctx context.Context) (*model.Warnings, error)
}
