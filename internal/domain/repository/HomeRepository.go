package repository

import (
	"context"
	"modulegue/internal/domain/model"
)

type Repository interface {
	GetProfile(ctx context.Context, userID int64) (*model.Profile, error)
	GetSummary(ctx context.Context, userID int64) (*model.Summary, error)
	GetJukirSummary(ctx context.Context, userID int64) (*model.JukirSummary, error)
	GetRecentEventsAndNews(ctx context.Context, limit int, offset int) ([]model.EventOrNews, error)
	GetWarnings(ctx context.Context, userID int64) (*model.Warnings, error)
}
