package service

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestSendDailyBillingNotification 測試發送每日扣款通知
func TestSendDailyBillingNotification(t *testing.T) {
	// 建立 mock Discord server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service := NewDiscordService()

	// 準備測試資料
	date := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)
	result := &DailyBillingResult{
		Date:              date,
		SubscriptionCount: 2,
		InstallmentCount:  1,
		TotalAmount:       1500.0,
		SubscriptionResult: &BillingResult{
			ProcessedCount: 2,
			FailedCount:    0,
			CreatedCashFlows: []*models.CashFlow{
				{
					ID:          uuid.New(),
					Description: "Netflix - 訂閱扣款",
					Amount:      390.0,
				},
				{
					ID:          uuid.New(),
					Description: "Spotify - 訂閱扣款",
					Amount:      149.0,
				},
			},
			Errors: []BillingError{},
		},
		InstallmentResult: &BillingResult{
			ProcessedCount: 1,
			FailedCount:    0,
			CreatedCashFlows: []*models.CashFlow{
				{
					ID:          uuid.New(),
					Description: "iPhone 15 - 分期付款 (3/12)",
					Amount:      961.0,
				},
			},
			Errors: []BillingError{},
		},
	}

	// 執行測試
	err := service.SendDailyBillingNotification(server.URL, result)

	// 驗證結果
	assert.NoError(t, err)
}

// TestSendDailyBillingNotification_NoData 測試沒有扣款資料時不發送通知
func TestSendDailyBillingNotification_NoData(t *testing.T) {
	service := NewDiscordService()

	// 準備測試資料（沒有扣款）
	date := time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC)
	result := &DailyBillingResult{
		Date:              date,
		SubscriptionCount: 0,
		InstallmentCount:  0,
		TotalAmount:       0,
		SubscriptionResult: &BillingResult{
			ProcessedCount:   0,
			FailedCount:      0,
			CreatedCashFlows: []*models.CashFlow{},
			Errors:           []BillingError{},
		},
		InstallmentResult: &BillingResult{
			ProcessedCount:   0,
			FailedCount:      0,
			CreatedCashFlows: []*models.CashFlow{},
			Errors:           []BillingError{},
		},
	}

	// 執行測試（不應該發送請求）
	err := service.SendDailyBillingNotification("http://example.com", result)

	// 驗證結果（應該沒有錯誤，因為直接返回）
	assert.NoError(t, err)
}

// TestSendSubscriptionExpiryNotification 測試發送訂閱到期通知
func TestSendSubscriptionExpiryNotification(t *testing.T) {
	// 建立 mock Discord server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service := NewDiscordService()

	// 準備測試資料
	endDate1 := time.Now().AddDate(0, 0, 5)
	endDate2 := time.Now().AddDate(0, 0, 3)
	subscriptions := []*models.Subscription{
		{
			ID:      uuid.New(),
			Name:    "Netflix",
			Amount:  390.0,
			EndDate: &endDate1,
		},
		{
			ID:      uuid.New(),
			Name:    "Spotify",
			Amount:  149.0,
			EndDate: &endDate2,
		},
	}

	// 執行測試
	err := service.SendSubscriptionExpiryNotification(server.URL, subscriptions, 7)

	// 驗證結果
	assert.NoError(t, err)
}

// TestSendSubscriptionExpiryNotification_NoData 測試沒有到期訂閱時不發送通知
func TestSendSubscriptionExpiryNotification_NoData(t *testing.T) {
	service := NewDiscordService()

	// 執行測試（空列表）
	err := service.SendSubscriptionExpiryNotification("http://example.com", []*models.Subscription{}, 7)

	// 驗證結果（應該沒有錯誤，因為直接返回）
	assert.NoError(t, err)
}

// TestSendInstallmentCompletionNotification 測試發送分期完成通知
func TestSendInstallmentCompletionNotification(t *testing.T) {
	// 建立 mock Discord server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	service := NewDiscordService()

	// 準備測試資料
	installments := []*models.Installment{
		{
			ID:                uuid.New(),
			Name:              "iPhone 15",
			TotalAmount:       30000.0,
			InstallmentCount:  12,
			InstallmentAmount: 2500.0,
			PaidCount:         10,
		},
		{
			ID:                uuid.New(),
			Name:              "MacBook Pro",
			TotalAmount:       60000.0,
			InstallmentCount:  24,
			InstallmentAmount: 2500.0,
			PaidCount:         23,
		},
	}

	// 執行測試
	err := service.SendInstallmentCompletionNotification(server.URL, installments, 3)

	// 驗證結果
	assert.NoError(t, err)
}

// TestSendInstallmentCompletionNotification_NoData 測試沒有即將完成的分期時不發送通知
func TestSendInstallmentCompletionNotification_NoData(t *testing.T) {
	service := NewDiscordService()

	// 執行測試（空列表）
	err := service.SendInstallmentCompletionNotification("http://example.com", []*models.Installment{}, 3)

	// 驗證結果（應該沒有錯誤，因為直接返回）
	assert.NoError(t, err)
}

