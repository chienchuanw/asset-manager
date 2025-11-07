package integration

import (
	"testing"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/stretchr/testify/assert"
)

// TestDiscordMessageFormat 測試 Discord 訊息格式
func TestDiscordMessageFormat(t *testing.T) {
	t.Run("Monthly Report Message Format", func(t *testing.T) {
		// 建立測試用的月度摘要
		summary := &models.MonthlyCashFlowSummary{
			Year:         2024,
			Month:        12,
			TotalIncome:  60000,
			TotalExpense: 20000,
			NetCashFlow:  40000,
			IncomeCount:  2,
			ExpenseCount: 2,
			IncomeCategoryBreakdown: []*models.CategorySummary{
				{CategoryName: "薪資", Amount: 50000, Count: 1},
				{CategoryName: "獎金", Amount: 10000, Count: 1},
			},
			ExpenseCategoryBreakdown: []*models.CategorySummary{
				{CategoryName: "飲食", Amount: 15000, Count: 1},
				{CategoryName: "交通", Amount: 5000, Count: 1},
			},
			ComparisonToPrev: &models.MonthComparison{
				PreviousYear:      2024,
				PreviousMonth:     11,
				IncomeChange:      12000,
				IncomeChangePct:   25.0,
				ExpenseChange:     8000,
				ExpenseChangePct:  66.67,
				NetCashFlowChange: 4000,
			},
		}

		// 建立 Discord Service
		discordService := service.NewDiscordService()
		message := discordService.FormatMonthlyCashFlowReport(summary)

		// 驗證訊息格式
		assert.NotNil(t, message)
		assert.NotEmpty(t, message.Content)

		content := message.Content
		assert.Contains(t, content, "📊 【2024年12月 現金流報告】")
		assert.Contains(t, content, "💰 收入：NT$ 60,000")
		assert.Contains(t, content, "💸 支出：NT$ 20,000")
		assert.Contains(t, content, "📈 淨現金流：NT$ 40,000")
		assert.Contains(t, content, "📊 與上月（2024年11月）比較")
		assert.Contains(t, content, "收入：+NT$ 12,000")
		assert.Contains(t, content, "支出：+NT$ 8,000")
		assert.Contains(t, content, "淨現金流：+NT$ 4,000")

		t.Logf("Monthly Report Message:\n%s", content)
	})

	t.Run("Yearly Report Message Format", func(t *testing.T) {
		// 建立測試用的年度摘要
		summary := &models.YearlyCashFlowSummary{
			Year:         2024,
			TotalIncome:  720000,
			TotalExpense: 480000,
			NetCashFlow:  240000,
			IncomeCount:  24,
			ExpenseCount: 48,
			IncomeCategoryBreakdown: []*models.CategorySummary{
				{CategoryName: "薪資", Amount: 600000, Count: 12},
				{CategoryName: "獎金", Amount: 120000, Count: 12},
			},
			ExpenseCategoryBreakdown: []*models.CategorySummary{
				{CategoryName: "飲食", Amount: 180000, Count: 24},
				{CategoryName: "交通", Amount: 60000, Count: 12},
				{CategoryName: "娛樂", Amount: 120000, Count: 12},
			},
			ComparisonToPrev: &models.YearComparison{
				PreviousYear:      2023,
				IncomeChange:      50000,
				IncomeChangePct:   7.46,
				ExpenseChange:     -20000,
				ExpenseChangePct:  -4.0,
				NetCashFlowChange: 70000,
			},
		}

		// 建立 Discord Service
		discordService := service.NewDiscordService()
		message := discordService.FormatYearlyCashFlowReport(summary)

		// 驗證訊息格式
		assert.NotNil(t, message)
		assert.NotEmpty(t, message.Content)

		content := message.Content
		assert.Contains(t, content, "📊 【2024年度 現金流報告】")
		assert.Contains(t, content, "💰 年度收入：NT$ 720,000")
		assert.Contains(t, content, "💸 年度支出：NT$ 480,000")
		assert.Contains(t, content, "📈 年度淨現金流：NT$ 240,000")
		assert.Contains(t, content, "📊 與去年（2023年）比較")
		assert.Contains(t, content, "收入：+NT$ 50,000")
		assert.Contains(t, content, "支出：NT$ -20,000") // 注意：負數格式是 "NT$ -20,000"
		assert.Contains(t, content, "淨現金流：+NT$ 70,000")

		t.Logf("Yearly Report Message:\n%s", content)
	})
}

