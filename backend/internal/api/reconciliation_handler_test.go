package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockReconciliationService struct{ mock.Mock }

func (m *mockReconciliationService) ReconcileBankAccount(id uuid.UUID, in *models.ReconcileBankAccountInput) (*models.ReconcileResult, error) {
	args := m.Called(id, in)
	if v := args.Get(0); v != nil {
		return v.(*models.ReconcileResult), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockReconciliationService) ReconcileCreditCard(id uuid.UUID, in *models.ReconcileCreditCardInput) (*models.ReconcileResult, error) {
	args := m.Called(id, in)
	if v := args.Get(0); v != nil {
		return v.(*models.ReconcileResult), args.Error(1)
	}
	return nil, args.Error(1)
}
func (m *mockReconciliationService) ReconcileBatch(in *models.ReconcileBatchInput) ([]*models.ReconcileResult, error) {
	args := m.Called(in)
	if v := args.Get(0); v != nil {
		return v.([]*models.ReconcileResult), args.Error(1)
	}
	return nil, args.Error(1)
}

func setupReconciliationRouter(svc *mockReconciliationService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewReconciliationHandler(svc)
	r.POST("/api/bank-accounts/:id/reconcile", h.ReconcileBankAccount)
	r.POST("/api/credit-cards/:id/reconcile", h.ReconcileCreditCard)
	r.POST("/api/user-management/reconcile-batch", h.ReconcileBatch)
	return r
}

func TestReconcileBankAccountHandler_Success(t *testing.T) {
	svc := &mockReconciliationService{}
	cfID := uuid.New()
	id := uuid.New()
	bal := 45000.0
	svc.On("ReconcileBankAccount", id, mock.Anything).Return(&models.ReconcileResult{
		TargetType: models.SourceTypeBankAccount, TargetID: id,
		Delta: 15000, CashFlowID: &cfID, NewBalance: &bal,
	}, nil)

	body, _ := json.Marshal(models.ReconcileBankAccountInput{NewBalance: 45000, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/bank-accounts/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestReconcileBankAccountHandler_InvalidUUID(t *testing.T) {
	svc := &mockReconciliationService{}
	body, _ := json.Marshal(models.ReconcileBankAccountInput{NewBalance: 1, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/bank-accounts/not-a-uuid/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReconcileBankAccountHandler_ServiceError(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	svc.On("ReconcileBankAccount", id, mock.Anything).Return(nil, errors.New("boom"))

	body, _ := json.Marshal(models.ReconcileBankAccountInput{NewBalance: 1, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/bank-accounts/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestReconcileBankAccountHandler_NotFoundMapsTo404(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	svc.On("ReconcileBankAccount", id, mock.Anything).Return(nil, service.ErrBankAccountNotFound)

	body, _ := json.Marshal(models.ReconcileBankAccountInput{NewBalance: 1, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/bank-accounts/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReconcileCreditCardHandler_NotFoundMapsTo404(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	used := 1.0
	svc.On("ReconcileCreditCard", id, mock.Anything).Return(nil, service.ErrCreditCardNotFound)

	body, _ := json.Marshal(models.ReconcileCreditCardInput{NewUsedCredit: &used, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/credit-cards/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReconcileCreditCardHandler_UsedExceedsLimitMapsTo400(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	used := 99999.0
	svc.On("ReconcileCreditCard", id, mock.Anything).Return(nil, service.ErrUsedExceedsLimit)

	body, _ := json.Marshal(models.ReconcileCreditCardInput{NewUsedCredit: &used, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/credit-cards/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReconcileBatchHandler_InvalidInputMapsTo400(t *testing.T) {
	svc := &mockReconciliationService{}
	bal := 1.0
	id := uuid.New()
	svc.On("ReconcileBatch", mock.Anything).Return(nil, service.ErrInvalidInput)

	body, _ := json.Marshal(models.ReconcileBatchInput{
		Date: time.Now(),
		Items: []models.ReconcileBatchItem{
			{TargetType: models.SourceTypeBankAccount, TargetID: id, NewBalance: &bal},
		},
	})
	req := httptest.NewRequest("POST", "/api/user-management/reconcile-batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestReconcileCreditCardHandler_Success(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	used := 8000.0
	svc.On("ReconcileCreditCard", id, mock.Anything).Return(&models.ReconcileResult{
		TargetType: models.SourceTypeCreditCard, TargetID: id, Delta: 3000, NewUsedCredit: &used,
	}, nil)

	body, _ := json.Marshal(models.ReconcileCreditCardInput{NewUsedCredit: &used, Date: time.Now()})
	req := httptest.NewRequest("POST", "/api/credit-cards/"+id.String()+"/reconcile", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReconcileBatchHandler_Success(t *testing.T) {
	svc := &mockReconciliationService{}
	id := uuid.New()
	bal := 100.0
	svc.On("ReconcileBatch", mock.Anything).Return([]*models.ReconcileResult{
		{TargetType: models.SourceTypeBankAccount, TargetID: id, Delta: 0, NewBalance: &bal},
	}, nil)

	body, _ := json.Marshal(models.ReconcileBatchInput{
		Date: time.Now(),
		Items: []models.ReconcileBatchItem{
			{TargetType: models.SourceTypeBankAccount, TargetID: id, NewBalance: &bal},
		},
	})
	req := httptest.NewRequest("POST", "/api/user-management/reconcile-batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestReconcileBatchHandler_EmptyItems(t *testing.T) {
	svc := &mockReconciliationService{}
	body := []byte(`{"date":"2026-04-29T00:00:00Z","items":[]}`)
	req := httptest.NewRequest("POST", "/api/user-management/reconcile-batch", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupReconciliationRouter(svc).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
