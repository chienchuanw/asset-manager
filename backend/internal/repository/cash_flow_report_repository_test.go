package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// cleanupCashFlowData 清理現金流相關測試資料
func cleanupCashFlowData(db *sql.DB) error {
	// 按照外鍵依賴順序清理
	tables := []string{
		"cash_flows",
		"cash_flow_categories",
	}

	for _, table := range tables {
		_, err := db.Exec("DELETE FROM " + table)
		if err != nil {
			return err
		}
	}

	return nil
}

// TestGetMonthlySummary 測試取得月度摘要
func TestGetMonthlySummary(t *testing.T) {
	// 設定測試資料庫
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupCashFlowData(db)
	require.NoError(t, err)

	repo := NewCashFlowRepository(db)
	categoryRepo := NewCategoryRepository(db)

	// 建立測試分類
	incomeCategory, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	expenseCategory, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 建立測試現金流記錄（2024年10月）
	testDate := time.Date(2024, 10, 15, 0, 0, 0, 0, time.Local)

	// 收入記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        testDate,
		Type:        models.CashFlowTypeIncome,
		CategoryID:  incomeCategory.ID,
		Amount:      50000,
		Description: "十月薪資",
	})
	require.NoError(t, err)

	// 支出記錄
	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        testDate,
		Type:        models.CashFlowTypeExpense,
		CategoryID:  expenseCategory.ID,
		Amount:      10000,
		Description: "餐飲支出",
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        testDate.AddDate(0, 0, 1),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  expenseCategory.ID,
		Amount:      5000,
		Description: "餐飲支出2",
	})
	require.NoError(t, err)

	// 測試取得月度摘要
	summary, err := repo.GetMonthlySummary(2024, 10)
	require.NoError(t, err)
	assert.NotNil(t, summary)

	// 驗證基本統計
	assert.Equal(t, 2024, summary.Year)
	assert.Equal(t, 10, summary.Month)
	assert.Equal(t, 50000.0, summary.TotalIncome)
	assert.Equal(t, 15000.0, summary.TotalExpense)
	assert.Equal(t, 35000.0, summary.NetCashFlow)
	assert.Equal(t, 1, summary.IncomeCount)
	assert.Equal(t, 2, summary.ExpenseCount)

	// 驗證分類摘要
	assert.Len(t, summary.IncomeCategoryBreakdown, 1)
	assert.Equal(t, "薪資", summary.IncomeCategoryBreakdown[0].CategoryName)
	assert.Equal(t, 50000.0, summary.IncomeCategoryBreakdown[0].Amount)

	assert.Len(t, summary.ExpenseCategoryBreakdown, 1)
	assert.Equal(t, "餐飲", summary.ExpenseCategoryBreakdown[0].CategoryName)
	assert.Equal(t, 15000.0, summary.ExpenseCategoryBreakdown[0].Amount)

	// 驗證前 10 大支出
	assert.Len(t, summary.TopExpenses, 2)
	assert.Equal(t, 10000.0, summary.TopExpenses[0].Amount)
	assert.Equal(t, 5000.0, summary.TopExpenses[1].Amount)
}

// TestGetYearlySummary 測試取得年度摘要
func TestGetYearlySummary(t *testing.T) {
	// 設定測試資料庫
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupCashFlowData(db)
	require.NoError(t, err)

	repo := NewCashFlowRepository(db)
	categoryRepo := NewCategoryRepository(db)

	// 建立測試分類
	incomeCategory, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "薪資",
		Type: models.CashFlowTypeIncome,
	})
	require.NoError(t, err)

	expenseCategory, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 建立多個月份的測試資料
	for month := 1; month <= 12; month++ {
		testDate := time.Date(2024, time.Month(month), 15, 0, 0, 0, 0, time.Local)

		// 每月收入
		_, err = repo.Create(&models.CreateCashFlowInput{
			Date:        testDate,
			Type:        models.CashFlowTypeIncome,
			CategoryID:  incomeCategory.ID,
			Amount:      50000,
			Description: "月薪",
		})
		require.NoError(t, err)

		// 每月支出
		_, err = repo.Create(&models.CreateCashFlowInput{
			Date:        testDate,
			Type:        models.CashFlowTypeExpense,
			CategoryID:  expenseCategory.ID,
			Amount:      30000,
			Description: "月支出",
		})
		require.NoError(t, err)
	}

	// 測試取得年度摘要
	summary, err := repo.GetYearlySummary(2024)
	require.NoError(t, err)
	assert.NotNil(t, summary)

	// 驗證基本統計
	assert.Equal(t, 2024, summary.Year)
	assert.Equal(t, 600000.0, summary.TotalIncome)  // 50000 * 12
	assert.Equal(t, 360000.0, summary.TotalExpense) // 30000 * 12
	assert.Equal(t, 240000.0, summary.NetCashFlow)
	assert.Equal(t, 12, summary.IncomeCount)
	assert.Equal(t, 12, summary.ExpenseCount)

	// 驗證分類摘要
	assert.Len(t, summary.IncomeCategoryBreakdown, 1)
	assert.Equal(t, "薪資", summary.IncomeCategoryBreakdown[0].CategoryName)
	assert.Equal(t, 600000.0, summary.IncomeCategoryBreakdown[0].Amount)

	assert.Len(t, summary.ExpenseCategoryBreakdown, 1)
	assert.Equal(t, "餐飲", summary.ExpenseCategoryBreakdown[0].CategoryName)
	assert.Equal(t, 360000.0, summary.ExpenseCategoryBreakdown[0].Amount)

	// 驗證月度細分
	assert.Len(t, summary.MonthlyBreakdown, 12)
	for i, breakdown := range summary.MonthlyBreakdown {
		assert.Equal(t, i+1, breakdown.Month)
		assert.Equal(t, 50000.0, breakdown.Income)
		assert.Equal(t, 30000.0, breakdown.Expense)
		assert.Equal(t, 20000.0, breakdown.NetCashFlow)
	}
}

// TestGetCategorySummary 測試取得分類摘要
func TestGetCategorySummary(t *testing.T) {
	// 設定測試資料庫
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupCashFlowData(db)
	require.NoError(t, err)

	repo := NewCashFlowRepository(db)
	categoryRepo := NewCategoryRepository(db)

	// 建立測試分類
	category1, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "餐飲",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	category2, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "交通",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 建立測試現金流記錄
	startDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2024, 10, 31, 23, 59, 59, 0, time.Local)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2024, 10, 15, 0, 0, 0, 0, time.Local),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  category1.ID,
		Amount:      10000,
		Description: "餐飲1",
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2024, 10, 16, 0, 0, 0, 0, time.Local),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  category1.ID,
		Amount:      5000,
		Description: "餐飲2",
	})
	require.NoError(t, err)

	_, err = repo.Create(&models.CreateCashFlowInput{
		Date:        time.Date(2024, 10, 17, 0, 0, 0, 0, time.Local),
		Type:        models.CashFlowTypeExpense,
		CategoryID:  category2.ID,
		Amount:      3000,
		Description: "交通",
	})
	require.NoError(t, err)

	// 測試取得分類摘要
	summaries, err := repo.GetCategorySummary(startDate, endDate, models.CashFlowTypeExpense)
	require.NoError(t, err)
	assert.Len(t, summaries, 2)

	// 驗證排序（金額由大到小）
	assert.Equal(t, "餐飲", summaries[0].CategoryName)
	assert.Equal(t, 15000.0, summaries[0].Amount)
	assert.Equal(t, 2, summaries[0].Count)

	assert.Equal(t, "交通", summaries[1].CategoryName)
	assert.Equal(t, 3000.0, summaries[1].Amount)
	assert.Equal(t, 1, summaries[1].Count)
}

// TestGetTopExpenses 測試取得前 N 大支出
func TestGetTopExpenses(t *testing.T) {
	// 設定測試資料庫
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupCashFlowData(db)
	require.NoError(t, err)

	repo := NewCashFlowRepository(db)
	categoryRepo := NewCategoryRepository(db)

	// 建立測試分類
	category, err := categoryRepo.Create(&models.CreateCategoryInput{
		Name: "購物",
		Type: models.CashFlowTypeExpense,
	})
	require.NoError(t, err)

	// 建立測試現金流記錄
	startDate := time.Date(2024, 10, 1, 0, 0, 0, 0, time.Local)
	endDate := time.Date(2024, 10, 31, 23, 59, 59, 0, time.Local)

	amounts := []float64{10000, 5000, 3000, 2000, 1000}
	for i, amount := range amounts {
		_, err = repo.Create(&models.CreateCashFlowInput{
			Date:        time.Date(2024, 10, i+1, 0, 0, 0, 0, time.Local),
			Type:        models.CashFlowTypeExpense,
			CategoryID:  category.ID,
			Amount:      amount,
			Description: "購物",
		})
		require.NoError(t, err)
	}

	// 測試取得前 3 大支出
	topExpenses, err := repo.GetTopExpenses(startDate, endDate, 3)
	require.NoError(t, err)
	assert.Len(t, topExpenses, 3)

	// 驗證排序（金額由大到小）
	assert.Equal(t, 10000.0, topExpenses[0].Amount)
	assert.Equal(t, 5000.0, topExpenses[1].Amount)
	assert.Equal(t, 3000.0, topExpenses[2].Amount)
}

