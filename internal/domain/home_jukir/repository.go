package home_jukir

import (
	"context"
	// "time"
)

type Repository interface {
	GetProfile(ctx context.Context, userID int64) (*Profile, error)
	GetSummary(ctx context.Context, userID int64) (*Summary, error)
	GetRecentEventsAndNews(ctx context.Context, limit int, offset int) ([]EventOrNews, error)
	GetWarnings(ctx context.Context, userID int64) (*Warnings, error)
}
