package services

import (
	"context"
	"point-exchange-api/internal/db"
	"point-exchange-api/models"
)

type SwapService struct {
	Repo db.SwapRepository
}

func (s *SwapService) ListSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	return s.Repo.ListSwapsBySourcePartnerID(ctx, partnerID)
}

func (s *SwapService) ListSwapsByTargetPartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	return s.Repo.ListSwapsByTargetPartnerID(ctx, partnerID)
}

func (s *SwapService) ListSwapsWithFilter(ctx context.Context, status, sourcePartnerID, targetPartnerID, from, to string) ([]*models.SwapLedger, error) {
	return s.Repo.ListSwapsWithFilter(ctx, status, sourcePartnerID, targetPartnerID, from, to)
}

func (s *SwapService) CreateSwap(ctx context.Context, req *models.SwapRequest) (string, error) {
	// Map request to ledger, set initial status and timestamps
	ledger := &models.SwapLedger{
		SourcePartnerID:   req.SourcePartnerID,
		SourceExternalRef: req.SourceExternalID,
		SourceCustomerID:  req.SourceCustomerID,
		SourcePoints:      req.SourcePoints,
		TargetPartnerID:   req.TargetPartnerID,
		TargetCustomerID:  req.TargetCustomerID,
		Status:            "PENDING",
		// USDValue, ExchangeRateAtTime, CommissionUSD, TargetPoints can be set by business logic if needed
	}
	return s.Repo.CreateSwap(ctx, ledger)
}

func (s *SwapService) GetSwap(ctx context.Context, id string) (*models.SwapLedger, error) {
	return s.Repo.GetSwapByID(ctx, id)
}

func (s *SwapService) ClaimSwaps(ctx context.Context, sourcePartnerID string) ([]*models.SwapLedger, error) {
	return s.Repo.GetPendingSwapsBySourcePartnerID(ctx, sourcePartnerID)
}

// ConfirmSwap updates the swap status to COMPLETED and sets completed_at
func (s *SwapService) ConfirmSwap(ctx context.Context, id string) error {
	return s.Repo.ConfirmSwap(ctx, id)
}
