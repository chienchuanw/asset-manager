package service

import (
	"fmt"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCreditCardRepository 信用卡 repository 的 mock
type MockCreditCardRepository struct {
	mock.Mock
}

func (m *MockCreditCardRepository) Create(input *models.CreateCreditCardInput) (*models.CreditCard, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetAll() ([]*models.CreditCard, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetByBillingDay(day int) ([]*models.CreditCard, error) {
	args := m.Called(day)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetByPaymentDueDay(day int) ([]*models.CreditCard, error) {
	args := m.Called(day)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetUpcomingBilling(daysAhead int) ([]*models.CreditCard, error) {
	args := m.Called(daysAhead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) GetUpcomingPayment(daysAhead int) ([]*models.CreditCard, error) {
	args := m.Called(daysAhead)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) Update(id uuid.UUID, input *models.UpdateCreditCardInput) (*models.CreditCard, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CreditCard), args.Error(1)
}

func (m *MockCreditCardRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestCreditCardService_CreateCreditCard 測試建立信用卡
func TestCreditCardService_CreateCreditCard(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	input := &models.CreateCreditCardInput{
		IssuingBank:     "玉山銀行",
		CardName:        "Pi 拍錢包信用卡",
		CardNumberLast4: "1234",
		BillingDay:      5,
		PaymentDueDay:   20,
		CreditLimit:     100000,
		UsedCredit:      30000,
	}

	expectedCard := &models.CreditCard{
		ID:              uuid.New(),
		IssuingBank:     input.IssuingBank,
		CardName:        input.CardName,
		CardNumberLast4: input.CardNumberLast4,
		BillingDay:      input.BillingDay,
		PaymentDueDay:   input.PaymentDueDay,
		CreditLimit:     input.CreditLimit,
		UsedCredit:      input.UsedCredit,
	}

	mockRepo.On("Create", input).Return(expectedCard, nil)

	result, err := service.CreateCreditCard(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard.CardName, result.CardName)
	assert.Equal(t, expectedCard.CardNumberLast4, result.CardNumberLast4)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_CreateCreditCard_InvalidInput 測試建立信用卡 - 無效輸入
func TestCreditCardService_CreateCreditCard_InvalidInput(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	tests := []struct {
		name  string
		input *models.CreateCreditCardInput
	}{
		{
			name: "卡號後四碼長度不正確",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "123", // 只有 3 位
				BillingDay:      5,
				PaymentDueDay:   20,
				CreditLimit:     100000,
				UsedCredit:      30000,
			},
		},
		{
			name: "帳單日超出範圍",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "1234",
				BillingDay:      32, // 超過 31
				PaymentDueDay:   20,
				CreditLimit:     100000,
				UsedCredit:      30000,
			},
		},
		{
			name: "繳款截止日超出範圍",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "1234",
				BillingDay:      5,
				PaymentDueDay:   0, // 小於 1
				CreditLimit:     100000,
				UsedCredit:      30000,
			},
		},
		{
			name: "信用額度為 0",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "1234",
				BillingDay:      5,
				PaymentDueDay:   20,
				CreditLimit:     0, // 必須大於 0
				UsedCredit:      0,
			},
		},
		{
			name: "已使用額度超過信用額度",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "1234",
				BillingDay:      5,
				PaymentDueDay:   20,
				CreditLimit:     100000,
				UsedCredit:      150000, // 超過信用額度
			},
		},
		{
			name: "已使用額度為負數",
			input: &models.CreateCreditCardInput{
				IssuingBank:     "玉山銀行",
				CardName:        "Pi 拍錢包信用卡",
				CardNumberLast4: "1234",
				BillingDay:      5,
				PaymentDueDay:   20,
				CreditLimit:     100000,
				UsedCredit:      -1000, // 負數
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateCreditCard(tt.input)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

// TestCreditCardService_GetCreditCard 測試取得信用卡
func TestCreditCardService_GetCreditCard(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	expectedCard := &models.CreditCard{
		ID:              cardID,
		IssuingBank:     "玉山銀行",
		CardName:        "Pi 拍錢包信用卡",
		CardNumberLast4: "1234",
		BillingDay:      5,
		PaymentDueDay:   20,
		CreditLimit:     100000,
		UsedCredit:      30000,
	}

	mockRepo.On("GetByID", cardID).Return(expectedCard, nil)

	result, err := service.GetCreditCard(cardID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCard.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_GetCreditCard_NotFound 測試取得信用卡 - 找不到
func TestCreditCardService_GetCreditCard_NotFound(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	mockRepo.On("GetByID", cardID).Return(nil, fmt.Errorf("credit card not found"))

	result, err := service.GetCreditCard(cardID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_ListCreditCards 測試列出所有信用卡
func TestCreditCardService_ListCreditCards(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	expectedCards := []*models.CreditCard{
		{
			ID:              uuid.New(),
			IssuingBank:     "玉山銀行",
			CardName:        "Pi 拍錢包信用卡",
			CardNumberLast4: "1234",
			BillingDay:      5,
			PaymentDueDay:   20,
			CreditLimit:     100000,
			UsedCredit:      30000,
		},
		{
			ID:              uuid.New(),
			IssuingBank:     "國泰世華銀行",
			CardName:        "CUBE 卡",
			CardNumberLast4: "5678",
			BillingDay:      10,
			PaymentDueDay:   25,
			CreditLimit:     150000,
			UsedCredit:      50000,
		},
	}

	mockRepo.On("GetAll").Return(expectedCards, nil)

	result, err := service.ListCreditCards()

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_GetUpcomingBilling 測試取得即將到來的帳單日信用卡
func TestCreditCardService_GetUpcomingBilling(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	daysAhead := 7
	expectedCards := []*models.CreditCard{
		{
			ID:              uuid.New(),
			IssuingBank:     "玉山銀行",
			CardName:        "Pi 拍錢包信用卡",
			CardNumberLast4: "1234",
			BillingDay:      5,
			PaymentDueDay:   20,
			CreditLimit:     100000,
			UsedCredit:      30000,
		},
	}

	mockRepo.On("GetUpcomingBilling", daysAhead).Return(expectedCards, nil)

	result, err := service.GetUpcomingBilling(daysAhead)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_GetUpcomingPayment 測試取得即將到來的繳款截止日信用卡
func TestCreditCardService_GetUpcomingPayment(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	daysAhead := 7
	expectedCards := []*models.CreditCard{
		{
			ID:              uuid.New(),
			IssuingBank:     "玉山銀行",
			CardName:        "Pi 拍錢包信用卡",
			CardNumberLast4: "1234",
			BillingDay:      5,
			PaymentDueDay:   20,
			CreditLimit:     100000,
			UsedCredit:      30000,
		},
	}

	mockRepo.On("GetUpcomingPayment", daysAhead).Return(expectedCards, nil)

	result, err := service.GetUpcomingPayment(daysAhead)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_UpdateCreditCard 測試更新信用卡
func TestCreditCardService_UpdateCreditCard(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	newUsedCredit := 40000.0
	input := &models.UpdateCreditCardInput{
		UsedCredit: &newUsedCredit,
	}

	expectedCard := &models.CreditCard{
		ID:              cardID,
		IssuingBank:     "玉山銀行",
		CardName:        "Pi 拍錢包信用卡",
		CardNumberLast4: "1234",
		BillingDay:      5,
		PaymentDueDay:   20,
		CreditLimit:     100000,
		UsedCredit:      newUsedCredit,
	}

	mockRepo.On("Update", cardID, input).Return(expectedCard, nil)

	result, err := service.UpdateCreditCard(cardID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newUsedCredit, result.UsedCredit)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_UpdateCreditCard_InvalidInput 測試更新信用卡 - 無效輸入
func TestCreditCardService_UpdateCreditCard_InvalidInput(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	negativeUsedCredit := -1000.0
	input := &models.UpdateCreditCardInput{
		UsedCredit: &negativeUsedCredit,
	}

	result, err := service.UpdateCreditCard(cardID, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestCreditCardService_DeleteCreditCard 測試刪除信用卡
func TestCreditCardService_DeleteCreditCard(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	mockRepo.On("Delete", cardID).Return(nil)

	err := service.DeleteCreditCard(cardID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestCreditCardService_DeleteCreditCard_NotFound 測試刪除信用卡 - 找不到
func TestCreditCardService_DeleteCreditCard_NotFound(t *testing.T) {
	mockRepo := new(MockCreditCardRepository)
	service := NewCreditCardService(mockRepo)

	cardID := uuid.New()
	mockRepo.On("Delete", cardID).Return(fmt.Errorf("credit card not found"))

	err := service.DeleteCreditCard(cardID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

