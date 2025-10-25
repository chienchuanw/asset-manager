package api

import (
	"net/http"
	"strconv"

	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
)

// AllocationHandler 資產配置 Handler
type AllocationHandler struct {
	service service.AllocationService
}

// NewAllocationHandler 建立資產配置 Handler
func NewAllocationHandler(service service.AllocationService) *AllocationHandler {
	return &AllocationHandler{
		service: service,
	}
}

// GetCurrentAllocation 取得當前資產配置摘要
// @Summary 取得當前資產配置摘要
// @Description 取得當前所有持倉的資產配置摘要，包含按資產類型和個別資產的分類
// @Tags allocation
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=models.AllocationSummary}
// @Failure 500 {object} APIResponse
// @Router /api/allocation/current [get]
func (h *AllocationHandler) GetCurrentAllocation(c *gin.Context) {
	summary, err := h.service.GetCurrentAllocation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_CURRENT_ALLOCATION_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: summary,
	})
}

// GetAllocationByType 取得按資產類型的配置
// @Summary 取得按資產類型的配置
// @Description 取得按資產類型分類的資產配置
// @Tags allocation
// @Accept json
// @Produce json
// @Success 200 {object} APIResponse{data=[]models.AllocationByType}
// @Failure 500 {object} APIResponse
// @Router /api/allocation/by-type [get]
func (h *AllocationHandler) GetAllocationByType(c *gin.Context) {
	allocations, err := h.service.GetAllocationByType()
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_ALLOCATION_BY_TYPE_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: allocations,
	})
}

// GetAllocationByAsset 取得按個別資產的配置
// @Summary 取得按個別資產的配置
// @Description 取得按個別資產分類的資產配置
// @Tags allocation
// @Accept json
// @Produce json
// @Param limit query int false "回傳數量限制" default(20)
// @Success 200 {object} APIResponse{data=[]models.AllocationByAsset}
// @Failure 500 {object} APIResponse
// @Router /api/allocation/by-asset [get]
func (h *AllocationHandler) GetAllocationByAsset(c *gin.Context) {
	// 取得 limit 參數，預設為 20
	limit := 20
	if limitStr := c.Query("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	allocations, err := h.service.GetAllocationByAsset(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Error: &APIError{
				Code:    "GET_ALLOCATION_BY_ASSET_FAILED",
				Message: err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Data: allocations,
	})
}

