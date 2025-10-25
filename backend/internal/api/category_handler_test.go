package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryService 模擬的 CategoryService
type MockCategoryService struct {
	mock.Mock
}

func (m *MockCategoryService) CreateCategory(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryService) GetCategory(id uuid.UUID) (*models.CashFlowCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryService) ListCategories(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	args := m.Called(flowType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryService) UpdateCategory(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryService) DeleteCategory(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// setupCategoryTestRouter 設定測試用的 router
func setupCategoryTestRouter(handler *CategoryHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		categories := api.Group("/categories")
		{
			categories.POST("", handler.CreateCategory)
			categories.GET("", handler.ListCategories)
			categories.GET("/:id", handler.GetCategory)
			categories.PUT("/:id", handler.UpdateCategory)
			categories.DELETE("/:id", handler.DeleteCategory)
		}
	}

	return router
}

// TestCreateCategory_Success 測試成功建立分類
func TestCreateCategory_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	input := models.CreateCategoryInput{
		Name: "投資收入",
		Type: models.CashFlowTypeIncome,
	}

	expectedCategory := &models.CashFlowCategory{
		ID:       uuid.New(),
		Name:     input.Name,
		Type:     input.Type,
		IsSystem: false,
	}

	mockService.On("CreateCategory", &input).Return(expectedCategory, nil)

	// 準備請求
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(body))
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

// TestCreateCategory_InvalidInput 測試無效的輸入資料
func TestCreateCategory_InvalidInput(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	// 無效的 JSON
	invalidJSON := []byte(`{"invalid": json}`)

	req, _ := http.NewRequest("POST", "/api/categories", bytes.NewBuffer(invalidJSON))
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

// TestGetCategory_Success 測試成功取得分類
func TestGetCategory_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	categoryID := uuid.New()
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	}

	mockService.On("GetCategory", categoryID).Return(expectedCategory, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/categories/%s", categoryID), nil)
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

// TestGetCategory_InvalidID 測試無效的 ID
func TestGetCategory_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/categories/invalid-id", nil)
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

// TestListCategories_Success 測試成功取得分類列表
func TestListCategories_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	expectedCategories := []*models.CashFlowCategory{
		{
			ID:       uuid.New(),
			Name:     "薪資",
			Type:     models.CashFlowTypeIncome,
			IsSystem: true,
		},
		{
			ID:       uuid.New(),
			Name:     "獎金",
			Type:     models.CashFlowTypeIncome,
			IsSystem: true,
		},
	}

	mockService.On("ListCategories", (*models.CashFlowType)(nil)).Return(expectedCategories, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/categories", nil)
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

// TestDeleteCategory_Success 測試成功刪除分類
func TestDeleteCategory_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCategoryService)
	handler := NewCategoryHandler(mockService)
	router := setupCategoryTestRouter(handler)

	categoryID := uuid.New()
	mockService.On("DeleteCategory", categoryID).Return(nil)

	// 準備請求
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/categories/%s", categoryID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

