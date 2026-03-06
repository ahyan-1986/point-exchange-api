package handlers_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"point-exchange-api/internal/handlers"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

type mockAdminService struct{}

func (m *mockAdminService) ListSwapLedger(ctx context.Context) ([]*models.SwapLedger, error) {
	return []*models.SwapLedger{{ID: "mock-swap-id", Status: "PENDING", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}

func setupAdminRouter() *gin.Engine {
	handlers.AdminService = &mockAdminService{}
	r := gin.Default()
	r.GET("/v1/admin/ledger", handlers.AdminLedger)
	return r
}

func TestAdminLedger(t *testing.T) {
	r := setupAdminRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/admin/ledger", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}
