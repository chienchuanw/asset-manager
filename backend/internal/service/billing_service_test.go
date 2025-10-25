package service

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// TestBillingService_ProcessSubscriptionBilling 測試處理訂閱扣款
func TestBillingService_ProcessSubscriptionBilling(t *testing.T) {
	mockSubscriptionRepo := new(MockSubscriptionRepository)
	mockInstallmentRepo := new(MockInstallmentRepository)
	mockCashFlowRepo := new(MockCashFlowRepository)

	service := NewBillingService(mockSubscriptionRepo, mockInstallmentRepo, mockCashFlowRepo)

	today := time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)
	categoryID := uuid.New()

	// 準備測試資料
	subscriptions := []*models.Subscription{
		{
			ID:           uuid.New(),
			Name:         "Netflix",
			Amount:       390,
			Currency:     models.CurrencyTWD,
			BillingCycle: models.BillingCycleMonthly,
			BillingDay:   15,
			CategoryID:   categoryID,
			StartDate:    time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			Status:       models.SubscriptionStatusActive,
		},
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        today,
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      390,
		Currency:    models.CurrencyTWD,
		Description: "Netflix - 訂閱扣款",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 設定 mock 期望
	mockSubscriptionRepo.On("GetDueBillings", today).Return(subscriptions, nil)
	mockCashFlowRepo.On("Create", mock.AnythingOfType("*models.CreateCashFlowInput")).Return(expectedCashFlow, nil)

	// 執行測試
	result, err := service.ProcessSubscriptionBilling(today)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ProcessedCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.CreatedCashFlows, 1)
	mockSubscriptionRepo.AssertExpectations(t)
	mockCashFlowRepo.AssertExpectations(t)
}

// TestBillingService_ProcessInstallmentBilling 測試處理分期扣款
func TestBillingService_ProcessInstallmentBilling(t *testing.T) {
	mockSubscriptionRepo := new(MockSubscriptionRepository)
	mockInstallmentRepo := new(MockInstallmentRepository)
	mockCashFlowRepo := new(MockCashFlowRepository)

	service := NewBillingService(mockSubscriptionRepo, mockInstallmentRepo, mockCashFlowRepo)

	today := time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)
	categoryID := uuid.New()
	installmentID := uuid.New()

	// 準備測試資料
	installments := []*models.Installment{
		{
			ID:                installmentID,
			Name:              "iPhone 15 Pro",
			TotalAmount:       36000,
			Currency:          models.CurrencyTWD,
			InstallmentCount:  12,
			InstallmentAmount: 3000,
			InterestRate:      0,
			TotalInterest:     0,
			PaidCount:         5,
			BillingDay:        15,
			CategoryID:        categoryID,
			StartDate:         time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
			Status:            models.InstallmentStatusActive,
		},
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        today,
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      3000,
		Currency:    models.CurrencyTWD,
		Description: "iPhone 15 Pro - 分期付款 (6/12)",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	updatedInstallment := &models.Installment{
		ID:                installmentID,
		Name:              "iPhone 15 Pro",
		TotalAmount:       36000,
		Currency:          models.CurrencyTWD,
		InstallmentCount:  12,
		InstallmentAmount: 3000,
		InterestRate:      0,
		TotalInterest:     0,
		PaidCount:         6, // 已付期數增加
		BillingDay:        15,
		CategoryID:        categoryID,
		StartDate:         time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
		Status:            models.InstallmentStatusActive,
	}

	// 設定 mock 期望
	mockInstallmentRepo.On("GetDueBillings", today).Return(installments, nil)
	mockCashFlowRepo.On("Create", mock.AnythingOfType("*models.CreateCashFlowInput")).Return(expectedCashFlow, nil)
	mockInstallmentRepo.On("Update", installmentID, mock.AnythingOfType("*models.UpdateInstallmentInput")).Return(updatedInstallment, nil)

	// 執行測試
	result, err := service.ProcessInstallmentBilling(today)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ProcessedCount)
	assert.Equal(t, 0, result.FailedCount)
	assert.Len(t, result.CreatedCashFlows, 1)
	mockInstallmentRepo.AssertExpectations(t)
	mockCashFlowRepo.AssertExpectations(t)
}

// TestBillingService_ProcessDailyBilling 測試處理每日扣款
func TestBillingService_ProcessDailyBilling(t *testing.T) {
	mockSubscriptionRepo := new(MockSubscriptionRepository)
	mockInstallmentRepo := new(MockInstallmentRepository)
	mockCashFlowRepo := new(MockCashFlowRepository)

	service := NewBillingService(mockSubscriptionRepo, mockInstallmentRepo, mockCashFlowRepo)

	today := time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC)

	// 設定 mock 期望
	mockSubscriptionRepo.On("GetDueBillings", today).Return([]*models.Subscription{}, nil)
	mockInstallmentRepo.On("GetDueBillings", today).Return([]*models.Installment{}, nil)

	// 執行測試
	result, err := service.ProcessDailyBilling(today)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 0, result.SubscriptionCount)
	assert.Equal(t, 0, result.InstallmentCount)
	mockSubscriptionRepo.AssertExpectations(t)
	mockInstallmentRepo.AssertExpectations(t)
}

