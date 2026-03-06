package db

import (
	"context"
	"database/sql"
	"point-exchange-api/models"
)

type PartnerRepository interface {
	CreatePartner(ctx context.Context, partner *models.Partner) (string, error)
	ListPartners(ctx context.Context) ([]*models.Partner, error)
	GetPartnerByID(ctx context.Context, id string) (*models.Partner, error)
	UpdatePartnerActive(ctx context.Context, id string, isActive bool) error
}

type partnerRepo struct {
	db *sql.DB
}

func NewPartnerRepository(db *sql.DB) PartnerRepository {
	return &partnerRepo{db: db}
}

func (r *partnerRepo) CreatePartner(ctx context.Context, partner *models.Partner) (string, error) {
	query := `INSERT INTO partners (id, name, api_key, api_secret, is_active, created_at, rate) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, partner.ID, partner.Name, partner.APIKey, partner.APISecret, partner.IsActive, partner.CreatedAt, partner.Rate)
	if err != nil {
		return "", err
	}
	return partner.ID, nil
}

func (r *partnerRepo) ListPartners(ctx context.Context) ([]*models.Partner, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, name, api_key, api_secret, is_active, created_at, rate FROM partners`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var partners []*models.Partner
	for rows.Next() {
		p := &models.Partner{}
		err := rows.Scan(&p.ID, &p.Name, &p.APIKey, &p.APISecret, &p.IsActive, &p.CreatedAt, &p.Rate)
		if err != nil {
			return nil, err
		}
		partners = append(partners, p)
	}
	return partners, nil
}

func (r *partnerRepo) GetPartnerByID(ctx context.Context, id string) (*models.Partner, error) {
	row := r.db.QueryRowContext(ctx, `SELECT id, name, api_key, api_secret, is_active, created_at, rate FROM partners WHERE id = $1`, id)
	p := &models.Partner{}
	err := row.Scan(&p.ID, &p.Name, &p.APIKey, &p.APISecret, &p.IsActive, &p.CreatedAt, &p.Rate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *partnerRepo) UpdatePartnerActive(ctx context.Context, id string, isActive bool) error {
	_, err := r.db.ExecContext(ctx, `UPDATE partners SET is_active = $1 WHERE id = $2`, isActive, id)
	return err
}
