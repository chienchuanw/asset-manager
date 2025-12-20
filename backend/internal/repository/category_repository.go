package repository

import (
	"database/sql"
	"fmt"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
)

// CategoryRepository 現金流分類資料存取介面
type CategoryRepository interface {
	Create(input *models.CreateCategoryInput) (*models.CashFlowCategory, error)
	GetByID(id uuid.UUID) (*models.CashFlowCategory, error)
	GetAll(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error)
	Update(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error)
	Delete(id uuid.UUID) error
	IsInUse(id uuid.UUID) (bool, error)
	Reorder(input *models.ReorderCategoryInput) error
	GetMaxSortOrder(flowType models.CashFlowType) (int, error)
}

// categoryRepository 現金流分類資料存取實作
type categoryRepository struct {
	db *sql.DB
}

// NewCategoryRepository 建立新的分類 repository
func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// Create 建立新的分類
func (r *categoryRepository) Create(input *models.CreateCategoryInput) (*models.CashFlowCategory, error) {
	query := `
		INSERT INTO cash_flow_categories (name, type, is_system)
		VALUES ($1, $2, false)
		RETURNING id, name, type, is_system, created_at, updated_at
	`

	category := &models.CashFlowCategory{}
	err := r.db.QueryRow(
		query,
		input.Name,
		input.Type,
	).Scan(
		&category.ID,
		&category.Name,
		&category.Type,
		&category.IsSystem,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// GetByID 根據 ID 取得分類
func (r *categoryRepository) GetByID(id uuid.UUID) (*models.CashFlowCategory, error) {
	query := `
		SELECT id, name, type, is_system, created_at, updated_at
		FROM cash_flow_categories
		WHERE id = $1
	`

	category := &models.CashFlowCategory{}
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Type,
		&category.IsSystem,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

// GetAll 取得所有分類（可選擇性篩選類型）
func (r *categoryRepository) GetAll(flowType *models.CashFlowType) ([]*models.CashFlowCategory, error) {
	query := `
		SELECT id, name, type, is_system, created_at, updated_at
		FROM cash_flow_categories
		WHERE 1=1
	`

	args := []interface{}{}

	// 如果有指定類型，加入篩選條件
	if flowType != nil {
		query += " AND type = $1"
		args = append(args, *flowType)
	}

	// 排序：系統分類優先，然後按名稱排序
	query += " ORDER BY is_system DESC, name ASC"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	categories := []*models.CashFlowCategory{}
	for rows.Next() {
		category := &models.CashFlowCategory{}
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Type,
			&category.IsSystem,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// Update 更新分類名稱
func (r *categoryRepository) Update(id uuid.UUID, input *models.UpdateCategoryInput) (*models.CashFlowCategory, error) {
	query := `
		UPDATE cash_flow_categories
		SET name = $1
		WHERE id = $2 AND is_system = false
		RETURNING id, name, type, is_system, created_at, updated_at
	`

	category := &models.CashFlowCategory{}
	err := r.db.QueryRow(query, input.Name, id).Scan(
		&category.ID,
		&category.Name,
		&category.Type,
		&category.IsSystem,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("category not found or is a system category")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// Delete 刪除分類（僅限非系統分類）
func (r *categoryRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM cash_flow_categories WHERE id = $1 AND is_system = false`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found or is a system category")
	}

	return nil
}

// IsInUse 檢查分類是否被現金流記錄使用
func (r *categoryRepository) IsInUse(id uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM cash_flows WHERE category_id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if category is in use: %w", err)
	}

	return exists, nil
}
