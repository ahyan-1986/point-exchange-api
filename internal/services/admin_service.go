package services

import (
	"context"
	"point-exchange-api/internal/db"
	"point-exchange-api/models"
)

type AdminService struct {
	Repo db.AdminRepository
}

func (s *AdminService) ListSwapLedger(ctx context.Context) ([]*models.SwapLedger, error) {
	return s.Repo.ListSwapLedger(ctx)
}
