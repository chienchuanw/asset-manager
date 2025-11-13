package service

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupAssetSnapshotServiceTest 設定測試環境
func setupAssetSnapshotServiceTest(t *testing.T) (*sql.DB, AssetSnapshotService, func()) {
	db, err := setupTestDB()
	require.NoError(t, err)

	// 清理測試資料
	_, err = db.Exec("DELETE FROM asset_snapshots")
	require.NoError(t, err)

	// 建立 repositories
	assetSnapshotRepo := repository.NewAssetSnapshotRepository(db)

	// 建立 service
	service := NewAssetSnapshotService(assetSnapshotRepo)

	cleanup := func() {
		db.Close()
	}

	return db, service, cleanup
}

// setupTestDB 設定測試資料庫連線
func setupTestDB() (*sql.DB, error) {
	// 使用 repository 包的私有函式需要透過反射或直接複製邏輯
	// 這裡我們直接使用 repository 的公開方法
	db, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=asset_manager_test sslmode=disable")
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

// TestAssetSnapshotService_CreateSnapshot 測試建立資產快照
func TestAssetSnapshotService_CreateSnapshot(t *testing.T) {
	_, service, cleanup := setupAssetSnapshotServiceTest(t)
	defer cleanup()

	tests := []struct {
		name    string
		input   *models.CreateAssetSnapshotInput
		wantErr bool
		errMsg  string
	}{
		{
			name: "成功建立總資產快照",
			input: &models.CreateAssetSnapshotInput{
				SnapshotDate: time.Now().Truncate(24 * time.Hour),
				AssetType:    models.SnapshotAssetTypeTotal,
				ValueTWD:     1000000.50,
			},
			wantErr: false,
		},
		{
			name: "成功建立台股快照",
			input: &models.CreateAssetSnapshotInput{
				SnapshotDate: time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour),
				AssetType:    models.SnapshotAssetTypeTWStock,
				ValueTWD:     500000.00,
			},
			wantErr: false,
		},
		{
			name: "失敗 - 缺少必要欄位",
			input: &models.CreateAssetSnapshotInput{
				SnapshotDate: time.Now().Truncate(24 * time.Hour).Add(48 * time.Hour),
				AssetType:    "",
				ValueTWD:     1000000.50,
			},
			wantErr: true,
			errMsg:  "asset_type is required",
		},
		{
			name: "失敗 - 負數金額",
			input: &models.CreateAssetSnapshotInput{
				SnapshotDate: time.Now().Truncate(24 * time.Hour).Add(72 * time.Hour),
				AssetType:    models.SnapshotAssetTypeTotal,
				ValueTWD:     -1000.00,
			},
			wantErr: true,
			errMsg:  "value_twd must be non-negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := service.CreateSnapshot(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, snapshot)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, snapshot)
				assert.NotZero(t, snapshot.ID)
				assert.Equal(t, tt.input.AssetType, snapshot.AssetType)
				assert.Equal(t, tt.input.ValueTWD, snapshot.ValueTWD)
			}
		})
	}
}

// TestAssetSnapshotService_GetSnapshotByDate 測試根據日期取得快照
func TestAssetSnapshotService_GetSnapshotByDate(t *testing.T) {
	db, service, cleanup := setupAssetSnapshotServiceTest(t)
	defer cleanup()

	// 準備測試資料
	today := time.Now().Truncate(24 * time.Hour)
	repo := repository.NewAssetSnapshotRepository(db)

	snapshot1, err := repo.Create(&models.CreateAssetSnapshotInput{
		SnapshotDate: today,
		AssetType:    models.SnapshotAssetTypeTotal,
		ValueTWD:     1000000.00,
	})
	require.NoError(t, err)

	tests := []struct {
		name      string
		date      time.Time
		assetType models.SnapshotAssetType
		wantErr   bool
		wantID    uuid.UUID
	}{
		{
			name:      "成功取得快照",
			date:      today,
			assetType: models.SnapshotAssetTypeTotal,
			wantErr:   false,
			wantID:    snapshot1.ID,
		},
		{
			name:      "找不到快照",
			date:      today.Add(-7 * 24 * time.Hour),
			assetType: models.SnapshotAssetTypeTotal,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshot, err := service.GetSnapshotByDate(tt.date, tt.assetType)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, snapshot)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, snapshot)
				assert.Equal(t, tt.wantID, snapshot.ID)
			}
		})
	}
}

// TestAssetSnapshotService_GetSnapshotsByDateRange 測試根據日期範圍取得快照
func TestAssetSnapshotService_GetSnapshotsByDateRange(t *testing.T) {
	db, service, cleanup := setupAssetSnapshotServiceTest(t)
	defer cleanup()

	// 準備測試資料 - 建立 7 天的快照
	today := time.Now().Truncate(24 * time.Hour)
	repo := repository.NewAssetSnapshotRepository(db)

	for i := 0; i < 7; i++ {
		date := today.Add(time.Duration(-i) * 24 * time.Hour)
		_, err := repo.Create(&models.CreateAssetSnapshotInput{
			SnapshotDate: date,
			AssetType:    models.SnapshotAssetTypeTotal,
			ValueTWD:     1000000.00 + float64(i*10000),
		})
		require.NoError(t, err)
	}

	tests := []struct {
		name      string
		startDate time.Time
		endDate   time.Time
		assetType models.SnapshotAssetType
		wantCount int
		wantErr   bool
	}{
		{
			name:      "取得最近 7 天的快照",
			startDate: today.Add(-6 * 24 * time.Hour),
			endDate:   today,
			assetType: models.SnapshotAssetTypeTotal,
			wantCount: 7,
			wantErr:   false,
		},
		{
			name:      "取得最近 3 天的快照",
			startDate: today.Add(-2 * 24 * time.Hour),
			endDate:   today,
			assetType: models.SnapshotAssetTypeTotal,
			wantCount: 3,
			wantErr:   false,
		},
		{
			name:      "日期範圍內沒有資料",
			startDate: today.Add(-30 * 24 * time.Hour),
			endDate:   today.Add(-20 * 24 * time.Hour),
			assetType: models.SnapshotAssetTypeTotal,
			wantCount: 0,
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshots, err := service.GetSnapshotsByDateRange(tt.startDate, tt.endDate, tt.assetType)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, snapshots, tt.wantCount)
			}
		})
	}
}

// TestAssetSnapshotService_GetLatestSnapshot 測試取得最新快照
func TestAssetSnapshotService_GetLatestSnapshot(t *testing.T) {
	db, service, cleanup := setupAssetSnapshotServiceTest(t)
	defer cleanup()

	// 準備測試資料
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)
	repo := repository.NewAssetSnapshotRepository(db)

	// 建立昨天的快照
	_, err := repo.Create(&models.CreateAssetSnapshotInput{
		SnapshotDate: yesterday,
		AssetType:    models.SnapshotAssetTypeTotal,
		ValueTWD:     900000.00,
	})
	require.NoError(t, err)

	// 建立今天的快照
	latestSnapshot, err := repo.Create(&models.CreateAssetSnapshotInput{
		SnapshotDate: today,
		AssetType:    models.SnapshotAssetTypeTotal,
		ValueTWD:     1000000.00,
	})
	require.NoError(t, err)

	// 測試取得最新快照
	snapshot, err := service.GetLatestSnapshot(models.SnapshotAssetTypeTotal)
	assert.NoError(t, err)
	assert.NotNil(t, snapshot)
	assert.Equal(t, latestSnapshot.ID, snapshot.ID)
	assert.Equal(t, today.Unix(), snapshot.SnapshotDate.Unix())
	assert.Equal(t, 1000000.00, snapshot.ValueTWD)
}

// TestAssetSnapshotService_CreateDailySnapshots_NoDoubleCurrencyConversion 測試建立每日快照時不會重複轉換台幣
// 這個測試驗證美股和加密貨幣的 MarketValue（已經是 TWD）不會被二次轉換
func TestAssetSnapshotService_CreateDailySnapshots_NoDoubleCurrencyConversion(t *testing.T) {
	db, err := setupTestDB()
	require.NoError(t, err)
	defer db.Close()

	// 清理測試資料
	_, err = db.Exec("DELETE FROM asset_snapshots")
	require.NoError(t, err)

	// 建立 mock HoldingService
	mockHoldingService := new(MockHoldingService)

	// 建立 repository 和 service
	assetSnapshotRepo := repository.NewAssetSnapshotRepository(db)
	service := NewAssetSnapshotServiceWithDeps(assetSnapshotRepo, mockHoldingService)

	// 準備測試資料：模擬三種資產類型的持倉
	// 重點：MarketValue 已經是 TWD，不應該再被轉換
	testHoldings := []*models.Holding{
		{
			Symbol:      "2330.TW",
			AssetType:   models.AssetTypeTWStock,
			Quantity:    100,
			TotalCost:   50000.0,  // TWD
			MarketValue: 60000.0,  // TWD (已轉換)
			Currency:    models.CurrencyTWD,
		},
		{
			Symbol:      "AAPL",
			AssetType:   models.AssetTypeUSStock,
			Quantity:    10,
			TotalCost:   50000.0,  // TWD (已轉換)
			MarketValue: 63000.0,  // TWD (已轉換，原本 USD 2000 * 31.5 = 63000)
			Currency:    models.CurrencyUSD, // 注意：這裡標記為 USD，但 MarketValue 已經是 TWD
		},
		{
			Symbol:      "BTC",
			AssetType:   models.AssetTypeCrypto,
			Quantity:    0.5,
			TotalCost:   500000.0, // TWD (已轉換)
			MarketValue: 630000.0, // TWD (已轉換，原本 USD 20000 * 31.5 = 630000)
			Currency:    models.CurrencyUSD, // 注意：這裡標記為 USD，但 MarketValue 已經是 TWD
		},
	}

	// 設定 mock 回傳值
	mockHoldingService.On("GetAllHoldings", models.HoldingFilters{}).Return(testHoldings, nil)

	// 執行建立每日快照
	err = service.CreateDailySnapshots()
	require.NoError(t, err)

	// 驗證快照是否正確建立
	today := time.Now().Truncate(24 * time.Hour)

	// 驗證總資產快照
	totalSnapshot, err := service.GetSnapshotByDate(today, models.SnapshotAssetTypeTotal)
	require.NoError(t, err)
	assert.NotNil(t, totalSnapshot)
	// 預期值：60000 + 63000 + 630000 = 753000 TWD
	// 如果有二次轉換的 bug，會是：60000 + (63000 * 31.5) + (630000 * 31.5) = 21,904,500 TWD
	assert.Equal(t, 753000.0, totalSnapshot.ValueTWD, "總資產應該是 753000 TWD，不應該有二次轉換")

	// 驗證台股快照
	twStockSnapshot, err := service.GetSnapshotByDate(today, models.SnapshotAssetTypeTWStock)
	require.NoError(t, err)
	assert.NotNil(t, twStockSnapshot)
	assert.Equal(t, 60000.0, twStockSnapshot.ValueTWD, "台股應該是 60000 TWD")

	// 驗證美股快照
	usStockSnapshot, err := service.GetSnapshotByDate(today, models.SnapshotAssetTypeUSStock)
	require.NoError(t, err)
	assert.NotNil(t, usStockSnapshot)
	// 預期值：63000 TWD（已轉換）
	// 如果有二次轉換的 bug，會是：63000 * 31.5 = 1,984,500 TWD
	assert.Equal(t, 63000.0, usStockSnapshot.ValueTWD, "美股應該是 63000 TWD，不應該有二次轉換")

	// 驗證加密貨幣快照
	cryptoSnapshot, err := service.GetSnapshotByDate(today, models.SnapshotAssetTypeCrypto)
	require.NoError(t, err)
	assert.NotNil(t, cryptoSnapshot)
	// 預期值：630000 TWD（已轉換）
	// 如果有二次轉換的 bug，會是：630000 * 31.5 = 19,845,000 TWD
	assert.Equal(t, 630000.0, cryptoSnapshot.ValueTWD, "加密貨幣應該是 630000 TWD，不應該有二次轉換")

	// 驗證 mock 被正確呼叫
	mockHoldingService.AssertExpectations(t)
}

