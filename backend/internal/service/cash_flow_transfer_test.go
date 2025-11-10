package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCreateTransferIn_Success 測試成功建立存入記錄
func TestCreateTransferIn_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferIn,
		CategoryID:  categoryID,
		Amount:      10000,
		Description: "存入現金",
		SourceType:  &sourceType,
		SourceID:    &accountID,
	}

	// 模擬分類資料（轉帳類型的分類）
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferIn,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  5000, // 原始餘額
		Currency: models.CurrencyTWD,
	}

	// 模擬更新後的帳戶資料
	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  15000, // 存入後餘額：5000 + 10000
		Currency: models.CurrencyTWD,
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
	}

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(10000)).Return(updatedAccount, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)
	assert.Equal(t, models.CashFlowTypeTransferIn, result.Type)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCreateTransferOut_Success 測試成功建立轉出記錄
func TestCreateTransferOut_Success(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      3000,
		Description: "提領現金",
		SourceType:  &sourceType,
		SourceID:    &accountID,
	}

	// 模擬分類資料（轉帳類型的分類）
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferOut,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  10000, // 原始餘額
		Currency: models.CurrencyTWD,
	}

	// 模擬更新後的帳戶資料
	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  7000, // 提領後餘額：10000 - 3000
		Currency: models.CurrencyTWD,
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
	}

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-3000)).Return(updatedAccount, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)
	assert.Equal(t, models.CashFlowTypeTransferOut, result.Type)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestCreateTransfer_AllowNegativeBalance 測試轉出允許帳戶餘額變成負數
func TestCreateTransfer_AllowNegativeBalance(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      15000, // 提領金額大於帳戶餘額
		Description: "大額提領",
		SourceType:  &sourceType,
		SourceID:    &accountID,
	}

	// 模擬分類資料
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferOut,
	}

	// 模擬銀行帳戶資料（餘額不足）
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  5000, // 原始餘額小於提領金額
		Currency: models.CurrencyTWD,
	}

	// 模擬更新後的帳戶資料（允許負餘額）
	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  -10000, // 提領後餘額：5000 - 15000 = -10000
		Currency: models.CurrencyTWD,
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
	}

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-15000)).Return(updatedAccount, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestDeleteTransfer_RevertBalance 測試刪除轉帳記錄時回復餘額
func TestDeleteTransfer_RevertBalance(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	cashFlowID := uuid.New()
	accountID := uuid.New()
	sourceType := models.SourceTypeBankAccount

	// 模擬要刪除的轉入記錄
	existingCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Type:        models.CashFlowTypeTransferIn,
		Amount:      5000,
		Description: "存入現金",
		SourceType:  &sourceType,
		SourceID:    &accountID,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	// 模擬刪除後回復餘額的帳戶資料
	revertedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  5000, // 回復餘額：10000 - 5000
		Currency: models.CurrencyTWD,
	}

	// 設定 Mock 期望
	mockRepo.On("GetByID", cashFlowID).Return(existingCashFlow, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	// 刪除轉入記錄時，需要減少餘額（相當於回復之前的增加）
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-5000)).Return(revertedAccount, nil)
	mockRepo.On("Delete", cashFlowID).Return(nil)

	// Act
	err := service.DeleteCashFlow(cashFlowID)

	// Assert
	assert.NoError(t, err)

	// 驗證所有 Mock 都被正確呼叫
	mockRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
}

// TestPayCreditCard_PaymentLessThanUsedCredit 測試繳款金額小於已使用額度
func TestPayCreditCard_PaymentLessThanUsedCredit(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	cardID := uuid.New()
	sourceType := models.SourceTypeBankAccount
	targetType := models.SourceTypeCreditCard

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      3000, // 繳款金額
		Description: "繳信用卡費",
		SourceType:  &sourceType,
		SourceID:    &accountID,
		TargetType:  &targetType,
		TargetID:    &cardID,
	}

	// 模擬分類資料
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferOut,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  7000, // 轉出後餘額：10000 - 3000
		Currency: models.CurrencyTWD,
	}

	// 模擬信用卡資料（已使用額度 5000）
	expectedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  5000, // 已使用額度
	}

	updatedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  2000, // 繳款後已使用額度：5000 - 3000
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

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-3000)).Return(updatedAccount, nil)
	mockCreditCardRepo.On("GetByID", cardID).Return(expectedCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", cardID, float64(-3000)).Return(updatedCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestPayCreditCard_PaymentEqualsUsedCredit 測試繳款金額等於已使用額度（全額繳清）
func TestPayCreditCard_PaymentEqualsUsedCredit(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	cardID := uuid.New()
	sourceType := models.SourceTypeBankAccount
	targetType := models.SourceTypeCreditCard

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      5000, // 繳款金額等於已使用額度
		Description: "繳信用卡費（全額繳清）",
		SourceType:  &sourceType,
		SourceID:    &accountID,
		TargetType:  &targetType,
		TargetID:    &cardID,
	}

	// 模擬分類資料
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferOut,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  5000, // 轉出後餘額：10000 - 5000
		Currency: models.CurrencyTWD,
	}

	// 模擬信用卡資料（已使用額度 5000）
	expectedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  5000, // 已使用額度
	}

	updatedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  0, // 全額繳清後已使用額度歸零
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

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-5000)).Return(updatedAccount, nil)
	mockCreditCardRepo.On("GetByID", cardID).Return(expectedCard, nil)
	mockCreditCardRepo.On("UpdateUsedCredit", cardID, float64(-5000)).Return(updatedCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

// TestPayCreditCard_PaymentExceedsUsedCredit 測試繳款金額大於已使用額度（超額繳款）
func TestPayCreditCard_PaymentExceedsUsedCredit(t *testing.T) {
	// Arrange
	mockRepo := new(MockCashFlowRepository)
	mockCategoryRepo := new(MockCategoryRepository)
	mockBankAccountRepo := new(MockBankAccountRepository)
	mockCreditCardRepo := new(MockCreditCardRepository)
	service := NewCashFlowService(mockRepo, mockCategoryRepo, mockBankAccountRepo, mockCreditCardRepo)

	categoryID := uuid.New()
	accountID := uuid.New()
	cardID := uuid.New()
	sourceType := models.SourceTypeBankAccount
	targetType := models.SourceTypeCreditCard

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeTransferOut,
		CategoryID:  categoryID,
		Amount:      8000, // 繳款金額大於已使用額度
		Description: "繳信用卡費（超額繳款）",
		SourceType:  &sourceType,
		SourceID:    &accountID,
		TargetType:  &targetType,
		TargetID:    &cardID,
	}

	// 模擬分類資料
	expectedCategory := &models.CashFlowCategory{
		ID:   categoryID,
		Name: "移轉",
		Type: models.CashFlowTypeTransferOut,
	}

	// 模擬銀行帳戶資料
	expectedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  10000,
		Currency: models.CurrencyTWD,
	}

	updatedAccount := &models.BankAccount{
		ID:       accountID,
		BankName: "測試銀行",
		Balance:  2000, // 轉出後餘額：10000 - 8000
		Currency: models.CurrencyTWD,
	}

	// 模擬信用卡資料（已使用額度 5000）
	expectedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  5000, // 已使用額度
	}

	updatedCard := &models.CreditCard{
		ID:          cardID,
		IssuingBank: "測試銀行",
		CardName:    "測試卡",
		CreditLimit: 100000,
		UsedCredit:  0, // 超額繳款後已使用額度歸零（不會變成負數）
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

	// 設定 Mock 期望
	mockCategoryRepo.On("GetByID", categoryID).Return(expectedCategory, nil)
	mockBankAccountRepo.On("GetByID", accountID).Return(expectedAccount, nil)
	mockBankAccountRepo.On("UpdateBalance", accountID, float64(-8000)).Return(updatedAccount, nil)
	mockCreditCardRepo.On("GetByID", cardID).Return(expectedCard, nil)
	// 關鍵：超額繳款時，只減少到 0，所以實際減少的金額是 -5000 而不是 -8000
	mockCreditCardRepo.On("UpdateUsedCredit", cardID, float64(-5000)).Return(updatedCard, nil)
	mockRepo.On("Create", input).Return(expectedCashFlow, nil)

	// Act
	result, err := service.CreateCashFlow(input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedCashFlow.ID, result.ID)
	assert.Equal(t, input.Amount, result.Amount)

	// 驗證所有 Mock 都被正確呼叫
	mockCategoryRepo.AssertExpectations(t)
	mockBankAccountRepo.AssertExpectations(t)
	mockCreditCardRepo.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}
