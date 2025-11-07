package service

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetMonthlySummaryWithComparison 測試取得月度摘要（有前一個月資料）
func TestGetMonthlySummaryWithComparison(t *testing.T) {
	mockRepo := new(MockCashFlowRepository)

	// 當月摘要
	currentSummary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        10,
		TotalIncome:  50000,
		TotalExpense: 15000,
		NetCashFlow:  35000,
	}

	// 前一個月摘要
	prevSummary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        9,
		TotalIncome:  45000,
		TotalExpense: 17000,
		NetCashFlow:  28000,
	}

	mockRepo.On("GetMonthlySummary", 2024, 10).Return(currentSummary, nil)
	mockRepo.On("GetMonthlySummary", 2024, 9).Return(prevSummary, nil)

	// 建立 service（需要完整的依賴）
	service := &cashFlowService{
		repo: mockRepo,
	}

	// 執行測試
	result, err := service.GetMonthlySummaryWithComparison(2024, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 驗證基本資料
	assert.Equal(t, 2024, result.Year)
	assert.Equal(t, 10, result.Month)
	assert.Equal(t, 50000.0, result.TotalIncome)
	assert.Equal(t, 15000.0, result.TotalExpense)
	assert.Equal(t, 35000.0, result.NetCashFlow)

	// 驗證比較資料
	assert.NotNil(t, result.ComparisonToPrev)
	assert.Equal(t, 9, result.ComparisonToPrev.PreviousMonth)
	assert.Equal(t, 2024, result.ComparisonToPrev.PreviousYear)
	assert.Equal(t, 5000.0, result.ComparisonToPrev.IncomeChange)
	assert.InDelta(t, 11.11, result.ComparisonToPrev.IncomeChangePct, 0.01)
	assert.Equal(t, -2000.0, result.ComparisonToPrev.ExpenseChange)
	assert.InDelta(t, -11.76, result.ComparisonToPrev.ExpenseChangePct, 0.01)
	assert.Equal(t, 7000.0, result.ComparisonToPrev.NetCashFlowChange)

	mockRepo.AssertExpectations(t)
}

// TestGetMonthlySummaryWithComparison_NoPreviousMonth 測試取得月度摘要（無前一個月資料）
func TestGetMonthlySummaryWithComparison_NoPreviousMonth(t *testing.T) {
	mockRepo := new(MockCashFlowRepository)

	// 當月摘要
	currentSummary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        10,
		TotalIncome:  50000,
		TotalExpense: 15000,
		NetCashFlow:  35000,
	}

	mockRepo.On("GetMonthlySummary", 2024, 10).Return(currentSummary, nil)
	mockRepo.On("GetMonthlySummary", 2024, 9).Return(nil, assert.AnError)

	service := &cashFlowService{
		repo: mockRepo,
	}

	// 執行測試
	result, err := service.GetMonthlySummaryWithComparison(2024, 10)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 驗證基本資料
	assert.Equal(t, 2024, result.Year)
	assert.Equal(t, 10, result.Month)

	// 驗證沒有比較資料
	assert.Nil(t, result.ComparisonToPrev)

	mockRepo.AssertExpectations(t)
}

// TestGetMonthlySummaryWithComparison_CrossYear 測試跨年情況（1月取得12月資料）
func TestGetMonthlySummaryWithComparison_CrossYear(t *testing.T) {
	mockRepo := new(MockCashFlowRepository)

	// 當月摘要（2024年1月）
	currentSummary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        1,
		TotalIncome:  50000,
		TotalExpense: 15000,
		NetCashFlow:  35000,
	}

	// 前一個月摘要（2023年12月）
	prevSummary := &models.MonthlyCashFlowSummary{
		Year:         2023,
		Month:        12,
		TotalIncome:  45000,
		TotalExpense: 17000,
		NetCashFlow:  28000,
	}

	mockRepo.On("GetMonthlySummary", 2024, 1).Return(currentSummary, nil)
	mockRepo.On("GetMonthlySummary", 2023, 12).Return(prevSummary, nil)

	service := &cashFlowService{
		repo: mockRepo,
	}

	// 執行測試
	result, err := service.GetMonthlySummaryWithComparison(2024, 1)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 驗證比較資料正確處理跨年
	assert.NotNil(t, result.ComparisonToPrev)
	assert.Equal(t, 12, result.ComparisonToPrev.PreviousMonth)
	assert.Equal(t, 2023, result.ComparisonToPrev.PreviousYear)

	mockRepo.AssertExpectations(t)
}

// TestGetYearlySummaryWithComparison 測試取得年度摘要
func TestGetYearlySummaryWithComparison(t *testing.T) {
	mockRepo := new(MockCashFlowRepository)

	// 當年摘要
	currentSummary := &models.YearlyCashFlowSummary{
		Year:         2024,
		TotalIncome:  600000,
		TotalExpense: 360000,
		NetCashFlow:  240000,
	}

	// 前一年摘要
	prevSummary := &models.YearlyCashFlowSummary{
		Year:         2023,
		TotalIncome:  550000,
		TotalExpense: 380000,
		NetCashFlow:  170000,
	}

	mockRepo.On("GetYearlySummary", 2024).Return(currentSummary, nil)
	mockRepo.On("GetYearlySummary", 2023).Return(prevSummary, nil)

	service := &cashFlowService{
		repo: mockRepo,
	}

	// 執行測試
	result, err := service.GetYearlySummaryWithComparison(2024)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// 驗證基本資料
	assert.Equal(t, 2024, result.Year)
	assert.Equal(t, 600000.0, result.TotalIncome)
	assert.Equal(t, 360000.0, result.TotalExpense)
	assert.Equal(t, 240000.0, result.NetCashFlow)

	// 驗證比較資料
	assert.NotNil(t, result.ComparisonToPrev)
	assert.Equal(t, 2023, result.ComparisonToPrev.PreviousYear)
	assert.Equal(t, 50000.0, result.ComparisonToPrev.IncomeChange)
	assert.InDelta(t, 9.09, result.ComparisonToPrev.IncomeChangePct, 0.01)
	assert.Equal(t, -20000.0, result.ComparisonToPrev.ExpenseChange)
	assert.InDelta(t, -5.26, result.ComparisonToPrev.ExpenseChangePct, 0.01)
	assert.Equal(t, 70000.0, result.ComparisonToPrev.NetCashFlowChange)

	mockRepo.AssertExpectations(t)
}

// TestFormatMonthlyCashFlowReport 測試月度報告格式化
func TestFormatMonthlyCashFlowReport(t *testing.T) {
	service := NewDiscordService()

	summary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        10,
		TotalIncome:  50000,
		TotalExpense: 15000,
		NetCashFlow:  35000,
		IncomeCount:  1,
		ExpenseCount: 2,
		ComparisonToPrev: &models.MonthComparison{
			PreviousMonth:     9,
			PreviousYear:      2024,
			IncomeChange:      5000,
			IncomeChangePct:   11.11,
			ExpenseChange:     -2000,
			ExpenseChangePct:  -11.76,
			NetCashFlowChange: 7000,
		},
	}

	message := service.FormatMonthlyCashFlowReport(summary)
	assert.NotNil(t, message)
	assert.Contains(t, message.Content, "2024年10月 現金流報告")
	assert.Contains(t, message.Content, "50,000")
	assert.Contains(t, message.Content, "15,000")
	assert.Contains(t, message.Content, "35,000")
	assert.Contains(t, message.Content, "與上月")
}

// TestFormatYearlyCashFlowReport 測試年度報告格式化
func TestFormatYearlyCashFlowReport(t *testing.T) {
	service := NewDiscordService()

	summary := &models.YearlyCashFlowSummary{
		Year:         2024,
		TotalIncome:  600000,
		TotalExpense: 360000,
		NetCashFlow:  240000,
		IncomeCount:  12,
		ExpenseCount: 24,
		ComparisonToPrev: &models.YearComparison{
			PreviousYear:      2023,
			IncomeChange:      50000,
			IncomeChangePct:   9.09,
			ExpenseChange:     -20000,
			ExpenseChangePct:  -5.26,
			NetCashFlowChange: 70000,
		},
	}

	message := service.FormatYearlyCashFlowReport(summary)
	assert.NotNil(t, message)
	assert.Contains(t, message.Content, "2024年度 現金流報告")
	assert.Contains(t, message.Content, "600,000")
	assert.Contains(t, message.Content, "360,000")
	assert.Contains(t, message.Content, "240,000")
	assert.Contains(t, message.Content, "與去年")
}

