package db

import (
	"context"
	"database/sql"
	"point-exchange-api/models"
	"strconv"
)

type SwapRepository interface {
	CreateSwap(ctx context.Context, swap *models.SwapLedger) (string, error)
	GetSwapByID(ctx context.Context, id string) (*models.SwapLedger, error)
	GetPendingSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error)
	ListSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error)
	ListSwapsByTargetPartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error)
	ConfirmSwap(ctx context.Context, id string) error
	ListSwapsWithFilter(ctx context.Context, status, sourcePartnerID, targetPartnerID, from, to string) ([]*models.SwapLedger, error)
}

type swapRepo struct {
	db *sql.DB
}

func (r *swapRepo) CreateSwap(ctx context.Context, swap *models.SwapLedger) (string, error) {
	query := `INSERT INTO swap_ledger (
		source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW()
	) RETURNING id`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		swap.SourcePartnerID,
		swap.SourceExternalRef,
		swap.SourceCustomerID,
		swap.SourcePoints,
		swap.USDValue,
		swap.ExchangeRateAtTime,
		swap.CommissionUSD,
		swap.TargetPartnerID,
		swap.TargetCustomerID,
		swap.TargetPoints,
		swap.Status,
	).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *swapRepo) ListSwapsWithFilter(ctx context.Context, status, sourcePartnerID, targetPartnerID, from, to string) ([]*models.SwapLedgerWithPartnerNames, error) {
	query := `SELECT s.id, s.source_partner_id, s.source_external_ref, s.source_customer_id, s.source_points, s.usd_value, s.exchange_rate_at_time, s.commission_usd, s.target_partner_id, s.target_customer_id, s.target_points, s.status, s.created_at, s.updated_at, s.claimed_at, s.completed_at,
		sp.name as source_partner_name, tp.name as target_partner_name
		FROM swap_ledger s
		LEFT JOIN partners sp ON s.source_partner_id = sp.id
		LEFT JOIN partners tp ON s.target_partner_id = tp.id
		WHERE 1=1`
	args := []interface{}{}
	idx := 1
	if status != "" {
		query += ` AND s.status = $` + strconv.Itoa(idx)
		args = append(args, status)
		idx++
	}
	if sourcePartnerID != "" {
		query += ` AND s.source_partner_id = $` + strconv.Itoa(idx)
		args = append(args, sourcePartnerID)
		idx++
	}
	if targetPartnerID != "" {
		query += ` AND s.target_partner_id = $` + strconv.Itoa(idx)
		args = append(args, targetPartnerID)
		idx++
	}
	if from != "" {
		query += ` AND s.created_at >= $` + strconv.Itoa(idx)
		args = append(args, from)
		idx++
	}
	if to != "" {
		query += ` AND s.created_at <= $` + strconv.Itoa(idx)
		args = append(args, to)
		idx++
	}
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var swaps []*models.SwapLedgerWithPartnerNames
	for rows.Next() {
		swap := &models.SwapLedgerWithPartnerNames{}
		var exchangeRate, commission sql.NullFloat64
		var claimedAt, completedAt sql.NullTime
		err := rows.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &exchangeRate, &commission, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt, &swap.SourcePartnerName, &swap.TargetPartnerName)
		if err != nil {
			return nil, err
		}
		if exchangeRate.Valid {
			swap.ExchangeRateAtTime = exchangeRate.Float64
		} else {
			swap.ExchangeRateAtTime = 0
		}
		if commission.Valid {
			swap.CommissionUSD = commission.Float64
		} else {
			swap.CommissionUSD = 0
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

func (r *swapRepo) ConfirmSwap(ctx context.Context, id string) error {
	query := `UPDATE swap_ledger SET status = 'COMPLETED', completed_at = NOW(), updated_at = NOW() WHERE id = $1 AND status != 'COMPLETED'`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *swapRepo) GetSwapByID(ctx context.Context, id string) (*models.SwapLedger, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at, claimed_at, completed_at FROM swap_ledger WHERE id = $1`, id)
	swap := &models.SwapLedger{}
	var claimedAt, completedAt sql.NullTime
	var exchangeRate, commission sql.NullFloat64
	err := row.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &exchangeRate, &commission, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if exchangeRate.Valid {
		swap.ExchangeRateAtTime = &exchangeRate.Float64
	} else {
		swap.ExchangeRateAtTime = nil
	}
	if commission.Valid {
		swap.CommissionUSD = &commission.Float64
	} else {
		swap.CommissionUSD = nil
	}
	if claimedAt.Valid {
		t := claimedAt.Time
		swap.ClaimedAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time
		swap.CompletedAt = &t
	}
	return swap, nil
}

func (r *swapRepo) GetPendingSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at, claimed_at, completed_at FROM swap_ledger WHERE source_partner_id = $1 AND status = 'PENDING'`, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var swaps []*models.SwapLedger
	for rows.Next() {
		swap := &models.SwapLedger{}
		var exchangeRate, commission sql.NullFloat64
		var claimedAt, completedAt sql.NullTime
		err := rows.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &exchangeRate, &commission, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if exchangeRate.Valid {
			swap.ExchangeRateAtTime = &exchangeRate.Float64
		} else {
			swap.ExchangeRateAtTime = nil
		}
		if commission.Valid {
			swap.CommissionUSD = &commission.Float64
		} else {
			swap.CommissionUSD = nil
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

func (r *swapRepo) ListSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at, claimed_at, completed_at FROM swap_ledger WHERE source_partner_id = $1 AND status = 'PENDING'`, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var swaps []*models.SwapLedger
	for rows.Next() {
		swap := &models.SwapLedger{}
		var exchangeRate, commission sql.NullFloat64
		var claimedAt, completedAt sql.NullTime
		err := rows.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &exchangeRate, &commission, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if exchangeRate.Valid {
			swap.ExchangeRateAtTime = &exchangeRate.Float64
		} else {
			swap.ExchangeRateAtTime = nil
		}
		if commission.Valid {
			swap.CommissionUSD = &commission.Float64
		} else {
			swap.CommissionUSD = nil
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

func (r *swapRepo) ListSwapsByTargetPartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, source_partner_id, source_external_ref, source_customer_id, source_points, usd_value, exchange_rate_at_time, commission_usd, target_partner_id, target_customer_id, target_points, status, created_at, updated_at, claimed_at, completed_at FROM swap_ledger WHERE target_partner_id = $1 AND status = 'PENDING'`, partnerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var swaps []*models.SwapLedger
	for rows.Next() {
		swap := &models.SwapLedger{}
		var exchangeRate, commission sql.NullFloat64
		var claimedAt, completedAt sql.NullTime
		err := rows.Scan(&swap.ID, &swap.SourcePartnerID, &swap.SourceExternalRef, &swap.SourceCustomerID, &swap.SourcePoints, &swap.USDValue, &exchangeRate, &commission, &swap.TargetPartnerID, &swap.TargetCustomerID, &swap.TargetPoints, &swap.Status, &swap.CreatedAt, &swap.UpdatedAt, &claimedAt, &completedAt)
		if err != nil {
			return nil, err
		}
		if exchangeRate.Valid {
			swap.ExchangeRateAtTime = &exchangeRate.Float64
		} else {
			swap.ExchangeRateAtTime = nil
		}
		if commission.Valid {
			swap.CommissionUSD = &commission.Float64
		} else {
			swap.CommissionUSD = nil
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

func NewSwapRepository(db *sql.DB) SwapRepository {
	return &swapRepo{db: db}
}
