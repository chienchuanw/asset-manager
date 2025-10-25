package service

import (
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
)

// CategoryService 現金流分類業務邏輯介面
type CategoryService interface {
	CreateCategory(input *models.CreateCategoryInput) (*models.CashFlowCategory, error)
	GetCategory(id uuid.UUID) (*models.CashFlowCategory, error)
	ListCategories(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error)
	UpdateCategory(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error)
	DeleteCategory(id uuid.UUID) error
}

// categoryService 現金流分類業務邏輯實作
type categoryService struct {
	repo repository.CategoryRepository
}

// NewCategoryService 建立新的分類 service
func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

// CreateCategory 建立新的分類
func (s *categoryService) CreateCategory(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	// 驗證現金流類型
	if !input.Type.Validate() {
		return nil, fmt.Errorf("invalid cash flow type: %s", input.Type)
	}

	// 驗證分類名稱
	if input.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	if len(input.Name) > 100 {
		return nil, fmt.Errorf("category name must not exceed 100 characters")
	}

	// 呼叫 repository 建立分類
	category, err := s.repo.Create(input)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetCategory 取得單筆分類
func (s *categoryService) GetCategory(id uuid.UUID) (*models.CashFlowCategory, error) {
	return s.repo.GetByID(id)
}

// ListCategories 取得分類列表
func (s *categoryService) ListCategories(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	// 驗證篩選條件
	if flowType != nil && !flowType.Validate() {
		return nil, fmt.Errorf("invalid cash flow type filter: %s", *flowType)
	}

	return s.repo.GetAll(flowType)
}

// UpdateCategory 更新分類
func (s *categoryService) UpdateCategory(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	// 驗證分類名稱
	if input.Name == "" {
		return nil, fmt.Errorf("category name is required")
	}

	if len(input.Name) > 100 {
		return nil, fmt.Errorf("category name must not exceed 100 characters")
	}

	// 呼叫 repository 更新分類
	category, err := s.repo.Update(id, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory 刪除分類
func (s *categoryService) DeleteCategory(id uuid.UUID) error {
	// 呼叫 repository 刪除分類
	// Repository 層會檢查是否為系統分類
	err := s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

