package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// DiscordService Discord 服務介面
type DiscordService interface {
	// SendMessage 發送訊息到 Discord Webhook
	SendMessage(webhookURL string, message *models.DiscordMessage) error

	// FormatDailyReport 格式化每日報告
	FormatDailyReport(data *models.DailyReportData) *models.DiscordMessage
}

// discordService Discord 服務實作
type discordService struct {
	httpClient *http.Client
}

// NewDiscordService 建立新的 Discord 服務
func NewDiscordService() DiscordService {
	return &discordService{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SendMessage 發送訊息到 Discord Webhook
func (s *discordService) SendMessage(webhookURL string, message *models.DiscordMessage) error {
	// 將訊息轉換為 JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 建立 HTTP 請求
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// 設定 Content-Type
	req.Header.Set("Content-Type", "application/json")

	// 發送請求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 檢查回應狀態碼
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}

// FormatDailyReport 格式化每日報告
func (s *discordService) FormatDailyReport(data *models.DailyReportData) *models.DiscordMessage {
	// 計算顏色（綠色=獲利，紅色=虧損）
	color := 0x00FF00 // 綠色
	if data.TotalUnrealizedPL < 0 {
		color = 0xFF0000 // 紅色
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "📊 每日資產報告",
		Description: fmt.Sprintf("報告日期：%s", data.Date.Format("2006-01-02")),
		Color:       color,
		Fields:      []models.DiscordEmbedField{},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &models.DiscordEmbedFooter{
			Text: "Asset Manager",
		},
	}

	// 總資產概況
	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "💰 總資產",
		Value:  fmt.Sprintf("NT$ %s", formatNumber(data.TotalMarketValue)),
		Inline: true,
	})

	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "💵 總成本",
		Value:  fmt.Sprintf("NT$ %s", formatNumber(data.TotalCost)),
		Inline: true,
	})

	plSymbol := "📈"
	if data.TotalUnrealizedPL < 0 {
		plSymbol = "📉"
	}
	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   fmt.Sprintf("%s 未實現損益", plSymbol),
		Value:  fmt.Sprintf("NT$ %s (%+.2f%%)", formatNumber(data.TotalUnrealizedPL), data.TotalUnrealizedPct),
		Inline: true,
	})

	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "📦 持倉數量",
		Value:  fmt.Sprintf("%d 個標的", data.HoldingCount),
		Inline: true,
	})

	// 各資產類型表現
	if len(data.ByAssetType) > 0 {
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "\n📊 各資產類型表現",
			Value:  "\u200B", // 空白字元
			Inline: false,
		})

		for assetType, perf := range data.ByAssetType {
			typeLabel := getAssetTypeLabel(assetType)
			embed.Fields = append(embed.Fields, models.DiscordEmbedField{
				Name: typeLabel,
				Value: fmt.Sprintf("市值: NT$ %s\n損益: NT$ %s (%+.2f%%)",
					formatNumber(perf.MarketValue),
					formatNumber(perf.UnrealizedPL),
					perf.UnrealizedPct,
				),
				Inline: true,
			})
		}
	}

	// 前 5 大持倉
	if len(data.TopHoldings) > 0 {
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "\n🏆 前 5 大持倉",
			Value:  "\u200B",
			Inline: false,
		})

		for i, holding := range data.TopHoldings {
			if i >= 5 {
				break
			}
			embed.Fields = append(embed.Fields, models.DiscordEmbedField{
				Name: fmt.Sprintf("%d. %s (%s)", i+1, holding.Name, holding.Symbol),
				Value: fmt.Sprintf("市值: NT$ %s\n損益: %+.2f%%",
					formatNumber(holding.MarketValue),
					holding.UnrealizedPLPct,
				),
				Inline: true,
			})
		}
	}

	// 再平衡檢查
	if data.RebalanceCheck != nil && data.RebalanceCheck.NeedsRebalance {
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "\n⚠️ 再平衡提醒",
			Value:  "\u200B",
			Inline: false,
		})

		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "狀態",
			Value:  fmt.Sprintf("⚠️ 需要再平衡（閾值: %.1f%%）", data.RebalanceCheck.Threshold),
			Inline: false,
		})

		// 顯示偏離情況
		deviationText := ""
		for _, deviation := range data.RebalanceCheck.Deviations {
			if deviation.ExceedsThreshold {
				typeLabel := getAssetTypeLabel(deviation.AssetType)
				symbol := "📈"
				if deviation.Deviation < 0 {
					symbol = "📉"
				}
				deviationText += fmt.Sprintf("%s %s: %.1f%% → %.1f%% (%s%.1f%%)\n",
					symbol,
					typeLabel,
					deviation.TargetPercent,
					deviation.CurrentPercent,
					getDeviationSign(deviation.Deviation),
					deviation.DeviationAbs,
				)
			}
		}
		if deviationText != "" {
			embed.Fields = append(embed.Fields, models.DiscordEmbedField{
				Name:   "偏離情況",
				Value:  deviationText,
				Inline: false,
			})
		}

		// 顯示建議（最多 3 個）
		if len(data.RebalanceCheck.Suggestions) > 0 {
			suggestionText := ""
			for i, suggestion := range data.RebalanceCheck.Suggestions {
				if i >= 3 {
					break
				}
				actionSymbol := "🔴"
				if suggestion.Action == "buy" {
					actionSymbol = "🟢"
				}
				typeLabel := getAssetTypeLabel(suggestion.AssetType)
				suggestionText += fmt.Sprintf("%s %s: NT$ %s\n",
					actionSymbol,
					typeLabel,
					formatNumber(suggestion.Amount),
				)
			}
			embed.Fields = append(embed.Fields, models.DiscordEmbedField{
				Name:   "建議操作",
				Value:  suggestionText,
				Inline: false,
			})
		}
	}

	return &models.DiscordMessage{
		Embeds: []models.DiscordEmbed{embed},
	}
}

// formatNumber 格式化數字（加上千分位）
func formatNumber(num float64) string {
	// 簡單的千分位格式化
	str := fmt.Sprintf("%.2f", num)
	// 這裡可以加上更複雜的千分位邏輯
	return str
}

// getAssetTypeLabel 取得資產類型標籤
func getAssetTypeLabel(assetType string) string {
	switch assetType {
	case "tw-stock":
		return "🇹🇼 台股"
	case "us-stock":
		return "🇺🇸 美股"
	case "crypto":
		return "₿ 加密貨幣"
	default:
		return assetType
	}
}

// getDeviationSign 取得偏離符號
func getDeviationSign(deviation float64) string {
	if deviation > 0 {
		return "+"
	}
	return ""
}
