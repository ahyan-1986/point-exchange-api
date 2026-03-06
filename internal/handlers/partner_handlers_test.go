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

// mockPartnerService implements the PartnerService interface for testing
// Only implements RegisterPartner for now

type mockPartnerService struct{}

func (m *mockPartnerService) RegisterPartner(ctx context.Context, partner *models.Partner) (string, error) {
	partner.ID = "mockid"
	partner.APIKey = "mockkey"
	partner.APISecret = "mocksecret"
	partner.IsActive = true
	partner.CreatedAt = time.Now()
	return partner.ID, nil
}

func (m *mockPartnerService) ListPartners(ctx context.Context) ([]*models.Partner, error) { return nil, nil }
func (m *mockPartnerService) GetPartner(ctx context.Context, id string) (*models.Partner, error) { return nil, nil }
func (m *mockPartnerService) ActivatePartner(ctx context.Context, id string, isActive bool) error { return nil }

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/v1/partners", handlers.RegisterPartner)
	r.GET("/v1/partners", handlers.ListPartners)
	r.GET("/v1/partners/:id", handlers.GetPartner)
	r.PATCH("/v1/partners/:id/activate", handlers.ActivatePartner)
	return r
}

func TestRegisterPartner(t *testing.T) {
	handlers.SetPartnerService(&mockPartnerService{})
	r := setupRouter()
	partner := map[string]string{"name": "TestPartner"}
	body, _ := json.Marshal(partner)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/v1/partners", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	var resp models.Partner
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid response: %v", err)
	}
	if resp.Name != "TestPartner" {
		t.Errorf("expected name TestPartner, got %s", resp.Name)
	}
	if resp.ID == "" || resp.APIKey == "" || resp.APISecret == "" {
		t.Error("missing ID, APIKey, or APISecret in response")
	}
	if !resp.IsActive {
		t.Error("expected IsActive true")
	}
	if time.Since(resp.CreatedAt) > time.Minute {
		t.Error("unexpected CreatedAt")
	}
}
