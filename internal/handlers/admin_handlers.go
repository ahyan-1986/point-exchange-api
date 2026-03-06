package handlers

import (
	"context"
	"net/http"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

// AdminServiceInterface abstracts the service for testability
type AdminServiceInterface interface {
	ListSwapLedger(ctx context.Context) ([]*models.SwapLedger, error)
}

// AdminService is injected at startup
var AdminService AdminServiceInterface

// AdminLedger godoc
// @Summary      Query swap ledger
// @Description  Admin: Query swap ledger with filters
// @Tags         admin
// @Produce      json
// @Success      200 {array} map[string]interface{}
// @Router       /v1/admin/ledger [get]
func AdminLedger(c *gin.Context) {
	ledgers, err := AdminService.ListSwapLedger(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ledgers)
}
