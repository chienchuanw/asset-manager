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
}

// cashFlowService 現金流記錄業務邏輯實作
type cashFlowService struct {
	repo         repository.CashFlowRepository
	categoryRepo repository.CategoryRepository
}

// NewCashFlowService 建立新的現金流記錄 service
func NewCashFlowService(
	repo repository.CashFlowRepository,
	categoryRepo repository.CategoryRepository,
) CashFlowService {
	return &cashFlowService{
		repo:         repo,
		categoryRepo: categoryRepo,
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

	// 呼叫 repository 建立現金流記錄
	cashFlow, err := s.repo.Create(input)
	if err != nil {
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

	// 如果要更新分類，需要驗證分類是否存在
	if input.CategoryID != nil {
		// 先取得原始記錄以確認類型
		original, err := s.repo.GetByID(id)
		if err != nil {
			return nil, fmt.Errorf("cash flow not found: %w", err)
		}

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

	// 呼叫 repository 更新記錄
	cashFlow, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update cash flow: %w", err)
	}

	return cashFlow, nil
}

// DeleteCashFlow 刪除現金流記錄
func (s *cashFlowService) DeleteCashFlow(id uuid.UUID) error {
	return s.repo.Delete(id)
}

// GetSummary 取得指定日期區間的現金流摘要
func (s *cashFlowService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	// 驗證日期範圍
	if startDate.After(endDate) {
		return nil, fmt.Errorf("start date must be before or equal to end date")
	}

	return s.repo.GetSummary(startDate, endDate)
}

