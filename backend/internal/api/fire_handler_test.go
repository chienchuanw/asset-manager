package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFireService 模擬的 FireService
type MockFireService struct {
	mock.Mock
}

func (m *MockFireService) GetProjection(input models.FireProjectionInput) (*models.FireProjectionResult, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FireProjectionResult), args.Error(1)
}

func setupFireTestRouter(handler *FireHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	api := router.Group("/api")
	{
		fire := api.Group("/fire")
		{
			fire.GET("/projection", handler.GetProjection)
		}
	}
	return router
}

func TestFireHandler_GetProjection_HappyPath(t *testing.T) {
	mockService := new(MockFireService)
	handler := NewFireHandler(mockService)
	router := setupFireTestRouter(handler)

	years := 12
	mockService.On("GetProjection", mock.Anything).Return(&models.FireProjectionResult{
		FireNumber: 15_000_000,
		YearsToFI:  &years,
		OnTrack:    true,
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/fire/projection", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp APIResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.NotNil(t, resp.Data)
}

func TestFireHandler_GetProjection_ParsesQueryOverrides(t *testing.T) {
	mockService := new(MockFireService)
	handler := NewFireHandler(mockService)
	router := setupFireTestRouter(handler)

	var captured models.FireProjectionInput
	mockService.On("GetProjection", mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(0).(models.FireProjectionInput)
		}).
		Return(&models.FireProjectionResult{}, nil)

	req := httptest.NewRequest(http.MethodGet,
		"/api/fire/projection?netWorth=1000000&annualExpenses=600000&annualSavings=500000&expectedReturn=0.06&withdrawalRate=0.035",
		nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotNil(t, captured.NetWorth)
	assert.Equal(t, 1_000_000.0, *captured.NetWorth)
	assert.NotNil(t, captured.AnnualExpenses)
	assert.Equal(t, 600_000.0, *captured.AnnualExpenses)
	assert.NotNil(t, captured.AnnualSavings)
	assert.Equal(t, 500_000.0, *captured.AnnualSavings)
	assert.NotNil(t, captured.ExpectedReturn)
	assert.InDelta(t, 0.06, *captured.ExpectedReturn, 1e-9)
	assert.NotNil(t, captured.WithdrawalRate)
	assert.InDelta(t, 0.035, *captured.WithdrawalRate, 1e-9)
}

func TestFireHandler_GetProjection_InvalidWithdrawalRateReturns400(t *testing.T) {
	mockService := new(MockFireService)
	handler := NewFireHandler(mockService)
	router := setupFireTestRouter(handler)

	mockService.On("GetProjection", mock.Anything).
		Return(nil, service.ErrInvalidWithdrawalRate)

	req := httptest.NewRequest(http.MethodGet, "/api/fire/projection?withdrawalRate=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestFireHandler_GetProjection_IgnoresUnparseableQueryParams(t *testing.T) {
	mockService := new(MockFireService)
	handler := NewFireHandler(mockService)
	router := setupFireTestRouter(handler)

	var captured models.FireProjectionInput
	mockService.On("GetProjection", mock.Anything).
		Run(func(args mock.Arguments) {
			captured = args.Get(0).(models.FireProjectionInput)
		}).
		Return(&models.FireProjectionResult{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/fire/projection?netWorth=abc", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Nil(t, captured.NetWorth)
}
