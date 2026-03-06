package handlers

import (
	"context"
	"net/http"
	"point-exchange-api/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// PartnerServiceInterface abstracts the service for testability
type PartnerServiceInterface interface {
	RegisterPartner(ctx context.Context, partner *models.Partner) (string, error)
	ListPartners(ctx context.Context) ([]*models.Partner, error)
	GetPartner(ctx context.Context, id string) (*models.Partner, error)
	ActivatePartner(ctx context.Context, id string, isActive bool) error
}

// PartnerService is injected at startup
var PartnerService PartnerServiceInterface

// SetPartnerService allows test code to inject a mock service
func SetPartnerService(s PartnerServiceInterface) {
	PartnerService = s
}

type RegisterPartnerRequest struct {
	Name      string  `json:"name" binding:"required"`
	APIKey    string  `json:"api_key"`
	APISecret string  `json:"api_secret"`
	IsActive  *bool   `json:"is_active"`
	Rate      float64 `json:"rate"`
}

// RegisterPartner godoc
// @Summary      Register a new partner
// @Description  Register a new partner (admin only)
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        partner body RegisterPartnerRequest true "Partner info"
// @Success      201 {object} models.Partner
// @Failure      400 {object} map[string]string
// @Router       /v1/partners [post]
func RegisterPartner(c *gin.Context) {
	var req RegisterPartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate UUID for ID
	id := uuid.New().String()

	// Use provided or generate API key/secret
	apiKey := req.APIKey
	if apiKey == "" {
		apiKey = RandString(32)
	}
	apiSecret := req.APISecret
	if apiSecret == "" {
		apiSecret = RandString(64)
	}
	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	partner := models.Partner{
		ID:        id,
		Name:      req.Name,
		APIKey:    apiKey,
		APISecret: apiSecret,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		Rate:      req.Rate,
	}
	_, err := PartnerService.RegisterPartner(c.Request.Context(), &partner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, partner)
}

// RandString generates a random string of n characters (for demo only)
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
