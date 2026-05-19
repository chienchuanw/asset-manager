package models

// FIRE (Financial Independence / Retire Early) 計算預設值
const (
	DefaultExpectedReturn = 0.05 // 預設預期年化報酬率 5%
	DefaultWithdrawalRate = 0.04 // 預設安全提領率 4%
	MaxProjectionYears    = 60   // 投影年數上限
)

// FireProjectionInput FIRE 投影輸入；所有欄位皆為可選覆寫值，
// 未提供時由服務從應用資料自動推導。
type FireProjectionInput struct {
	NetWorth       *float64
	AnnualExpenses *float64
	AnnualSavings  *float64
	ExpectedReturn *float64
	WithdrawalRate *float64
}

// FireProjectionPoint 投影圖表上的單一資料點
type FireProjectionPoint struct {
	Year              int     `json:"year"`
	ProjectedNetWorth float64 `json:"projected_net_worth"`
	FireTarget        float64 `json:"fire_target"`
}

// FireProjectionResult FIRE 投影結果，回傳已套用的輸入值供前端回填表單。
type FireProjectionResult struct {
	CurrentNetWorth float64               `json:"current_net_worth"`
	AnnualExpenses  float64               `json:"annual_expenses"`
	AnnualSavings   float64               `json:"annual_savings"`
	ExpectedReturn  float64               `json:"expected_return"`
	WithdrawalRate  float64               `json:"withdrawal_rate"`
	FireNumber      float64               `json:"fire_number"`
	ProgressPct     float64               `json:"progress_pct"`
	YearsToFI       *int                  `json:"years_to_fi"`
	OnTrack         bool                  `json:"on_track"`
	Projection      []FireProjectionPoint `json:"projection"`
}
