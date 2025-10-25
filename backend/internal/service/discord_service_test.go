package service

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSendMessage_Success 測試發送訊息成功
func TestSendMessage_Success(t *testing.T) {
	// 建立 mock Discord server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 驗證請求方法
		assert.Equal(t, http.MethodPost, r.Method)

		// 驗證 Content-Type
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		// 解析請求 body
		var msg models.DiscordMessage
		err := json.NewDecoder(r.Body).Decode(&msg)
		require.NoError(t, err)

		// 驗證訊息內容
		assert.Equal(t, "Test message", msg.Content)

		// 回傳成功
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// 建立 Discord Service
	service := NewDiscordService()

	// 建立測試訊息
	message := &models.DiscordMessage{
		Content: "Test message",
	}

	// 發送訊息
	err := service.SendMessage(server.URL, message)

	// 驗證結果
	assert.NoError(t, err)
}

// TestSendMessage_InvalidWebhookURL 測試無效的 Webhook URL
func TestSendMessage_InvalidWebhookURL(t *testing.T) {
	service := NewDiscordService()

	message := &models.DiscordMessage{
		Content: "Test message",
	}

	// 使用無效的 URL
	err := service.SendMessage("invalid-url", message)

	// 應該回傳錯誤
	assert.Error(t, err)
}

// TestSendMessage_ServerError 測試伺服器錯誤
func TestSendMessage_ServerError(t *testing.T) {
	// 建立 mock Discord server（回傳錯誤）
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "Internal Server Error"}`))
	}))
	defer server.Close()

	service := NewDiscordService()

	message := &models.DiscordMessage{
		Content: "Test message",
	}

	// 發送訊息
	err := service.SendMessage(server.URL, message)

	// 應該回傳錯誤
	assert.Error(t, err)
}

// TestFormatDailyReport_Success 測試格式化每日報告
func TestFormatDailyReport_Success(t *testing.T) {
	service := NewDiscordService()

	// 建立測試資料
	reportData := &models.DailyReportData{
		Date:               time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		TotalMarketValue:   2524906.22,
		TotalCost:          2501130.52,
		TotalUnrealizedPL:  23775.70,
		TotalUnrealizedPct: 0.9506,
		HoldingCount:       28,
		TopHoldings: []*models.Holding{
			{
				Symbol:          "BTC",
				Name:            "Bitcoin",
				AssetType:       "crypto",
				Quantity:        0.5,
				AvgCost:         50000,
				CurrentPrice:    60000,
				MarketValue:     30000,
				UnrealizedPL:    5000,
				UnrealizedPLPct: 20.0,
			},
		},
		ByAssetType: map[string]*models.AssetTypePerformance{
			"tw-stock": {
				AssetType:     "tw-stock",
				MarketValue:   1000000,
				Cost:          1000245.81,
				UnrealizedPL:  -245.81,
				UnrealizedPct: -0.0098,
				HoldingCount:  16,
			},
		},
	}

	// 格式化報告
	message := service.FormatDailyReport(reportData)

	// 驗證結果
	require.NotNil(t, message)
	require.Len(t, message.Embeds, 1)

	embed := message.Embeds[0]
	assert.Equal(t, "📊 每日資產報告", embed.Title)
	assert.NotEmpty(t, embed.Description)
	assert.Greater(t, len(embed.Fields), 0)
	assert.NotNil(t, embed.Footer)
}

