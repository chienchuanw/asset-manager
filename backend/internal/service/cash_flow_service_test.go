package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCashFlowService_CreateCashFlow 測試建立現金流記錄
func TestCashFlowService_CreateCashFlow(t *testing.T) {
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	}

	note := "十月薪資"
	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
		Note:        &note,
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		Note:        input.Note,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// 執行測試
	result, err := service.CreateCashFlow(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)
	mockCategoryRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCreateCashFlow_Success 測試成功建立現金流記錄
func TestCreateCashFlow_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	cashFlowID := uuid.New()

	// 建立一個現金交易記錄（manual 類型，不需要回復餘額）
	sourceType := models.SourceTypeManual
	cashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Type:        models.CashFlowTypeExpense,
		Amount:      100.0,
		SourceType:  &sourceType,
		SourceID:    nil,
	}

	mockRepo.On("GetByID", cashFlowID).Return(cashFlow, nil)
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
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

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

// TestCreateCashFlow_CreditCardPayment_Success 測試成功建立信用卡繳款記錄
func TestCreateCashFlow_CreditCardPayment_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	sourceType := models.SourceTypeBankAccount
	targetType := models.SourceTypeCreditCard

	// 建立轉帳分類
	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "轉帳",
		Type: models.CashFlowTypeTransferOut,
	}

	// 建立銀行帳戶
	bankAccount := &models.BankAccount{
		ID:      bankAccountID,
		Balance: 50000,
	}

	// 建立信用卡
	creditCard := &models.CreditCard{
		ID:          creditCardID,
		CreditLimit: 100000,
		UsedCredit:  30000,
	}

	// 繳款金額
	paymentAmount := 30000.0

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      paymentAmount,
		Description: "繳信用卡費",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
		TargetType:  &targetType,
		TargetID:    &creditCardID,
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
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedBankAccount := &models.BankAccount{
		ID:      bankAccountID,
		Balance: 20000, // 50000 - 30000
	}

	updatedCreditCard := &models.CreditCard{
		ID:          creditCardID,
		CreditLimit: 100000,
		UsedCredit:  0, // 30000 - 30000
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(bankAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, -paymentAmount).Return(updatedBankAccount, nil)
	mockCreditCardRepo.On("GetByID", creditCardID).Return(creditCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", creditCardID, -paymentAmount).Return(updatedCreditCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, expectedCashFlow.Amount, result.Amount)
	assert.Equal(t, models.CashFlowTypeTransferOut, result.Type)
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCreateCashFlow_CreditCardPayment_PartialPayment 測試部分繳款
func TestCreateCashFlow_CreditCardPayment_PartialPayment(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	bankAccountID := uuid.New()
	creditCardID := uuid.New()
	sourceType := models.SourceTypeBankAccount
	targetType := models.SourceTypeCreditCard

	category := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "轉帳",
		Type: models.CashFlowTypeTransferOut,
	}

	bankAccount := &models.BankAccount{
		ID:      bankAccountID,
		Balance: 50000,
	}

	creditCard := &models.CreditCard{
		ID:          creditCardID,
		CreditLimit: 100000,
		UsedCredit:  30000,
	}

	// 部分繳款金額
	paymentAmount := 15000.0

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 1, 20, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      paymentAmount,
		Description: "部分繳信用卡費",
		SourceType:  &sourceType,
		SourceID:    &bankAccountID,
		TargetType:  &targetType,
		TargetID:    &creditCardID,
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
		TargetType:  input.TargetType,
		TargetID:    input.TargetID,
	}

	updatedBankAccount := &models.BankAccount{
		ID:      bankAccountID,
		Balance: 35000, // 50000 - 15000
	}

	updatedCreditCard := &models.CreditCard{
		ID:          creditCardID,
		CreditLimit: 100000,
		UsedCredit:  15000, // 30000 - 15000 (還有 15000 未繳)
	}

	// 設定 mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(category, nil)
	mockBankAccountRepo.On("GetByID", bankAccountID).Return(bankAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", bankAccountID, -paymentAmount).Return(updatedBankAccount, nil)
	mockCreditCardRepo.On("GetByID", creditCardID).Return(creditCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", creditCardID, -paymentAmount).Return(updatedCreditCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, paymentAmount, result.Amount)
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
