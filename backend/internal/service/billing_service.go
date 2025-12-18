package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// BillingService 扣款服務介面
type BillingService interface {
	// ProcessSubscriptionBilling 處理訂閱扣款
	ProcessSubscriptionBilling(date time.Time) (*BillingResult, error)
	// ProcessInstallmentBilling 處理分期扣款
	ProcessInstallmentBilling(date time.Time) (*BillingResult, error)
	// ProcessDailyBilling 處理每日扣款（訂閱 + 分期）
	ProcessDailyBilling(date time.Time) (*DailyBillingResult, error)
}

// BillingResult 扣款結果
type BillingResult struct {
	ProcessedCount   int                `json:"processed_count"`   // 處理的數量
	FailedCount      int                `json:"failed_count"`      // 失敗的數量
	CreatedCashFlows []*models.CashFlow `json:"created_cash_flows"` // 建立的現金流記錄
	Errors           []BillingError     `json:"errors,omitempty"`   // 錯誤列表
}

// BillingError 扣款錯誤
type BillingError struct {
	ID      uuid.UUID `json:"id"`      // 訂閱或分期的 ID
	Name    string    `json:"name"`    // 訂閱或分期的名稱
	Message string    `json:"message"` // 錯誤訊息
}

// DailyBillingResult 每日扣款結果
type DailyBillingResult struct {
	Date              time.Time      `json:"date"`               // 扣款日期
	SubscriptionCount int            `json:"subscription_count"` // 訂閱扣款數量
	InstallmentCount  int            `json:"installment_count"`  // 分期扣款數量
	TotalAmount       float64        `json:"total_amount"`       // 總扣款金額
	SubscriptionResult *BillingResult `json:"subscription_result"` // 訂閱扣款結果
	InstallmentResult  *BillingResult `json:"installment_result"`  // 分期扣款結果
}

// billingService 扣款服務實作
type billingService struct {
	subscriptionRepo repository.SubscriptionRepository
	installmentRepo  repository.InstallmentRepository
	cashFlowRepo     repository.CashFlowRepository
}

// NewBillingService 建立新的扣款 service
func NewBillingService(
	subscriptionRepo repository.SubscriptionRepository,
	installmentRepo repository.InstallmentRepository,
	cashFlowRepo repository.CashFlowRepository,
) BillingService {
	return &billingService{
		subscriptionRepo: subscriptionRepo,
		installmentRepo:  installmentRepo,
		cashFlowRepo:     cashFlowRepo,
	}
}

// ProcessSubscriptionBilling 處理訂閱扣款
func (s *billingService) ProcessSubscriptionBilling(date time.Time) (*BillingResult, error) {
	// 取得當日需要扣款的訂閱
	subscriptions, err := s.subscriptionRepo.GetDueBillings(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get due subscriptions: %w", err)
	}

	result := &BillingResult{
		ProcessedCount:   0,
		FailedCount:      0,
		CreatedCashFlows: []*models.CashFlow{},
		Errors:           []BillingError{},
	}

	// 處理每個訂閱
	for _, subscription := range subscriptions {
		// 根據訂閱的付款方式設定現金流的來源類型
		sourceType := subscription.PaymentMethod.ToSourceType()
		var sourceID *uuid.UUID
		if subscription.AccountID != nil {
			sourceID = subscription.AccountID
		}

		// 建立現金流記錄
		cashFlowInput := &models.CreateCashFlowInput{
			Date:        date,
			Type:        models.CashFlowTypeExpense,
			CategoryID:  subscription.CategoryID,
			Amount:      subscription.Amount,
			Description: fmt.Sprintf("%s - 訂閱扣款", subscription.Name),
			SourceType:  &sourceType,
			SourceID:    sourceID,
		}

		if subscription.Note != nil {
			cashFlowInput.Note = subscription.Note
		}

		cashFlow, err := s.cashFlowRepo.Create(cashFlowInput)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BillingError{
				ID:      subscription.ID,
				Name:    subscription.Name,
				Message: fmt.Sprintf("failed to create cash flow: %v", err),
			})
			continue
		}

		result.ProcessedCount++
		result.CreatedCashFlows = append(result.CreatedCashFlows, cashFlow)
	}

	return result, nil
}

// ProcessInstallmentBilling 處理分期扣款
func (s *billingService) ProcessInstallmentBilling(date time.Time) (*BillingResult, error) {
	// 取得當日需要扣款的分期
	installments, err := s.installmentRepo.GetDueBillings(date)
	if err != nil {
		return nil, fmt.Errorf("failed to get due installments: %w", err)
	}

	result := &BillingResult{
		ProcessedCount:   0,
		FailedCount:      0,
		CreatedCashFlows: []*models.CashFlow{},
		Errors:           []BillingError{},
	}

	// 處理每個分期
	for _, installment := range installments {
		// 計算當前期數
		currentPeriod := installment.PaidCount + 1

		// 根據分期的付款方式設定現金流的來源類型
		sourceType := installment.PaymentMethod.ToSourceType()
		var sourceID *uuid.UUID
		if installment.AccountID != nil {
			sourceID = installment.AccountID
		}

		// 建立現金流記錄
		cashFlowInput := &models.CreateCashFlowInput{
			Date:        date,
			Type:        models.CashFlowTypeExpense,
			CategoryID:  installment.CategoryID,
			Amount:      installment.InstallmentAmount,
			Description: fmt.Sprintf("%s - 分期付款 (%d/%d)", installment.Name, currentPeriod, installment.InstallmentCount),
			SourceType:  &sourceType,
			SourceID:    sourceID,
		}

		if installment.Note != nil {
			cashFlowInput.Note = installment.Note
		}

		cashFlow, err := s.cashFlowRepo.Create(cashFlowInput)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, BillingError{
				ID:      installment.ID,
				Name:    installment.Name,
				Message: fmt.Sprintf("failed to create cash flow: %v", err),
			})
			continue
		}

		// 更新分期的已付期數
		newPaidCount := currentPeriod
		updateInput := &models.UpdateInstallmentInput{
			PaidCount: &newPaidCount,
		}

		// 如果已付完所有期數，更新狀態為已完成
		if currentPeriod >= installment.InstallmentCount {
			status := models.InstallmentStatusCompleted
			updateInput.Status = &status
		}

		_, err = s.installmentRepo.Update(installment.ID, updateInput)
		if err != nil {
			// 記錄錯誤但不影響現金流記錄的建立
			result.Errors = append(result.Errors, BillingError{
				ID:      installment.ID,
				Name:    installment.Name,
				Message: fmt.Sprintf("cash flow created but failed to update installment: %v", err),
			})
		}

		result.ProcessedCount++
		result.CreatedCashFlows = append(result.CreatedCashFlows, cashFlow)
	}

	return result, nil
}

// ProcessDailyBilling 處理每日扣款（訂閱 + 分期）
func (s *billingService) ProcessDailyBilling(date time.Time) (*DailyBillingResult, error) {
	// 處理訂閱扣款
	subscriptionResult, err := s.ProcessSubscriptionBilling(date)
	if err != nil {
		return nil, fmt.Errorf("failed to process subscription billing: %w", err)
	}

	// 處理分期扣款
	installmentResult, err := s.ProcessInstallmentBilling(date)
	if err != nil {
		return nil, fmt.Errorf("failed to process installment billing: %w", err)
	}

	// 計算總金額
	totalAmount := 0.0
	for _, cf := range subscriptionResult.CreatedCashFlows {
		totalAmount += cf.Amount
	}
	for _, cf := range installmentResult.CreatedCashFlows {
		totalAmount += cf.Amount
	}

	result := &DailyBillingResult{
		Date:               date,
		SubscriptionCount:  subscriptionResult.ProcessedCount,
		InstallmentCount:   installmentResult.ProcessedCount,
		TotalAmount:        totalAmount,
		SubscriptionResult: subscriptionResult,
		InstallmentResult:  installmentResult,
	}

	return result, nil
}

