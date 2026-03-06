package services

import (
	"context"
	"point-exchange-api/internal/db"
	"point-exchange-api/models"
)

type PartnerService struct {
	Repo db.PartnerRepository
}

func (s *PartnerService) RegisterPartner(ctx context.Context, partner *models.Partner) (string, error) {
	return s.Repo.CreatePartner(ctx, partner)
}

func (s *PartnerService) ListPartners(ctx context.Context) ([]*models.Partner, error) {
	return s.Repo.ListPartners(ctx)
}

func (s *PartnerService) GetPartner(ctx context.Context, id string) (*models.Partner, error) {
	return s.Repo.GetPartnerByID(ctx, id)
}

func (s *PartnerService) ActivatePartner(ctx context.Context, id string, isActive bool) error {
	return s.Repo.UpdatePartnerActive(ctx, id, isActive)
}
