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

// cleanupSubscriptions 清理測試資料庫中的訂閱記錄
func cleanupSubscriptions(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM subscriptions")
	return err
}

// TestSubscriptionRepository_Create 測試建立訂閱
func TestSubscriptionRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupSubscriptions(db))

	repo := NewSubscriptionRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 準備測試資料
	note := "Netflix 訂閱"
	input := &models.CreateSubscriptionInput{
		Name:          "Netflix",
		Amount:        390,
		BillingCycle:  models.BillingCycleMonthly,
		BillingDay:    15,
		CategoryID:    categoryID,
		PaymentMethod: models.PaymentMethodCash,
		StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		AutoRenew:     true,
		Note:          &note,
	}

	// 執行測試
	subscription, err := repo.Create(input)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.NotEqual(t, uuid.Nil, subscription.ID)
	assert.Equal(t, input.Name, subscription.Name)
	assert.Equal(t, input.Amount, subscription.Amount)
	assert.Equal(t, models.CurrencyTWD, subscription.Currency)
	assert.Equal(t, input.BillingCycle, subscription.BillingCycle)
	assert.Equal(t, input.BillingDay, subscription.BillingDay)
	assert.Equal(t, input.CategoryID, subscription.CategoryID)
	assert.Equal(t, input.StartDate.Format("2006-01-02"), subscription.StartDate.Format("2006-01-02"))
	assert.Nil(t, subscription.EndDate)
	assert.Equal(t, input.AutoRenew, subscription.AutoRenew)
	assert.Equal(t, models.SubscriptionStatusActive, subscription.Status)
	assert.Equal(t, *input.Note, *subscription.Note)
	assert.NotZero(t, subscription.CreatedAt)
	assert.NotZero(t, subscription.UpdatedAt)
}

// TestSubscriptionRepository_GetByID 測試根據 ID 取得訂閱
func TestSubscriptionRepository_GetByID(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupSubscriptions(db))

	repo := NewSubscriptionRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試訂閱
	note := "Spotify 訂閱"
	input := &models.CreateSubscriptionInput{
		Name:          "Spotify",
		Amount:        149,
		BillingCycle:  models.BillingCycleMonthly,
		BillingDay:    1,
		CategoryID:    categoryID,
		PaymentMethod: models.PaymentMethodCash,
		StartDate:     time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		AutoRenew:     true,
		Note:          &note,
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	// 執行測試
	subscription, err := repo.GetByID(created.ID)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, subscription)
	assert.Equal(t, created.ID, subscription.ID)
	assert.Equal(t, created.Name, subscription.Name)
	assert.NotNil(t, subscription.Category)
	assert.Equal(t, categoryID, subscription.Category.ID)
}

// TestSubscriptionRepository_List 測試取得訂閱列表
func TestSubscriptionRepository_List(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupSubscriptions(db))

	repo := NewSubscriptionRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立多個測試訂閱
	subscriptions := []struct {
		name   string
		amount float64
		status models.SubscriptionStatus
	}{
		{"Netflix", 390, models.SubscriptionStatusActive},
		{"Spotify", 149, models.SubscriptionStatusActive},
		{"YouTube Premium", 179, models.SubscriptionStatusCancelled},
	}

	for _, sub := range subscriptions {
		input := &models.CreateSubscriptionInput{
			Name:          sub.name,
			Amount:        sub.amount,
			BillingCycle:  models.BillingCycleMonthly,
			BillingDay:    15,
			CategoryID:    categoryID,
			PaymentMethod: models.PaymentMethodCash,
			StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			AutoRenew:     true,
		}
		created, err := repo.Create(input)
		require.NoError(t, err)

		// 如果是已取消狀態，需要更新
		if sub.status == models.SubscriptionStatusCancelled {
			_, err = db.Exec("UPDATE subscriptions SET status = $1 WHERE id = $2", sub.status, created.ID)
			require.NoError(t, err)
		}
	}

	// 測試：取得所有訂閱
	t.Run("get all subscriptions", func(t *testing.T) {
		filters := SubscriptionFilters{}
		result, err := repo.List(filters)

		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	// 測試：只取得進行中的訂閱
	t.Run("get active subscriptions only", func(t *testing.T) {
		status := models.SubscriptionStatusActive
		filters := SubscriptionFilters{
			Status: &status,
		}
		result, err := repo.List(filters)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		for _, sub := range result {
			assert.Equal(t, models.SubscriptionStatusActive, sub.Status)
		}
	})

	// 測試：分頁
	t.Run("pagination", func(t *testing.T) {
		filters := SubscriptionFilters{
			Limit:  2,
			Offset: 0,
		}
		result, err := repo.List(filters)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

// TestSubscriptionRepository_Update 測試更新訂閱
func TestSubscriptionRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupSubscriptions(db))

	repo := NewSubscriptionRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試訂閱
	input := &models.CreateSubscriptionInput{
		Name:          "Netflix",
		Amount:        390,
		BillingCycle:  models.BillingCycleMonthly,
		BillingDay:    15,
		CategoryID:    categoryID,
		PaymentMethod: models.PaymentMethodCash,
		StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		AutoRenew:     true,
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	// 準備更新資料
	newAmount := 490.0
	newNote := "價格調漲"
	updateInput := &models.UpdateSubscriptionInput{
		Amount: &newAmount,
		Note:   &newNote,
	}

	// 執行測試
	updated, err := repo.Update(created.ID, updateInput)

	// 驗證結果
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, created.ID, updated.ID)
	assert.Equal(t, newAmount, updated.Amount)
	assert.Equal(t, newNote, *updated.Note)
	assert.Equal(t, created.Name, updated.Name) // 未更新的欄位應保持不變
}

// TestSubscriptionRepository_Delete 測試刪除訂閱
func TestSubscriptionRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err, "Failed to setup test database")
	defer db.Close()

	// 清理測試資料
	require.NoError(t, cleanupSubscriptions(db))

	repo := NewSubscriptionRepository(db)

	// 取得測試用的分類 ID
	categoryID, err := getTestCategory(db, models.CashFlowTypeExpense)
	require.NoError(t, err, "Failed to get test category")

	// 建立測試訂閱
	input := &models.CreateSubscriptionInput{
		Name:          "Netflix",
		Amount:        390,
		BillingCycle:  models.BillingCycleMonthly,
		BillingDay:    15,
		CategoryID:    categoryID,
		PaymentMethod: models.PaymentMethodCash,
		StartDate:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
		AutoRenew:     true,
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

