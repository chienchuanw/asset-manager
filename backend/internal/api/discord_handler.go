package api

import (
	"net/http"
	"sort"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// DiscordHandler Discord API Handler
type DiscordHandler struct {
	discordService  service.DiscordService
	settingsService service.SettingsService
	holdingService  service.HoldingService
}

// NewDiscordHandler 建立新的 Discord Handler
func NewDiscordHandler(
	discordService service.DiscordService,
	settingsService service.SettingsService,
	holdingService service.HoldingService,
) *DiscordHandler {
	return &DiscordHandler{
		discordService:  discordService,
		settingsService: settingsService,
		holdingService:  holdingService,
	}
}

// TestDiscordInput 測試 Discord 輸入
type TestDiscordInput struct {
	Message string `json:"message" binding:"required"` // 測試訊息
}

// TestDiscord 測試 Discord 發送
// @Summary 測試 Discord 發送
// @Description 發送測試訊息到 Discord Webhook
// @Tags discord
// @Accept json
// @Produce json
// @Param input body TestDiscordInput true "測試訊息"
// @Success 200 {object} APIResponse[string]
// @Failure 400 {object} APIResponse[any]
// @Failure 500 {object} APIResponse[any]
// @Router /api/discord/test [post]
func (h *DiscordHandler) TestDiscord(c *gin.Context) {
	// 解析輸入
	var input TestDiscordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_INPUT",
				Message: err.Error(),
			},
		})
		return
	}

	// 取得 Discord 設定
	settings, err := h.settingsService.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_SETTINGS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "DISCORD_DISABLED",
				Message: "Discord is not enabled",
			},
		})
		return
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "WEBHOOK_URL_NOT_SET",
				Message: "Discord webhook URL is not set",
			},
		})
		return
	}

	// 建立測試訊息
	message := &models.DiscordMessage{
		Content: input.Message,
	}

	// 發送訊息
	if err := h.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "SEND_MESSAGE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Message sent successfully",
	})
}

// SendDailyReport 發送每日報告
// @Summary 發送每日報告
// @Description 發送每日資產報告到 Discord
// @Tags discord
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse[string]
// @Failure 500 {object} APIResponse[any]
// @Router /api/discord/daily-report [post]
func (h *DiscordHandler) SendDailyReport(c *gin.Context) {
	// 取得 Discord 設定
	settings, err := h.settingsService.GetSettings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_SETTINGS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 檢查 Discord 是否啟用
	if !settings.Discord.Enabled {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "DISCORD_DISABLED",
				Message: "Discord is not enabled",
			},
		})
		return
	}

	// 檢查 Webhook URL 是否設定
	if settings.Discord.WebhookURL == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "WEBHOOK_URL_NOT_SET",
				Message: "Discord webhook URL is not set",
			},
		})
		return
	}

	// 取得所有持倉
	holdings, err := h.holdingService.GetAllHoldings(models.HoldingFilters{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_HOLDINGS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 計算總資產資料
	var totalMarketValue, totalCost, totalUnrealizedPL float64
	byAssetType := make(map[string]*models.AssetTypePerformance)

	for _, holding := range holdings {
		totalMarketValue += holding.MarketValue
		totalCost += holding.TotalCost
		totalUnrealizedPL += holding.UnrealizedPL

		// 按資產類型分類
		assetTypeStr := string(holding.AssetType)
		if _, exists := byAssetType[assetTypeStr]; !exists {
			byAssetType[assetTypeStr] = &models.AssetTypePerformance{
				AssetType: assetTypeStr,
			}
		}
		perf := byAssetType[assetTypeStr]
		perf.MarketValue += holding.MarketValue
		perf.Cost += holding.TotalCost
		perf.UnrealizedPL += holding.UnrealizedPL
		perf.HoldingCount++
	}

	// 計算各資產類型的損益百分比
	for _, perf := range byAssetType {
		if perf.Cost > 0 {
			perf.UnrealizedPct = (perf.UnrealizedPL / perf.Cost) * 100
		}
	}

	// 計算總損益百分比
	totalUnrealizedPct := 0.0
	if totalCost > 0 {
		totalUnrealizedPct = (totalUnrealizedPL / totalCost) * 100
	}

	// 排序持倉（按市值降序）
	sort.Slice(holdings, func(i, j int) bool {
		return holdings[i].MarketValue > holdings[j].MarketValue
	})

	// 取前 5 大持倉
	topHoldings := holdings
	if len(topHoldings) > 5 {
		topHoldings = topHoldings[:5]
	}

	// 建立報告資料
	reportData := &models.DailyReportData{
		Date:               time.Now(),
		TotalMarketValue:   totalMarketValue,
		TotalCost:          totalCost,
		TotalUnrealizedPL:  totalUnrealizedPL,
		TotalUnrealizedPct: totalUnrealizedPct,
		HoldingCount:       len(holdings),
		TopHoldings:        topHoldings,
		ByAssetType:        byAssetType,
	}

	// 格式化報告
	message := h.discordService.FormatDailyReport(reportData)

	// 發送訊息
	if err := h.discordService.SendMessage(settings.Discord.WebhookURL, message); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "SEND_MESSAGE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: "Daily report sent successfully",
	})
}

