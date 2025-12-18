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

// MockSubscriptionRepository 訂閱 repository 的 mock
type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) Create(input *models.CreateSubscriptionInput) (*models.Subscription, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetByID(id uuid.UUID) (*models.Subscription, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) List(filters repository.SubscriptionFilters) ([]*models.Subscription, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Update(id uuid.UUID, input *models.UpdateSubscriptionInput) (*models.Subscription, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetDueBillings(date time.Time) ([]*models.Subscription, error) {
	args := m.Called(date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) GetExpiringSoon(days int) ([]*models.Subscription, error) {
	args := m.Called(days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Subscription), args.Error(1)
}

// TestSubscriptionService_CreateSubscription 測試建立訂閱
func TestSubscriptionService_CreateSubscription(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewSubscriptionService(mockRepo, mockCategoryRepo)

	categoryID := uuid.New()
	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "娛樂",
		Type: models.CashFlowTypeExpense,
	}

	input := &models.CreateSubscriptionInput{
		Name:          "Netflix",
		Amount:        390,
		BillingCycle:  models.BillingCycleMonthly,
		BillingDay:    15,
		CategoryID:    categoryID,
		PaymentMethod: models.PaymentMethodCash,
		StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		AutoRenew:     true,
	}

	expectedSubscription := &models.Subscription{
		ID:            uuid.New(),
		Name:          input.Name,
		Amount:        input.Amount,
		Currency:      models.CurrencyTWD,
		BillingCycle:  input.BillingCycle,
		BillingDay:    input.BillingDay,
		CategoryID:    input.CategoryID,
		PaymentMethod: input.PaymentMethod,
		StartDate:     input.StartDate,
		AutoRenew:     input.AutoRenew,
		Status:        models.SubscriptionStatusActive,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockRepo.On("Create", input).Return(expectedSubscription, nil)

	// 執行測試
	result, err := service.CreateSubscription(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSubscription.ID, result.ID)
	assert.Equal(t, expectedSubscription.Name, result.Name)
	mockCategoryRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestSubscriptionService_CreateSubscription_InvalidInput 測試建立訂閱時的輸入驗證
func TestSubscriptionService_CreateSubscription_InvalidInput(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewSubscriptionService(mockRepo, mockCategoryRepo)

	tests := []struct {
		name        string
		input       *models.CreateSubscriptionInput
		expectedErr string
	}{
		{
			name: "空白名稱",
			input: &models.CreateSubscriptionInput{
				Name:         "",
				Amount:       390,
				BillingCycle: models.BillingCycleMonthly,
				BillingDay:   15,
				CategoryID:   uuid.New(),
				StartDate:    time.Now(),
			},
			expectedErr: "subscription name is required",
		},
		{
			name: "金額為零",
			input: &models.CreateSubscriptionInput{
				Name:         "Netflix",
				Amount:       0,
				BillingCycle: models.BillingCycleMonthly,
				BillingDay:   15,
				CategoryID:   uuid.New(),
				StartDate:    time.Now(),
			},
			expectedErr: "amount must be greater than zero",
		},
		{
			name: "無效的計費週期",
			input: &models.CreateSubscriptionInput{
				Name:         "Netflix",
				Amount:       390,
				BillingCycle: "invalid",
				BillingDay:   15,
				CategoryID:   uuid.New(),
				StartDate:    time.Now(),
			},
			expectedErr: "invalid billing cycle",
		},
		{
			name: "無效的扣款日（小於 1）",
			input: &models.CreateSubscriptionInput{
				Name:         "Netflix",
				Amount:       390,
				BillingCycle: models.BillingCycleMonthly,
				BillingDay:   0,
				CategoryID:   uuid.New(),
				StartDate:    time.Now(),
			},
			expectedErr: "billing day must be between 1 and 31",
		},
		{
			name: "無效的扣款日（大於 31）",
			input: &models.CreateSubscriptionInput{
				Name:         "Netflix",
				Amount:       390,
				BillingCycle: models.BillingCycleMonthly,
				BillingDay:   32,
				CategoryID:   uuid.New(),
				StartDate:    time.Now(),
			},
			expectedErr: "billing day must be between 1 and 31",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateSubscription(tt.input)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

// TestSubscriptionService_GetSubscription 測試取得訂閱
func TestSubscriptionService_GetSubscription(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewSubscriptionService(mockRepo, mockCategoryRepo)

	subscriptionID := uuid.New()
	expectedSubscription := &models.Subscription{
		ID:     subscriptionID,
		Name:   "Netflix",
		Amount: 390,
		Status: models.SubscriptionStatusActive,
	}

	mockRepo.On("GetByID", subscriptionID).Return(expectedSubscription, nil)

	result, err := service.GetSubscription(subscriptionID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSubscription.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestSubscriptionService_CancelSubscription 測試取消訂閱
func TestSubscriptionService_CancelSubscription(t *testing.T) {
	mockRepo := new(MockSubscriptionRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	service := NewSubscriptionService(mockRepo, mockCategoryRepo)

	subscriptionID := uuid.New()
	endDate := time.Now()

	existingSubscription := &models.Subscription{
		ID:     subscriptionID,
		Name:   "Netflix",
		Amount: 390,
		Status: models.SubscriptionStatusActive,
	}

	updatedSubscription := &models.Subscription{
		ID:      subscriptionID,
		Name:    "Netflix",
		Amount:  390,
		Status:  models.SubscriptionStatusCancelled,
		EndDate: &endDate,
	}

	mockRepo.On("GetByID", subscriptionID).Return(existingSubscription, nil)
	mockRepo.On("Update", subscriptionID, mock.AnythingOfType("*models.UpdateSubscriptionInput")).Return(updatedSubscription, nil)

	result, err := service.CancelSubscription(subscriptionID, endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, models.SubscriptionStatusCancelled, result.Status)
	mockRepo.AssertExpectations(t)
}

