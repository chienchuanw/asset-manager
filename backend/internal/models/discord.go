package models

import "time"

// DiscordMessage Discord 訊息
type DiscordMessage struct {
	Content string         `json:"content,omitempty"` // 純文字訊息
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`  // Embed 訊息（更豐富的格式）
}

// DiscordEmbed Discord Embed 訊息格式
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`       // 標題
	Description string              `json:"description,omitempty"` // 描述
	Color       int                 `json:"color,omitempty"`       // 顏色（十進位）
	Fields      []DiscordEmbedField `json:"fields,omitempty"`      // 欄位列表
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`      // 頁尾
	Timestamp   string              `json:"timestamp,omitempty"`   // 時間戳記（ISO 8601 格式）
}

// DiscordEmbedField Discord Embed 欄位
type DiscordEmbedField struct {
	Name   string `json:"name"`             // 欄位名稱
	Value  string `json:"value"`            // 欄位值
	Inline bool   `json:"inline,omitempty"` // 是否並排顯示
}

// DiscordEmbedFooter Discord Embed 頁尾
type DiscordEmbedFooter struct {
	Text string `json:"text"`               // 頁尾文字
	Icon string `json:"icon_url,omitempty"` // 頁尾圖示 URL
}

// DailyReportData 每日報告資料
type DailyReportData struct {
	Date               time.Time                       // 報告日期
	TotalMarketValue   float64                         // 總市值（TWD）
	TotalCost          float64                         // 總成本（TWD）
	TotalUnrealizedPL  float64                         // 總未實現損益（TWD）
	TotalUnrealizedPct float64                         // 總未實現損益百分比
	HoldingCount       int                             // 持倉數量
	TopHoldings        []*Holding                      // 前 5 大持倉
	ByAssetType        map[string]*AssetTypePerformance // 按資產類型分類
}

// AssetTypePerformance 資產類型績效
type AssetTypePerformance struct {
	AssetType     string  // 資產類型
	MarketValue   float64 // 市值
	Cost          float64 // 成本
	UnrealizedPL  float64 // 未實現損益
	UnrealizedPct float64 // 未實現損益百分比
	HoldingCount  int     // 持倉數量
}

