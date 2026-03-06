
package models

import "time"

type Rate struct {
	ID                int64     `json:"id"`
	PartnerID         string    `json:"partner_id"`
	PointType         string    `json:"point_type"`
	PointsPerUSD      float64   `json:"points_per_usd"`
	MinExchangePoints float64   `json:"min_exchange_points"`
	EffectiveDate     time.Time `json:"effective_date"`
	CreatedAt         time.Time `json:"created_at"`
}
