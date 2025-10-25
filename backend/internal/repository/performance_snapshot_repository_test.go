package repository

import (
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPerformanceSnapshotRepository_Create(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM daily_performance_snapshots")
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewPerformanceSnapshotRepository(dbx)

	input := &models.CreateDailyPerformanceSnapshotInput{
		SnapshotDate:       time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		TotalMarketValue:   1000000,
		TotalCost:          900000,
		TotalUnrealizedPL:  100000,
		TotalUnrealizedPct: 11.11,
		TotalRealizedPL:    50000,
		TotalRealizedPct:   5.56,
		HoldingCount:       10,
		Currency:           "TWD",
		Details: []models.CreateDailyPerformanceSnapshotDetailInput{
			{
				AssetType:     models.AssetTypeTWStock,
				MarketValue:   500000,
				Cost:          450000,
				UnrealizedPL:  50000,
				UnrealizedPct: 11.11,
				RealizedPL:    25000,
				RealizedPct:   5.56,
				HoldingCount:  5,
			},
			{
				AssetType:     models.AssetTypeUSStock,
				MarketValue:   500000,
				Cost:          450000,
				UnrealizedPL:  50000,
				UnrealizedPct: 11.11,
				RealizedPL:    25000,
				RealizedPct:   5.56,
				HoldingCount:  5,
			},
		},
	}

	snapshot, err := repo.Create(input)
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.Equal(t, input.SnapshotDate.Format("2006-01-02"), snapshot.SnapshotDate.Format("2006-01-02"))
	assert.Equal(t, input.TotalMarketValue, snapshot.TotalMarketValue)
	assert.Equal(t, input.TotalCost, snapshot.TotalCost)
	assert.Equal(t, input.HoldingCount, snapshot.HoldingCount)
}

func TestPerformanceSnapshotRepository_GetByDate(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM daily_performance_snapshots")
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewPerformanceSnapshotRepository(dbx)

	// 建立測試資料
	input := &models.CreateDailyPerformanceSnapshotInput{
		SnapshotDate:       time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		TotalMarketValue:   1000000,
		TotalCost:          900000,
		TotalUnrealizedPL:  100000,
		TotalUnrealizedPct: 11.11,
		TotalRealizedPL:    50000,
		TotalRealizedPct:   5.56,
		HoldingCount:       10,
		Currency:           "TWD",
		Details:            []models.CreateDailyPerformanceSnapshotDetailInput{},
	}
	_, err = repo.Create(input)
	require.NoError(t, err)

	// 測試取得快照
	snapshot, err := repo.GetByDate(input.SnapshotDate)
	require.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.Equal(t, input.SnapshotDate.Format("2006-01-02"), snapshot.SnapshotDate.Format("2006-01-02"))
}

func TestPerformanceSnapshotRepository_GetByDateRange(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM daily_performance_snapshots")
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewPerformanceSnapshotRepository(dbx)

	// 建立多筆測試資料
	dates := []time.Time{
		time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		input := &models.CreateDailyPerformanceSnapshotInput{
			SnapshotDate:       date,
			TotalMarketValue:   1000000,
			TotalCost:          900000,
			TotalUnrealizedPL:  100000,
			TotalUnrealizedPct: 11.11,
			TotalRealizedPL:    50000,
			TotalRealizedPct:   5.56,
			HoldingCount:       10,
			Currency:           "TWD",
			Details:            []models.CreateDailyPerformanceSnapshotDetailInput{},
		}
		_, err = repo.Create(input)
		require.NoError(t, err)
	}

	// 測試取得日期範圍內的快照
	snapshots, err := repo.GetByDateRange(dates[0], dates[2])
	require.NoError(t, err)
	assert.Len(t, snapshots, 3)
}

func TestPerformanceSnapshotRepository_GetLatest(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM daily_performance_snapshots")
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewPerformanceSnapshotRepository(dbx)

	// 建立多筆測試資料
	dates := []time.Time{
		time.Date(2025, 10, 23, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
	}

	for _, date := range dates {
		input := &models.CreateDailyPerformanceSnapshotInput{
			SnapshotDate:       date,
			TotalMarketValue:   1000000,
			TotalCost:          900000,
			TotalUnrealizedPL:  100000,
			TotalUnrealizedPct: 11.11,
			TotalRealizedPL:    50000,
			TotalRealizedPct:   5.56,
			HoldingCount:       10,
			Currency:           "TWD",
			Details:            []models.CreateDailyPerformanceSnapshotDetailInput{},
		}
		_, err = repo.Create(input)
		require.NoError(t, err)
	}

	// 測試取得最新快照
	snapshots, err := repo.GetLatest(2)
	require.NoError(t, err)
	assert.Len(t, snapshots, 2)
	// 應該按日期降序排列
	assert.Equal(t, dates[2].Format("2006-01-02"), snapshots[0].SnapshotDate.Format("2006-01-02"))
	assert.Equal(t, dates[1].Format("2006-01-02"), snapshots[1].SnapshotDate.Format("2006-01-02"))
}

func TestPerformanceSnapshotRepository_GetDetailsBySnapshotID(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM daily_performance_snapshots")
	require.NoError(t, err)

	dbx := sqlx.NewDb(db, "postgres")
	repo := NewPerformanceSnapshotRepository(dbx)

	// 建立測試資料
	input := &models.CreateDailyPerformanceSnapshotInput{
		SnapshotDate:       time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		TotalMarketValue:   1000000,
		TotalCost:          900000,
		TotalUnrealizedPL:  100000,
		TotalUnrealizedPct: 11.11,
		TotalRealizedPL:    50000,
		TotalRealizedPct:   5.56,
		HoldingCount:       10,
		Currency:           "TWD",
		Details: []models.CreateDailyPerformanceSnapshotDetailInput{
			{
				AssetType:     models.AssetTypeTWStock,
				MarketValue:   500000,
				Cost:          450000,
				UnrealizedPL:  50000,
				UnrealizedPct: 11.11,
				RealizedPL:    25000,
				RealizedPct:   5.56,
				HoldingCount:  5,
			},
			{
				AssetType:     models.AssetTypeUSStock,
				MarketValue:   500000,
				Cost:          450000,
				UnrealizedPL:  50000,
				UnrealizedPct: 11.11,
				RealizedPL:    25000,
				RealizedPct:   5.56,
				HoldingCount:  5,
			},
		},
	}

	snapshot, err := repo.Create(input)
	require.NoError(t, err)

	// 測試取得明細
	details, err := repo.GetDetailsBySnapshotID(snapshot.ID)
	require.NoError(t, err)
	assert.Len(t, details, 2)
}

