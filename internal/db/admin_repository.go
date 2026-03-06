package db

import (
	"context"
	"database/sql"
	"point-exchange-api/models"
)

type AdminRepository interface {
	ListSwapLedger(ctx context.Context) ([]*models.SwapLedger, error)
}

type adminRepo struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) AdminRepository {
	return &adminRepo{db: db}
}

func (r *adminRepo) ListSwapLedger(ctx context.Context) ([]*models.SwapLedger, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at, claimed_at, completed_at FROM swap_ledger`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var swaps []*models.SwapLedger
	for rows.Next() {
		swap := &models.SwapLedger{}
		var claimedAt, completedAt sql.NullTime
		err := rows.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &swap.ExchangeRateAtTime, &swap.CommissionUSD, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if claimedAt.Valid {
			t := claimedAt.Time
			swap.ClaimedAt = &t
		}
		if completedAt.Valid {
			t := completedAt.Time
			swap.CompletedAt = &t
		}
		swaps = append(swaps, swap)
	}
	return swaps, nil
}
