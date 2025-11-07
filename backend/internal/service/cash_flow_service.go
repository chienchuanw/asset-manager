package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// CashFlowService 現金流記錄業務邏輯介面
type CashFlowService interface {
	CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error)
	GetCashFlow(id uuid.UUID) (*models.CashFlow, error)
	ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error)
	UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error)
	DeleteCashFlow(id uuid.UUID) error
	GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error)
	GetMonthlySummaryWithComparison(year, month int) (*models.MonthlyCashFlowSummary, error)
	GetYearlySummaryWithComparison(year int) (*models.YearlyCashFlowSummary, error)
}

// cashFlowService 現金流記錄業務邏輯實作
type cashFlowService struct {
	repo            repository.CashFlowRepository
	categoryRepo    repository.CategoryRepository
	bankAccountRepo repository.BankAccountRepository
	creditCardRepo  repository.CreditCardRepository
}

// NewCashFlowService 建立新的現金流記錄 service
func NewCashFlowService(
	repo repository.CashFlowRepository,
	categoryRepo repository.CategoryRepository,
	bankAccountRepo repository.BankAccountRepository,
	creditCardRepo repository.CreditCardRepository,
) CashFlowService {
	return &cashFlowService{
		repo:            repo,
		categoryRepo:    categoryRepo,
		bankAccountRepo: bankAccountRepo,
		creditCardRepo:  creditCardRepo,
	}
}

// CreateCashFlow 建立新的現金流記錄
func (s *cashFlowService) CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	// 驗證現金流類型
	if !input.Type.Validate() {
		return nil, fmt.Errorf("invalid cash flow type: %s", input.Type)
	}

	// 驗證金額
	if input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 驗證描述
	if input.Description == "" {
		return nil, fmt.Errorf("description is required")
	}

	if len(input.Description) > 500 {
		return nil, fmt.Errorf("description must not exceed 500 characters")
	}

	// 驗證分類是否存在且類型匹配
	category, err := s.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	// 確認分類類型與現金流類型一致
	if category.Type != input.Type {
		return nil, fmt.Errorf("category type (%s) does not match cash flow type (%s)", category.Type, input.Type)
	}

	// 對於轉帳類型，強制要求選擇銀行帳戶
	if input.Type == models.CashFlowTypeTransferIn || input.Type == models.CashFlowTypeTransferOut {
		if input.SourceType == nil || input.SourceID == nil {
			return nil, fmt.Errorf("bank account is required for transfer transactions")
		}
		if *input.SourceType != models.SourceTypeBankAccount {
			return nil, fmt.Errorf("only bank account is allowed for transfer transactions")
		}
	}

	// 驗證並處理付款方式 (source)
	if input.SourceType != nil && input.SourceID != nil {
		err := s.validateAndUpdateBalance(input.Type, *input.SourceType, *input.SourceID, input.Amount)
		if err != nil {
			return nil, fmt.Errorf("failed to update balance: %w", err)
		}
	}

	// 處理轉帳目標 (target) - 只在 transfer_out 類型時處理
	if input.Type == models.CashFlowTypeTransferOut && input.TargetType != nil && input.TargetID != nil {
		err := s.validateAndUpdateTarget(*input.TargetType, *input.TargetID, input.Amount)
		if err != nil {
			// 如果目標更新失敗，需要回復 source 的餘額變動
			if input.SourceType != nil && input.SourceID != nil {
				s.revertBalanceUpdate(input.Type, *input.SourceType, *input.SourceID, input.Amount)
			}
			return nil, fmt.Errorf("failed to update target balance: %w", err)
		}
	}

	// 呼叫 repository 建立現金流記錄
	cashFlow, err := s.repo.Create(input)
	if err != nil {
		// 如果建立失敗，需要回復所有餘額變動
		if input.SourceType != nil && input.SourceID != nil {
			s.revertBalanceUpdate(input.Type, *input.SourceType, *input.SourceID, input.Amount)
		}
		if input.Type == models.CashFlowTypeTransferOut && input.TargetType != nil && input.TargetID != nil {
			s.revertTargetUpdate(*input.TargetType, *input.TargetID, input.Amount)
		}
		return nil, fmt.Errorf("failed to create cash flow: %w", err)
	}

	return cashFlow, nil
}

// GetCashFlow 取得單筆現金流記錄
func (s *cashFlowService) GetCashFlow(id uuid.UUID) (*models.CashFlow, error) {
	return s.repo.GetByID(id)
}

// ListCashFlows 取得現金流記錄列表
func (s *cashFlowService) ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	// 驗證篩選條件
	if filters.Type != nil && !filters.Type.Validate() {
		return nil, fmt.Errorf("invalid cash flow type filter: %s", *filters.Type)
	}

	// 驗證日期範圍
	if filters.StartDate != nil && filters.EndDate != nil {
		if filters.StartDate.After(*filters.EndDate) {
			return nil, fmt.Errorf("start date must be before or equal to end date")
		}
	}

	return s.repo.GetAll(filters)
}

// UpdateCashFlow 更新現金流記錄
func (s *cashFlowService) UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	// 驗證金額
	if input.Amount != nil && *input.Amount <= 0 {
		return nil, fmt.Errorf("amount must be greater than zero")
	}

	// 驗證描述
	if input.Description != nil {
		if *input.Description == "" {
			return nil, fmt.Errorf("description cannot be empty")
		}
		if len(*input.Description) > 500 {
			return nil, fmt.Errorf("description must not exceed 500 characters")
		}
	}

	// 先取得原始記錄以處理餘額變動和分類驗證
	original, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("cash flow not found: %w", err)
	}

	// 如果要更新分類，需要驗證分類是否存在
	if input.CategoryID != nil {
		// 驗證新分類
		category, err := s.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}

		// 確認分類類型與現金流類型一致
		if category.Type != original.Type {
			return nil, fmt.Errorf("category type (%s) does not match cash flow type (%s)", category.Type, original.Type)
		}
	}

	// 處理付款方式變更時的餘額調整
	err = s.handlePaymentMethodChange(original, input)
	if err != nil {
		return nil, fmt.Errorf("failed to handle payment method change: %w", err)
	}

	// 呼叫 repository 更新記錄
	cashFlow, err := s.repo.Update(id, input)
	if err != nil {
		// 如果更新失敗，需要回復餘額變動
		s.revertPaymentMethodChange(original, input)
		return nil, fmt.Errorf("failed to update cash flow: %w", err)
	}

	return cashFlow, nil
}

// DeleteCashFlow 刪除現金流記錄
func (s *cashFlowService) DeleteCashFlow(id uuid.UUID) error {
	// 先取得現金流記錄以便回復餘額
	cashFlow, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get cash flow for deletion: %w", err)
	}

	// 回復餘額變動（在刪除記錄之前）
	if cashFlow.SourceType != nil && cashFlow.SourceID != nil {
		err = s.revertBalanceUpdateWithError(cashFlow.Type, *cashFlow.SourceType, *cashFlow.SourceID, cashFlow.Amount)
		if err != nil {
			return fmt.Errorf("failed to revert balance update: %w", err)
		}
	}

	// 刪除現金流記錄
	err = s.repo.Delete(id)
	if err != nil {
		// 如果刪除失敗，需要重新套用餘額變動
		if cashFlow.SourceType != nil && cashFlow.SourceID != nil {
			s.validateAndUpdateBalance(cashFlow.Type, *cashFlow.SourceType, *cashFlow.SourceID, cashFlow.Amount)
		}
		return fmt.Errorf("failed to delete cash flow: %w", err)
	}

	return nil
}

// GetSummary 取得指定日期區間的現金流摘要
func (s *cashFlowService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	// 驗證日期範圍
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	return s.repo.GetSummary(startDate, endDate)
}

// validateAndUpdateBalance 驗證付款方式並更新對應的餘額
func (s *cashFlowService) validateAndUpdateBalance(cashFlowType models.CashFlowType, sourceType models.SourceType, sourceID uuid.UUID, amount float64) error {
	// 驗證 SourceType 是否有效
	if !sourceType.Validate() {
		return fmt.Errorf("invalid source type: %s", sourceType)
	}

	// 根據 SourceType 處理不同的付款方式
	switch sourceType {
	case models.SourceTypeBankAccount:
		return s.updateBankAccountBalance(cashFlowType, sourceID, amount)
	case models.SourceTypeCreditCard:
		return s.updateCreditCardBalance(cashFlowType, sourceID, amount)
	case models.SourceTypeManual:
		// 現金交易，不需要更新任何餘額
		return nil
	default:
		// 其他類型（訂閱、分期）暫時不處理餘額更新
		return nil
	}
}

// revertBalanceUpdate 回復餘額變動（用於交易失敗時的回滾）
func (s *cashFlowService) revertBalanceUpdate(cashFlowType models.CashFlowType, sourceType models.SourceType, sourceID uuid.UUID, amount float64) {
	// 回復操作：將原本的金額變動反向操作
	switch sourceType {
	case models.SourceTypeBankAccount:
		s.updateBankAccountBalance(cashFlowType, sourceID, -amount)
	case models.SourceTypeCreditCard:
		s.updateCreditCardBalance(cashFlowType, sourceID, -amount)
	}
	// 忽略錯誤，因為這是回復操作
}

// revertBalanceUpdateWithError 回復餘額變動並返回錯誤（用於刪除操作）
func (s *cashFlowService) revertBalanceUpdateWithError(cashFlowType models.CashFlowType, sourceType models.SourceType, sourceID uuid.UUID, amount float64) error {
	// 回復操作：將原本的金額變動反向操作
	switch sourceType {
	case models.SourceTypeBankAccount:
		return s.updateBankAccountBalance(cashFlowType, sourceID, -amount)
	case models.SourceTypeCreditCard:
		return s.updateCreditCardBalance(cashFlowType, sourceID, -amount)
	default:
		return nil // manual 類型不需要回復餘額
	}
}

// updateBankAccountBalance 更新銀行帳戶餘額
func (s *cashFlowService) updateBankAccountBalance(cashFlowType models.CashFlowType, accountID uuid.UUID, amount float64) error {
	// 驗證銀行帳戶是否存在
	_, err := s.bankAccountRepo.GetByID(accountID)
	if err != nil {
		return fmt.Errorf("bank account not found: %w", err)
	}

	// 計算餘額變動
	var balanceChange float64
	switch cashFlowType {
	case models.CashFlowTypeIncome, models.CashFlowTypeTransferIn:
		// 收入或存入：增加銀行帳戶餘額
		balanceChange = amount
	case models.CashFlowTypeExpense, models.CashFlowTypeTransferOut:
		// 支出或轉出：減少銀行帳戶餘額
		balanceChange = -amount
	}

	// 更新餘額
	_, err = s.bankAccountRepo.UpdateBalance(accountID, balanceChange)
	if err != nil {
		return fmt.Errorf("failed to update bank account balance: %w", err)
	}

	return nil
}

// updateCreditCardBalance 更新信用卡已使用額度
func (s *cashFlowService) updateCreditCardBalance(cashFlowType models.CashFlowType, cardID uuid.UUID, amount float64) error {
	// 驗證信用卡是否存在
	_, err := s.creditCardRepo.GetByID(cardID)
	if err != nil {
		return fmt.Errorf("credit card not found: %w", err)
	}

	// 計算已使用額度變動
	var usedCreditChange float64
	if cashFlowType == models.CashFlowTypeIncome {
		// 收入：減少已使用額度（例如退款）
		usedCreditChange = -amount
	} else {
		// 支出：增加已使用額度
		usedCreditChange = amount
	}

	// 更新已使用額度
	_, err = s.creditCardRepo.UpdateUsedCredit(cardID, usedCreditChange)
	if err != nil {
		return fmt.Errorf("failed to update credit card used credit: %w", err)
	}

	return nil
}

// validateAndUpdateTarget 驗證並更新轉帳目標的餘額
func (s *cashFlowService) validateAndUpdateTarget(targetType models.SourceType, targetID uuid.UUID, amount float64) error {
	switch targetType {
	case models.SourceTypeCreditCard:
		// 驗證信用卡是否存在
		_, err := s.creditCardRepo.GetByID(targetID)
		if err != nil {
			return fmt.Errorf("credit card not found: %w", err)
		}

		// 繳款給信用卡 → 減少已使用額度
		_, err = s.creditCardRepo.UpdateUsedCredit(targetID, -amount)
		if err != nil {
			return fmt.Errorf("failed to update credit card used credit: %w", err)
		}

	case models.SourceTypeBankAccount:
		// 驗證銀行帳戶是否存在
		_, err := s.bankAccountRepo.GetByID(targetID)
		if err != nil {
			return fmt.Errorf("bank account not found: %w", err)
		}

		// 轉帳到銀行帳戶 → 增加餘額
		_, err = s.bankAccountRepo.UpdateBalance(targetID, amount)
		if err != nil {
			return fmt.Errorf("failed to update bank account balance: %w", err)
		}

	default:
		return fmt.Errorf("invalid target type: %s", targetType)
	}

	return nil
}

// revertTargetUpdate 回復轉帳目標的餘額變動
func (s *cashFlowService) revertTargetUpdate(targetType models.SourceType, targetID uuid.UUID, amount float64) {
	switch targetType {
	case models.SourceTypeCreditCard:
		// 回復信用卡已使用額度（加回去）
		s.creditCardRepo.UpdateUsedCredit(targetID, amount)

	case models.SourceTypeBankAccount:
		// 回復銀行帳戶餘額（減回去）
		s.bankAccountRepo.UpdateBalance(targetID, -amount)
	}
}

// handlePaymentMethodChange 處理付款方式變更時的餘額調整
func (s *cashFlowService) handlePaymentMethodChange(original *models.CashFlow, input *models.UpdateCashFlowInput) error {
	// 檢查是否有付款方式相關的變更
	sourceTypeChanged := input.SourceType != nil && (original.SourceType == nil || *input.SourceType != *original.SourceType)
	sourceIDChanged := input.SourceID != nil && (original.SourceID == nil || *input.SourceID != *original.SourceID)
	amountChanged := input.Amount != nil && *input.Amount != original.Amount

	// 如果沒有相關變更，直接返回
	if !sourceTypeChanged && !sourceIDChanged && !amountChanged {
		return nil
	}

	// 先回復原本的餘額變動
	if original.SourceType != nil && original.SourceID != nil {
		s.revertBalanceUpdate(original.Type, *original.SourceType, *original.SourceID, original.Amount)
	}

	// 計算新的金額
	newAmount := original.Amount
	if input.Amount != nil {
		newAmount = *input.Amount
	}

	// 計算新的付款方式
	newSourceType := original.SourceType
	if input.SourceType != nil {
		newSourceType = input.SourceType
	}

	newSourceID := original.SourceID
	if input.SourceID != nil {
		newSourceID = input.SourceID
	}

	// 套用新的餘額變動
	if newSourceType != nil && newSourceID != nil {
		err := s.validateAndUpdateBalance(original.Type, *newSourceType, *newSourceID, newAmount)
		if err != nil {
			// 如果新的餘額變動失敗，需要回復原本的餘額變動
			if original.SourceType != nil && original.SourceID != nil {
				s.validateAndUpdateBalance(original.Type, *original.SourceType, *original.SourceID, original.Amount)
			}
			return err
		}
	}

	return nil
}

// revertPaymentMethodChange 回復付款方式變更（用於更新失敗時的回滾）
func (s *cashFlowService) revertPaymentMethodChange(original *models.CashFlow, input *models.UpdateCashFlowInput) {
	// 檢查是否有付款方式相關的變更
	sourceTypeChanged := input.SourceType != nil && (original.SourceType == nil || *input.SourceType != *original.SourceType)
	sourceIDChanged := input.SourceID != nil && (original.SourceID == nil || *input.SourceID != *original.SourceID)
	amountChanged := input.Amount != nil && *input.Amount != original.Amount

	// 如果沒有相關變更，直接返回
	if !sourceTypeChanged && !sourceIDChanged && !amountChanged {
		return
	}

	// 計算新的金額
	newAmount := original.Amount
	if input.Amount != nil {
		newAmount = *input.Amount
	}

	// 計算新的付款方式
	newSourceType := original.SourceType
	if input.SourceType != nil {
		newSourceType = input.SourceType
	}

	newSourceID := original.SourceID
	if input.SourceID != nil {
		newSourceID = input.SourceID
	}

	// 回復新的餘額變動
	if newSourceType != nil && newSourceID != nil {
		s.revertBalanceUpdate(original.Type, *newSourceType, *newSourceID, newAmount)
	}

	// 重新套用原本的餘額變動
	if original.SourceType != nil && original.SourceID != nil {
		s.validateAndUpdateBalance(original.Type, *original.SourceType, *original.SourceID, original.Amount)
	}
}

// GetMonthlySummaryWithComparison 取得月度摘要（包含與前一個月的比較）
func (s *cashFlowService) GetMonthlySummaryWithComparison(year, month int) (*models.MonthlyCashFlowSummary, error) {
	// 取得當月摘要
	summary, err := s.repo.GetMonthlySummary(year, month)
	if err != nil {
		return nil, fmt.Errorf("failed to get monthly summary: %w", err)
	}

	// 計算前一個月的年份和月份
	prevYear := year
	prevMonth := month - 1
	if prevMonth < 1 {
		prevMonth = 12
		prevYear = year - 1
	}

	// 取得前一個月摘要
	prevSummary, err := s.repo.GetMonthlySummary(prevYear, prevMonth)
	if err != nil {
		// 如果前一個月沒有資料，不算錯誤，只是沒有比較資料
		return summary, nil
	}

	// 計算比較資料
	comparison := &models.MonthComparison{
		PreviousMonth: prevMonth,
		PreviousYear:  prevYear,
	}

	// 計算收入變化
	comparison.IncomeChange = summary.TotalIncome - prevSummary.TotalIncome
	if prevSummary.TotalIncome > 0 {
		comparison.IncomeChangePct = (comparison.IncomeChange / prevSummary.TotalIncome) * 100
	}

	// 計算支出變化
	comparison.ExpenseChange = summary.TotalExpense - prevSummary.TotalExpense
	if prevSummary.TotalExpense > 0 {
		comparison.ExpenseChangePct = (comparison.ExpenseChange / prevSummary.TotalExpense) * 100
	}

	// 計算淨現金流變化
	comparison.NetCashFlowChange = summary.NetCashFlow - prevSummary.NetCashFlow

	summary.ComparisonToPrev = comparison

	return summary, nil
}

// GetYearlySummaryWithComparison 取得年度摘要（包含與前一年的比較）
func (s *cashFlowService) GetYearlySummaryWithComparison(year int) (*models.YearlyCashFlowSummary, error) {
	// 取得當年摘要
	summary, err := s.repo.GetYearlySummary(year)
	if err != nil {
		return nil, fmt.Errorf("failed to get yearly summary: %w", err)
	}

	// 取得前一年摘要
	prevYear := year - 1
	prevSummary, err := s.repo.GetYearlySummary(prevYear)
	if err != nil {
		// 如果前一年沒有資料，不算錯誤，只是沒有比較資料
		return summary, nil
	}

	// 計算比較資料
	comparison := &models.YearComparison{
		PreviousYear: prevYear,
	}

	// 計算收入變化
	comparison.IncomeChange = summary.TotalIncome - prevSummary.TotalIncome
	if prevSummary.TotalIncome > 0 {
		comparison.IncomeChangePct = (comparison.IncomeChange / prevSummary.TotalIncome) * 100
	}

	// 計算支出變化
	comparison.ExpenseChange = summary.TotalExpense - prevSummary.TotalExpense
	if prevSummary.TotalExpense > 0 {
		comparison.ExpenseChangePct = (comparison.ExpenseChange / prevSummary.TotalExpense) * 100
	}

	// 計算淨現金流變化
	comparison.NetCashFlowChange = summary.NetCashFlow - prevSummary.NetCashFlow

	summary.ComparisonToPrev = comparison

	return summary, nil
}

