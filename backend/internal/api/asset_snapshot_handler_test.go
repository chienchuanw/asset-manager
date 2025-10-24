package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAssetSnapshotService 模擬 AssetSnapshotService
type MockAssetSnapshotService struct {
	mock.Mock
}

func (m *MockAssetSnapshotService) CreateSnapshot(input *models.CreateAssetSnapshotInput) (*models.AssetSnapshot, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetSnapshotByDate(date time.Time, assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	args := m.Called(date, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetSnapshotsByDateRange(startDate, endDate time.Time, assetType models.SnapshotAssetType) ([]*models.AssetSnapshot, error) {
	args := m.Called(startDate, endDate, assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) GetLatestSnapshot(assetType models.SnapshotAssetType) (*models.AssetSnapshot, error) {
	args := m.Called(assetType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) UpdateSnapshot(date time.Time, assetType models.SnapshotAssetType, valueTWD float64) (*models.AssetSnapshot, error) {
	args := m.Called(date, assetType, valueTWD)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.AssetSnapshot), args.Error(1)
}

func (m *MockAssetSnapshotService) DeleteSnapshot(date time.Time, assetType models.SnapshotAssetType) error {
	args := m.Called(date, assetType)
	return args.Error(0)
}

// setupAssetSnapshotHandlerTest 設定測試環境
func setupAssetSnapshotHandlerTest() (*gin.Engine, *MockAssetSnapshotService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(MockAssetSnapshotService)
	handler := NewAssetSnapshotHandler(mockService)

	// 註冊路由
	api := router.Group("/api")
	{
		snapshots := api.Group("/snapshots")
		{
			snapshots.POST("", handler.CreateSnapshot)
			snapshots.GET("/trend", handler.GetAssetTrend)
		}
	}

	return router, mockService
}

// TestAssetSnapshotHandler_CreateSnapshot 測試建立資產快照
func TestAssetSnapshotHandler_CreateSnapshot(t *testing.T) {
	router, mockService := setupAssetSnapshotHandlerTest()

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		mockSetup      func()
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name: "成功建立快照",
			requestBody: map[string]interface{}{
				"snapshot_date": "2024-01-15",
				"asset_type":    "total",
				"value_twd":     1000000.50,
			},
			mockSetup: func() {
				mockService.On("CreateSnapshot", mock.AnythingOfType("*models.CreateAssetSnapshotInput")).
					Return(&models.AssetSnapshot{
						SnapshotDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
						AssetType:    models.SnapshotAssetTypeTotal,
						ValueTWD:     1000000.50,
					}, nil)
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].(map[string]interface{})
				assert.Equal(t, "total", data["asset_type"])
				assert.Equal(t, 1000000.50, data["value_twd"])
			},
		},
		{
			name: "失敗 - 缺少必要欄位",
			requestBody: map[string]interface{}{
				"snapshot_date": "2024-01-15",
			},
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.NotNil(t, resp["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			tt.mockSetup()

			// 建立請求
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/snapshots", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// 執行請求
			router.ServeHTTP(w, req)

			// 驗證狀態碼
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 驗證回應
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			tt.checkResponse(t, response)
		})
	}
}

// TestAssetSnapshotHandler_GetAssetTrend 測試取得資產趨勢
func TestAssetSnapshotHandler_GetAssetTrend(t *testing.T) {
	router, mockService := setupAssetSnapshotHandlerTest()

	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func()
		expectedStatus int
		checkResponse  func(*testing.T, map[string]interface{})
	}{
		{
			name:        "成功取得 30 天趨勢",
			queryParams: "?days=30&asset_type=total",
			mockSetup: func() {
				snapshots := []*models.AssetSnapshot{
					{
						SnapshotDate: time.Now().Add(-1 * 24 * time.Hour),
						AssetType:    models.SnapshotAssetTypeTotal,
						ValueTWD:     1000000.00,
					},
					{
						SnapshotDate: time.Now().Add(-2 * 24 * time.Hour),
						AssetType:    models.SnapshotAssetTypeTotal,
						ValueTWD:     990000.00,
					},
				}
				mockService.On("GetSnapshotsByDateRange", mock.Anything, mock.Anything, models.SnapshotAssetTypeTotal).
					Return(snapshots, nil)
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				data := resp["data"].([]interface{})
				assert.Len(t, data, 2)
			},
		},
		{
			name:           "失敗 - 缺少必要參數",
			queryParams:    "",
			mockSetup:      func() {},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.NotNil(t, resp["error"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 重置 mock
			mockService.ExpectedCalls = nil
			mockService.Calls = nil

			tt.mockSetup()

			// 建立請求
			req := httptest.NewRequest(http.MethodGet, "/api/snapshots/trend"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			// 執行請求
			router.ServeHTTP(w, req)

			// 驗證狀態碼
			assert.Equal(t, tt.expectedStatus, w.Code)

			// 驗證回應
			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			tt.checkResponse(t, response)
		})
	}
}

