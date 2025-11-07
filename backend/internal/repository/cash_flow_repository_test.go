package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cleanupCashFlows 清理測試資料庫中的現金流記錄
func cleanupCashFlows(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM cash_flows")
	return err
}

// cleanupCategories 清理測試資料庫中的自訂分類（保留系統分類）
func cleanupCategories(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM cash_flow_categories WHERE is_system = false")
	return err
}

// getTestCategory 取得測試用的分類 ID（使用系統預設分類）
// 如果系統分類不存在，則建立一個
func getTestCategory(db *sql.DB, flowType models.CashFlowType) (uuid.UUID, error) {
	var categoryID uuid.UUID
	query := `SELECT id FROM cash_flow_categories WHERE type = $1 AND is_system = true LIMIT 1`
	err := db.QueryRow(query, flowType).Scan(&categoryID)

	// 如果找不到系統分類，建立一個
	if err == sql.ErrNoRows {
		// 根據類型決定分類名稱
		var categoryName string
		switch flowType {
		case models.CashFlowTypeIncome:
			categoryName = "測試收入"
		case models.CashFlowTypeExpense:
			categoryName = "測試支出"
		case models.CashFlowTypeTransferIn:
			categoryName = "測試轉入"
		case models.CashFlowTypeTransferOut:
			categoryName = "測試轉出"
		default:
			categoryName = "測試分類"
		}

		// 插入系統分類
		insertQuery := `
			INSERT INTO cash_flow_categories (name, type, is_system)
			VALUES ($1, $2, true)
			ON CONFLICT (name, type) DO UPDATE SET is_system = true
			RETURNING id
		`
		err = db.QueryRow(insertQuery, categoryName, flowType).Scan(&categoryID)
		if err != nil {
			return uuid.Nil, err
		}
		return categoryID, nil
	}

	return categoryID, err
}

// TestCashFlowRepository_Create 測試建立現金流記錄
func TestCashFlowRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeIncome)
	require.NoError(t, err, "Failed to get test category")

	// 準備測試資料
	note := "測試備註"
	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
		Note:        &note,
	}

	// 執行測試
	cashFlow, err := repo.Create(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, cashFlow)
	assert.NotEqual(t, uuid.Nil, cashFlow.ID)
	assert.Equal(t, input.Date.Format("2006-01-02"), cashFlow.Date.Format("2006-01-02"))
	assert.Equal(t, input.Type, cashFlow.Type)
	assert.Equal(t, input.CategoryID, cashFlow.CategoryID)
	assert.Equal(t, input.Amount, cashFlow.Amount)
	assert.Equal(t, models.CurrencyTWD, cashFlow.Currency)
	assert.Equal(t, input.Description, cashFlow.Description)
	assert.Equal(t, *input.Note, *cashFlow.Note)
	assert.NotZero(t, cashFlow.CreatedAt)
	assert.NotZero(t, cashFlow.UpdatedAt)
}

// TestCashFlowRepository_GetByID 測試根據 ID 取得現金流記錄
func TestCashFlowRepository_GetByID(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 建立測試資料
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err)

	input := &models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1200,
		Description: "午餐",
		Note:        nil,
	}

	created, err := repo.Create(input)
	require.NoError(t, err)

	// 執行測試
	cashFlow, err := repo.GetByID(created.ID)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, cashFlow)
	assert.Equal(t, created.ID, cashFlow.ID)
	assert.Equal(t, created.Amount, cashFlow.Amount)
	assert.NotNil(t, cashFlow.Category, "Category should be loaded")
	assert.Equal(t, categoryID, cashFlow.Category.ID)
}

// TestCashFlowRepository_GetByID_NotFound 測試取得不存在的記錄
func TestCashFlowRepository_GetByID_NotFound(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCashFlowRepository(db)

	// 使用不存在的 ID
	nonExistentID := uuid.New()

	// 執行測試
	cashFlow, err := repo.GetByID(nonExistentID)

	// 驗證結果
	assert.Error(t, err)
	assert.Nil(t, cashFlow)
	assert.Contains(t, err.Error(), "not found")
}

// TestCashFlowRepository_GetAll 測試取得所有現金流記錄
func TestCashFlowRepository_GetAll(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 建立多筆測試資料
	incomeCategoryID, err := getTestCategory(db, models.CashFlowTypeIncome)
	require.NoError(t, err)

	expenseCategoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err)

	// 建立收入記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  incomeCategoryID,
		Amount:      50000,
		Description: "薪資",
	})
	require.NoError(t, err)

	// 建立支出記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  expenseCategoryID,
		Amount:      1200,
		Description: "午餐",
	})
	require.NoError(t, err)

	// 測試：取得所有記錄
	t.Run("Get all without filters", func(t *testing.T) {
		cashFlows, err := repo.GetAll(CashFlowFilters{})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(cashFlows), 2)
	})

	// 測試：篩選收入
	t.Run("Filter by income type", func(t *testing.T) {
		incomeType := models.CashFlowTypeIncome
		cashFlows, err := repo.GetAll(CashFlowFilters{
			Type: &incomeType,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(cashFlows), 1)
		for _, cf := range cashFlows {
			assert.Equal(t, models.CashFlowTypeIncome, cf.Type)
		}
	})

	// 測試：篩選支出
	t.Run("Filter by expense type", func(t *testing.T) {
		expenseType := models.CashFlowTypeExpense
		cashFlows, err := repo.GetAll(CashFlowFilters{
			Type: &expenseType,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(cashFlows), 1)
		for _, cf := range cashFlows {
			assert.Equal(t, models.CashFlowTypeExpense, cf.Type)
		}
	})

	// 測試：日期範圍篩選
	t.Run("Filter by date range", func(t *testing.T) {
		startDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)
		cashFlows, err := repo.GetAll(CashFlowFilters{
			StartDate: &startDate,
			EndDate:   &endDate,
		})
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(cashFlows), 2)
	})

	// 測試：分頁
	t.Run("Pagination", func(t *testing.T) {
		cashFlows, err := repo.GetAll(CashFlowFilters{
			Limit:  1,
			Offset: 0,
		})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(cashFlows))
	})
}

// TestCashFlowRepository_Update 測試更新現金流記錄
func TestCashFlowRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 建立測試資料
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err)

	created, err := repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1000,
		Description: "原始描述",
	})
	require.NoError(t, err)

	// 準備更新資料
	newAmount := 1500.0
	newDescription := "更新後的描述"
	updateInput := &models.UpdateCashFlowInput{
		Amount:      &newAmount,
		Description: &newDescription,
	}

	// 執行測試
	updated, err := repo.Update(created.ID, updateInput)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, newAmount, updated.Amount)
	assert.Equal(t, newDescription, updated.Description)
}

// TestCashFlowRepository_Delete 測試刪除現金流記錄
func TestCashFlowRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 建立測試資料
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err)

	created, err := repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  categoryID,
		Amount:      1000,
		Description: "測試刪除",
	})
	require.NoError(t, err)

	// 執行刪除
	err = repo.Delete(created.ID)
	assert.NoError(t, err)

	// 驗證已刪除
	deleted, err := repo.GetByID(created.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}

// TestCashFlowRepository_GetSummary 測試取得現金流摘要
func TestCashFlowRepository_GetSummary(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCashFlows(db))

	repo := NewCashFlowRepository(db)

	// 建立測試資料
	incomeCategoryID, err := getTestCategory(db, models.CashFlowTypeIncome)
	require.NoError(t, err)

	expenseCategoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err)

	// 建立收入記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 5, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  incomeCategoryID,
		Amount:      50000,
		Description: "薪資",
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 10, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  incomeCategoryID,
		Amount:      5000,
		Description: "獎金",
	})
	require.NoError(t, err)

	// 建立支出記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 15, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  expenseCategoryID,
		Amount:      10000,
		Description: "房租",
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  expenseCategoryID,
		Amount:      5000,
		Description: "生活費",
	})
	require.NoError(t, err)

	// 執行測試
	startDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)
	summary, err := repo.GetSummary(startDate, endDate)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.Equal(t, 55000.0, summary.TotalIncome)
	assert.Equal(t, 15000.0, summary.TotalExpense)
	assert.Equal(t, 40000.0, summary.NetCashFlow)
}

