package models

import "time"

type SwapRequest struct {
	SourcePartnerID  string  `json:"source_partner_id" binding:"required"`
	SourceExternalID string  `json:"source_external_id" binding:"required"`
	SourceCustomerID string  `json:"source_customer_id" binding:"required"`
	SourcePoints     float64 `json:"source_points" binding:"required"`
	TargetPartnerID  string  `json:"target_partner_id" binding:"required"`
	TargetCustomerID string  `json:"target_customer_id" binding:"required"`
}

type SwapLedger struct {
	ID                 string     `json:"id"`
	SourcePartnerID    string     `json:"source_partner_id"`
	SourceExternalRef  string     `json:"source_external_ref"`
	SourceCustomerID   string     `json:"source_customer_id"`
	SourcePoints       float64    `json:"source_points"`
	USDValue           float64    `json:"usd_value"`
	ExchangeRateAtTime *float64   `json:"exchange_rate_at_time"`
	CommissionUSD      *float64   `json:"commission_usd"`
	TargetPartnerID    string     `json:"target_partner_id"`
	TargetCustomerID   string     `json:"target_customer_id"`
	TargetPoints       float64    `json:"target_points"`
	Status             string     `json:"status"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ClaimedAt          *time.Time `json:"claimed_at"`
	CompletedAt        *time.Time `json:"completed_at"`
}

// SwapLedgerWithPartnerNames is used for API responses with partner names
// It embeds SwapLedger and adds source/target partner names for display
// This is not stored in the DB, only for API output
//
type SwapLedgerWithPartnerNames struct {
	SwapLedger
	SourcePartnerName string `json:"source_partner_name"`
	TargetPartnerName string `json:"target_partner_name"`
}

type ConfirmSwapRequest struct {
	ID string `json:"id"`
}
