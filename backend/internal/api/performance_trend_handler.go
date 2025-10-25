package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// PerformanceTrendHandler 績效趨勢 Handler
type PerformanceTrendHandler struct {
	service service.PerformanceTrendService
}

// NewPerformanceTrendHandler 建立績效趨勢 Handler
func NewPerformanceTrendHandler(service service.PerformanceTrendService) *PerformanceTrendHandler {
	return &PerformanceTrendHandler{
		service: service,
	}
}

// CreateDailySnapshot 建立每日績效快照
// @Summary 建立每日績效快照
// @Description 建立當天的績效快照，包含總體和各資產類型的績效指標
// @Tags performance-trends
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=models.DailyPerformanceSnapshot}
// @Failure 500 {object} APIResponse
// @Router /api/performance-trends/snapshot [post]
func (h *PerformanceTrendHandler) CreateDailySnapshot(c *gin.Context) {
	snapshot, err := h.service.CreateDailySnapshot()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "CREATE_SNAPSHOT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: snapshot,
	})
}

// GetTrendByDateRange 取得日期範圍內的績效趨勢
// @Summary 取得日期範圍內的績效趨勢
// @Description 取得指定日期範圍內的績效趨勢資料
// @Tags performance-trends
// @Accept json
// @Produce json
// @Param start_date query string true "起始日期 (YYYY-MM-DD)"
// @Param end_date query string true "結束日期 (YYYY-MM-DD)"
// @Success 200 {object} APIResponse{data=models.PerformanceTrendSummary}
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /api/performance-trends/range [get]
func (h *PerformanceTrendHandler) GetTrendByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_PARAMETERS",
				Message: "start_date and end_date are required",
			},
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_DATE_FORMAT",
				Message: "start_date must be in YYYY-MM-DD format",
			},
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_DATE_FORMAT",
				Message: "end_date must be in YYYY-MM-DD format",
			},
		})
		return
	}

	summary, err := h.service.GetTrendByDateRange(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_TREND_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: summary,
	})
}

// GetLatestTrend 取得最新的績效趨勢
// @Summary 取得最新的績效趨勢
// @Description 取得最新 N 天的績效趨勢資料
// @Tags performance-trends
// @Accept json
// @Produce json
// @Param days query int false "天數" default(30)
// @Success 200 {object} APIResponse{data=[]models.PerformanceTrendPoint}
// @Failure 500 {object} APIResponse
// @Router /api/performance-trends/latest [get]
func (h *PerformanceTrendHandler) GetLatestTrend(c *gin.Context) {
	days := 30
	if daysStr := c.Query("days"); daysStr != "" {
		if parsedDays, err := strconv.Atoi(daysStr); err == nil && parsedDays > 0 {
			days = parsedDays
		}
	}

	data, err := h.service.GetLatestTrend(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_LATEST_TREND_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: data,
	})
}

