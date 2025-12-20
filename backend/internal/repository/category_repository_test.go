package repository

import (
	"database/sql"
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ensureSystemCategories 確保系統預設分類存在
func ensureSystemCategories(db *sql.DB) {
	// 收入分類
	incomeCategories := []string{"薪資", "獎金", "利息", "其他收入"}
	for _, name := range incomeCategories {
		_, _ = db.Exec(`
			INSERT INTO cash_flow_categories (name, type, is_system)
			VALUES ($1, $2, true)
			ON CONFLICT (name, type) DO NOTHING
		`, name, models.CashFlowTypeIncome)
	}

	// 支出分類
	expenseCategories := []string{"飲食", "交通", "娛樂", "醫療", "房租", "水電", "保險", "其他支出"}
	for _, name := range expenseCategories {
		_, _ = db.Exec(`
			INSERT INTO cash_flow_categories (name, type, is_system)
			VALUES ($1, $2, true)
			ON CONFLICT (name, type) DO NOTHING
		`, name, models.CashFlowTypeExpense)
	}

	// 轉帳分類
	_, _ = db.Exec(`
		INSERT INTO cash_flow_categories (name, type, is_system)
		VALUES ($1, $2, true)
		ON CONFLICT (name, type) DO NOTHING
	`, "移轉", models.CashFlowTypeTransferIn)

	_, _ = db.Exec(`
		INSERT INTO cash_flow_categories (name, type, is_system)
		VALUES ($1, $2, true)
		ON CONFLICT (name, type) DO NOTHING
	`, "移轉", models.CashFlowTypeTransferOut)

	// 提領分類
	_, _ = db.Exec(`
		INSERT INTO cash_flow_categories (name, type, is_system)
		VALUES ($1, $2, true)
		ON CONFLICT (name, type) DO NOTHING
	`, "提領", models.CashFlowTypeTransferOut)
}

// TestCategoryRepository_Create 測試建立分類
func TestCategoryRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 準備測試資料
	input := &models.CreateCategoryInput{
		Name: "測試分類",
		Type: models.CashFlowTypeIncome,
	}

	// 執行測試
	category, err := repo.Create(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.NotEqual(t, uuid.Nil, category.ID)
	assert.Equal(t, input.Name, category.Name)
	assert.Equal(t, input.Type, category.Type)
	assert.False(t, category.IsSystem, "Custom category should not be system category")
	assert.NotZero(t, category.CreatedAt)
	assert.NotZero(t, category.UpdatedAt)
}

// TestCategoryRepository_GetByID 測試根據 ID 取得分類
func TestCategoryRepository_GetByID(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 建立測試資料
	created, err := repo.Create(&models.CreateCategoryInput{
		Name: "測試分類",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 執行測試
	category, err := repo.GetByID(created.ID)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, category)
	assert.Equal(t, created.ID, category.ID)
	assert.Equal(t, created.Name, category.Name)
	assert.Equal(t, created.Type, category.Type)
}

// TestCategoryRepository_GetByID_NotFound 測試取得不存在的分類
func TestCategoryRepository_GetByID_NotFound(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCategoryRepository(db)

	// 使用不存在的 ID
	nonExistentID := uuid.New()

	// 執行測試
	category, err := repo.GetByID(nonExistentID)

	// 驗證結果
	assert.Error(t, err)
	assert.Nil(t, category)
	assert.Contains(t, err.Error(), "not found")
}

// TestCategoryRepository_GetAll 測試取得所有分類
func TestCategoryRepository_GetAll(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 確保系統預設分類存在
	ensureSystemCategories(db)

	// 建立測試資料
	_, err = repo.Create(&models.CreateCategoryInput{
		Name: "自訂收入分類",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCategoryInput{
		Name: "自訂支出分類",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 測試：取得所有分類
	t.Run("Get all categories", func(t *testing.T) {
		categories, err := repo.GetAll(nil)
		assert.NoError(t, err)
		// 應該包含系統預設分類 + 自訂分類
		assert.GreaterOrEqual(t, len(categories), 14) // 12 個系統分類 + 2 個自訂分類
	})

	// 測試：只取得收入分類
	t.Run("Get income categories only", func(t *testing.T) {
		incomeType := models.CashFlowTypeIncome
		categories, err := repo.GetAll(&incomeType)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 5) // 4 個系統收入分類 + 1 個自訂
		for _, cat := range categories {
			assert.Equal(t, models.CashFlowTypeIncome, cat.Type)
		}
	})

	// 測試：只取得支出分類
	t.Run("Get expense categories only", func(t *testing.T) {
		expenseType := models.CashFlowTypeExpense
		categories, err := repo.GetAll(&expenseType)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(categories), 9) // 8 個系統支出分類 + 1 個自訂
		for _, cat := range categories {
			assert.Equal(t, models.CashFlowTypeExpense, cat.Type)
		}
	})
}

// TestCategoryRepository_Update 測試更新分類
func TestCategoryRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 建立測試資料
	created, err := repo.Create(&models.CreateCategoryInput{
		Name: "原始名稱",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	// 準備更新資料
	updateInput := &models.UpdateCategoryInput{
		Name: "更新後的名稱",
	}

	// 執行測試
	updated, err := repo.Update(created.ID, updateInput)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, updateInput.Name, updated.Name)
}

// TestCategoryRepository_Update_SystemCategory 測試更新系統分類（應該失敗）
func TestCategoryRepository_Update_SystemCategory(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCategoryRepository(db)

	// 取得系統分類
	systemCategoryID, err := getTestCategory(db, models.CashFlowTypeIncome)
	require.NoError(t, err)

	// 嘗試更新系統分類
	updateInput := &models.UpdateCategoryInput{
		Name: "嘗試更新系統分類",
	}

	updated, err := repo.Update(systemCategoryID, updateInput)

	// 驗證結果：應該失敗
	assert.Error(t, err)
	assert.Nil(t, updated)
	assert.Contains(t, err.Error(), "system category")
}

// TestCategoryRepository_Delete 測試刪除分類
func TestCategoryRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 建立測試資料
	created, err := repo.Create(&models.CreateCategoryInput{
		Name: "待刪除分類",
		Type: models.CashFlowTypeExpense,
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

// TestCategoryRepository_Delete_SystemCategory 測試刪除系統分類（應該失敗）
func TestCategoryRepository_Delete_SystemCategory(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	repo := NewCategoryRepository(db)

	// 取得系統分類
	systemCategoryID, err := getTestCategory(db, models.CashFlowTypeIncome)
	require.NoError(t, err)

	// 嘗試刪除系統分類
	err = repo.Delete(systemCategoryID)

	// 驗證結果：應該失敗
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "system category")
}

// TestCategoryRepository_Reorder 測試批次更新分類排序
func TestCategoryRepository_Reorder(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 建立測試分類
	cat1, err := repo.Create(&models.CreateCategoryInput{
		Name: "分類A",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	cat2, err := repo.Create(&models.CreateCategoryInput{
		Name: "分類B",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	cat3, err := repo.Create(&models.CreateCategoryInput{
		Name: "分類C",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 執行重新排序：C -> A -> B
	reorderInput := &models.ReorderCategoryInput{
		Orders: []models.CategoryOrderItem{
			{ID: cat3.ID, SortOrder: 0},
			{ID: cat1.ID, SortOrder: 1},
			{ID: cat2.ID, SortOrder: 2},
		},
	}

	err = repo.Reorder(reorderInput)
	assert.NoError(t, err)

	// 驗證排序結果
	expenseType := models.CashFlowTypeExpense
	categories, err := repo.GetAll(&expenseType)
	require.NoError(t, err)

	// 過濾出自訂分類（非系統分類）
	var customCategories []*models.CashFlowCategory
	for _, c := range categories {
		if !c.IsSystem {
			customCategories = append(customCategories, c)
		}
	}

	require.Len(t, customCategories, 3)
	assert.Equal(t, cat3.ID, customCategories[0].ID, "第一個應該是分類C")
	assert.Equal(t, cat1.ID, customCategories[1].ID, "第二個應該是分類A")
	assert.Equal(t, cat2.ID, customCategories[2].ID, "第三個應該是分類B")
}

// TestCategoryRepository_GetAll_SortBySortOrder 測試 GetAll 按 sort_order 排序
func TestCategoryRepository_GetAll_SortBySortOrder(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))
	ensureSystemCategories(db)

	repo := NewCategoryRepository(db)

	// 建立測試分類並設定不同的 sort_order
	cat1, err := repo.Create(&models.CreateCategoryInput{
		Name: "第三個",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	cat2, err := repo.Create(&models.CreateCategoryInput{
		Name: "第一個",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	cat3, err := repo.Create(&models.CreateCategoryInput{
		Name: "第二個",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	// 重新排序
	err = repo.Reorder(&models.ReorderCategoryInput{
		Orders: []models.CategoryOrderItem{
			{ID: cat2.ID, SortOrder: 100}, // 第一個
			{ID: cat3.ID, SortOrder: 101}, // 第二個
			{ID: cat1.ID, SortOrder: 102}, // 第三個
		},
	})
	require.NoError(t, err)

	// 取得所有收入分類
	incomeType := models.CashFlowTypeIncome
	categories, err := repo.GetAll(&incomeType)
	require.NoError(t, err)

	// 過濾自訂分類
	var customCategories []*models.CashFlowCategory
	for _, c := range categories {
		if !c.IsSystem {
			customCategories = append(customCategories, c)
		}
	}

	require.Len(t, customCategories, 3)
	assert.Equal(t, "第一個", customCategories[0].Name)
	assert.Equal(t, "第二個", customCategories[1].Name)
	assert.Equal(t, "第三個", customCategories[2].Name)
}

// TestCategoryRepository_GetMaxSortOrder 測試取得最大 sort_order
func TestCategoryRepository_GetMaxSortOrder(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 空資料庫應該回傳 -1
	maxOrder, err := repo.GetMaxSortOrder(models.CashFlowTypeIncome)
	assert.NoError(t, err)
	assert.Equal(t, -1, maxOrder)

	// 建立分類
	_, err = repo.Create(&models.CreateCategoryInput{
		Name: "測試分類",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	// 應該回傳 0（新分類的預設 sort_order）
	maxOrder, err = repo.GetMaxSortOrder(models.CashFlowTypeIncome)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, maxOrder, 0)
}

// TestCategoryRepository_Create_WithSortOrder 測試新建分類自動設定 sort_order
func TestCategoryRepository_Create_WithSortOrder(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	require.NoError(t, cleanupCategories(db))

	repo := NewCategoryRepository(db)

	// 建立第一個分類
	cat1, err := repo.Create(&models.CreateCategoryInput{
		Name: "第一個分類",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)
	firstSortOrder := cat1.SortOrder

	// 建立第二個分類，sort_order 應該比第一個大
	cat2, err := repo.Create(&models.CreateCategoryInput{
		Name: "第二個分類",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	assert.Greater(t, cat2.SortOrder, firstSortOrder, "第二個分類的 sort_order 應該比第一個大")
}

