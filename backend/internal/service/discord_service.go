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

	// SendDailyBillingNotification 發送每日扣款通知
	SendDailyBillingNotification(webhookURL string, result *DailyBillingResult) error

	// SendSubscriptionExpiryNotification 發送訂閱到期通知
	SendSubscriptionExpiryNotification(webhookURL string, subscriptions []*models.Subscription, days int) error

	// SendInstallmentCompletionNotification 發送分期完成通知
	SendInstallmentCompletionNotification(webhookURL string, installments []*models.Installment, remainingCount int) error

	// SendCreditCardPaymentReminder 發送信用卡繳款提醒
	SendCreditCardPaymentReminder(webhookURL string, creditCards []*models.CreditCard) error
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

	// 檢查是否有價格資料問題
	priceWarning := s.checkPriceDataQuality(data)

	// 建立描述（包含日期和可能的警告）
	description := fmt.Sprintf("報告日期：%s", data.Date.Format("2006-01-02"))
	if priceWarning != "" {
		description += "\n\n" + priceWarning
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "📊 每日資產報告",
		Description: description,
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

// checkPriceDataQuality 檢查價格資料品質並返回警告訊息
func (s *discordService) checkPriceDataQuality(data *models.DailyReportData) string {
	if data == nil || len(data.TopHoldings) == 0 {
		return ""
	}

	staleCount := 0
	unavailableCount := 0

	for _, holding := range data.TopHoldings {
		if holding.PriceSource == "unavailable" {
			unavailableCount++
		} else if holding.IsPriceStale {
			staleCount++
		}
	}

	// 如果有價格問題，返回警告訊息
	if unavailableCount > 0 || staleCount > 0 {
		warning := "⚠️ **價格資料警告**\n"
		if unavailableCount > 0 {
			warning += fmt.Sprintf("• %d 個標的無法取得價格（使用成本價估算）\n", unavailableCount)
		}
		if staleCount > 0 {
			warning += fmt.Sprintf("• %d 個標的使用快取價格（API 達到速率限制）\n", staleCount)
		}
		return warning
	}

	return ""
}

// SendDailyBillingNotification 發送每日扣款通知
func (s *discordService) SendDailyBillingNotification(webhookURL string, result *DailyBillingResult) error {
	if result == nil {
		return fmt.Errorf("billing result is nil")
	}

	// 如果沒有任何扣款，不發送通知
	if result.SubscriptionCount == 0 && result.InstallmentCount == 0 {
		return nil
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "💳 每日扣款通知",
		Description: fmt.Sprintf("扣款日期：%s", result.Date.Format("2006-01-02")),
		Color:       0x3498db, // 藍色
		Fields:      []models.DiscordEmbedField{},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &models.DiscordEmbedFooter{
			Text: "Asset Manager - 訂閱分期管理",
		},
	}

	// 總覽
	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name: "📊 扣款總覽",
		Value: fmt.Sprintf("訂閱扣款：%d 筆\n分期扣款：%d 筆\n總金額：NT$ %.2f",
			result.SubscriptionCount,
			result.InstallmentCount,
			result.TotalAmount,
		),
		Inline: false,
	})

	// 訂閱扣款詳情
	if result.SubscriptionCount > 0 && len(result.SubscriptionResult.CreatedCashFlows) > 0 {
		subscriptionText := ""
		for i, cf := range result.SubscriptionResult.CreatedCashFlows {
			if i >= 5 { // 最多顯示 5 筆
				subscriptionText += fmt.Sprintf("...及其他 %d 筆\n", result.SubscriptionCount-5)
				break
			}
			subscriptionText += fmt.Sprintf("• %s - NT$ %.2f\n", cf.Description, cf.Amount)
		}
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "📅 訂閱扣款",
			Value:  subscriptionText,
			Inline: false,
		})
	}

	// 分期扣款詳情
	if result.InstallmentCount > 0 && len(result.InstallmentResult.CreatedCashFlows) > 0 {
		installmentText := ""
		for i, cf := range result.InstallmentResult.CreatedCashFlows {
			if i >= 5 { // 最多顯示 5 筆
				installmentText += fmt.Sprintf("...及其他 %d 筆\n", result.InstallmentCount-5)
				break
			}
			installmentText += fmt.Sprintf("• %s - NT$ %.2f\n", cf.Description, cf.Amount)
		}
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "💰 分期扣款",
			Value:  installmentText,
			Inline: false,
		})
	}

	// 錯誤訊息（如果有）
	totalErrors := len(result.SubscriptionResult.Errors) + len(result.InstallmentResult.Errors)
	if totalErrors > 0 {
		errorText := fmt.Sprintf("⚠️ 有 %d 筆扣款失敗，請檢查系統日誌", totalErrors)
		embed.Fields = append(embed.Fields, models.DiscordEmbedField{
			Name:   "錯誤",
			Value:  errorText,
			Inline: false,
		})
	}

	message := &models.DiscordMessage{
		Embeds: []models.DiscordEmbed{embed},
	}

	return s.SendMessage(webhookURL, message)
}

// SendSubscriptionExpiryNotification 發送訂閱到期通知
func (s *discordService) SendSubscriptionExpiryNotification(webhookURL string, subscriptions []*models.Subscription, days int) error {
	if len(subscriptions) == 0 {
		return nil
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "⏰ 訂閱到期提醒",
		Description: fmt.Sprintf("以下訂閱將在 %d 天內到期", days),
		Color:       0xf39c12, // 橘色
		Fields:      []models.DiscordEmbedField{},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &models.DiscordEmbedFooter{
			Text: "Asset Manager - 訂閱分期管理",
		},
	}

	// 訂閱列表
	subscriptionText := ""
	for i, sub := range subscriptions {
		if i >= 10 { // 最多顯示 10 筆
			subscriptionText += fmt.Sprintf("...及其他 %d 筆\n", len(subscriptions)-10)
			break
		}
		if sub.EndDate != nil {
			daysUntilExpiry := int(sub.EndDate.Sub(time.Now()).Hours() / 24)
			subscriptionText += fmt.Sprintf("• %s - NT$ %.2f/月 (剩餘 %d 天)\n",
				sub.Name,
				sub.Amount,
				daysUntilExpiry,
			)
		}
	}

	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "即將到期的訂閱",
		Value:  subscriptionText,
		Inline: false,
	})

	message := &models.DiscordMessage{
		Embeds: []models.DiscordEmbed{embed},
	}

	return s.SendMessage(webhookURL, message)
}

// SendInstallmentCompletionNotification 發送分期完成通知
func (s *discordService) SendInstallmentCompletionNotification(webhookURL string, installments []*models.Installment, remainingCount int) error {
	if len(installments) == 0 {
		return nil
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "🎉 分期即將完成",
		Description: fmt.Sprintf("以下分期剩餘 %d 期或更少", remainingCount),
		Color:       0x2ecc71, // 綠色
		Fields:      []models.DiscordEmbedField{},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &models.DiscordEmbedFooter{
			Text: "Asset Manager - 訂閱分期管理",
		},
	}

	// 分期列表
	installmentText := ""
	for i, inst := range installments {
		if i >= 10 { // 最多顯示 10 筆
			installmentText += fmt.Sprintf("...及其他 %d 筆\n", len(installments)-10)
			break
		}
		remaining := inst.RemainingCount()
		remainingAmount := inst.RemainingAmount()
		installmentText += fmt.Sprintf("• %s - 剩餘 %d/%d 期 (NT$ %.2f)\n",
			inst.Name,
			remaining,
			inst.InstallmentCount,
			remainingAmount,
		)
	}

	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "即將完成的分期",
		Value:  installmentText,
		Inline: false,
	})

	message := &models.DiscordMessage{
		Embeds: []models.DiscordEmbed{embed},
	}

	return s.SendMessage(webhookURL, message)
}

// SendCreditCardPaymentReminder 發送信用卡繳款提醒
func (s *discordService) SendCreditCardPaymentReminder(webhookURL string, creditCards []*models.CreditCard) error {
	if len(creditCards) == 0 {
		return nil
	}

	// 建立 Embed
	embed := models.DiscordEmbed{
		Title:       "🔔 信用卡繳款提醒",
		Description: "以下信用卡明天是繳款日，請記得準時繳款！",
		Color:       0xff9800, // 橘色
		Fields:      []models.DiscordEmbedField{},
		Timestamp:   time.Now().Format(time.RFC3339),
		Footer: &models.DiscordEmbedFooter{
			Text: "Asset Manager - 信用卡管理",
		},
	}

	// 信用卡列表
	cardText := ""
	for i, card := range creditCards {
		if i >= 10 { // 最多顯示 10 張
			cardText += fmt.Sprintf("...及其他 %d 張信用卡\n", len(creditCards)-10)
			break
		}

		availableCredit := card.CreditLimit - card.UsedCredit
		utilizationRate := (card.UsedCredit / card.CreditLimit) * 100

		cardText += fmt.Sprintf(
			"💳 **%s %s** (****%s)\n"+
				"   繳款日: 每月 %d 號\n"+
				"   目前已使用額度: NT$ %s\n"+
				"   可用額度: NT$ %s\n"+
				"   使用率: %.1f%%\n\n",
			card.IssuingBank,
			card.CardName,
			card.CardNumberLast4,
			card.PaymentDueDay,
			formatCurrency(card.UsedCredit),
			formatCurrency(availableCredit),
			utilizationRate,
		)
	}

	embed.Fields = append(embed.Fields, models.DiscordEmbedField{
		Name:   "明天需要繳款的信用卡",
		Value:  cardText,
		Inline: false,
	})

	message := &models.DiscordMessage{
		Embeds: []models.DiscordEmbed{embed},
	}

	return s.SendMessage(webhookURL, message)
}

// formatCurrency 格式化貨幣顯示（加入千分位逗號）
func formatCurrency(amount float64) string {
	// 將數字轉為整數字串
	intAmount := int64(amount)
	str := fmt.Sprintf("%d", intAmount)

	// 如果是負數，先處理符號
	negative := false
	if str[0] == '-' {
		negative = true
		str = str[1:]
	}

	// 加入千分位逗號
	n := len(str)
	if n <= 3 {
		if negative {
			return "-" + str
		}
		return str
	}

	// 從右到左每三位加一個逗號
	result := ""
	for i, digit := range str {
		if i > 0 && (n-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}

	if negative {
		return "-" + result
	}
	return result
}
