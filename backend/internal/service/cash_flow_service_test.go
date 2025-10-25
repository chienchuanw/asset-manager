package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCashFlowRepository 模擬的 CashFlowRepository
type MockCashFlowRepository struct {
	mock.Mock
}

func (m *MockCashFlowRepository) Create(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) GetByID(id uuid.UUID) (*models.CashFlow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) GetAll(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) Update(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashFlowRepository) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CashFlowSummary), args.Error(1)
}

// MockCategoryRepository 模擬的 CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) GetByID(id uuid.UUID) (*models.CashFlowCategory, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) GetAll(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	args := m.Called(flowType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) Update(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlowCategory), args.Error(1)
}

func (m *MockCategoryRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestCreateCashFlow_Success 測試成功建立現金流記錄
func TestCreateCashFlow_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	categoryID := uuid.New()
	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
	}

	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
	}

	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)
	mockCategoryRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCreateCashFlow_InvalidType 測試無效的現金流類型
func TestCreateCashFlow_InvalidType(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowType("invalid"),
		CategoryID:  uuid.New(),
		Amount:      50000,
		Description: "測試",
	}

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid cash flow type")
}

// TestCreateCashFlow_InvalidAmount 測試無效的金額
func TestCreateCashFlow_InvalidAmount(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	tests := []struct {
		name   string
		amount float64
	}{
		{"zero amount", 0},
		{"negative amount", -100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := &models.CreateCashFlowInput{
				Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
				Type:        models.CashFlowTypeIncome,
				CategoryID:  uuid.New(),
				Amount:      tt.amount,
				Description: "測試",
			}

			// Act
			result, err := service.CreateCashFlow(input)

			// Assert
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "amount must be greater than zero")
		})
	}
}

// TestCreateCashFlow_CategoryTypeMismatch 測試分類類型不匹配
func TestCreateCashFlow_CategoryTypeMismatch(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	categoryID := uuid.New()
	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "測試",
	}

	// 分類是支出類型，但現金流是收入類型
	wrongTypeCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "飲食",
		Type: models.CashFlowTypeExpense,
	}

	mockCategoryRepo.On("GetByID", categoryID).Return(wrongTypeCategory, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "does not match")
	mockCategoryRepo.AssertExpectations(t)
}

// TestGetCashFlow_Success 測試成功取得現金流記錄
func TestGetCashFlow_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	cashFlowID := uuid.New()
	expectedCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		Amount:      50000,
		Description: "薪資",
	}

	mockRepo.On("GetByID", cashFlowID).Return(expectedCashFlow, nil)

	// Act
	result, err := service.GetCashFlow(cashFlowID)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, cashFlowID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestListCashFlows_Success 測試成功取得現金流列表
func TestListCashFlows_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	filters := repository.CashFlowFilters{}
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

	mockRepo.On("GetAll", filters).Return(expectedCashFlows, nil)

	// Act
	result, err := service.ListCashFlows(filters)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestDeleteCashFlow_Success 測試成功刪除現金流記錄
func TestDeleteCashFlow_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	cashFlowID := uuid.New()
	mockRepo.On("Delete", cashFlowID).Return(nil)

	// Act
	err := service.DeleteCashFlow(cashFlowID)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestGetSummary_Success 測試成功取得摘要
func TestGetSummary_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo)

	startDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	expectedSummary := &repository.CashFlowSummary{
		TotalIncome:  55000,
		TotalExpense: 15000,
		NetCashFlow:  40000,
	}

	mockRepo.On("GetSummary", startDate, endDate).Return(expectedSummary, nil)

	// Act
	result, err := service.GetSummary(startDate, endDate)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 55000.0, result.TotalIncome)
	assert.Equal(t, 15000.0, result.TotalExpense)
	assert.Equal(t, 40000.0, result.NetCashFlow)
	mockRepo.AssertExpectations(t)
}
