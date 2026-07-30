package settings

import (
	"context"
	"database/sql"

	model "modulegue/internal/domain/mobile/model/settings"
	"modulegue/internal/domain/mobile/repository"
)

type SettingsRepositoryImpl struct {
	db *sql.DB
}

func NewSettingsRepositoryImpl(db *sql.DB) repository.SettingsRepository {
	return &SettingsRepositoryImpl{db: db}
}

func (r *SettingsRepositoryImpl) ChangeProfile(ctx context.Context, req model.SettingsModel) error {
	query := `
		UPDATE user_identity 
		SET 
			email = COALESCE(?, email), 
			phone_number = COALESCE(?, phone_number), 
			photo_url = COALESCE(?, photo_url), 
			updated_at = NOW() 
		WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, req.Email, req.NoTelp, req.FotoProfil, req.UserId)
	return err
}
