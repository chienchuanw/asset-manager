package service

import (
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// InstallmentService 分期業務邏輯介面
type InstallmentService interface {
	CreateInstallment(input *models.CreateInstallmentInput) (*models.Installment, error)
	GetInstallment(id uuid.UUID) (*models.Installment, error)
	ListInstallments(filters repository.InstallmentFilters) ([]*models.Installment, error)
	UpdateInstallment(id uuid.UUID, input *models.UpdateInstallmentInput) (*models.Installment, error)
	DeleteInstallment(id uuid.UUID) error
	GetDueBillings(date time.Time) ([]*models.Installment, error)
	GetCompletingSoon(remainingCount int) ([]*models.Installment, error)
}

// installmentService 分期業務邏輯實作
type installmentService struct {
	repo         repository.InstallmentRepository
	categoryRepo repository.CategoryRepository
}

// NewInstallmentService 建立新的分期 service
func NewInstallmentService(
	repo repository.InstallmentRepository,
	categoryRepo repository.CategoryRepository,
) InstallmentService {
	return &installmentService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

// CreateInstallment 建立新的分期
func (s *installmentService) CreateInstallment(input *models.CreateInstallmentInput) (*models.Installment, error) {
	// 驗證分期名稱
	if input.Name == "" {
		return nil, fmt.Errorf("installment name is required")
	}

	if len(input.Name) > 255 {
		return nil, fmt.Errorf("installment name must not exceed 255 characters")
	}

	// 驗證總金額
	if input.TotalAmount <= 0 {
		return nil, fmt.Errorf("total amount must be greater than zero")
	}

	// 驗證期數
	if input.InstallmentCount <= 0 {
		return nil, fmt.Errorf("installment count must be greater than zero")
	}

	// 驗證利率
	if input.InterestRate < 0 {
		return nil, fmt.Errorf("interest rate cannot be negative")
	}

	// 驗證扣款日
	if input.BillingDay < 1 || input.BillingDay > 31 {
		return nil, fmt.Errorf("billing day must be between 1 and 31")
	}

	// 驗證分類是否存在且為支出類型
	category, err := s.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("invalid category: %w", err)
	}

	if category.Type != models.CashFlowTypeExpense {
		return nil, fmt.Errorf("installment category must be expense type")
	}

	// 呼叫 repository 建立分期
	installment, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create installment: %w", err)
	}

	return installment, nil
}

// GetInstallment 取得單筆分期
func (s *installmentService) GetInstallment(id uuid.UUID) (*models.Installment, error) {
	return s.repo.GetByID(id)
}

// ListInstallments 取得分期列表
func (s *installmentService) ListInstallments(filters repository.InstallmentFilters) ([]*models.Installment, error) {
	// 驗證篩選條件
	if filters.Status != nil && !filters.Status.Validate() {
		return nil, fmt.Errorf("invalid installment status filter: %s", *filters.Status)
	}

	return s.repo.List(filters)
}

// UpdateInstallment 更新分期
func (s *installmentService) UpdateInstallment(id uuid.UUID, input *models.UpdateInstallmentInput) (*models.Installment, error) {
	// 驗證分期是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("installment not found: %w", err)
	}

	// 驗證分期名稱
	if input.Name != nil {
		if *input.Name == "" {
			return nil, fmt.Errorf("installment name cannot be empty")
		}
		if len(*input.Name) > 255 {
			return nil, fmt.Errorf("installment name must not exceed 255 characters")
		}
	}

	// 驗證分類
	if input.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(*input.CategoryID)
		if err != nil {
			return nil, fmt.Errorf("invalid category: %w", err)
		}
		if category.Type != models.CashFlowTypeExpense {
			return nil, fmt.Errorf("installment category must be expense type")
		}
	}

	// 呼叫 repository 更新分期
	installment, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update installment: %w", err)
	}

	return installment, nil
}

// DeleteInstallment 刪除分期
func (s *installmentService) DeleteInstallment(id uuid.UUID) error {
	// 驗證分期是否存在
	_, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("installment not found: %w", err)
	}

	// 呼叫 repository 刪除分期
	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete installment: %w", err)
	}

	return nil
}

// GetDueBillings 取得指定日期需要扣款的分期
func (s *installmentService) GetDueBillings(date time.Time) ([]*models.Installment, error) {
	return s.repo.GetDueBillings(date)
}

// GetCompletingSoon 取得即將完成的分期
func (s *installmentService) GetCompletingSoon(remainingCount int) ([]*models.Installment, error) {
	if remainingCount <= 0 {
		return nil, fmt.Errorf("remaining count must be greater than zero")
	}

	return s.repo.GetCompletingSoon(remainingCount)
}

