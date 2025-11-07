package models

import (
	"time"

	"github.com/google/uuid"
)

// CashFlowReportType 現金流報告類型
type CashFlowReportType string

const (
	CashFlowReportTypeMonthly CashFlowReportType = "monthly" // 月度報告
	CashFlowReportTypeYearly  CashFlowReportType = "yearly"  // 年度報告
)

// CashFlowReportLog 現金流報告發送記錄
type CashFlowReportLog struct {
	ID          uuid.UUID          `json:"id" db:"id"`
	ReportType  CashFlowReportType `json:"report_type" db:"report_type"`
	Year        int                `json:"year" db:"year"`
	Month       *int               `json:"month,omitempty" db:"month"` // 月度報告才有
	SentAt      time.Time          `json:"sent_at" db:"sent_at"`
	Success     bool               `json:"success" db:"success"`
	ErrorMsg    *string            `json:"error_msg,omitempty" db:"error_msg"`
	RetryCount  int                `json:"retry_count" db:"retry_count"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// CreateCashFlowReportLogInput 建立報告記錄的輸入
type CreateCashFlowReportLogInput struct {
	ReportType CashFlowReportType `json:"report_type"`
	Year       int                `json:"year"`
	Month      *int               `json:"month,omitempty"`
	Success    bool               `json:"success"`
	ErrorMsg   *string            `json:"error_msg,omitempty"`
	RetryCount int                `json:"retry_count"`
}

// UpdateCashFlowReportLogInput 更新報告記錄的輸入
type UpdateCashFlowReportLogInput struct {
	Success    *bool   `json:"success,omitempty"`
	ErrorMsg   *string `json:"error_msg,omitempty"`
	RetryCount *int    `json:"retry_count,omitempty"`
}

// CategorySummary 分類摘要（用於月度/年度報告）
type CategorySummary struct {
	CategoryID   uuid.UUID `json:"category_id" db:"category_id"`
	CategoryName string    `json:"category_name" db:"category_name"`
	Amount       float64   `json:"amount" db:"amount"`
	Count        int       `json:"count" db:"count"`
}

// MonthComparison 月份比較資料
type MonthComparison struct {
	PreviousMonth      int     `json:"previous_month"`
	PreviousYear       int     `json:"previous_year"`
	IncomeChange       float64 `json:"income_change"`        // 收入變化金額
	IncomeChangePct    float64 `json:"income_change_pct"`    // 收入變化百分比
	ExpenseChange      float64 `json:"expense_change"`       // 支出變化金額
	ExpenseChangePct   float64 `json:"expense_change_pct"`   // 支出變化百分比
	NetCashFlowChange  float64 `json:"net_cash_flow_change"` // 淨現金流變化
}

// MonthlyCashFlowSummary 月度現金流摘要
type MonthlyCashFlowSummary struct {
	Year                     int                `json:"year"`
	Month                    int                `json:"month"`
	TotalIncome              float64            `json:"total_income"`
	TotalExpense             float64            `json:"total_expense"`
	NetCashFlow              float64            `json:"net_cash_flow"`
	IncomeCount              int                `json:"income_count"`
	ExpenseCount             int                `json:"expense_count"`
	IncomeCategoryBreakdown  []*CategorySummary `json:"income_category_breakdown"`
	ExpenseCategoryBreakdown []*CategorySummary `json:"expense_category_breakdown"`
	TopExpenses              []*CashFlow        `json:"top_expenses"`
	ComparisonToPrev         *MonthComparison   `json:"comparison_to_prev,omitempty"`
}

// YearComparison 年度比較資料
type YearComparison struct {
	PreviousYear       int     `json:"previous_year"`
	IncomeChange       float64 `json:"income_change"`        // 收入變化金額
	IncomeChangePct    float64 `json:"income_change_pct"`    // 收入變化百分比
	ExpenseChange      float64 `json:"expense_change"`       // 支出變化金額
	ExpenseChangePct   float64 `json:"expense_change_pct"`   // 支出變化百分比
	NetCashFlowChange  float64 `json:"net_cash_flow_change"` // 淨現金流變化
}

// MonthlyBreakdown 月度細分（用於年度報告）
type MonthlyBreakdown struct {
	Month       int     `json:"month"`
	Income      float64 `json:"income"`
	Expense     float64 `json:"expense"`
	NetCashFlow float64 `json:"net_cash_flow"`
}

// YearlyCashFlowSummary 年度現金流摘要
type YearlyCashFlowSummary struct {
	Year                     int                 `json:"year"`
	TotalIncome              float64             `json:"total_income"`
	TotalExpense             float64             `json:"total_expense"`
	NetCashFlow              float64             `json:"net_cash_flow"`
	IncomeCount              int                 `json:"income_count"`
	ExpenseCount             int                 `json:"expense_count"`
	IncomeCategoryBreakdown  []*CategorySummary  `json:"income_category_breakdown"`
	ExpenseCategoryBreakdown []*CategorySummary  `json:"expense_category_breakdown"`
	MonthlyBreakdown         []*MonthlyBreakdown `json:"monthly_breakdown"`
	TopExpenses              []*CashFlow         `json:"top_expenses"`
	ComparisonToPrev         *YearComparison     `json:"comparison_to_prev,omitempty"`
}

