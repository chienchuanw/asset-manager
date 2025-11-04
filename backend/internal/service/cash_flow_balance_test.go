package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCashFlowService_CreateCashFlow_WithBankAccount 測試使用銀行帳戶建立現金流記錄
func TestCashFlowService_CreateCashFlow_WithBankAccount(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	bankAccountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	}

	bankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		SourceType:  input.SourceType,
		SourceID:    input.SourceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedBankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  60000, // 原本 10000 + 收入 50000
		Currency: models.CurrencyTWD,
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(bankAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, float64(50000)).Return(updatedBankAccount, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)

	// 驗證所有 mock 期望都被呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCashFlowService_CreateCashFlow_WithBankAccount_Expense 測試使用銀行帳戶建立支出記錄
func TestCashFlowService_CreateCashFlow_WithBankAccount_Expense(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	bankAccountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	}

	bankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1500,
		Description: "午餐",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		SourceType:  input.SourceType,
		SourceID:    input.SourceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedBankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  8500, // 原本 10000 - 支出 1500
		Currency: models.CurrencyTWD,
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(bankAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, float64(-1500)).Return(updatedBankAccount, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)

	// 驗證所有 mock 期望都被呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCashFlowService_CreateCashFlow_WithCreditCard 測試使用信用卡建立支出記錄
func TestCashFlowService_CreateCashFlow_WithCreditCard(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	creditCardID := uuid.New()
	sourceType := models.SourceTypeCreditCard

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "購物",
		Type: models.CashFlowTypeExpense,
	}

	creditCard := &models.CreditCard{
		ID:          creditCardID,
		IssuingBank: "玉山銀行",
		CardName:    "Pi 拍錢包信用卡",
		CreditLimit: 100000,
		UsedCredit:  20000,
	}

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      5000,
		Description: "網購",
		SourceType:  &sourceType,
		SourceID:    &creditCardID,
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		SourceType:  input.SourceType,
		SourceID:    input.SourceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedCreditCard := &models.CreditCard{
		ID:          creditCardID,
		IssuingBank: "玉山銀行",
		CardName:    "Pi 拍錢包信用卡",
		CreditLimit: 100000,
		UsedCredit:  25000, // 原本 20000 + 支出 5000
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockCreditCardRepo.On("GetByID", creditCardID).Return(creditCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", creditCardID, float64(5000)).Return(updatedCreditCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)

	// 驗證所有 mock 期望都被呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCashFlowService_CreateCashFlow_BankAccountNotFound 測試銀行帳戶不存在的情況
func TestCashFlowService_CreateCashFlow_BankAccountNotFound(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	bankAccountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	}

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1500,
		Description: "午餐",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(nil, assert.AnError)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "bank account not found")

	// 驗證所有 mock 期望都被呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
}

// TestCashFlowService_CreateCashFlow_CashTransaction 測試現金交易（不影響餘額）
func TestCashFlowService_CreateCashFlow_CashTransaction(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	sourceType := models.SourceTypeManual

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	}

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1500,
		Description: "現金午餐",
		SourceType:  &sourceType,
		SourceID:    nil, // 現金交易不需要 SourceID
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		SourceType:  input.SourceType,
		SourceID:    input.SourceID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)

	// 驗證所有 mock 期望都被呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
	// 注意：現金交易不應該呼叫銀行帳戶或信用卡的方法
}

// TestCashFlowService_DeleteCashFlow_WithBankAccount 測試刪除銀行帳戶現金流記錄並回復餘額
func TestCashFlowService_DeleteCashFlow_WithBankAccount(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	cashFlowID := uuid.New()
	bankAccountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	// 模擬要刪除的現金流記錄（支出）
	existingCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		Amount:      1500,
		Description: "午餐",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
	}

	updatedBankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  11500, // 回復餘額：原本減少了 1500，現在加回來
		Currency: models.CurrencyTWD,
	}

	// 設定 mock 期望
	mockRepo.On("GetByID", cashFlowID).Return(existingCashFlow, nil)
	mockRepo.On("Delete", cashFlowID).Return(nil)
	// 回復餘額時需要先驗證銀行帳戶存在
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(updatedBankAccount, nil)
	// 回復餘額：支出記錄刪除時，應該增加銀行帳戶餘額（+1500）
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, float64(1500)).Return(updatedBankAccount, nil)

	// Act
	err := service.DeleteCashFlow(cashFlowID)

	// Assert
	assert.NoError(t, err)

	// 驗證所有 mock 期望都被呼叫
	mockRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
}

// TestCashFlowService_DeleteCashFlow_WithCreditCard 測試刪除信用卡現金流記錄並回復額度
func TestCashFlowService_DeleteCashFlow_WithCreditCard(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	cashFlowID := uuid.New()
	creditCardID := uuid.New()
	sourceType := models.SourceTypeCreditCard

	// 模擬要刪除的現金流記錄（支出）
	existingCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		Amount:      5000,
		Description: "網購",
		SourceType:  &sourceType,
		SourceID:    &creditCardID,
	}

	updatedCreditCard := &models.CreditCard{
		ID:          creditCardID,
		IssuingBank: "玉山銀行",
		CardName:    "Pi 拍錢包信用卡",
		CreditLimit: 100000,
		UsedCredit:  15000, // 回復額度：原本增加了 5000，現在減回來
	}

	// 設定 mock 期望
	mockRepo.On("GetByID", cashFlowID).Return(existingCashFlow, nil)
	mockRepo.On("Delete", cashFlowID).Return(nil)
	// 回復額度時需要先驗證信用卡存在
	mockCreditCardRepo.On("GetByID", creditCardID).Return(updatedCreditCard, nil)
	// 回復額度：支出記錄刪除時，應該減少已使用額度（-5000）
	mockCreditCardRepo.On("UpdateUsedCredit", creditCardID, float64(-5000)).Return(updatedCreditCard, nil)

	// Act
	err := service.DeleteCashFlow(cashFlowID)

	// Assert
	assert.NoError(t, err)

	// 驗證所有 mock 期望都被呼叫
	mockRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
}

// TestCashFlowService_UpdateCashFlow_ChangePaymentMethod 測試更改付款方式
func TestCashFlowService_UpdateCashFlow_ChangePaymentMethod(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	cashFlowID := uuid.New()
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	oldSourceType := models.SourceTypeBankAccount
	newSourceType := models.SourceTypeCreditCard

	// 原始現金流記錄（銀行帳戶支出）
	originalCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		Amount:      1500,
		Description: "午餐",
		SourceType:  &oldSourceType,
		SourceID:    &bankAccountID,
	}

	// 更新輸入（改為信用卡）
	input := &models.UpdateCashFlowInput{
		SourceType: &newSourceType,
		SourceID:   &creditCardID,
	}

	// 更新後的現金流記錄
	updatedCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		Amount:      1500,
		Description: "午餐",
		SourceType:  &newSourceType,
		SourceID:    &creditCardID,
	}

	bankAccount := &models.BankAccount{
		ID:       bankAccountID,
		BankName: "台灣銀行",
		Balance:  11500, // 回復餘額：原本減少了 1500，現在加回來
		Currency: models.CurrencyTWD,
	}

	creditCard := &models.CreditCard{
		ID:          creditCardID,
		IssuingBank: "玉山銀行",
		CardName:    "Pi 拍錢包信用卡",
		CreditLimit: 100000,
		UsedCredit:  25000, // 新增使用額度：增加 1500
	}

	// 設定 mock 期望
	mockRepo.On("GetByID", cashFlowID).Return(originalCashFlow, nil)
	// 回復原本的銀行帳戶餘額時需要先驗證帳戶存在
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(bankAccount, nil)
	// 回復原本的銀行帳戶餘額（支出 1500 -> 加回 1500）
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, float64(1500)).Return(bankAccount, nil)
	// 驗證信用卡存在並更新額度
	mockCreditCardRepo.On("GetByID", creditCardID).Return(creditCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", creditCardID, float64(1500)).Return(creditCard, nil)
	// 更新記錄
	mockRepo.On("Update", cashFlowID, input).Return(updatedCashFlow, nil)

	// Act
	result, err := service.UpdateCashFlow(cashFlowID, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, updatedCashFlow.ID, result.ID)
	assert.Equal(t, newSourceType, *result.SourceType)
	assert.Equal(t, creditCardID, *result.SourceID)

	// 驗證所有 mock 期望都被呼叫
	mockRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
}
