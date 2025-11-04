package service

import (
	"fmt"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBankAccountRepository 銀行帳戶 repository 的 mock
type MockBankAccountRepository struct {
	mock.Mock
}

func (m *MockBankAccountRepository) Create(input *models.CreateBankAccountInput) (*models.BankAccount, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetByID(id uuid.UUID) (*models.BankAccount, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) GetAll(currency *models.Currency) ([]*models.BankAccount, error) {
	args := m.Called(currency)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) Update(id uuid.UUID, input *models.UpdateBankAccountInput) (*models.BankAccount, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) UpdateBalance(id uuid.UUID, amount float64) (*models.BankAccount, error) {
	args := m.Called(id, amount)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.BankAccount), args.Error(1)
}

func (m *MockBankAccountRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestBankAccountService_CreateBankAccount 測試建立銀行帳戶
func TestBankAccountService_CreateBankAccount(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	input := &models.CreateBankAccountInput{
		BankName:           "台灣銀行",
		AccountType:        "活期存款",
		AccountNumberLast4: "1234",
		Currency:           models.CurrencyTWD,
		Balance:            50000,
	}

	expectedAccount := &models.BankAccount{
		ID:                 uuid.New(),
		BankName:           input.BankName,
		AccountType:        input.AccountType,
		AccountNumberLast4: input.AccountNumberLast4,
		Currency:           input.Currency,
		Balance:            input.Balance,
	}

	mockRepo.On("Create", input).Return(expectedAccount, nil)

	result, err := service.CreateBankAccount(input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAccount.BankName, result.BankName)
	assert.Equal(t, expectedAccount.AccountNumberLast4, result.AccountNumberLast4)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_CreateBankAccount_InvalidInput 測試建立銀行帳戶 - 無效輸入
func TestBankAccountService_CreateBankAccount_InvalidInput(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	tests := []struct {
		name  string
		input *models.CreateBankAccountInput
	}{
		{
			name: "帳號後四碼長度不正確",
			input: &models.CreateBankAccountInput{
				BankName:           "台灣銀行",
				AccountType:        "活期存款",
				AccountNumberLast4: "123", // 只有 3 位
				Currency:           models.CurrencyTWD,
				Balance:            50000,
			},
		},
		{
			name: "餘額為負數",
			input: &models.CreateBankAccountInput{
				BankName:           "台灣銀行",
				AccountType:        "活期存款",
				AccountNumberLast4: "1234",
				Currency:           models.CurrencyTWD,
				Balance:            -1000, // 負數
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.CreateBankAccount(tt.input)
			assert.Error(t, err)
			assert.Nil(t, result)
		})
	}
}

// TestBankAccountService_GetBankAccount 測試取得銀行帳戶
func TestBankAccountService_GetBankAccount(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	expectedAccount := &models.BankAccount{
		ID:                 accountID,
		BankName:           "台灣銀行",
		AccountType:        "活期存款",
		AccountNumberLast4: "1234",
		Currency:           models.CurrencyTWD,
		Balance:            50000,
	}

	mockRepo.On("GetByID", accountID).Return(expectedAccount, nil)

	result, err := service.GetBankAccount(accountID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedAccount.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_GetBankAccount_NotFound 測試取得銀行帳戶 - 找不到
func TestBankAccountService_GetBankAccount_NotFound(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	mockRepo.On("GetByID", accountID).Return(nil, fmt.Errorf("bank account not found"))

	result, err := service.GetBankAccount(accountID)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_ListBankAccounts 測試列出所有銀行帳戶
func TestBankAccountService_ListBankAccounts(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	expectedAccounts := []*models.BankAccount{
		{
			ID:                 uuid.New(),
			BankName:           "台灣銀行",
			AccountType:        "活期存款",
			AccountNumberLast4: "1234",
			Currency:           models.CurrencyTWD,
			Balance:            50000,
		},
		{
			ID:                 uuid.New(),
			BankName:           "玉山銀行",
			AccountType:        "外幣帳戶",
			AccountNumberLast4: "5678",
			Currency:           models.CurrencyUSD,
			Balance:            1000,
		},
	}

	mockRepo.On("GetAll", (*models.Currency)(nil)).Return(expectedAccounts, nil)

	result, err := service.ListBankAccounts(nil)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_ListBankAccounts_FilterByCurrency 測試列出銀行帳戶 - 依幣別篩選
func TestBankAccountService_ListBankAccounts_FilterByCurrency(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	currency := models.CurrencyTWD
	expectedAccounts := []*models.BankAccount{
		{
			ID:                 uuid.New(),
			BankName:           "台灣銀行",
			AccountType:        "活期存款",
			AccountNumberLast4: "1234",
			Currency:           models.CurrencyTWD,
			Balance:            50000,
		},
	}

	mockRepo.On("GetAll", &currency).Return(expectedAccounts, nil)

	result, err := service.ListBankAccounts(&currency)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, models.CurrencyTWD, result[0].Currency)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_UpdateBankAccount 測試更新銀行帳戶
func TestBankAccountService_UpdateBankAccount(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	newBalance := 60000.0
	input := &models.UpdateBankAccountInput{
		Balance: &newBalance,
	}

	expectedAccount := &models.BankAccount{
		ID:                 accountID,
		BankName:           "台灣銀行",
		AccountType:        "活期存款",
		AccountNumberLast4: "1234",
		Currency:           models.CurrencyTWD,
		Balance:            newBalance,
	}

	mockRepo.On("Update", accountID, input).Return(expectedAccount, nil)

	result, err := service.UpdateBankAccount(accountID, input)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newBalance, result.Balance)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_UpdateBankAccount_InvalidInput 測試更新銀行帳戶 - 無效輸入
func TestBankAccountService_UpdateBankAccount_InvalidInput(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	negativeBalance := -1000.0
	input := &models.UpdateBankAccountInput{
		Balance: &negativeBalance,
	}

	result, err := service.UpdateBankAccount(accountID, input)

	assert.Error(t, err)
	assert.Nil(t, result)
}

// TestBankAccountService_DeleteBankAccount 測試刪除銀行帳戶
func TestBankAccountService_DeleteBankAccount(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	mockRepo.On("Delete", accountID).Return(nil)

	err := service.DeleteBankAccount(accountID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

// TestBankAccountService_DeleteBankAccount_NotFound 測試刪除銀行帳戶 - 找不到
func TestBankAccountService_DeleteBankAccount_NotFound(t *testing.T) {
	mockRepo := new(MockBankAccountRepository)
	service := NewBankAccountService(mockRepo)

	accountID := uuid.New()
	mockRepo.On("Delete", accountID).Return(fmt.Errorf("bank account not found"))

	err := service.DeleteBankAccount(accountID)

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

