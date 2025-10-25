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

// cleanupInstallments 清理測試資料庫中的分期記錄
func cleanupInstallments(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM installments")
	return err
}

// TestInstallmentRepository_Create 測試建立分期
func TestInstallmentRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 準備測試資料
	note := "iPhone 15 Pro 分期"
	input := &models.CreateInstallmentInput{
		Name:             "iPhone 15 Pro",
		TotalAmount:      36000,
		InstallmentCount: 12,
		InterestRate:     0,
		BillingDay:       15,
		CategoryID:       categoryID,
		StartDate:        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		Note:             &note,
	}

	// 執行測試
	installment, err := repo.Create(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, installment)
	assert.NotEqual(t, uuid.Nil, installment.ID)
	assert.Equal(t, input.Name, installment.Name)
	assert.Equal(t, input.TotalAmount, installment.TotalAmount)
	assert.Equal(t, models.CurrencyTWD, installment.Currency)
	assert.Equal(t, input.InstallmentCount, installment.InstallmentCount)
	assert.Equal(t, 3000.0, installment.InstallmentAmount) // 36000 / 12
	assert.Equal(t, input.InterestRate, installment.InterestRate)
	assert.Equal(t, 0.0, installment.TotalInterest)
	assert.Equal(t, 0, installment.PaidCount)
	assert.Equal(t, input.BillingDay, installment.BillingDay)
	assert.Equal(t, input.CategoryID, installment.CategoryID)
	assert.Equal(t, input.StartDate.Format("2006-01-02"), installment.StartDate.Format("2006-01-02"))
	assert.Equal(t, models.InstallmentStatusActive, installment.Status)
	assert.Equal(t, *input.Note, *installment.Note)
	assert.NotZero(t, installment.CreatedAt)
	assert.NotZero(t, installment.UpdatedAt)
}

// TestInstallmentRepository_Create_WithInterest 測試建立有利息的分期
func TestInstallmentRepository_Create_WithInterest(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 準備測試資料（10% 利率）
	input := &models.CreateInstallmentInput{
		Name:             "MacBook Pro",
		TotalAmount:      60000,
		InstallmentCount: 12,
		InterestRate:     10,
		BillingDay:       1,
		CategoryID:       categoryID,
		StartDate:        time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	// 執行測試
	installment, err := repo.Create(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, installment)
	assert.Equal(t, 6000.0, installment.TotalInterest)      // 60000 * 0.1
	assert.Equal(t, 5500.0, installment.InstallmentAmount)  // (60000 + 6000) / 12
}

// TestInstallmentRepository_GetByID 測試根據 ID 取得分期
func TestInstallmentRepository_GetByID(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試分期
	input := &models.CreateInstallmentInput{
		Name:             "iPad Pro",
		TotalAmount:      30000,
		InstallmentCount: 6,
		InterestRate:     0,
		BillingDay:       10,
		CategoryID:       categoryID,
		StartDate:        time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	// 執行測試
	installment, err := repo.GetByID(created.ID)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, installment)
	assert.Equal(t, created.ID, installment.ID)
	assert.Equal(t, created.Name, installment.Name)
	assert.NotNil(t, installment.Category)
	assert.Equal(t, categoryID, installment.Category.ID)
}

// TestInstallmentRepository_List 測試取得分期列表
func TestInstallmentRepository_List(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立多個測試分期
	installments := []struct {
		name   string
		amount float64
		status models.InstallmentStatus
	}{
		{"iPhone", 36000, models.InstallmentStatusActive},
		{"MacBook", 60000, models.InstallmentStatusActive},
		{"iPad", 30000, models.InstallmentStatusCompleted},
	}

	for _, inst := range installments {
		input := &models.CreateInstallmentInput{
			Name:             inst.name,
			TotalAmount:      inst.amount,
			InstallmentCount: 12,
			InterestRate:     0,
			BillingDay:       15,
			CategoryID:       categoryID,
			StartDate:        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		}
		created, err := repo.Create(input)
		require.NoError(t, err)

		// 如果是已完成狀態，需要更新
		if inst.status == models.InstallmentStatusCompleted {
			_, err = db.Exec("UPDATE installments SET status = $1, paid_count = installment_count WHERE id = $2", inst.status, created.ID)
			require.NoError(t, err)
		}
	}

	// 測試：取得所有分期
	t.Run("get all installments", func(t *testing.T) {
		filters := InstallmentFilters{}
		result, err := repo.List(filters)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	// 測試：只取得進行中的分期
	t.Run("get active installments only", func(t *testing.T) {
		status := models.InstallmentStatusActive
		filters := InstallmentFilters{
			Status: &status,
		}
		result, err := repo.List(filters)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		for _, inst := range result {
			assert.Equal(t, models.InstallmentStatusActive, inst.Status)
		}
	})
}

// TestInstallmentRepository_Update 測試更新分期
func TestInstallmentRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試分期
	input := &models.CreateInstallmentInput{
		Name:             "iPhone",
		TotalAmount:      36000,
		InstallmentCount: 12,
		InterestRate:     0,
		BillingDay:       15,
		CategoryID:       categoryID,
		StartDate:        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	// 準備更新資料
	newNote := "已付 3 期"
	updateInput := &models.UpdateInstallmentInput{
		Note: &newNote,
	}

	// 執行測試
	updated, err := repo.Update(created.ID, updateInput)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, newNote, *updated.Note)
	assert.Equal(t, created.Name, updated.Name) // 未更新的欄位應保持不變
}

// TestInstallmentRepository_Delete 測試刪除分期
func TestInstallmentRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupInstallments(db))

	repo := NewInstallmentRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試分期
	input := &models.CreateInstallmentInput{
		Name:             "iPhone",
		TotalAmount:      36000,
		InstallmentCount: 12,
		InterestRate:     0,
		BillingDay:       15,
		CategoryID:       categoryID,
		StartDate:        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	// 執行測試
	err = repo.Delete(created.ID)

	// 驗證結果
	assert.NoError(t, err)

	// 確認已刪除
	_, err = repo.GetByID(created.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

