package models

import (
	"time"

	"github.com/google/uuid"
)

// Setting 設定資料模型
type Setting struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Key         string    `json:"key" db:"key"`
	Value       string    `json:"value" db:"value"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// UpdateSettingInput 更新設定輸入
type UpdateSettingInput struct {
	Value string `json:"value" binding:"required"`
}

// SettingsGroup 設定群組（用於批次取得和更新）
type SettingsGroup struct {
	Discord      DiscordSettings      `json:"discord"`
	Allocation   AllocationSettings   `json:"allocation"`
	Notification NotificationSettings `json:"notification"`
}

// DiscordSettings Discord 設定
type DiscordSettings struct {
	WebhookURL            string `json:"webhook_url"`
	Enabled               bool   `json:"enabled"`
	ReportTime            string `json:"report_time"`              // HH:MM 格式
	MonthlyReportEnabled  bool   `json:"monthly_report_enabled"`   // 月度現金流報告開關
	MonthlyReportDay      int    `json:"monthly_report_day"`       // 每月幾號發送 (1-10)
	YearlyReportEnabled   bool   `json:"yearly_report_enabled"`    // 年度現金流報告開關
	YearlyReportMonth     int    `json:"yearly_report_month"`      // 每年幾月發送 (1-12)
	YearlyReportDay       int    `json:"yearly_report_day"`        // 每年幾號發送 (1-10)
}

// AllocationSettings 資產配置設定
type AllocationSettings struct {
	TWStock           float64 `json:"tw_stock"`            // 台股目標配置百分比
	USStock           float64 `json:"us_stock"`            // 美股目標配置百分比
	Crypto            float64 `json:"crypto"`              // 加密貨幣目標配置百分比
	RebalanceThreshold float64 `json:"rebalance_threshold"` // 再平衡閾值百分比
}

// NotificationSettings 通知設定
type NotificationSettings struct {
	DailyBilling            bool `json:"daily_billing"`             // 每日扣款通知
	SubscriptionExpiry      bool `json:"subscription_expiry"`       // 訂閱到期通知
	InstallmentCompletion   bool `json:"installment_completion"`    // 分期完成通知
	ExpiryDays              int  `json:"expiry_days"`               // 到期提醒天數
}

// UpdateSettingsGroupInput 更新設定群組輸入
type UpdateSettingsGroupInput struct {
	Discord      *DiscordSettings      `json:"discord,omitempty"`
	Allocation   *AllocationSettings   `json:"allocation,omitempty"`
	Notification *NotificationSettings `json:"notification,omitempty"`
}

