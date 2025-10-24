package api

import (
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// AnalyticsHandler 分析 API Handler
type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

// NewAnalyticsHandler 建立新的 AnalyticsHandler
func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

// GetSummary 取得分析摘要
// @Summary 取得分析摘要
// @Description 取得指定時間範圍的已實現損益摘要
// @Tags analytics
// @Accept json
// @Produce json
// @Param time_range query string false "時間範圍 (week, month, quarter, year, all)" default(month)
// @Success 200 {object} models.AnalyticsSummary
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/summary [get]
func (h *AnalyticsHandler) GetSummary(c *gin.Context) {
	// 取得時間範圍參數
	timeRangeStr := c.DefaultQuery("time_range", "month")
	timeRange := models.TimeRange(timeRangeStr)

	// 呼叫 service
	summary, err := h.analyticsService.GetSummary(timeRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_TIME_RANGE",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: summary,
	})
}

// GetPerformance 取得各資產類型績效
// @Summary 取得各資產類型績效
// @Description 取得指定時間範圍內各資產類型的已實現損益績效
// @Tags analytics
// @Accept json
// @Produce json
// @Param time_range query string false "時間範圍 (week, month, quarter, year, all)" default(month)
// @Success 200 {array} models.PerformanceData
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/performance [get]
func (h *AnalyticsHandler) GetPerformance(c *gin.Context) {
	// 取得時間範圍參數
	timeRangeStr := c.DefaultQuery("time_range", "month")
	timeRange := models.TimeRange(timeRangeStr)

	// 呼叫 service
	performance, err := h.analyticsService.GetPerformance(timeRange)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_TIME_RANGE",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: performance,
	})
}

// GetTopAssets 取得最佳/最差表現資產
// @Summary 取得最佳/最差表現資產
// @Description 取得指定時間範圍內表現最佳的資產（按已實現損益排序）
// @Tags analytics
// @Accept json
// @Produce json
// @Param time_range query string false "時間範圍 (week, month, quarter, year, all)" default(month)
// @Param limit query int false "回傳數量限制" default(5)
// @Success 200 {array} models.TopAsset
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/analytics/top-assets [get]
func (h *AnalyticsHandler) GetTopAssets(c *gin.Context) {
	// 取得時間範圍參數
	timeRangeStr := c.DefaultQuery("time_range", "month")
	timeRange := models.TimeRange(timeRangeStr)

	// 取得 limit 參數
	limitStr := c.DefaultQuery("limit", "5")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 5 // 預設 5 筆
	}

	// 呼叫 service
	topAssets, err := h.analyticsService.GetTopAssets(timeRange, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_TIME_RANGE",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: topAssets,
	})
}

