package api

import (
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// UnrealizedAnalyticsHandler 未實現損益分析 API Handler
type UnrealizedAnalyticsHandler struct {
	service service.UnrealizedAnalyticsService
}

// NewUnrealizedAnalyticsHandler 建立新的 UnrealizedAnalyticsHandler
func NewUnrealizedAnalyticsHandler(service service.UnrealizedAnalyticsService) *UnrealizedAnalyticsHandler {
	return &UnrealizedAnalyticsHandler{
		service: service,
	}
}

// GetSummary 取得未實現損益摘要
// @Summary 取得未實現損益摘要
// @Description 取得當前所有持倉的未實現損益摘要
// @Tags unrealized-analytics
// @Accept json
// @Produce json
// @Success 200 {object} models.UnrealizedSummary
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/unrealized/summary [get]
func (h *UnrealizedAnalyticsHandler) GetSummary(c *gin.Context) {
	// 呼叫 service
	summary, err := h.service.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_SUMMARY_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: summary,
	})
}

// GetPerformance 取得各資產類型未實現績效
// @Summary 取得各資產類型未實現績效
// @Description 取得當前各資產類型的未實現損益績效
// @Tags unrealized-analytics
// @Accept json
// @Produce json
// @Success 200 {array} models.UnrealizedPerformance
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/unrealized/performance [get]
func (h *UnrealizedAnalyticsHandler) GetPerformance(c *gin.Context) {
	// 呼叫 service
	performance, err := h.service.GetPerformance()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_PERFORMANCE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: performance,
	})
}

// GetTopAssets 取得 Top 未實現損益資產
// @Summary 取得 Top 未實現損益資產
// @Description 取得未實現損益最佳的資產（按未實現損益排序）
// @Tags unrealized-analytics
// @Accept json
// @Produce json
// @Param limit query int false "回傳數量限制" default(10)
// @Success 200 {array} models.UnrealizedTopAsset
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/unrealized/top-assets [get]
func (h *UnrealizedAnalyticsHandler) GetTopAssets(c *gin.Context) {
	// 取得 limit 參數
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // 預設 10 筆
	}

	// 呼叫 service
	topAssets, err := h.service.GetTopAssets(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_TOP_ASSETS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: topAssets,
	})
}

