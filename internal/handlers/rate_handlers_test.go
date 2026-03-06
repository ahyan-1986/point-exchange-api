package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"point-exchange-api/internal/handlers"
	"point-exchange-api/models"

	"github.com/gin-gonic/gin"
)

type mockRateService struct{}

func (m *mockRateService) AddOrUpdateRate(ctx context.Context, partnerID string, req *models.AddOrUpdateRateRequest) error {
	return nil
}
func (m *mockRateService) ListRates(ctx context.Context, partnerID string) ([]*models.Rate, error) {
	return []*models.Rate{{ID: 1, PartnerID: partnerID, PointType: "A", PointsPerUSD: 1000, MinExchangePoints: 10}}, nil
}

func setupRateRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/v1/partners/:id/rates", handlers.AddOrUpdateRate)
	r.GET("/v1/partners/:id/rates", handlers.ListRates)
	return r
}

func TestAddOrUpdateRate(t *testing.T) {
	handlers.RateService = &mockRateService{}
	r := setupRateRouter()
	body, _ := json.Marshal(models.AddOrUpdateRateRequest{PointType: "A", PointsPerUSD: 1000, MinExchangePoints: 10})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/partners/1/rates", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
}

func TestListRates(t *testing.T) {
	handlers.RateService = &mockRateService{}
	r := setupRateRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/partners/1/rates", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	var resp []models.Rate
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if len(resp) != 1 || resp[0].PartnerID != "1" {
		t.Errorf("unexpected response: %+v", resp)
	}
}
