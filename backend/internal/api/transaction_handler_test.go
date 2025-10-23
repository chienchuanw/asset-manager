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

// MockTransactionService 模擬的 service
type MockTransactionService struct {
	mock.Mock
}

func (m *MockTransactionService) CreateTransaction(input *models.CreateTransactionInput) (*models.Transaction, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) GetTransaction(id uuid.UUID) (*models.Transaction, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) ListTransactions(filters repository.TransactionFilters) ([]*models.Transaction, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) UpdateTransaction(id uuid.UUID, input *models.UpdateTransactionInput) (*models.Transaction, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Transaction), args.Error(1)
}

func (m *MockTransactionService) DeleteTransaction(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// setupTestRouter 設定測試用的 router
func setupTestRouter(handler *TransactionHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		transactions := api.Group("/transactions")
		{
			transactions.POST("", handler.CreateTransaction)
			transactions.GET("", handler.ListTransactions)
			transactions.GET("/:id", handler.GetTransaction)
			transactions.PUT("/:id", handler.UpdateTransaction)
			transactions.DELETE("/:id", handler.DeleteTransaction)
		}
	}

	return router
}

// TestCreateTransaction_Success 測試成功建立交易記錄
func TestCreateTransaction_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	fee := 28.0
	input := models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
	}

	expectedTransaction := &models.Transaction{
		ID:              uuid.New(),
		Date:            input.Date,
		AssetType:       input.AssetType,
		Symbol:          input.Symbol,
		Name:            input.Name,
		TransactionType: input.TransactionType,
		Quantity:        input.Quantity,
		Price:           input.Price,
		Amount:          input.Amount,
		Fee:             input.Fee,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	mockService.On("CreateTransaction", &input).Return(expectedTransaction, nil)

	// 準備請求
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(body))
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

// TestCreateTransaction_InvalidInput 測試無效的輸入資料
func TestCreateTransaction_InvalidInput(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	// 無效的 JSON
	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(invalidJSON))
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

// TestGetTransaction_Success 測試成功取得交易記錄
func TestGetTransaction_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	transactionID := uuid.New()
	expectedTransaction := &models.Transaction{
		ID:              transactionID,
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
	}

	mockService.On("GetTransaction", transactionID).Return(expectedTransaction, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/transactions/%s", transactionID), nil)
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

// TestGetTransaction_InvalidID 測試無效的 ID 格式
func TestGetTransaction_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/transactions/invalid-id", nil)
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

// TestListTransactions_Success 測試成功取得交易記錄列表
func TestListTransactions_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	expectedTransactions := []*models.Transaction{
		{
			ID:              uuid.New(),
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
		},
		{
			ID:              uuid.New(),
			Symbol:          "ETH",
			Name:            "Ethereum",
			TransactionType: models.TransactionTypeBuy,
		},
	}

	mockService.On("ListTransactions", mock.AnythingOfType("repository.TransactionFilters")).
		Return(expectedTransactions, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/transactions", nil)
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

// TestDeleteTransaction_Success 測試成功刪除交易記錄
func TestDeleteTransaction_Success(t *testing.T) {
	// Arrange
	mockService := new(MockTransactionService)
	handler := NewTransactionHandler(mockService)
	router := setupTestRouter(handler)

	transactionID := uuid.New()
	mockService.On("DeleteTransaction", transactionID).Return(nil)

	// 準備請求
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/transactions/%s", transactionID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)

	mockService.AssertExpectations(t)
}

