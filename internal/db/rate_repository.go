package db

import (
	"context"
	"database/sql"
	"point-exchange-api/models"
)

type RateRepository interface {
	AddOrUpdateRate(ctx context.Context, partnerID string, req *models.AddOrUpdateRateRequest) error
	ListRates(ctx context.Context, partnerID string) ([]*models.Rate, error)
}

type rateRepo struct {
	db *sql.DB
}

func NewRateRepository(db *sql.DB) RateRepository {
	return &rateRepo{db: db}
}

func (r *rateRepo) AddOrUpdateRate(ctx context.Context, partnerID string, req *models.AddOrUpdateRateRequest) error {
	// Upsert logic: try update, if not found then insert
	res, err := r.db.ExecContext(ctx, `UPDATE rates SET points_per_usd = $1, min_exchange_points = $2 WHERE partner_id = $3 AND point_type = $4`, req.PointsPerUSD, req.MinExchangePoints, partnerID, req.PointType)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		_, err = r.db.ExecContext(ctx, `INSERT INTO rates (partner_id, point_type, points_per_usd, min_exchange_points, effective_date, created_at) VALUES ($1, $2, $3, $4, NOW(), NOW())`, partnerID, req.PointType, req.PointsPerUSD, req.MinExchangePoints)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *rateRepo) ListRates(ctx context.Context, partnerID string) ([]*models.Rate, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, partner_id, point_type, points_per_usd, min_exchange_points, effective_date, created_at FROM rates WHERE partner_id = $1`, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []*models.Rate
	for rows.Next() {
		rate := &models.Rate{}
		err := rows.Scan(&rate.ID, &rate.PartnerID, &rate.PointType, &rate.PointsPerUSD, &rate.MinExchangePoints, &rate.EffectiveDate, &rate.CreatedAt)
		if err != nil {
			return nil, err
		}
		rates = append(rates, rate)
	}
	return rates, nil
}
