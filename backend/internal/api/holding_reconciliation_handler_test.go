package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockHoldingReconcileService struct{ mock.Mock }

func (m *mockHoldingReconcileService) ReconcileHoldings(_ context.Context, items []models.HoldingReconcileItem, dryRun bool) (*models.HoldingReconcilePreview, error) {
	args := m.Called(items, dryRun)
	if v := args.Get(0); v != nil {
		return v.(*models.HoldingReconcilePreview), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupHoldingReconcileRouter(svc *mockHoldingReconcileService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewHoldingReconciliationHandler(svc)
	r.POST("/api/holdings/reconcile/preview", h.Preview)
	r.POST("/api/holdings/reconcile", h.Execute)
	r.POST("/api/holdings/:symbol/reconcile", h.ExecuteSingle)
	return r
}

func TestHoldingReconcileHandler_Preview_OK(t *testing.T) {
	svc := &mockHoldingReconcileService{}
	preview := &models.HoldingReconcilePreview{Items: []models.HoldingReconcilePreviewItem{{Action: models.ReconcileActionUpdate}}}
	svc.On("ReconcileHoldings", mock.Anything, true).Return(preview, nil)

	body, _ := json.Marshal(models.HoldingReconcileBatchInput{Items: []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180},
	}})
	req := httptest.NewRequest("POST", "/api/holdings/reconcile/preview", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupHoldingReconcileRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHoldingReconcileHandler_Execute_ValidationErrorMapsTo400(t *testing.T) {
	svc := &mockHoldingReconcileService{}
	svc.On("ReconcileHoldings", mock.Anything, false).Return(nil, models.ErrInvalidReconcileQuantity)

	body, _ := json.Marshal(models.HoldingReconcileBatchInput{Items: []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: -1, TargetAvgCost: 100},
	}})
	req := httptest.NewRequest("POST", "/api/holdings/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupHoldingReconcileRouter(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHoldingReconcileHandler_Execute_GenericErrorMapsTo500(t *testing.T) {
	svc := &mockHoldingReconcileService{}
	svc.On("ReconcileHoldings", mock.Anything, false).Return(nil, errors.New("db down"))

	body, _ := json.Marshal(models.HoldingReconcileBatchInput{Items: []models.HoldingReconcileItem{
		{AssetType: models.AssetTypeUSStock, Symbol: "AAPL", Currency: models.CurrencyUSD, TargetQuantity: 150, TargetAvgCost: 180},
	}})
	req := httptest.NewRequest("POST", "/api/holdings/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupHoldingReconcileRouter(svc).ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestHoldingReconcileHandler_ExecuteSingle_PathSymbolOverridesBody(t *testing.T) {
	svc := &mockHoldingReconcileService{}
	svc.On("ReconcileHoldings",
		mock.MatchedBy(func(items []models.HoldingReconcileItem) bool {
			return len(items) == 1 && items[0].Symbol == "TSLA"
		}),
		false,
	).Return(&models.HoldingReconcilePreview{Items: []models.HoldingReconcilePreviewItem{{Action: models.ReconcileActionUpdate}}}, nil)

	body, _ := json.Marshal(models.HoldingReconcileItem{
		AssetType: models.AssetTypeUSStock, Symbol: "WRONG", Currency: models.CurrencyUSD,
		TargetQuantity: 50, TargetAvgCost: 250,
	})
	req := httptest.NewRequest("POST", "/api/holdings/TSLA/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupHoldingReconcileRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestHoldingReconcileHandler_BadJSON_400(t *testing.T) {
	svc := &mockHoldingReconcileService{}
	req := httptest.NewRequest("POST", "/api/holdings/reconcile", bytes.NewReader([]byte("{not-json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupHoldingReconcileRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
