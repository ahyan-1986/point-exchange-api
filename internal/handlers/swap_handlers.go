package handlers

import (
	"net/http"
	"point-exchange-api/models"

	"context"

	"github.com/gin-gonic/gin"
)

// SwapServiceInterface abstracts the service for testability
type SwapServiceInterface interface {
	CreateSwap(ctx context.Context, req *models.SwapRequest) (string, error)
	GetSwap(ctx context.Context, id string) (*models.SwapLedger, error)
	ClaimSwaps(ctx context.Context, sourcePartnerID string) ([]*models.SwapLedger, error)
	ListSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error)
	ListSwapsByTargetPartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error)
	ConfirmSwap(ctx context.Context, id string) error
	ListSwapsWithFilter(ctx context.Context, status, sourcePartnerID, targetPartnerID, from, to string) ([]*models.SwapLedgerWithPartnerNames, error)
}

// SwapService is injected at startup
var SwapService SwapServiceInterface

// CreateDeposit godoc
// @Summary      Deposit (initiate a point swap)
// @Description  Partner A pushes a point deduction request
// @Tags         swaps
// @Accept       json
// @Produce      json
// @Param        deposit body models.SwapRequest true "Swap request"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Router       /v1/swap/deposit [post]
func CreateDeposit(c *gin.Context) {
	var req models.SwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := SwapService.CreateSwap(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transaction_id": id, "status": "PENDING"})
}

// GetSwap godoc
// @Summary      Get swap transaction details
// @Description  Get details for a specific swap
// @Tags         swaps
// @Produce      json
// @Param        id path string true "Swap ID"
// @Success      200 {object} models.SwapLedger
// @Failure      404 {object} map[string]string
// @Router       /v1/swap/{id} [get]
func GetSwap(c *gin.Context) {
	id := c.Param("id")
	swap, err := SwapService.GetSwap(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if swap == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swap not found"})
		return
	}
	c.JSON(http.StatusOK, swap)
}

// ClaimSwaps godoc
// @Summary      Claim PENDING swaps
// @Description  Partner B pulls all PENDING swaps to be claimed
// @Tags         swaps
// @Produce      json
// @Param        source_partner_id query string true "Source Partner ID"
// @Success      200 {array} models.SwapLedger
// @Router       /v1/swap/claims [get]
func ClaimSwaps(c *gin.Context) {
	partnerID := c.Query("source_partner_id")
	if partnerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source_partner_id is required"})
		return
	}
	swaps, err := SwapService.ClaimSwaps(c.Request.Context(), partnerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, swaps)
}

// ConfirmSwap godoc
// @Summary      Confirm a swap
// @Description  Partner B confirms a swap (settlement)
// @Tags         swaps
// @Accept       json
// @Produce      json
// @Param        confirm body models.ConfirmSwapRequest true "Swap confirmation"
// @Success      200 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Router       /v1/swap/confirm [post]
func ConfirmSwap(c *gin.Context) {
	var req models.ConfirmSwapRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Update the swap status to COMPLETED
	err := SwapService.ConfirmSwap(c.Request.Context(), req.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": req.ID, "status": "COMPLETED"})
}

// ListSwaps godoc
// @Summary      List swap records with filters
// @Description  List swap records filtered by status, partner, or date range
// @Tags         swaps
// @Produce      json
// @Param        status query string false "Status (PENDING, COMPLETED, etc.)"
// @Param        source_partner_id query string false "Source Partner ID"
// @Param        target_partner_id query string false "Target Partner ID"
// @Param        from query string false "From date (YYYY-MM-DD)"
// @Param        to query string false "To date (YYYY-MM-DD)"
// @Success      200 {array} models.SwapLedger
// @Router       /v1/swaps [get]
func ListSwaps(c *gin.Context) {
	status := c.Query("status")
	sourcePartnerID := c.Query("source_partner_id")
	targetPartnerID := c.Query("target_partner_id")
	from := c.Query("from")
	to := c.Query("to")

	swaps, err := SwapService.ListSwapsWithFilter(c.Request.Context(), status, sourcePartnerID, targetPartnerID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, swaps)
}

// ListSwapsWithFilter godoc
// @Summary      List swap records with custom filters
// @Description  List swap records filtered by custom criteria
// @Tags         swaps
// @Produce      json
// @Param        status query string false "Status (PENDING, COMPLETED, etc.)"
// @Param        source_partner_id query string false "Source Partner ID"
// @Param        target_partner_id query string false "Target Partner ID"
// @Param        from query string false "From date (YYYY-MM-DD)"
// @Param        to query string false "To date (YYYY-MM-DD)"
// @Success      200 {array} models.SwapLedger
// @Router       /v1/swaps [get]
func ListSwapsWithFilter(c *gin.Context) {
	status := c.Query("status")
	sourcePartnerID := c.Query("source_partner_id")
	targetPartnerID := c.Query("target_partner_id")
	from := c.Query("from")
	to := c.Query("to")

	swaps, err := SwapService.ListSwapsWithFilter(c.Request.Context(), status, sourcePartnerID, targetPartnerID, from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, swaps)
}
