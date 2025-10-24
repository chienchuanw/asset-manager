package repository

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssetSnapshotRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupAssetSnapshots(db)
	require.NoError(t, err)

	repo := NewAssetSnapshotRepository(db)

	t.Run("成功建立資產快照", func(t *testing.T) {
		input := &models.CreateAssetSnapshotInput{
			SnapshotDate: time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
			AssetType:    models.SnapshotAssetTypeTWStock,
			ValueTWD:     1000000.50,
		}

		snapshot, err := repo.Create(input)

		require.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.NotEqual(t, "", snapshot.ID.String())
		assert.Equal(t, input.SnapshotDate.Format("2006-01-02"), snapshot.SnapshotDate.Format("2006-01-02"))
		assert.Equal(t, input.AssetType, snapshot.AssetType)
		assert.Equal(t, input.ValueTWD, snapshot.ValueTWD)
		assert.False(t, snapshot.CreatedAt.IsZero())
		assert.False(t, snapshot.UpdatedAt.IsZero())
	})

	t.Run("建立多個不同類型的快照", func(t *testing.T) {
		date := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)

		inputs := []*models.CreateAssetSnapshotInput{
			{SnapshotDate: date, AssetType: models.SnapshotAssetTypeTWStock, ValueTWD: 500000},
			{SnapshotDate: date, AssetType: models.SnapshotAssetTypeUSStock, ValueTWD: 300000},
			{SnapshotDate: date, AssetType: models.SnapshotAssetTypeCrypto, ValueTWD: 200000},
			{SnapshotDate: date, AssetType: models.SnapshotAssetTypeTotal, ValueTWD: 1000000},
		}

		for _, input := range inputs {
			snapshot, err := repo.Create(input)
			require.NoError(t, err)
			assert.Equal(t, input.AssetType, snapshot.AssetType)
			assert.Equal(t, input.ValueTWD, snapshot.ValueTWD)
		}
	})

	t.Run("相同日期和類型的快照應該失敗（UNIQUE 約束）", func(t *testing.T) {
		date := time.Date(2025, 10, 26, 0, 0, 0, 0, time.UTC)
		input := &models.CreateAssetSnapshotInput{
			SnapshotDate: date,
			AssetType:    models.SnapshotAssetTypeTWStock,
			ValueTWD:     100000,
		}

		// 第一次建立應該成功
		_, err := repo.Create(input)
		require.NoError(t, err)

		// 第二次建立應該失敗
		_, err = repo.Create(input)
		assert.Error(t, err)
	})
}

func TestAssetSnapshotRepository_GetByDateAndType(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupAssetSnapshots(db)
	require.NoError(t, err)

	repo := NewAssetSnapshotRepository(db)

	// 準備測試資料
	date := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	input := &models.CreateAssetSnapshotInput{
		SnapshotDate: date,
		AssetType:    models.SnapshotAssetTypeTWStock,
		ValueTWD:     1500000,
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	t.Run("成功取得快照", func(t *testing.T) {
		snapshot, err := repo.GetByDateAndType(date, models.SnapshotAssetTypeTWStock)

		require.NoError(t, err)
		assert.NotNil(t, snapshot)
		assert.Equal(t, created.ID, snapshot.ID)
		assert.Equal(t, created.ValueTWD, snapshot.ValueTWD)
	})

	t.Run("查詢不存在的快照", func(t *testing.T) {
		snapshot, err := repo.GetByDateAndType(date, models.SnapshotAssetTypeUSStock)

		assert.Error(t, err)
		assert.Nil(t, snapshot)
	})
}

func TestAssetSnapshotRepository_GetByDateRange(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupAssetSnapshots(db)
	require.NoError(t, err)

	repo := NewAssetSnapshotRepository(db)

	// 準備測試資料：建立 7 天的快照
	baseDate := time.Date(2025, 10, 20, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 7; i++ {
		date := baseDate.AddDate(0, 0, i)
		for _, assetType := range []models.SnapshotAssetType{
			models.SnapshotAssetTypeTWStock,
			models.SnapshotAssetTypeUSStock,
			models.SnapshotAssetTypeCrypto,
			models.SnapshotAssetTypeTotal,
		} {
			input := &models.CreateAssetSnapshotInput{
				SnapshotDate: date,
				AssetType:    assetType,
				ValueTWD:     float64((i + 1) * 100000),
			}
			_, err := repo.Create(input)
			require.NoError(t, err)
		}
	}

	t.Run("取得指定日期範圍的所有快照", func(t *testing.T) {
		startDate := baseDate.AddDate(0, 0, 2) // 10/22
		endDate := baseDate.AddDate(0, 0, 4)   // 10/24
		filters := models.AssetSnapshotFilters{
			StartDate: &startDate,
			EndDate:   &endDate,
		}

		snapshots, err := repo.GetByDateRange(filters)

		require.NoError(t, err)
		// 3 天 × 4 種類型 = 12 筆
		assert.Equal(t, 12, len(snapshots))
	})

	t.Run("取得指定日期範圍和資產類型的快照", func(t *testing.T) {
		startDate := baseDate
		endDate := baseDate.AddDate(0, 0, 6)
		assetType := models.SnapshotAssetTypeTWStock
		filters := models.AssetSnapshotFilters{
			StartDate: &startDate,
			EndDate:   &endDate,
			AssetType: &assetType,
		}

		snapshots, err := repo.GetByDateRange(filters)

		require.NoError(t, err)
		// 7 天 × 1 種類型 = 7 筆
		assert.Equal(t, 7, len(snapshots))
		for _, snapshot := range snapshots {
			assert.Equal(t, models.SnapshotAssetTypeTWStock, snapshot.AssetType)
		}
	})

	t.Run("只指定開始日期", func(t *testing.T) {
		startDate := baseDate.AddDate(0, 0, 5) // 10/25
		filters := models.AssetSnapshotFilters{
			StartDate: &startDate,
		}

		snapshots, err := repo.GetByDateRange(filters)

		require.NoError(t, err)
		// 2 天 × 4 種類型 = 8 筆
		assert.Equal(t, 8, len(snapshots))
	})

	t.Run("只指定結束日期", func(t *testing.T) {
		endDate := baseDate.AddDate(0, 0, 1) // 10/21
		filters := models.AssetSnapshotFilters{
			EndDate: &endDate,
		}

		snapshots, err := repo.GetByDateRange(filters)

		require.NoError(t, err)
		// 2 天 × 4 種類型 = 8 筆
		assert.Equal(t, 8, len(snapshots))
	})

	t.Run("沒有篩選條件時取得所有快照", func(t *testing.T) {
		filters := models.AssetSnapshotFilters{}

		snapshots, err := repo.GetByDateRange(filters)

		require.NoError(t, err)
		// 7 天 × 4 種類型 = 28 筆
		assert.Equal(t, 28, len(snapshots))
	})
}

func TestAssetSnapshotRepository_Update(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupAssetSnapshots(db)
	require.NoError(t, err)

	repo := NewAssetSnapshotRepository(db)

	// 準備測試資料
	date := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	input := &models.CreateAssetSnapshotInput{
		SnapshotDate: date,
		AssetType:    models.SnapshotAssetTypeTWStock,
		ValueTWD:     1000000,
	}
	created, err := repo.Create(input)
	require.NoError(t, err)

	t.Run("成功更新快照價值", func(t *testing.T) {
		newValue := 1500000.75

		updated, err := repo.Update(date, models.SnapshotAssetTypeTWStock, newValue)

		require.NoError(t, err)
		assert.NotNil(t, updated)
		assert.Equal(t, created.ID, updated.ID)
		assert.Equal(t, newValue, updated.ValueTWD)
		assert.True(t, updated.UpdatedAt.After(created.UpdatedAt))
	})

	t.Run("更新不存在的快照", func(t *testing.T) {
		_, err := repo.Update(date, models.SnapshotAssetTypeUSStock, 500000)

		assert.Error(t, err)
	})
}

func TestAssetSnapshotRepository_Delete(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	err = cleanupAssetSnapshots(db)
	require.NoError(t, err)

	repo := NewAssetSnapshotRepository(db)

	// 準備測試資料
	date := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
	input := &models.CreateAssetSnapshotInput{
		SnapshotDate: date,
		AssetType:    models.SnapshotAssetTypeTWStock,
		ValueTWD:     1000000,
	}
	_, err = repo.Create(input)
	require.NoError(t, err)

	t.Run("成功刪除快照", func(t *testing.T) {
		err := repo.Delete(date, models.SnapshotAssetTypeTWStock)

		require.NoError(t, err)

		// 驗證已刪除
		snapshot, err := repo.GetByDateAndType(date, models.SnapshotAssetTypeTWStock)
		assert.Error(t, err)
		assert.Nil(t, snapshot)
	})

	t.Run("刪除不存在的快照", func(t *testing.T) {
		err := repo.Delete(date, models.SnapshotAssetTypeUSStock)

		assert.Error(t, err)
	})
}

