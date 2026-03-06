package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListSwapsRequested godoc
// @Summary      List swaps requested by partner
// @Description  List all swaps where this partner is the source (requested swaps)
// @Tags         swaps
// @Produce      json
// @Param        id path string true "Partner ID"
// @Success      200 {array} models.SwapLedger
// @Router       /v1/partners/{id}/swaps/requested [get]
// @Success      200 {array} models.SwapLedgerWithPartnerNames
func ListSwapsRequested(c *gin.Context) {
	partnerID := c.Param("id")
	swaps, err := SwapService.ListSwapsBySourcePartnerID(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, swaps)
}

// ListSwapsReceived godoc
// @Summary      List swaps received by partner
// @Description  List all swaps where this partner is the target (received swaps)
// @Tags         swaps
// @Produce      json
// @Param        id path string true "Partner ID"
// @Success      200 {array} models.SwapLedger
// @Router       /v1/partners/{id}/swaps/received [get]
// @Success      200 {array} models.SwapLedgerWithPartnerNames
func ListSwapsReceived(c *gin.Context) {
	partnerID := c.Param("id")
	swaps, err := SwapService.ListSwapsByTargetPartnerID(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, swaps)
}
