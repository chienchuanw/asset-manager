package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// AssetSnapshotHandler 資產快照 API Handler
type AssetSnapshotHandler struct {
	service service.AssetSnapshotService
}

// NewAssetSnapshotHandler 建立新的資產快照 Handler
func NewAssetSnapshotHandler(service service.AssetSnapshotService) *AssetSnapshotHandler {
	return &AssetSnapshotHandler{
		service: service,
	}
}

// CreateSnapshotRequest 建立快照請求
type CreateSnapshotRequest struct {
	SnapshotDate string                    `json:"snapshot_date" binding:"required"` // 格式: YYYY-MM-DD
	AssetType    models.SnapshotAssetType  `json:"asset_type" binding:"required"`
	ValueTWD     float64                   `json:"value_twd" binding:"required,gte=0"`
}

// GetAssetTrendRequest 取得資產趨勢請求
type GetAssetTrendRequest struct {
	Days      int                       `form:"days" binding:"required,gte=1,lte=365"`
	AssetType models.SnapshotAssetType  `form:"asset_type" binding:"required"`
}

// AssetTrendResponse 資產趨勢回應
type AssetTrendResponse struct {
	Date     string  `json:"date"`      // 日期 (YYYY-MM-DD)
	ValueTWD float64 `json:"value_twd"` // 資產價值 (TWD)
}

// CreateSnapshot 建立資產快照
// @Summary 建立資產快照
// @Tags snapshots
// @Accept json
// @Produce json
// @Param request body CreateSnapshotRequest true "建立快照請求"
// @Success 201 {object} APIResponse{data=models.AssetSnapshot}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/snapshots [post]
func (h *AssetSnapshotHandler) CreateSnapshot(c *gin.Context) {
	var req CreateSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// 解析日期
	snapshotDate, err := time.Parse("2006-01-02", req.SnapshotDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_DATE_FORMAT",
				Message: "snapshot_date must be in YYYY-MM-DD format",
			},
		})
		return
	}

	// 建立快照
	input := &models.CreateAssetSnapshotInput{
		SnapshotDate: snapshotDate,
		AssetType:    req.AssetType,
		ValueTWD:     req.ValueTWD,
	}

	snapshot, err := h.service.CreateSnapshot(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "CREATE_SNAPSHOT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Data: snapshot,
	})
}

// GetAssetTrend 取得資產價值趨勢
// @Summary 取得資產價值趨勢
// @Tags snapshots
// @Accept json
// @Produce json
// @Param days query int true "天數" minimum(1) maximum(365)
// @Param asset_type query string true "資產類型" Enums(total, tw-stock, us-stock, crypto)
// @Success 200 {object} APIResponse{data=[]AssetTrendResponse}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/snapshots/trend [get]
func (h *AssetSnapshotHandler) GetAssetTrend(c *gin.Context) {
	var req GetAssetTrendRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		})
		return
	}

	// 計算日期範圍
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.Add(-time.Duration(req.Days-1) * 24 * time.Hour)

	// 取得快照列表
	snapshots, err := h.service.GetSnapshotsByDateRange(startDate, endDate, req.AssetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_SNAPSHOTS_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	// 轉換為回應格式
	trendData := make([]AssetTrendResponse, 0, len(snapshots))
	for _, snapshot := range snapshots {
		trendData = append(trendData, AssetTrendResponse{
			Date:     snapshot.SnapshotDate.Format("2006-01-02"),
			ValueTWD: snapshot.ValueTWD,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: trendData,
	})
}

// GetLatestSnapshot 取得最新快照
// @Summary 取得最新快照
// @Tags snapshots
// @Accept json
// @Produce json
// @Param asset_type query string true "資產類型" Enums(total, tw-stock, us-stock, crypto)
// @Success 200 {object} APIResponse{data=models.AssetSnapshot}
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /api/snapshots/latest [get]
func (h *AssetSnapshotHandler) GetLatestSnapshot(c *gin.Context) {
	assetTypeStr := c.Query("asset_type")
	if assetTypeStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "asset_type is required",
			},
		})
		return
	}

	assetType := models.SnapshotAssetType(assetTypeStr)

	snapshot, err := h.service.GetLatestSnapshot(assetType)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Error: &APIError{
				Code:    "SNAPSHOT_NOT_FOUND",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: snapshot,
	})
}

// UpdateSnapshot 更新快照
// @Summary 更新快照
// @Tags snapshots
// @Accept json
// @Produce json
// @Param date query string true "日期 (YYYY-MM-DD)"
// @Param asset_type query string true "資產類型"
// @Param value_twd query number true "資產價值 (TWD)"
// @Success 200 {object} APIResponse{data=models.AssetSnapshot}
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/snapshots [put]
func (h *AssetSnapshotHandler) UpdateSnapshot(c *gin.Context) {
	dateStr := c.Query("date")
	assetTypeStr := c.Query("asset_type")
	valueTWDStr := c.Query("value_twd")

	if dateStr == "" || assetTypeStr == "" || valueTWDStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "date, asset_type, and value_twd are required",
			},
		})
		return
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_DATE_FORMAT",
				Message: "date must be in YYYY-MM-DD format",
			},
		})
		return
	}

	// 解析金額
	valueTWD, err := strconv.ParseFloat(valueTWDStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_VALUE",
				Message: "value_twd must be a valid number",
			},
		})
		return
	}

	assetType := models.SnapshotAssetType(assetTypeStr)

	snapshot, err := h.service.UpdateSnapshot(date, assetType, valueTWD)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "UPDATE_SNAPSHOT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: snapshot,
	})
}

// DeleteSnapshot 刪除快照
// @Summary 刪除快照
// @Tags snapshots
// @Accept json
// @Produce json
// @Param date query string true "日期 (YYYY-MM-DD)"
// @Param asset_type query string true "資產類型"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /api/snapshots [delete]
func (h *AssetSnapshotHandler) DeleteSnapshot(c *gin.Context) {
	dateStr := c.Query("date")
	assetTypeStr := c.Query("asset_type")

	if dateStr == "" || assetTypeStr == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_REQUEST",
				Message: "date and asset_type are required",
			},
		})
		return
	}

	// 解析日期
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Error: &APIError{
				Code:    "INVALID_DATE_FORMAT",
				Message: "date must be in YYYY-MM-DD format",
			},
		})
		return
	}

	assetType := models.SnapshotAssetType(assetTypeStr)

	err = h.service.DeleteSnapshot(date, assetType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "DELETE_SNAPSHOT_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: map[string]string{
			"message": "Snapshot deleted successfully",
		},
	})
}

// TriggerDailySnapshots 手動觸發每日快照建立（用於測試或手動執行）
// @Summary 手動觸發每日快照建立
// @Description 立即執行每日快照建立任務，計算當前所有持倉的價值並建立快照
// @Tags snapshots
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse "成功建立快照"
// @Failure 500 {object} APIResponse "伺服器錯誤"
// @Router /api/snapshots/trigger [post]
func (h *AssetSnapshotHandler) TriggerDailySnapshots(c *gin.Context) {
	// 呼叫 Service 層
	if err := h.service.CreateDailySnapshots(); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "INTERNAL_ERROR",
				Message: err.Error(),
			},
		})
		return
	}

	// 返回成功結果
	c.JSON(http.StatusOK, APIResponse{
		Data: map[string]string{
			"message": "Daily snapshots created successfully",
		},
	})
}
