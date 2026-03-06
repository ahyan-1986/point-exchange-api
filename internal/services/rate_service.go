package services

import (
	"context"
	"point-exchange-api/internal/db"
	"point-exchange-api/models"
)

type RateService struct {
	Repo db.RateRepository
}

func (s *RateService) AddOrUpdateRate(ctx context.Context, partnerID string, req *models.AddOrUpdateRateRequest) error {
	return s.Repo.AddOrUpdateRate(ctx, partnerID, req)
}

func (s *RateService) ListRates(ctx context.Context, partnerID string) ([]*models.Rate, error) {
	return s.Repo.ListRates(ctx, partnerID)
}
