package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestCashFlowReportType 測試報告類型常數
func TestCashFlowReportType(t *testing.T) {
	assert.Equal(t, CashFlowReportType("monthly"), CashFlowReportTypeMonthly)
	assert.Equal(t, CashFlowReportType("yearly"), CashFlowReportTypeYearly)
}

// TestCashFlowReportLog 測試報告記錄模型
func TestCashFlowReportLog(t *testing.T) {
	month := 10
	errorMsg := "test error"

	log := CashFlowReportLog{
		ID:         uuid.New(),
		ReportType: CashFlowReportTypeMonthly,
		Year:       2024,
		Month:      &month,
		SentAt:     time.Now(),
		Success:    true,
		ErrorMsg:   &errorMsg,
		RetryCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// 驗證所有欄位都能正確設定
	assert.NotEqual(t, uuid.Nil, log.ID)
	assert.Equal(t, CashFlowReportTypeMonthly, log.ReportType)
	assert.Equal(t, 2024, log.Year)
	assert.NotNil(t, log.Month)
	assert.Equal(t, 10, *log.Month)
	assert.True(t, log.Success)
	assert.NotNil(t, log.ErrorMsg)
	assert.Equal(t, "test error", *log.ErrorMsg)
	assert.Equal(t, 0, log.RetryCount)
}

// TestCreateCashFlowReportLogInput 測試建立報告記錄輸入
func TestCreateCashFlowReportLogInput(t *testing.T) {
	month := 10
	errorMsg := "test error"

	input := CreateCashFlowReportLogInput{
		ReportType: CashFlowReportTypeMonthly,
		Year:       2024,
		Month:      &month,
		Success:    false,
		ErrorMsg:   &errorMsg,
		RetryCount: 1,
	}

	assert.Equal(t, CashFlowReportTypeMonthly, input.ReportType)
	assert.Equal(t, 2024, input.Year)
	assert.NotNil(t, input.Month)
	assert.Equal(t, 10, *input.Month)
	assert.False(t, input.Success)
	assert.NotNil(t, input.ErrorMsg)
	assert.Equal(t, "test error", *input.ErrorMsg)
	assert.Equal(t, 1, input.RetryCount)
}

// TestUpdateCashFlowReportLogInput 測試更新報告記錄輸入
func TestUpdateCashFlowReportLogInput(t *testing.T) {
	success := true
	errorMsg := "updated error"
	retryCount := 2

	input := UpdateCashFlowReportLogInput{
		Success:    &success,
		ErrorMsg:   &errorMsg,
		RetryCount: &retryCount,
	}

	assert.NotNil(t, input.Success)
	assert.True(t, *input.Success)
	assert.NotNil(t, input.ErrorMsg)
	assert.Equal(t, "updated error", *input.ErrorMsg)
	assert.NotNil(t, input.RetryCount)
	assert.Equal(t, 2, *input.RetryCount)
}

// TestMonthlyCashFlowSummary 測試月度摘要模型
func TestMonthlyCashFlowSummary(t *testing.T) {
	summary := MonthlyCashFlowSummary{
		Year:         2024,
		Month:        10,
		TotalIncome:  50000,
		TotalExpense: 30000,
		NetCashFlow:  20000,
		IncomeCount:  5,
		ExpenseCount: 15,
		IncomeCategoryBreakdown: []*CategorySummary{
			{
				CategoryID:   uuid.New(),
				CategoryName: "薪資",
				Amount:       50000,
				Count:        1,
			},
		},
		ExpenseCategoryBreakdown: []*CategorySummary{
			{
				CategoryID:   uuid.New(),
				CategoryName: "餐飲",
				Amount:       10000,
				Count:        10,
			},
		},
		TopExpenses: []*CashFlow{},
		ComparisonToPrev: &MonthComparison{
			PreviousMonth:     9,
			PreviousYear:      2024,
			IncomeChange:      5000,
			IncomeChangePct:   11.11,
			ExpenseChange:     -2000,
			ExpenseChangePct:  -6.25,
			NetCashFlowChange: 7000,
		},
	}

	assert.Equal(t, 2024, summary.Year)
	assert.Equal(t, 10, summary.Month)
	assert.Equal(t, 50000.0, summary.TotalIncome)
	assert.Equal(t, 30000.0, summary.TotalExpense)
	assert.Equal(t, 20000.0, summary.NetCashFlow)
	assert.Equal(t, 5, summary.IncomeCount)
	assert.Equal(t, 15, summary.ExpenseCount)
	assert.Len(t, summary.IncomeCategoryBreakdown, 1)
	assert.Len(t, summary.ExpenseCategoryBreakdown, 1)
	assert.NotNil(t, summary.ComparisonToPrev)
	assert.Equal(t, 9, summary.ComparisonToPrev.PreviousMonth)
}

// TestYearlyCashFlowSummary 測試年度摘要模型
func TestYearlyCashFlowSummary(t *testing.T) {
	summary := YearlyCashFlowSummary{
		Year:         2024,
		TotalIncome:  600000,
		TotalExpense: 360000,
		NetCashFlow:  240000,
		IncomeCount:  60,
		ExpenseCount: 180,
		IncomeCategoryBreakdown: []*CategorySummary{
			{
				CategoryID:   uuid.New(),
				CategoryName: "薪資",
				Amount:       600000,
				Count:        12,
			},
		},
		ExpenseCategoryBreakdown: []*CategorySummary{
			{
				CategoryID:   uuid.New(),
				CategoryName: "餐飲",
				Amount:       120000,
				Count:        120,
			},
		},
		MonthlyBreakdown: []*MonthlyBreakdown{
			{
				Month:       1,
				Income:      50000,
				Expense:     30000,
				NetCashFlow: 20000,
			},
		},
		TopExpenses: []*CashFlow{},
		ComparisonToPrev: &YearComparison{
			PreviousYear:      2023,
			IncomeChange:      60000,
			IncomeChangePct:   11.11,
			ExpenseChange:     -24000,
			ExpenseChangePct:  -6.25,
			NetCashFlowChange: 84000,
		},
	}

	assert.Equal(t, 2024, summary.Year)
	assert.Equal(t, 600000.0, summary.TotalIncome)
	assert.Equal(t, 360000.0, summary.TotalExpense)
	assert.Equal(t, 240000.0, summary.NetCashFlow)
	assert.Equal(t, 60, summary.IncomeCount)
	assert.Equal(t, 180, summary.ExpenseCount)
	assert.Len(t, summary.IncomeCategoryBreakdown, 1)
	assert.Len(t, summary.ExpenseCategoryBreakdown, 1)
	assert.Len(t, summary.MonthlyBreakdown, 1)
	assert.NotNil(t, summary.ComparisonToPrev)
	assert.Equal(t, 2023, summary.ComparisonToPrev.PreviousYear)
}

