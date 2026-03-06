package handlers

import (
	"context"
	"net/http"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

// RateServiceInterface abstracts the service for testability
type RateServiceInterface interface {
	AddOrUpdateRate(ctx context.Context, partnerID string, req *models.AddOrUpdateRateRequest) error
	ListRates(ctx context.Context, partnerID string) ([]*models.Rate, error)
}

// RateService is injected at startup
var RateService RateServiceInterface

// AddOrUpdateRate godoc
// @Summary      Add or update a partner's point rate
// @Description  Add or update a point rate for a partner (admin only)
// @Tags         rates
// @Accept       json
// @Produce      json
// @Param        id path string true "Partner ID"
// @Param        rate body models.AddOrUpdateRateRequest true "Rate info"
// @Success      200 {object} map[string]string
// @Router       /v1/partners/{id}/rates [post]
func AddOrUpdateRate(c *gin.Context) {
	var req models.AddOrUpdateRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	partnerID := c.Param("id")
	err := RateService.AddOrUpdateRate(c.Request.Context(), partnerID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Rate added/updated"})
}

// ListRates godoc
// @Summary      List all rates for a partner
// @Description  Get all point rates for a partner
// @Tags         rates
// @Produce      json
// @Param        id path string true "Partner ID"
// @Success      200 {array} models.Rate
// @Router       /v1/partners/{id}/rates [get]
func ListRates(c *gin.Context) {
	partnerID := c.Param("id")
	rates, err := RateService.ListRates(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rates)
}
