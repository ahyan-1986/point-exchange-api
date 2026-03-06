package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"point-exchange-api/internal/handlers"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

type mockSwapService struct{}

func (m *mockSwapService) CreateSwap(ctx context.Context, req *models.SwapRequest) (string, error) {
	return "mock-swap-id", nil
}
func (m *mockSwapService) GetSwap(ctx context.Context, id string) (*models.SwapLedger, error) {
	if id == "mock-swap-id" {
		return &models.SwapLedger{ID: id, Status: "PENDING", CreatedAt: time.Now(), UpdatedAt: time.Now()}, nil
	}
	return nil, nil
}
func (m *mockSwapService) ClaimSwaps(ctx context.Context, sourcePartnerID string) ([]*models.SwapLedger, error) {
	return []*models.SwapLedger{{ID: "mock-swap-id", SourcePartnerID: sourcePartnerID, Status: "PENDING", CreatedAt: time.Now(), UpdatedAt: time.Now()}}, nil
}
func (m *mockSwapService) ListSwapsBySourcePartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	return []*models.SwapLedger{}, nil
}
func (m *mockSwapService) ListSwapsByTargetPartnerID(ctx context.Context, partnerID string) ([]*models.SwapLedger, error) {
	return []*models.SwapLedger{}, nil
}
func (m *mockSwapService) ConfirmSwap(ctx context.Context, id string) error {
	return nil
}
func (m *mockSwapService) ListSwapsWithFilter(ctx context.Context, status, sourcePartnerID, targetPartnerID, from, to string) ([]*models.SwapLedger, error) {
	return []*models.SwapLedger{}, nil
}

func setupSwapRouter() *gin.Engine {
	handlers.SwapService = &mockSwapService{}
	r := gin.Default()
	r.POST("/v1/swap/deposit", handlers.CreateDeposit)
	r.GET("/v1/swap/:id", handlers.GetSwap)
	r.GET("/v1/swap/claims", handlers.ClaimSwaps)
	return r
}

func TestCreateDeposit(t *testing.T) {
	r := setupSwapRouter()
	body, _ := json.Marshal(models.SwapRequest{SourcePartnerID: "1", SourceExternalID: "ext", SourceCustomerID: "cust", SourcePoints: 100, TargetPartnerID: "2", TargetCustomerID: "tcust"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/swap/deposit", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestGetSwap(t *testing.T) {
	r := setupSwapRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swap/mock-swap-id", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestClaimSwaps(t *testing.T) {
	r := setupSwapRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/swap/claims?source_partner_id=1", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}
