package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCashFlowService 模擬的 CashFlowService
type MockCashFlowService struct {
	mock.Mock
}

func (m *MockCashFlowService) CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) GetCashFlow(id uuid.UUID) (*models.CashFlow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) DeleteCashFlow(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashFlowService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CashFlowSummary), args.Error(1)
}

// setupCashFlowTestRouter 設定測試用的 router
func setupCashFlowTestRouter(handler *CashFlowHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		cashFlows := api.Group("/cash-flows")
		{
			cashFlows.POST("", handler.CreateCashFlow)
			cashFlows.GET("", handler.ListCashFlows)
			cashFlows.GET("/summary", handler.GetSummary)
			cashFlows.GET("/:id", handler.GetCashFlow)
			cashFlows.PUT("/:id", handler.UpdateCashFlow)
			cashFlows.DELETE("/:id", handler.DeleteCashFlow)
		}
	}

	return router
}

// TestCreateCashFlow_Success 測試成功建立現金流記錄
func TestCreateCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	categoryID := uuid.New()
	input := models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateCashFlow", &input).Return(expectedCashFlow, nil)

	// 準備請求
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/cash-flows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestCreateCashFlow_InvalidInput 測試無效的輸入資料
func TestCreateCashFlow_InvalidInput(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 無效的 JSON
	invalidJSON := []byte(`{"invalid": json}`)

	req, _ := http.NewRequest("POST", "/api/cash-flows", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}

// TestGetCashFlow_Success 測試成功取得現金流記錄
func TestGetCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	cashFlowID := uuid.New()
	expectedCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		Amount:      50000,
		Description: "薪資",
	}

	mockService.On("GetCashFlow", cashFlowID).Return(expectedCashFlow, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/cash-flows/%s", cashFlowID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetCashFlow_InvalidID 測試無效的 ID
func TestGetCashFlow_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows/invalid-id", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INVALID_ID", response.Error.Code)
}

// TestListCashFlows_Success 測試成功取得現金流列表
func TestListCashFlows_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	expectedCashFlows := []*models.CashFlow{
		{
			ID:          uuid.New(),
			Type:        models.CashFlowTypeIncome,
			Amount:      50000,
			Description: "薪資",
		},
		{
			ID:          uuid.New(),
			Type:        models.CashFlowTypeExpense,
			Amount:      1200,
			Description: "午餐",
		},
	}

	mockService.On("ListCashFlows", mock.AnythingOfType("repository.CashFlowFilters")).Return(expectedCashFlows, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestDeleteCashFlow_Success 測試成功刪除現金流記錄
func TestDeleteCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	cashFlowID := uuid.New()
	mockService.On("DeleteCashFlow", cashFlowID).Return(nil)

	// 準備請求
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/cash-flows/%s", cashFlowID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetSummary_Success 測試成功取得摘要
func TestGetSummary_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	startDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	expectedSummary := &repository.CashFlowSummary{
		TotalIncome:  55000,
		TotalExpense: 15000,
		NetCashFlow:  40000,
	}

	mockService.On("GetSummary", startDate, endDate).Return(expectedSummary, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows/summary?start_date=2025-10-01&end_date=2025-10-31", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetSummary_MissingParameters 測試缺少參數
func TestGetSummary_MissingParameters(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 準備請求（缺少 end_date）
	req, _ := http.NewRequest("GET", "/api/cash-flows/summary?start_date=2025-10-01", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "MISSING_PARAMETERS", response.Error.Code)
}
