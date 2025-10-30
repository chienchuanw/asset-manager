package api

import (
	"net/http"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// ExchangeRateHandler 匯率 API Handler
type ExchangeRateHandler struct {
	service service.ExchangeRateService
}

// NewExchangeRateHandler 建立新的匯率 Handler
func NewExchangeRateHandler(service service.ExchangeRateService) *ExchangeRateHandler {
	return &ExchangeRateHandler{
		service: service,
	}
}

// ExchangeRateResponse 匯率回應
type ExchangeRateResponse struct {
	FromCurrency string    `json:"from_currency"` // 來源幣別
	ToCurrency   string    `json:"to_currency"`   // 目標幣別
	Rate         float64   `json:"rate"`          // 匯率
	Date         string    `json:"date"`          // 日期 (YYYY-MM-DD)
	UpdatedAt    time.Time `json:"updated_at"`    // 更新時間
	Source       string    `json:"source"`        // 資料來源
}

// RefreshExchangeRate 更新今日匯率
// @Summary 更新今日匯率
// @Description 從 ExchangeRate-API 更新今日的 USD/TWD 匯率
// @Tags exchange-rates
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=ExchangeRateResponse}
// @Failure 500 {object} APIResponse
// @Router /api/exchange-rates/refresh [post]
func (h *ExchangeRateHandler) RefreshExchangeRate(c *gin.Context) {
	// 呼叫 service 更新匯率
	if err := h.service.RefreshTodayRate(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "REFRESH_RATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 取得更新後的匯率記錄
	today := time.Now().Truncate(24 * time.Hour)
	rateRecord, err := h.service.GetRateRecord(models.CurrencyUSD, models.CurrencyTWD, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_RATE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 建立回應
	response := ExchangeRateResponse{
		FromCurrency: string(rateRecord.FromCurrency),
		ToCurrency:   string(rateRecord.ToCurrency),
		Rate:         rateRecord.Rate,
		Date:         rateRecord.Date.Format("2006-01-02"),
		UpdatedAt:    rateRecord.UpdatedAt,
		Source:       "ExchangeRate-API",
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: response,
	})
}

