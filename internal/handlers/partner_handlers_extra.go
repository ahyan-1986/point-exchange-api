package handlers

import (
	"net/http"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

// ListPartners godoc
// @Summary      List all partners
// @Description  Get a list of all partners (admin only)
// @Tags         partners
// @Produce      json
// @Success      200 {array} models.Partner
// @Router       /v1/partners [get]
func ListPartners(c *gin.Context) {
	partners, err := PartnerService.ListPartners(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, partners)
}

// GetPartner godoc
// @Summary      Get partner details
// @Description  Get details for a specific partner
// @Tags         partners
// @Produce      json
// @Param        id path string true "Partner ID"
// @Success      200 {object} models.Partner
// @Failure      404 {object} map[string]string
// @Router       /v1/partners/{id} [get]
func GetPartner(c *gin.Context) {
	id := c.Param("id")
	partner, err := PartnerService.GetPartner(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if partner == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}
	c.JSON(http.StatusOK, partner)
}

// ActivatePartner godoc
// @Summary      Activate or deactivate a partner
// @Description  Set a partner's active status (admin only)
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        id path string true "Partner ID"
// @Param        body body models.ActivatePartnerRequest true "Active status"
// @Success      200 {object} models.Partner
// @Failure      404 {object} map[string]string
// @Router       /v1/partners/{id}/activate [patch]
func ActivatePartner(c *gin.Context) {
	id := c.Param("id")
	var req models.ActivatePartnerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := PartnerService.ActivatePartner(c.Request.Context(), id, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	partner, err := PartnerService.GetPartner(c.Request.Context(), id)
	if err != nil || partner == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Partner not found after update"})
		return
	}
	c.JSON(http.StatusOK, partner)
}
