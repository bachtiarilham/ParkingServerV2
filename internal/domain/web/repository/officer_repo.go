package repository

import (
	"context"
	"modulegue/internal/delivery/web/dto/request"
	"modulegue/internal/domain/web/officer"
)

type officerRepo interface {
	UpdateOfficerStatus(ctx context.Context, adminID int64, officerID string, status string) (officer.ParkingOfficerOption, error)
	ApplyOfficerMutation(ctx context.Context, adminID int64, req request.ApplyOfficerMutationRequest) (officer.ParkingOfficerOption, error)
}
