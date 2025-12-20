package service

import (
	"fmt"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateCategory_Success 測試成功建立分類
func TestCreateCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.CreateCategoryInput{
		Name: "投資收入",
		Type: models.CashFlowTypeIncome,
	}

	expectedCategory := &models.CashFlowCategory{
		ID:       uuid.New(),
		Name:     input.Name,
		Type:     input.Type,
		IsSystem: false,
	}

	mockRepo.On("Create", input).Return(expectedCategory, nil)

	// Act
	result, err := service.CreateCategory(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCategory.ID, result.ID)
	assert.Equal(t, input.Name, result.Name)
	assert.Equal(t, input.Type, result.Type)
	mockRepo.AssertExpectations(t)
}

// TestCreateCategory_InvalidType 測試無效的分類類型
func TestCreateCategory_InvalidType(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.CreateCategoryInput{
		Name: "測試分類",
		Type: models.CashFlowType("invalid"),
	}

	// Act
	result, err := service.CreateCategory(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid cash flow type")
}

// TestCreateCategory_EmptyName 測試空白分類名稱
func TestCreateCategory_EmptyName(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.CreateCategoryInput{
		Name: "",
		Type: models.CashFlowTypeIncome,
	}

	// Act
	result, err := service.CreateCategory(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
}

// TestCreateCategory_NameTooLong 測試分類名稱過長
func TestCreateCategory_NameTooLong(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.CreateCategoryInput{
		Name: "這是一個超過二十個字元的分類名稱測試",
		Type: models.CashFlowTypeIncome,
	}

	// Act
	result, err := service.CreateCategory(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "must not exceed 20 characters")
}

// TestGetCategory_Success 測試成功取得分類
func TestGetCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	}

	mockRepo.On("GetByID", categoryID).Return(expectedCategory, nil)

	// Act
	result, err := service.GetCategory(categoryID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, categoryID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestListCategories_Success 測試成功取得分類列表
func TestListCategories_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

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

	mockRepo.On("GetAll", (*models.CashFlowType)(nil)).Return(expectedCategories, nil)

	// Act
	result, err := service.ListCategories(nil)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestListCategories_WithFilter 測試篩選分類列表
func TestListCategories_WithFilter(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	incomeType := models.CashFlowTypeIncome
	expectedCategories := []*models.CashFlowCategory{
		{
			ID:       uuid.New(),
			Name:     "薪資",
			Type:     models.CashFlowTypeIncome,
			IsSystem: true,
		},
	}

	mockRepo.On("GetAll", &incomeType).Return(expectedCategories, nil)

	// Act
	result, err := service.ListCategories(&incomeType)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, models.CashFlowTypeIncome, result[0].Type)
	mockRepo.AssertExpectations(t)
}

// TestUpdateCategory_Success 測試成功更新分類
func TestUpdateCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	input := &models.UpdateCategoryInput{
		Name: "更新後的名稱",
	}

	expectedCategory := &models.CashFlowCategory{
		ID:       categoryID,
		Name:     input.Name,
		Type:     models.CashFlowTypeIncome,
		IsSystem: false,
	}

	mockRepo.On("Update", categoryID, input).Return(expectedCategory, nil)

	// Act
	result, err := service.UpdateCategory(categoryID, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, input.Name, result.Name)
	mockRepo.AssertExpectations(t)
}

// TestUpdateCategory_EmptyName 測試更新為空白名稱
func TestUpdateCategory_EmptyName(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	input := &models.UpdateCategoryInput{
		Name: "",
	}

	// Act
	result, err := service.UpdateCategory(categoryID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "name is required")
}

// TestUpdateCategory_NameTooLong 測試更新名稱過長
func TestUpdateCategory_NameTooLong(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	input := &models.UpdateCategoryInput{
		Name: "這是一個超過二十個字元的分類名稱測試",
	}

	// Act
	result, err := service.UpdateCategory(categoryID, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "must not exceed 20 characters")
}

// TestDeleteCategory_Success 測試成功刪除分類
func TestDeleteCategory_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	mockRepo.On("IsInUse", categoryID).Return(false, nil)
	mockRepo.On("Delete", categoryID).Return(nil)

	// Act
	err := service.DeleteCategory(categoryID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestDeleteCategory_InUse 測試刪除正在使用的分類（應該失敗）
func TestDeleteCategory_InUse(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	mockRepo.On("IsInUse", categoryID).Return(true, nil)

	// Act
	err := service.DeleteCategory(categoryID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "being used by existing cash flow records")
	mockRepo.AssertExpectations(t)
}

// TestDeleteCategory_SystemCategory 測試刪除系統分類（應該失敗）
func TestDeleteCategory_SystemCategory(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	mockRepo.On("IsInUse", categoryID).Return(false, nil)
	mockRepo.On("Delete", categoryID).Return(fmt.Errorf("category not found or is a system category"))

	// Act
	err := service.DeleteCategory(categoryID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "system category")
	mockRepo.AssertExpectations(t)
}

// TestDeleteCategory_CheckInUseFailed 測試檢查分類使用狀態失敗
func TestDeleteCategory_CheckInUseFailed(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	categoryID := uuid.New()
	mockRepo.On("IsInUse", categoryID).Return(false, fmt.Errorf("database error"))

	// Act
	err := service.DeleteCategory(categoryID)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to check if category is in use")
	mockRepo.AssertExpectations(t)
}

// TestReorderCategories_Success 測試成功重新排序分類
func TestReorderCategories_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.ReorderCategoryInput{
		Orders: []models.CategoryOrderItem{
			{ID: uuid.New(), SortOrder: 0},
			{ID: uuid.New(), SortOrder: 1},
			{ID: uuid.New(), SortOrder: 2},
		},
	}

	mockRepo.On("Reorder", input).Return(nil)

	// Act
	err := service.ReorderCategories(input)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestReorderCategories_EmptyOrders 測試空的排序列表
func TestReorderCategories_EmptyOrders(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.ReorderCategoryInput{
		Orders: []models.CategoryOrderItem{},
	}

	// Act
	err := service.ReorderCategories(input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one category order is required")
}

// TestReorderCategories_RepositoryError 測試 repository 錯誤
func TestReorderCategories_RepositoryError(t *testing.T) {
	// Arrange
	mockRepo := new(MockCategoryRepository)
	service := NewCategoryService(mockRepo)

	input := &models.ReorderCategoryInput{
		Orders: []models.CategoryOrderItem{
			{ID: uuid.New(), SortOrder: 0},
		},
	}

	mockRepo.On("Reorder", input).Return(fmt.Errorf("database error"))

	// Act
	err := service.ReorderCategories(input)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to reorder categories")
	mockRepo.AssertExpectations(t)
}
